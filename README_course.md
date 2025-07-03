# Atlas Platform

🚀 **Complete Kubernetes-native platform for automated application deployment and management**

Atlas Platform provides a comprehensive solution for managing application deployments across multiple environments with automated promotion workflows, health monitoring, and declarative configuration.

## 🌟 Overview

The Atlas Platform consists of two main components that work together to provide a complete deployment automation solution:

### 🎯 [Atlas Controller](./atlas-controller/README.md)
A Kubernetes controller that automates application deployments with intelligent promotion workflows:
- **Automated Promotion**: dev → stage → prod with approval gates
- **Health Monitoring**: Continuous application health checks
- **Declarative Management**: Custom Resource Definitions (CRDs)
- **Production Safety**: Manual approval required for production deployments

### 🔍 [atlasctl](./atlasctl/README.md)
A command-line tool for observing and managing Atlas applications:
- **Multi-environment View**: See all deployments across clusters
- **Real-time Status**: Current deployment status and health
- **Flexible Output**: Table, JSON, and YAML formats
- **Kubeconfig Support**: Connect to any Kubernetes cluster

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           Atlas Platform Architecture                            │
│                                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────────┐  │
│  │     Dev     │    │    Stage    │    │    Prod     │    │    atlasctl     │  │
│  │             │    │             │    │             │    │                 │  │
│  │  AtlasApp   │───▶│  AtlasApp   │───▶│  AtlasApp   │    │  ┌─────────────┐ │  │
│  │  v1.21.0    │    │  v1.21.0    │    │  v1.21.0    │    │  │ CLI Tool    │ │  │
│  │  ready ✓    │    │  ready ✓    │    │  pending... │    │  │             │ │  │
│  │             │    │             │    │             │    │  │ List/Watch  │ │  │
│  │  ┌────────┐ │    │  ┌────────┐ │    │  ┌────────┐ │    │  │ Monitor     │ │  │
│  │  │Deploy  │ │    │  │Deploy  │ │    │  │Deploy  │ │    │  │ Status      │ │  │
│  │  │        │ │    │  │        │ │    │  │        │ │    │  └─────────────┘ │  │
│  │  └────────┘ │    │  └────────┘ │    │  └────────┘ │    └─────────────────┘  │
│  └─────────────┘    └─────────────┘    └─────────────┘                         │
│           │                  │                  │                              │
│           └──────────────────┼──────────────────┘                              │
│                              │                                                 │
│                    ┌─────────▼──────────┐                                      │
│                    │  Atlas Controller  │                                      │
│                    │                    │                                      │
│                    │  ┌──────────────┐  │                                      │
│                    │  │ Reconciler   │  │                                      │
│                    │  │ Health Check │  │                                      │
│                    │  │ Promotion    │  │                                      │
│                    │  │ Logic        │  │                                      │
│                    │  └──────────────┘  │                                      │
│                    └────────────────────┘                                      │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### 1. Prerequisites
- Kubernetes cluster (1.19+)
- kubectl configured
- Go 1.21+ (for building from source)

### 2. Install Atlas Controller
```bash
# Create namespaces
kubectl create namespace dev
kubectl create namespace stage
kubectl create namespace prod
kubectl create namespace atlas-system

# Install the controller
cd atlas-controller
kubectl apply -f config/crd/atlasapp.yaml
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/manager/manager.yaml
```

### 3. Install atlasctl
```bash
# Build from source
cd atlasctl
go build -o atlasctl main.go

# Or download binary (when releases are available)
curl -L https://github.com/your-org/atlas-platform/releases/latest/download/atlasctl-linux-amd64 -o atlasctl
chmod +x atlasctl
```

### 4. Deploy Your First Application
```bash
# Create an AtlasApp resource
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

### 5. Monitor with atlasctl
```bash
# Watch your applications
./atlasctl list --watch

# Example output:
┼───────────┼───────┼─────────┼──────────────┼─────────┼──────────┼─────────────────────┼─────┼
│ NAMESPACE │  APP  │ VERSION │ MIGRATION ID │ STATUS  │ REPLICAS │     LAST UPDATE     │ AGE │
┼───────────┼───────┼─────────┼──────────────┼─────────┼──────────┼─────────────────────┼─────┼
│ dev       │ atlas │ 1.21.0  │ 5            │ Running │ 2/2      │ 2025-07-03 02:00:00 │ 2m  │
│ stage     │ atlas │ 1.21.0  │ 5            │ Running │ 2/2      │ 2025-07-03 02:01:00 │ 1m  │
```

## 🎯 Use Cases

### 🔄 Automated DevOps Pipeline
Perfect for teams wanting to automate their deployment pipeline:
- **Developers** push code → **Atlas Controller** handles promotion
- **Automatic testing** in each environment before promotion
- **Manual approval gates** for production deployments

### 📊 Multi-Environment Management
Ideal for organizations managing multiple environments:
- **Consistent deployments** across dev, stage, and prod
- **Version tracking** and migration management
- **Centralized monitoring** with atlasctl

### 🛡️ Production Safety
Built-in safety mechanisms for production deployments:
- **Approval workflows** prevent accidental production deployments
- **Health checks** ensure applications are ready before promotion
- **Rollback capabilities** for quick recovery

## 📋 Components Deep Dive

### Atlas Controller Features
- 🔄 **Automatic Promotion**: Intelligent workflow management
- 📊 **Custom Resources**: Kubernetes-native approach
- 🔍 **Health Monitoring**: Continuous application health checks
- 🛡️ **Production Gates**: Manual approval for critical deployments
- 📈 **Status Tracking**: Real-time deployment status
- 🎯 **Namespace Isolation**: Multi-tenant support

[📖 Read the complete Atlas Controller documentation →](./atlas-controller/README.md)

### atlasctl Features
- 📊 **Multi-environment View**: See all deployments at once
- 🔍 **Real-time Monitoring**: Live status updates
- 🎯 **Flexible Output**: Table, JSON, YAML formats
- 🔧 **Kubeconfig Support**: Connect to any cluster
- 📈 **Migration Tracking**: Database version monitoring
- 🚀 **Integration Ready**: Works with or without Atlas Controller

[📖 Read the complete atlasctl documentation →](./atlasctl/README.md)

## 🔧 Configuration Examples

### Basic Development Workflow
```yaml
# dev-app.yaml
apiVersion: atlas.io/v1
kind: AtlasApp
metadata:
  name: myapp-dev
  namespace: dev
spec:
  environment: dev
  version: "1.22.0"
  migrationId: 6
  replicas: 2
  autoPromote: true
  nextEnvironment: stage
```

### Production Deployment
```yaml
# prod-app.yaml
apiVersion: atlas.io/v1
kind: AtlasApp
metadata:
  name: myapp-prod
  namespace: prod
spec:
  environment: prod
  version: "1.22.0"
  migrationId: 6
  replicas: 10
  autoPromote: false
  requireApproval: true
  healthCheckPath: "/health"
```

### Multi-cluster Monitoring
```bash
# Monitor development cluster
./atlasctl list --kubeconfig ~/.kube/dev-cluster

# Monitor production cluster
./atlasctl list --kubeconfig ~/.kube/prod-cluster

# Watch all environments
./atlasctl list --watch --kubeconfig ~/.kube/all-clusters
```

## 🛠️ Development

### Repository Structure
```
atlas-platform/
├── atlas-controller/          # Kubernetes controller
│   ├── api/v1/                # CRD definitions
│   ├── internal/controller/   # Controller logic
│   ├── config/               # Kubernetes manifests
│   └── README.md
├── atlasctl/                 # CLI tool
│   ├── cmd/                  # CLI commands
│   ├── pkg/                  # Core logic
│   └── README.md
├── docs/                     # Additional documentation
├── examples/                 # Usage examples
└── README.md                 # This file
```

### Building from Source
```bash
# Build Atlas Controller
cd atlas-controller
go build -o bin/manager main.go
docker build -t atlas-controller:latest .

# Build atlasctl
cd ../atlasctl
go build -o atlasctl main.go
```

### Running Tests
```bash
# Test Atlas Controller
cd atlas-controller
go test ./...

# Test atlasctl
cd ../atlasctl
go test ./...
```

## 🎯 Roadmap

### Phase 1: Core Features ✅
- [x] Atlas Controller with auto-promotion
- [x] atlasctl for monitoring
- [x] Basic health checks
- [x] Production approval gates

### Phase 2: Enhanced Features 🚧
- [ ] GitOps integration (ArgoCD/Flux)
- [ ] Webhook validation
- [ ] Prometheus metrics
- [ ] Slack/Teams notifications

### Phase 3: Advanced Features 🔮
- [ ] Blue/Green deployments
- [ ] Canary releases
- [ ] Multi-cluster support
- [ ] Policy engine
- [ ] Advanced rollback capabilities

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### Development Guidelines
- Follow Go best practices and idioms
- Add comprehensive tests for new features
- Update documentation for API changes
- Ensure backward compatibility
- Use conventional commit messages

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Kubebuilder](https://kubebuilder.io/) for the controller framework
- Uses [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Inspired by GitOps principles and Kubernetes best practices
- Thanks to the Kubernetes community for excellent tooling and documentation

## 📞 Support

- 📖 **Documentation**: Check component-specific READMEs
- 🐛 **Issues**: Report bugs via GitHub Issues
- 💬 **Discussions**: Join GitHub Discussions for questions
- 📧 **Contact**: Reach out to maintainers

---

**Ready to automate your deployments?** 🚀

Get started with the [Atlas Controller](./atlas-controller/README.md) or [atlasctl](./atlasctl/README.md) today!
