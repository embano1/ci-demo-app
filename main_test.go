package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
	"gotest.tools/assert"
	"knative.dev/pkg/logging"
)

func TestServer(t *testing.T) {
	type response struct {
		code int
		body []byte
	}
	tests := []struct {
		name   string
		method string
		path   string
		want   response
	}{
		{
			name:   "200 on /healthz",
			method: http.MethodGet,
			path:   "/healthz",
			want:   response{code: 200, body: []byte(`{"status":"ok"}`)},
		},
		{
			name:   "404 on root",
			method: http.MethodGet,
			path:   "/",
			want:   response{code: 404, body: []byte("404 page not found\n")},
		},
	}
	for _, tt := range tests {
		ctx := logging.WithLogger(context.Background(), zaptest.NewLogger(t).Sugar())

		t.Run(tt.name, func(t *testing.T) {
			srv := newServer(ctx)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			srv.Handler.ServeHTTP(rec, req)

			res := rec.Result()
			b, err := io.ReadAll(res.Body)
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.NilError(t, err)
			assert.Equal(t, string(b), string(tt.want.body))
		})
	}
}
