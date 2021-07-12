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

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

const (
	// http timeouts
	timeout     = time.Second * 5
	healthzPath = "/healthz"
	defaultPort = "8080"
)

// set at compile time
var (
	buildVersion = "unknown"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

func main() {
	// print version information
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("version: %s\n", buildVersion)
		fmt.Printf("commit: %s\n", buildCommit)
		fmt.Printf("date: %s\n", buildDate)
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
	ctx = logging.WithLogger(ctx, logger.Sugar().Named("ci-demo-app").With("commit", buildCommit))

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
	mux := http.NewServeMux()
	mux.HandleFunc(healthzPath, requestLogger(ctx, healthZHandler(ctx)))
	mux.HandleFunc("/", requestLogger(ctx, greeterHandler(ctx)))

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
	} else {
		return defaultPort
	}
}

func requestLogger(ctx context.Context, next http.HandlerFunc) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logging.FromContext(ctx).Debugw("new request", "method", req.Method, "path", html.EscapeString(req.URL.Path), "client", req.RemoteAddr)
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
