package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

const (
	timeout     = time.Second * 5
	healthzPath = "/healthz"
	defaultPort = "8080"
	maxRPS      = time.Millisecond * 100 // 10 rps

	requestID = "req-correlation-id"
)

// set at compile time
var (
	buildVersion = "unknown"
	buildCommit  = "unknown"
)

func main() {
	// print version information
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("version: %s\n", buildVersion)
		fmt.Printf("commit: %s\n", buildCommit)
		os.Exit(0)
	}

	var (
		logger *zap.Logger
		err    error
	)
	jsonCfg := os.Getenv("ZAP_CONFIG")

	if jsonCfg != "" {
		var cfg zap.Config
		b := []byte(jsonCfg)

		err = json.Unmarshal(b, &cfg)
		if err != nil {
			panic(fmt.Errorf("unmarshal ZAP JSON config: %v", err).Error())
		}
		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}

	ctx := signals.NewContext()
	ctx = logging.WithLogger(ctx, logger.Sugar().Named("ci-demo-app").With("commit", buildCommit, "version", buildVersion))

	if err = run(ctx); !errors.Is(err, http.ErrServerClosed) {
		logging.FromContext(ctx).Fatalf("run server: %v", err)
	}
}

func run(ctx context.Context) error {
	srv := newServer(ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		<-ctx.Done()
		logging.FromContext(ctx).Info("got signal, attempting graceful shutdown")
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		return srv.Shutdown(timeoutCtx)
	})

	eg.Go(func() error {
		logging.FromContext(ctx).Infow("running server", "address", srv.Addr)
		return srv.ListenAndServe()
	})

	return eg.Wait()
}

func newServer(ctx context.Context) *http.Server {
	rl := rate.NewLimiter(rate.Every(maxRPS), 10)

	mux := http.NewServeMux()
	mux.HandleFunc(healthzPath, requestLogger(ctx, healthZHandler(ctx)))
	mux.HandleFunc("/", requestLogger(ctx, rateLimiter(ctx, rl, greeterHandler(ctx))))

	port := getPort()
	addr := fmt.Sprintf(":%s", port)
	srv := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return &srv
}

func getPort() string {
	// Knative injected PORT
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return defaultPort
}

func rateLimiter(ctx context.Context, rl *rate.Limiter, next http.HandlerFunc) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if !rl.Allow() {
			var id string
			id = req.Header.Get(requestID)

			if id == "" {
				id = "undefined"
			}

			logging.FromContext(ctx).With(zap.String("id", id)).Debug("rate limited")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next(w, req)
	}
}

func requestLogger(ctx context.Context, next http.HandlerFunc) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		// inject correlation ID
		id := uuid.New().String()
		req.Header.Del(requestID) // clear if exist
		req.Header.Add(requestID, id)

		logging.FromContext(ctx).With(zap.String("id", id)).Debugw("new request", "method", req.Method, "path", html.EscapeString(req.URL.Path), "client", req.RemoteAddr)
		next(w, req)
	}
}

func healthZHandler(ctx context.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			logging.FromContext(ctx).Errorf("write response: %v", err)
		}
	}
}

func greeterHandler(ctx context.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := "Stranger"
		if param := req.URL.Query().Get("name"); param != "" {
			// https://codeql.github.com/codeql-query-help/go/go-reflected-xss/
			name = html.EscapeString(param)
		}

		_, err := w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
		if err != nil {
			logging.FromContext(ctx).Errorf("write response: %v", err)
		}
	}
}
