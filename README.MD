## Project Status 🚦

[![Build](https://github.com/embano1/ci-demo-app/actions/workflows/build.yaml/badge.svg)](https://github.com/embano1/ci-demo-app/actions/workflows/build.yaml)
[![Tests](https://github.com/embano1/ci-demo-app/actions/workflows/go-unit-tests.yaml/badge.svg)](https://github.com/embano1/ci-demo-app/actions/workflows/go-unit-tests.yaml)
[![Latest Release](https://img.shields.io/github/release/embano1/ci-demo-app.svg?logo=github&style=flat-square)](https://github.com/embano1/ci-demo-app/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/embano1/ci-demo-app)](https://goreportcard.com/report/github.com/embano1/ci-demo-app)
[![codecov](https://codecov.io/gh/embano1/ci-demo-app/branch/main/graph/badge.svg?token=TC7MW723JO)](https://codecov.io/gh/embano1/ci-demo-app)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/embano1/ci-demo-app)](https://github.com/embano1/ci-demo-app)


# Motivation 💪

I created this repository for learning and education purposes to help everyone
to level up their (open source) projects with good documentation and
code/release automation. I spend a lot of time learning and experimenting with
different tools and integrations to manage the projects I'm working, thus trying
to give some inspiration to save you precious time 😀

Of course there's many ways to automate the release process for your (Go)
applications, so consider the examples shared in this repository **ONE** but not
**THE** way.

Even though this repository puts a focus on building and releasing Go
applications, the Github Actions workflows and explanations provided apply to
any language/environment.

💡 **Note:** Check out this online `README` [generator](https://readme.so/) to
bootstrap a nice landing page for your users. 

# About 💡

This repository shows an end-to-end release management approach for
[Go](https://golang.org/) applications using [Github
Actions](https://github.com/features/actions) as a pipeline for **continuous
integration** (CI), [`goreleaser`](https://goreleaser.com/) for **artifact and
release management**, [`ko`](https://github.com/google/ko) to build OCI-compliant
**container images**, and Kubernetes in Docker
([`kind`](https://github.com/kubernetes-sigs/kind)) for **end-to-end testing**.

Another goodie: depending on your repository settings and usage, you are likely
to stay within the generous Github Actions [**free
tier**](https://docs.github.com/en/github/setting-up-and-managing-billing-and-payments-on-github/about-billing-for-github-actions).

Github Actions is also deeply integrated into [Github
Packages](https://docs.github.com/en/packages/learn-github-packages/introduction-to-github-packages),
a software package hosting service that allows you to **host your software
packages**, e.g. container images.

Even if you don't use Github Actions (yet), I hope the repository still provides
value since the tools used here can be used standalone, too.

## Features 🤩

Based on a trivial Go demo application this repository was built to show the
power of Github Actions and related Go tools with the following highlights:

💫 &nbsp; Add Continuous Integration and **end-to-end release automation** with
[Github Actions](https://github.com/embano1/ci-demo-app/actions) ("workflows")  
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
[releases](https://github.com/embano1/ci-demo-app/releases) page.

### Homebrew 🍺

As part of this project, a Homebrew formula is maintained and continuously
updated as part of the release pipeline. Detailed installation instructions can
be found in the dedicated
[repository](https://github.com/embano1/homebrew-ci-demo-app).

Quick install via `brew`:

```console
brew install embano1/ci-demo-app/demo-app
```

### Container Images 🐳

Multi-arch (platform) OCI-compliant container images are pushed to Github
Container Registry ([Github Packages](https://github.com/features/packages)) as
part of the release pipeline. Images are also publicly signed using `cosign` and
the [sigstore](https://www.sigstore.dev/) project.

Visit the [package
page](https://github.com/embano1/ci-demo-app/pkgs/container/ci-demo-app) for
details.

Example how to run the demo app via `docker`:

```console
# print version information
docker run ghcr.io/embano1/ci-demo-app:latest version
version: v0.4.0
commit: c59a1872

# expose default port 8080 to 8080 on the Docker host
# interrupt with CTRL-C
docker run -p 8080:8080 ghcr.io/embano1/ci-demo-app:latest
2022-04-22T08:34:09.432Z        INFO    ci-demo-app     ci-demo-app/main.go:86  running server  {"commit": "c59a1872", "version": "v0.4.0", "address": ":8080"}
2022-04-22T08:34:09.561Z        DEBUG   ci-demo-app     ci-demo-app/main.go:119 new request     {"commit": "c59a1872", "version": "v0.4.0", "method": "GET", "path": "/api", "client": "172.17.0.1:63376"}
2022-04-22T08:34:15.841Z        INFO    ci-demo-app     ci-demo-app/main.go:79  got signal, attempting graceful shutdown     {"commit": "c59a1872", "version": "v0.4.0"}
```

### Source Tarballs 👾
See [Github Release](https://github.com/embano1/ci-demo-app/releases) page. A
full [`CHANGELOG`](./CHANGELOG.md) is also provided.

# Navigating this Repository 🧭

If you want to see the Github Action workflows used in this repository *in
action* please see this [document](./ACTIONS.md).

If you are interested in the individual tools used here or you want to learn
more about Github Actions, please see this [document](DETAILS.md).

# Out of Scope (for now) ❌

- Golang deep dive, advanced testing, etc. (this repo isn't about teaching Go)
- Advanced 3rd party integration, e.g. Slack, bots to manage issues, comments,
  etc.
- External (manual) workflow dispatch triggers
- Branch protection and role based access control
- Organization-wide workflows and templates
- Detailed security and repository hardening and guidelines (see [notes](DETAILS.md) here)
- Detailed guidelines for writing your custom Github Actions
- Creating/distributing software artifacts other than Go binaries/Homebrew
  formulas (see `goreleaser` section in [DETAILS](./DETAILS.md) to get some
  inspiration)

# Credits 💐

A BIG THANK YOU to the [Knative](https://knative.dev/) community for their endless
inspiration 🙏