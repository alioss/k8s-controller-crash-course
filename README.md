# Golang Kubernetes Controller Tutorial

The goal of this repository is to demonstrate how to build a Kubernetes controller step by step. Starting from a basic CLI application, the tutorial progressively adds features like informers, reconcilers, metrics, and leader election to create a fully functional, production-ready Kubernetes operator using Go and controller-runtime.

🎆 **Latest Achievement: Complete Platform as a Service with Port.io Integration!**

> **🚀 NEW: Port.io Integration** - Our controller now includes a full Internal Developer Portal integration with Port.io, transforming it into a complete Platform as a Service with beautiful self-service UI, real-time sync, and developer workflows. See [README-port.md](./README-port.md) for complete setup guide.

> **📚 Reference Repository:** The complete reference implementation for this tutorial can be found at:
> **https://github.com/den-vasyliev/k8s-controller-tutorial-ref**

## 📋 Steps Overview

| Step | Branch | Status |
|------|--------|---------|
| **Step 1** | [feature/step1-cobra-cli](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step1-cobra-cli) | ✅ Done |
| **Step 2** | [feature/step2-zerolog](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step2-zerolog-logging) | ✅ Done |
| **Step 3** | [feature/step3-pflag](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step3-pflag-loglevel) | ✅ Done |
| **Step 4** | [feature/step4-fasthttp-server](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step4-fasthttp-server) | ✅ Done |
| **Step 5** | [feature/step5-makefile-dockerfile-workflow](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step5-makefile-docker-ci) | ✅ Done |
| **Step 6** | [feature/step6-list-deployments](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step6-list-deployments) | ✅ Done |
| **Step 7** | [feature/step7-deployment-informer](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step7-informer) | ✅ Done |
| **Step 8** | [feature/step8-deployments-api-endpoint](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step8-api-handler) | ✅ Done |
| **Step 9** | [feature/step9-controller-runtime](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step9-controller-runtime) | ✅ Done |
| **Step 10** | [feature/step10-leader-election](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step10-leader-election) | ✅ Done |
| **Step 11** | [feature/step11-frontendpage-crd](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step11-frontendpage-crd) | ✅ Done |
| **Step 12** | [feature/step12-platform-api](https://github.com/alioss/k8s-controller-crash-course/blob/feature/step12-platform-api/README-port.md) | ✅ Done |
| **Step 13** | [feature/step13-mcp-integration](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step13-mcp-integration) | ⭕ Todo |
| **Step 14** | [feature/step14-jwt-auth](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step14-jwt-auth) | ⭕ Todo |
| **Step 15** | [feature/step15-opentelemetry](https://github.com/alioss/k8s-controller-crash-course/tree/feature/step15-opentelemetry) | ⭕ Todo |

**Progress: 12/15 (80%) Complete** 🚀

## Progress Status

### ✅ Completed Steps

- [x] **Step 1: Golang CLI Application using Cobra** — Initialize a CLI app with cobra-cli
  * [x] Set up basic Cobra CLI structure
  * [x] Create go-basic command with Kubernetes struct
  * [x] Add version flag (--version, -v)
  * [x] Create status command with cluster information
  * [x] Add unit tests for struct methods
  * [x] Update project documentation

- [x] **Step 2: Zerolog for Log Levels** — Add structured logging with zerolog
  * [x] Replace fmt.Println with structured logging in go-basic command
  * [x] Configure different log levels (trace, debug, info, warn, error)
  * [x] Add JSON log formatting for all commands
  * [x] Implement structured fields (cluster_name, version, username, etc.)
  * [x] Demonstrate zerolog fluent API with method chaining

- [x] **Step 3: pflag for Log Level Flags** — Integrate pflag for CLI log level flags
  * [x] Add --log-level flag to control zerolog output level
  * [x] Implement flag validation for valid log levels (trace, debug, info, warn, error)
  * [x] Configure dynamic log level switching with zerolog.SetGlobalLevel()
  * [x] Add different output formats based on log level (console vs JSON)
  * [x] Use PersistentFlags for global log level control across all commands
  * [x] Demonstrate comprehensive logging at all verbosity levels

- [x] **Step 4: FastHTTP Server Command** — Add a server command with configurable port and log level
  * [x] Add server command using FastHTTP instead of net/http
  * [x] Implement configurable port and host flags (--port)
  * [x] Add structured logging for HTTP requests and responses
  * [x] Create simple routing system with switch statement
  * [x] Add /health endpoint returning JSON status
  * [x] Implement 404 handling with JSON error responses
  * [x] Set appropriate Content-Type headers for different endpoints

- [x] **Step 5: Makefile, Dockerfile, and GitHub Workflow** — Introduce build automation, secure containerization, CI/CD, and tests
  * [x] Create Makefile for build automation (build, test, run, docker-build targets)
  * [x] Add multi-stage Dockerfile for optimized containerization
  * [x] Setup GitHub Actions workflow for comprehensive CI/CD pipeline
  * [x] Implement Docker security scanning with Trivy
  * [x] Add Helm chart packaging and artifact upload
  * [x] Configure cross-platform build support
  * [x] Add automated testing and code quality checks

- [x] **Step 6: List Kubernetes Deployments with client-go** — List deployments in the default namespace
  * [x] Add client-go dependency for Kubernetes API access
  * [x] Implement deployment listing functionality with detailed information
  * [x] Add kubeconfig handling and cluster connectivity
  * [x] Create colorful status indicators with emoji (✅ ⚠️ ❌ ⏸️)
  * [x] Display replica status, container images, and deployment age
  * [x] Add comprehensive summary statistics (ready deployments, running pods)
  * [x] Implement structured logging for debugging and troubleshooting
  * [x] Add professional deployment overview with human-readable formatting

- [x] **Step 7: Deployment Informer with client-go** — Watch and log Deployment events
  * [x] Implement Kubernetes Informer pattern for real-time events
  * [x] Add event watching and logging for deployment changes (ADD, UPDATE, DELETE)
  * [x] Create event handlers with structured logging via zerolog
  * [x] Integrate informer with FastHTTP server for concurrent operation
  * [x] Support both kubeconfig and in-cluster authentication methods
  * [x] Implement automatic cache synchronization with 30-second resync
  * [x] Verify real-time event detection and proper informer lifecycle
  * [x] Establish foundation for production Kubernetes controller development

- [x] **Step 8: /deployments JSON API Endpoint** — Serve deployment names as JSON from the informer cache
  * [x] Create HTTP endpoint to expose cached deployment data
  * [x] Integrate informer cache with REST API responses
  * [x] Add JSON serialization for deployment information
  * [x] Implement request ID tracking with UUID for better debugging
  * [x] Add /nodes endpoint for cluster-wide node information
  * [x] Create dual informer setup (deployments + nodes) running concurrently
  * [x] Add structured logging for API requests with contextual information
  * [x] Demonstrate efficient API responses using in-memory informer cache

- [x] **Step 9: controller-runtime Deployment Controller** — Reconcile Deployments and log events
  * [x] Implement controller-runtime framework integration with manager
  * [x] Create Deployment Reconciler with structured logging
  * [x] Add concurrent controller-runtime manager alongside informers
  * [x] Integrate controller-runtime with existing FastHTTP server architecture
  * [x] Support both kubeconfig and in-cluster authentication for controller-runtime
  * [x] Implement real-time reconciliation for deployment creation, updates, deletion
  * [x] **BONUS: Event Monitoring System** — Added comprehensive Kubernetes Event tracking
  * [x] Add Event Informer to monitor all cluster events in real-time
  * [x] Create /events JSON API endpoint for event inspection
  * [x] Implement structured event logging with emoji indicators (📅)
  * [x] Add proper JSON marshalling for complex event data
  * [x] Establish foundation for production Kubernetes operator development

- [x] **Step 10: Leader Election and Metrics** — Add HA and metrics endpoint to the controller manager
  * [x] Implement Leader Election for High Availability controller deployment
  * [x] Add Prometheus metrics endpoint with controller-runtime integration
  * [x] Configure proper logging for controller-runtime with zap logger

- [x] **Step 11: FrontendPage CRD and Advanced Controller** — Define a custom resource and manage Deployments/ConfigMaps
  * [x] Implement Custom Resource Definition (CRD) for FrontendPage with frontendpage.alios.io domain
  * [x] Create Go types and deepcopy generation for custom resources
  * [x] Build advanced FrontendPage Controller with sophisticated reconciliation logic
  * [x] Add ConfigMap management for HTML content storage
  * [x] Implement Deployment creation and updates with nginx containers
  * [x] Configure OwnerReferences for automatic cascade deletion
  * [x] Handle replica scaling and container image updates
  * [x] Add comprehensive Mermaid architecture diagram showing CRD workflow
  * [x] Establish foundation for production Kubernetes operators with custom resources

- [x] **Step 12: Platform API + Port.io Integration** — Transform controller into Platform as a Service with Port.io
  * [x] Implement RESTful CRUD API for FrontendPage resources
  * [x] Add Swagger UI documentation with interactive API explorer
  * [x] **Port.io Integration**: Connect controller with Internal Developer Portal
  * [x] Create Port.io Blueprint for FrontendPage schema definition
  * [x] Build Self-Service Actions for developer self-service workflows
  * [x] Implement bi-directional real-time sync (Port.io ↔ Kubernetes)
  * [x] Add automatic status updates and metrics synchronization
  * [x] Support webhook-driven deployments from Port.io UI
  * [x] Create comprehensive Port.io documentation with troubleshooting
  * [x] **Result**: Complete Platform as a Service with beautiful developer portal

### 🔄 In Progress

- [ ] **Step 13: MCP Integration** — Integrate MCP server for multi-cluster management
- [ ] **Step 14: JWT Authentication** — Secure API endpoints with JWT
- [ ] **Step 15: OpenTelemetry Instrumentation** — Add distributed tracing with OpenTelemetry