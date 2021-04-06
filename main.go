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
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

const (
	// http timeouts
	timeout = time.Second * 3
)

// set at compile time
var (
	buildVersion = "unknown"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

var (
	defaultPort = "8080"
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

	if err := run(ctx); !errors.Is(err, http.ErrServerClosed) {
		logging.FromContext(ctx).Fatalf("run server: %v", err)
	}
}

func run(ctx context.Context) error {
	srv := newServer(ctx)

	go func() {
		<-ctx.Done()
		logging.FromContext(ctx).Info("got signal, attempting graceful shutdown")
		_ = srv.Shutdown(ctx)
	}()

	logging.FromContext(ctx).Infow("running server", "address", srv.Addr)
	return srv.ListenAndServe()
}

func newServer(ctx context.Context) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", requestLogger(ctx, healthZHandler(ctx)))

	// Knative Serving injects PORT
	var port string
	if p := os.Getenv("PORT"); p != "" {
		port = p
	} else {
		port = defaultPort
	}

	addr := fmt.Sprintf(":%s", port)
	srv := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return &srv
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
