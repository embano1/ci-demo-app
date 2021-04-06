package main

import (
	"encoding/json"
	"flag"
	"io"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

const (
	timeout = time.Second * 30
)

var (
	url string
)

func main() {
	flag.StringVar(&url, "url", "http://demo-app.default.svc.cluster.local", "target URL")
	flag.Parse()

	ctx := signals.NewContext()

	// don't block test code too long
	go func() {
		<-time.After(timeout)
		logging.FromContext(ctx).Fatal("timed out waiting for successful response")
	}()

	resp, err := retryablehttp.Get(url)
	if err != nil {
		logging.FromContext(ctx).Fatalf("get response: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	if err != nil {
		logging.FromContext(ctx).Fatalf("read body: %v", err)
	}

	logging.FromContext(ctx).Infof("received response: %s", string(body))
	var response struct {
		Status string `json:"status"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logging.FromContext(ctx).Fatalf("unmarshal body: %v", err)
	}

	if response.Status != "ok" {
		logging.FromContext(ctx).Fatalf("unexpected status received: %+v", response)
	}

	logging.FromContext(ctx).Infof("received successful response: %+v", response)
}
