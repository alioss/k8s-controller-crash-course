# Atlas Controller

🚀 **Kubernetes Controller for automating Atlas application deployments with auto-promotion between environments**

## ✨ Features

- 🔄 **Auto-promotion**: dev → stage → prod (with approval)
- 📊 **Version Management**: Control application versions and migrations
- 🛡️ **Production Approval**: Requires manual approval for production deployments
- 🔍 **Health Checks**: Application health monitoring
- 📈 **Status Tracking**: Real-time deployment status tracking
- 🎯 **Custom Resources**: Declarative approach with Kubernetes CRDs

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                Atlas Controller Workflow                    │
│                                                             │
│  dev ────auto────▶ stage ────auto────▶ prod (approval)     │
│   ↓                  ↓                   ↓                 │
│ health             health              health              │
│ check              check               check               │
│   ↓                  ↓                   ↓                 │
│ ready              ready               manual              │
│                                       approval             │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Kubernetes cluster (1.19+)
- kubectl configured
- Docker (for building)
- Go 1.21+ (for development)

### 1. Setup Namespaces
```bash
kubectl create namespace dev
kubectl create namespace stage
kubectl create namespace prod
kubectl create namespace atlas-system
```

### 2. Install Atlas Controller

#### Option A: Using Pre-built Image
```bash
# Install CRD
kubectl apply -f config/crd/atlasapp.yaml

# Setup RBAC
kubectl apply -f config/rbac/role.yaml

# Create image pull secret (if using private registry)
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password=YOUR_GITHUB_TOKEN \
  --docker-email=YOUR_EMAIL \
  -n atlas-system

# Deploy controller
kubectl apply -f config/manager/manager.yaml
```

#### Option B: Build from Source
```bash
# Clone repository
git clone https://github.com/your-org/atlas-controller.git
cd atlas-controller

# Build for Linux (if building on Mac/Windows)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/manager main.go

# Build Docker image
docker build --platform linux/amd64 -t your-registry/atlas-controller:latest .
docker push your-registry/atlas-controller:latest

# Update manager.yaml with your image
# Then apply manifests as in Option A
```

### 3. Create Your First AtlasApp

```bash
# Deploy to dev environment
kubectl apply -f - <<EOF
apiVersion: atlas.io/v1
kind: AtlasApp
metadata:
  name: atlas-dev
  namespace: dev
spec:
  environment: dev
  version: "1.21.0"
  migrationId: 5
  replicas: 2
  autoPromote: true
  nextEnvironment: stage
  healthCheckPath: "/"
EOF
```

### 4. Watch the Magic Happen! ✨

```bash
# Watch AtlasApp resources
kubectl get atlasapp -A -w

# Check controller logs
kubectl logs -n atlas-system deployment/atlas-controller -f

# View deployment status
kubectl get deployments -A
```

## 📋 AtlasApp Custom Resource

### Spec Fields
```yaml
apiVersion: atlas.io/v1
kind: AtlasApp
metadata:
  name: atlas-dev
  namespace: dev
spec:
  environment: dev           # Environment: dev/stage/prod
  version: "1.21.0"         # Application version
  migrationId: 5            # Database migration ID
  replicas: 2               # Number of replicas
  autoPromote: true         # Enable auto-promotion
  nextEnvironment: stage    # Next environment for promotion
  healthCheckPath: "/"      # Health check endpoint
  requireApproval: false    # Require manual approval
```

### Status Fields
```yaml
status:
  phase: Ready              # Current phase
  ready: true               # Readiness status
  readyReplicas: 2         # Ready replica count
  totalReplicas: 2         # Total replica count
  lastUpdate: "2025-07-03T02:00:00Z"
  approvalRequired: false   # Approval needed
  promotionPending: false   # Promotion waiting
  message: "Application is healthy and ready"
```

## 🔄 Promotion Workflows

### Automatic Promotion (dev → stage)
```bash
# Update dev version
kubectl patch atlasapp atlas-dev -n dev --type='json' \
  -p='[{"op": "replace", "path": "/spec/version", "value": "1.22.0"}]'

# Controller will:
# 1. Update dev deployment
# 2. Wait for ready status
# 3. Automatically create/update stage AtlasApp
# 4. Deploy to stage environment
```

### Manual Approval (stage → prod)
```bash
# Check if stage is ready for promotion
kubectl get atlasapp atlas-stage -n stage -o yaml

# If promotionPending: true, approve by creating prod AtlasApp
kubectl apply -f - <<EOF
apiVersion: atlas.io/v1
kind: AtlasApp
metadata:
  name: atlas-prod
  namespace: prod
spec:
  environment: prod
  version: "1.22.0"
  migrationId: 6
  replicas: 5
  autoPromote: false
  requireApproval: true
EOF
```

## 📊 Monitoring & Observability

### Check Application Status
```bash
# All AtlasApp resources
kubectl get atlasapp -A

# Detailed status
kubectl describe atlasapp atlas-dev -n dev

# YAML output
kubectl get atlasapp atlas-dev -n dev -o yaml
```

### Controller Logs
```bash
# Follow controller logs
kubectl logs -n atlas-system deployment/atlas-controller -f

# Recent events
kubectl get events -A --sort-by='.firstTimestamp'
```

### Integration with atlasctl
The controller works seamlessly with the existing `atlasctl` CLI:

```bash
# Your existing atlasctl continues to work
./atlasctl list --kubeconfig ~/.kube/config

# Output shows managed applications:
# NAMESPACE │ APP   │ VERSION │ MIGRATION ID │ STATUS  │ REPLICAS │ LAST UPDATE │ AGE
# dev       │ atlas │ 1.22.0  │ 6            │ Running │ 2/2      │ 2 min ago   │ 5m
# stage     │ atlas │ 1.22.0  │ 6            │ Running │ 3/3      │ 1 min ago   │ 3m
# prod      │ atlas │ 1.21.0  │ 5            │ Running │ 5/5      │ 1 hour ago  │ 2d
```

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Docker
- kubectl
- A Kubernetes cluster for testing

### Build & Test
```bash
# Install dependencies
go mod tidy

# Generate code
make generate

# Update manifests
make manifests

# Build binary
go build -o bin/manager main.go

# Run tests
go test ./...

# Build Docker image
docker build -t atlas-controller:latest .
```

### Run Locally
```bash
# Run controller locally (connects to cluster via kubeconfig)
go run main.go

# In another terminal, test with AtlasApp resources
kubectl apply -f examples/dev-deployment.yaml
```

## 📚 Configuration

### Environment Variables
```yaml
env:
- name: CONTROLLER_LOG_LEVEL
  value: "info"
- name: METRICS_ADDR
  value: ":8080"
- name: HEALTH_PROBE_ADDR
  value: ":8081"
```

### RBAC Permissions
The controller requires the following permissions:
- `atlasapps`: Full access for managing AtlasApp resources
- `deployments`: CRUD operations for application deployments
- `services`: CRUD operations for service resources
- `leases`: Leader election coordination

### Health Checks
```yaml
# Custom health check path
spec:
  healthCheckPath: "/api/health"
  
# Controller will check: http://service-name/api/health
```

## 🔧 Troubleshooting

### Common Issues

#### Controller Won't Start
```bash
# Check RBAC permissions
kubectl auth can-i create atlasapps --as=system:serviceaccount:atlas-system:atlas-controller-sa

# Verify CRD installation
kubectl get crd atlasapps.atlas.io

# Check controller logs
kubectl logs -n atlas-system deployment/atlas-controller
```

#### AtlasApp Not Creating Deployments
```bash
# Check AtlasApp status
kubectl describe atlasapp atlas-dev -n dev

# Verify namespace permissions
kubectl auth can-i create deployments --as=system:serviceaccount:atlas-system:atlas-controller-sa -n dev

# Check for resource conflicts
kubectl get events -n dev --field-selector involvedObject.name=atlas
```

#### Auto-promotion Not Working
```bash
# Verify autoPromote setting
kubectl get atlasapp atlas-dev -n dev -o jsonpath='{.spec.autoPromote}'

# Check ready status
kubectl get atlasapp atlas-dev -n dev -o jsonpath='{.status.ready}'

# Review controller logs for promotion events
kubectl logs -n atlas-system deployment/atlas-controller | grep -i promotion
```

#### Image Pull Errors
```bash
# For private registries, create image pull secret
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=USERNAME \
  --docker-password=TOKEN \
  -n atlas-system

# Update manager.yaml to reference the secret
```

## 🎯 Roadmap

- [ ] **GitOps Integration**: ArgoCD/Flux compatibility
- [ ] **Webhook Validation**: Admission controllers for validation
- [ ] **Metrics & Monitoring**: Prometheus metrics integration
- [ ] **Advanced Routing**: Canary and blue/green deployments
- [ ] **Multi-cluster Support**: Cross-cluster promotions
- [ ] **Policy Engine**: Custom promotion rules and policies
- [ ] **Slack/Teams Integration**: Approval notifications
- [ ] **Rollback Capabilities**: Automatic rollback on failures

## 📈 Performance & Scaling

### Resource Requirements
```yaml
resources:
  requests:
    cpu: 10m
    memory: 64Mi
  limits:
    cpu: 500m
    memory: 128Mi
```

### Scaling Considerations
- Controller supports multiple replicas with leader election
- Horizontal scaling of managed applications via `spec.replicas`
- Namespace-based isolation for multi-tenant environments

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation for API changes
- Ensure backward compatibility

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Kubebuilder](https://kubebuilder.io/)
- Uses [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- Inspired by GitOps principles and Kubernetes best practices

---

**Ready to automate your deployments?** 🚀 [Get started](#quick-start) or [contribute](#contributing) to make Atlas Controller even better!
