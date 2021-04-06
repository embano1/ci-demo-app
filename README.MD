[![Build](https://github.com/embano1/ci-demo-app/actions/workflows/build.yaml/badge.svg)](https://github.com/embano1/ci-demo-app/actions/workflows/build.yaml)
[![Tests](https://github.com/embano1/ci-demo-app/actions/workflows/go-unit-tests.yaml/badge.svg)](https://github.com/embano1/ci-demo-app/actions/workflows/go-unit-tests.yaml)
[![Latest Release](https://img.shields.io/github/release/embano1/ci-demo-app.svg?logo=github&style=flat-square)](https://github.com/embano1/ci-demo-app/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/embano1/ci-demo-app)](https://goreportcard.com/report/github.com/embano1/ci-demo-app)
[![codecov](https://codecov.io/gh/embano1/ci-demo-app/branch/main/graph/badge.svg?token=TC7MW723JO)](https://codecov.io/gh/embano1/ci-demo-app)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/embano1/ci-demo-app)](https://github.com/embano1/ci-demo-app)

# About

TODO




# TODO

## Custom Logging

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: demo-app
spec:
  template:
    spec:
      containers:
        - image: ko://github.com/embano1/ci-demo-app
          readinessProbe:
            httpGet:
              path: /healthz
          env:
            - name: ZAP_CONFIG
              value: '{"level":"info","encoding":"json","outputPaths": ["stdout"],"errorOutputPaths": ["stderr"],"initialFields":{"app":"demo-app","logConfig":"custom"},"encoderConfig":{"messageKey":"message","levelKey":"level","levelEncoder":"lowercase"}}'
```

💡 **Note:** For complex logging configurations use a `ConfigMap` and project its data into the manifest using `valueFrom` semantics as described [here](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#define-container-environment-variables-using-configmap-data).


## E2E


kind create cluster --name e2e

export KIND_CLUSTER_NAME=e2e
export KO_DOCKER_REPO=kind.local
export KNATIVE_VERSION=v0.22.0


	KubeConfig string `envconfig:"KUBECONFIG"`
	Namespace  string `default:"default"`
	DockerRepo string `default:"kind.local"`
	DockerTag  string `default:"latest"`
	URLPath    string `default:"/healthz"`
	ServerName string `default:"demo-app"`

kubectl apply -f https://github.com/knative/serving/releases/download/${KNATIVE_VERSION}/serving-crds.yaml
kubectl wait --timeout=20s --for=condition=Established crd/services.serving.knative.dev

kubectl apply -f https://github.com/knative/serving/releases/download/${KNATIVE_VERSION}/serving-core.yaml
kubectl -n knative-serving wait --timeout=2m --for=condition=Available deploy/webhook

kubectl apply -f https://github.com/knative/net-kourier/releases/download/${KNATIVE_VERSION}/kourier.yaml
kubectl -n kourier-system wait --timeout=2m --for=condition=Available deploy/3scale-kourier-gateway

kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'

./test/upload_test_images.sh

ko resolve -BR -f config/ | k apply -f -
kubectl wait --timeout=1m --for=condition=Ready ksvc/demo-app

go test -timeout 1m -v -tags=e2e -count 1 -race ./test/...



# Out of Scope

Issue/PR templates, CONTRIB and community management
Advanced 3rd party integration, e.g. Slack, bots
External/Workflow Dispatch triggers
advanced branch protection and RBAC
Security and repo hardening guidelines
