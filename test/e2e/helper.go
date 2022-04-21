//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"knative.dev/pkg/ptr"

	"github.com/embano1/ci-demo-app/test"
)

const (
	clientImageName = "http" // name of test http getter folder
	pollInterval    = 2 * time.Second
	pollTimeout     = 30 * time.Second
	letterBytes     = "abcdefghijklmnopqrstuvwxyz"
	randSuffixLen   = 8
	sep             = "-"
)

type (
	waitFunc    func()
	cleanupFunc func()
)

func createHttpClient(t *testing.T, client *kubernetes.Clientset) (waitFunc, cleanupFunc) {
	t.Helper()

	var env test.Environment
	if err := envconfig.Process("", &env); err != nil {
		t.Fatal(err)
	}

	name := fmt.Sprintf("%s%s%s", clientImageName, sep, randomString())
	job := v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: env.Namespace,
		},
		Spec: v1.JobSpec{
			Completions: ptr.Int32(1),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            clientImageName,
							Image:           imageName(t, clientImageName),
							Args:            []string{"-url", serverURL(t, env.ServerName)},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			TTLSecondsAfterFinished: nil,
		},
	}

	ctx := context.TODO()
	t.Logf("creating job %q", name)
	_, err := client.BatchV1().Jobs(env.Namespace).Create(ctx, &job, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	waiter := func() {
		// Wait for the Job to report a successful execution.
		waitErr := wait.PollImmediate(pollInterval, pollTimeout, func() (bool, error) {
			j, err := client.BatchV1().Jobs(env.Namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					t.Logf("job not found: %v", err)
					return false, nil
				}
				return true, err
			}

			t.Logf("Active=%d, Failed=%d, Succeeded=%d", j.Status.Active, j.Status.Failed, j.Status.Succeeded)

			// Check for successful completions.
			return j.Status.Succeeded > 0, nil
		})
		if waitErr != nil {
			t.Fatalf("waiting for Job to complete successfully: %v", waitErr)
		}
	}

	cleanup := func() {
		if err := client.BatchV1().Jobs(env.Namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			t.Errorf("cleanup job: %v", err)
		}
	}

	return waiter, cleanup
}

func imageName(t *testing.T, name string) string {
	t.Helper()

	var env test.Environment
	if err := envconfig.Process("", &env); err != nil {
		t.Fatal(err)
	}

	return fmt.Sprintf("%s/%s:%s", env.DockerRepo, name, env.DockerTag)
}

func serverURL(t *testing.T, name string) string {
	t.Helper()

	var env test.Environment
	if err := envconfig.Process("", &env); err != nil {
		t.Fatal(err)
	}

	return fmt.Sprintf("http://%s.%s.svc.cluster.local%s", name, env.Namespace, env.URLPath)
}

// randomString will generate a random string.
func randomString() string {
	rand.Seed(time.Now().UnixNano())
	suffix := make([]byte, randSuffixLen)
	for i := range suffix {
		suffix[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(suffix)
}
