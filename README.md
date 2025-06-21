# Golang Kubernetes Controller Tutorial

This project is a step-by-step tutorial for DevOps and SRE engineers to learn about building Golang applications and Kubernetes controllers. Each step is implemented as a feature branch and includes a README section with explanations and command history.

## Roadmap

This tutorial guides you through building a production-grade Golang Kubernetes controller, step by step:

1. **Golang CLI Application using Cobra** — Initialize a CLI app with cobra-cli.
2. **Zerolog for Log Levels** — Add structured logging with zerolog.
3. **pflag for Log Level Flags** — Integrate pflag for CLI log level flags.
4. **FastHTTP Server Command** — Add a server command with configurable port and log level.
5. **Makefile, Dockerfile, and GitHub Workflow** — Introduce build automation, secure containerization, CI/CD, and tests.
6. **List Kubernetes Deployments with client-go** — List deployments in the default namespace.
7. **Deployment Informer with client-go** — Watch and log Deployment events.
8. **/deployments JSON API Endpoint** — Serve deployment names as JSON from the informer cache.
9. **controller-runtime Deployment Controller** — Reconcile Deployments and log events.
10. **Leader Election and Metrics** — Add HA and metrics endpoint to the controller manager.
11. **FrontendPage CRD and Advanced Controller** — Define a custom resource and manage Deployments/ConfigMaps.
12. **Platform API (CRUD + Swagger)** — Add RESTful CRUD API and Swagger UI.
13. **MCP Integration** — Integrate MCP server for multi-cluster management.
14. **JWT Authentication** — Secure API endpoints with JWT.
15. **OpenTelemetry Instrumentation** — Add distributed tracing with OpenTelemetry.

---
![logo](docs/img/ctrllogo.png)

## Dev Environment (GitHub Codespaces)
### Recommended way:
 - 1. develop and test your code with envtest k8s control plane first (we use it on step7)
 - 2. test controller with local kind cluster (find example)
 - 3. deploy controller as helm chart in production cluster (use helm chart from step10)

```sh
#### Install brew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
### Install kind
brew install kind
### Install K9S for cluster management
curl -sS https://webi.sh/k9s | sh
### Create kind cluster
kind create cluster
### Add some usefull aliases
alias kk="EDITOR='code --wait' k9s"
alias k=kubectl
## Run complete controller setup
 go run main.go
```

## Progress Status

### ✅ Completed Steps

- [x] **Step 1: Golang CLI Application using Cobra** — Initialize a CLI app with cobra-cli
  * [x] Set up basic Cobra CLI structure
  * [x] Create go-basic command with Kubernetes struct
  * [x] Add version flag (--version, -v)
  * [x] Create status command with cluster information
  * [x] Update project documentation

### 🔄 In Progress

- [ ] **Step 2: Zerolog for Log Levels** — Add structured logging with zerolog
  * [ ] Replace fmt.Println with structured logging
  * [ ] Configure different log levels (trace, debug, info, warn, error)
  * [ ] Add JSON log formatting

### 📋 Todo

- [ ] **Step 3: pflag for Log Level Flags** — Integrate pflag for CLI log level flags
- [ ] **Step 4: FastHTTP Server Command** — Add a server command with configurable port and log level
- [ ] **Step 5: Makefile, Dockerfile, and GitHub Workflow** — Introduce build automation, secure containerization, CI/CD, and tests
- [ ] **Step 6: List Kubernetes Deployments with client-go** — List deployments in the default namespace
- [ ] **Step 7: Deployment Informer with client-go** — Watch and log Deployment events
- [ ] **Step 8: /deployments JSON API Endpoint** — Serve deployment names as JSON from the informer cache
- [ ] **Step 9: controller-runtime Deployment Controller** — Reconcile Deployments and log events
- [ ] **Step 10: Leader Election and Metrics** — Add HA and metrics endpoint to the controller manager
- [ ] **Step 11: FrontendPage CRD and Advanced Controller** — Define a custom resource and manage Deployments/ConfigMaps
- [ ] **Step 12: Platform API (CRUD + Swagger)** — Add RESTful CRUD API and Swagger UI
- [ ] **Step 13: MCP Integration** — Integrate MCP server for multi-cluster management
- [ ] **Step 14: JWT Authentication** — Secure API endpoints with JWT
- [ ] **Step 15: OpenTelemetry Instrumentation** — Add distributed tracing with OpenTelemetry
