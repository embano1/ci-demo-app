// +build e2e

package e2e

import (
	"testing"

	"github.com/embano1/ci-demo-app/test"
)

func TestServer(t *testing.T) {
	client, err := test.GetKubeClient(t)
	if err != nil {
		t.Fatalf("get kube client: %v", err)
	}

	awaitJob, deleteFn := createHttpClient(t, client)
	t.Cleanup(func() {
		deleteFn()
	})

	// wait for http client to finish
	awaitJob()
}
