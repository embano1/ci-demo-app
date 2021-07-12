package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
			name:   "no name param set on /",
			method: http.MethodGet,
			path:   "/",
			want:   response{code: 200, body: []byte("Hello Stranger!")},
		},
		{
			name:   "name param set on /",
			method: http.MethodGet,
			path:   "/?name=Michael",
			want:   response{code: 200, body: []byte("Hello Michael!")},
		},
		{
			name:   "prevent XSS",
			method: http.MethodGet,
			path:   "/?name=<script>alert('xss')</script>",
			want:   response{code: 200, body: []byte("Hello &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;!")},
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

func Test_getPort(t *testing.T) {
	type fields struct {
		osEnvPort string
	}
	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name:   "env PORT not set",
			want:   defaultPort,
			fields: fields{},
		},
		{
			name: "env PORT set to 8081",
			want: "8081",
			fields: fields{
				osEnvPort: "8081",
			},
		},
	}
	for _, tt := range tests {
		if f := tt.fields.osEnvPort; f != "" {
			err := os.Setenv("PORT", f)
			assert.NilError(t, err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got := getPort()
			assert.Equal(t, got, tt.want)
		})
	}
}
