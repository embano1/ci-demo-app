//go:build e2e
// +build e2e

package test

import (
	"path/filepath"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Environment struct {
	KubeConfig string `envconfig:"KUBECONFIG"`
	Namespace  string `default:"default"`
	DockerRepo string `default:"kind.local"`
	DockerTag  string `default:"latest"`
	URLPath    string `default:"/healthz"`
	ServerName string `default:"demo-app"`
}

func GetKubeClient(t *testing.T) (*kubernetes.Clientset, error) {
	t.Helper()

	var env Environment
	if err := envconfig.Process("", &env); err != nil {
		t.Fatal(err)
	}

	var kubeconfig string
	if env.KubeConfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			t.Fatal("kubernetes configuration not found")
		}
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	return clientset, nil
}
