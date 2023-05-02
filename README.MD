# Policy as code Demo 💪
Demo showcases Scribe Policy As Code capabilities.

## Policy Modules
Policy providers can provide OPA based Policy modules.
Modules are signed by provider and can be verified by consumers.
Using the policy modules, consumers can generate compliance reports fitted to the providers regulations.

See https://github.com/scribe-security/policy-modules repo for details.

## Policies
* .valint.yaml - Policies select the modules and configuration required by a specific report.
* .valint_provider - Policy used to verify the Policy bundle trust (who is the provider).

## Evidence collection
* Checkout stage:
  * Sign sources
  * Collect source SBOM
* Image Release stage:
  * Sign image
  * Collect image SLSA prov
## Evaluation
* Check release tagged sources are signed.
* Check release image is signed.
* SBOM collected on sources.
* SLSA prov collected on img.
* Check Code owners (Who is modifying what part of the SCM)

# Demo app💡

This repository shows an end-to-end release management approach for
[Go](https://golang.org/) applications using [Github
Actions](https://github.com/features/actions) as a pipeline for **continuous
integration** (CI), [`goreleaser`](https://goreleaser.com/) for **artifact and
release management**, [`ko`](https://github.com/google/ko) to build OCI-compliant
**container images**, and Kubernetes in Docker
([`kind`](https://github.com/kubernetes-sigs/kind)) for **end-to-end testing**.

## Project Status 🚦

[![Build](https://github.com/scribe-security/ci-demo-app/actions/workflows/build.yaml/badge.svg)](https://github.com/scribe-security/ci-demo-app/actions/workflows/build.yaml)
[![Tests](https://github.com/scribe-security/ci-demo-app/actions/workflows/go-unit-tests.yaml/badge.svg)](https://github.com/scribe-security/ci-demo-app/actions/workflows/go-unit-tests.yaml)
[![Latest Release](https://img.shields.io/github/release/scribe-security/ci-demo-app.svg?logo=github&style=flat-square)](https://github.com/scribe-security/ci-demo-app/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/scribe-security/ci-demo-app)](https://goreportcard.com/report/github.com/scribe-security/ci-demo-app)
[![codecov](https://codecov.io/gh/scribe-security/ci-demo-app/branch/main/graph/badge.svg?token=TC7MW723JO)](https://codecov.io/gh/scribe-security/ci-demo-app)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/scribe-security/ci-demo-app)](https://github.com/scribe-security/ci-demo-app)

## Features 🤩

Based on a trivial Go demo application this repository was built to show the
power of Github Actions and related Go tools with the following highlights:

💫 &nbsp; Add Continuous Integration and **end-to-end release automation** with
[Github Actions](https://github.com/scribe-security/ci-demo-app/actions) ("workflows")  
💫 &nbsp; Add [**status badges**](#project-status-) to your project `README.md`  
💫 &nbsp; **Manage your repository** (issues) and contributions (PRs) with
automatic greetings and code quality checks  
💫 &nbsp; **Produce [artifacts](#artifacts-)** for different platforms and software management
tools, e.g. Homebrew, Snaps, etc.  
💫 &nbsp; **Automate** Github Releases incl. a (templatized) [`CHANGELOG.md`](CHANGELOG.md)  
💫 &nbsp; Use `ko` to **not worry** about building container images and
`Dockerfiles`  
💫 &nbsp; Use `ko` with `kind` for **super fast** development and end-to-end
testing against real APIs (Kubernetes)  

## Artifacts 💎

You can grab the resulting artifacts produced by the Github release
[workflow](./.github/workflows/release.yaml) from the
[releases](https://github.com/scribe-security/ci-demo-app/releases) page.

### Docker 🐳

Multi-arch (platform) Docker images are pushed to Scribe Artifactory as part of the
release pipeline. 
Run the app via `docker`:

```console
$ # expose default port 8080 to 8080 on the Docker host
$ # interrupt with CTRL-C
$ docker run -it --rm -p 8080:8080 scribesecuriy.jfrog.io/scribe-docker-public-local/ci-demo-app:latest
2021-05-04T14:31:00.658Z      INFO    ci-demo-app     ci-demo-app/main.go:85  running server  {"commit": "e616df52", "address": ":8080"}
2021-05-04T14:31:02.933Z      INFO    ci-demo-app     ci-demo-app/main.go:81  got signal, attempting graceful shutdown        {"commit": "e616df52"}
```




















































