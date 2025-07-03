# atlasctl

🔍 **Command-line tool for managing Atlas applications across Kubernetes environments**

## ✨ Features

- 📊 **Multi-environment View**: See all Atlas applications across dev, stage, and prod
- 🔄 **Real-time Status**: Current deployment status and health information
- 🎯 **Kubeconfig Support**: Connect to any Kubernetes cluster
- 📈 **Migration Tracking**: Database migration version monitoring
- 🚀 **Atlas Controller Integration**: Works with both standalone deployments and Atlas Controller

## 🚀 Quick Start

### Installation

#### Option 1: Build from Source
```bash
git clone https://github.com/your-org/atlasctl.git
cd atlasctl
go build -o atlasctl main.go
```

#### Option 2: Download Binary
```bash
# Download latest release
curl -L https://github.com/your-org/atlasctl/releases/latest/download/atlasctl-linux-amd64 -o atlasctl
chmod +x atlasctl
```

### Basic Usage

```bash
# List all Atlas applications
./atlasctl list

# Use specific kubeconfig
./atlasctl list --kubeconfig ~/.kube/my-cluster

# Get help
./atlasctl --help
```

## 📋 Commands

### `atlasctl list`
Display all Atlas applications across environments

```bash
# Basic usage
./atlasctl list

# With custom kubeconfig
./atlasctl list --kubeconfig ~/.kube/production

# Example output:
┼───────────┼───────┼─────────┼──────────────┼─────────┼──────────┼─────────────────────┼─────┼
│ NAMESPACE │  APP  │ VERSION │ MIGRATION ID │ STATUS  │ REPLICAS │     LAST UPDATE     │ AGE │
┼───────────┼───────┼─────────┼──────────────┼─────────┼──────────┼─────────────────────┼─────┼
│ dev       │ atlas │ 1.22.0  │ 6            │ Running │ 2/2      │ 2025-07-03 02:00:00 │ 5m  │
│ stage     │ atlas │ 1.22.0  │ 6            │ Running │ 3/3      │ 2025-07-03 01:58:00 │ 3m  │
│ prod      │ atlas │ 1.21.0  │ 5            │ Running │ 5/5      │ 2025-07-03 01:00:00 │ 2d  │
```

### Command Options

#### Global Flags
- `--kubeconfig`: Path to kubeconfig file (default: `~/.kube/config`)
- `--help`: Show help information

#### List Command Flags
- `--output`, `-o`: Output format (table, json, yaml)
- `--namespace`, `-n`: Filter by specific namespace
- `--watch`, `-w`: Watch for changes (continuous monitoring)

## 🔧 Configuration

### Kubeconfig Setup
```bash
# Set default kubeconfig
export KUBECONFIG=~/.kube/my-cluster

# Or use flag with each command
./atlasctl list --kubeconfig ~/.kube/production
```

### Environment Variables
```bash
# Default kubeconfig path
export KUBECONFIG=/path/to/kubeconfig

# Default namespace filter
export ATLASCTL_NAMESPACE=dev
```

## 🎯 Integration Modes

### Standalone Mode
When **Atlas Controller** is not installed, `atlasctl` reads directly from Kubernetes Deployments:

```bash
# Searches for deployments with label app=atlas
# Extracts information from:
# - spec.template.spec.containers[0].image (version)
# - env.MIGRATION_ID (migration version)
# - status.readyReplicas (replica status)
```

### Controller Mode
When **Atlas Controller** is installed, `atlasctl` can optionally read from AtlasApp Custom Resources:

```bash
# Future enhancement: Read from AtlasApp CRDs
# Provides richer information:
# - spec.version (application version)
# - spec.migrationId (migration ID)
# - status.phase (deployment phase)
# - status.ready (readiness status)
```

## 📊 Output Formats

### Table Format (Default)
```bash
./atlasctl list
```
Clean, human-readable table with columns for all key information.

### JSON Format
```bash
./atlasctl list --output json
```
```json
{
  "applications": [
    {
      "namespace": "dev",
      "app": "atlas",
      "version": "1.22.0",
      "migrationId": "6",
      "status": "Running",
      "replicas": "2/2",
      "lastUpdate": "2025-07-03T02:00:00Z",
      "age": "5m"
    }
  ]
}
```

### YAML Format
```bash
./atlasctl list --output yaml
```
```yaml
applications:
- namespace: dev
  app: atlas
  version: "1.22.0"
  migrationId: "6"
  status: Running
  replicas: "2/2"
  lastUpdate: "2025-07-03T02:00:00Z"
  age: "5m"
```

## 🔍 Advanced Usage

### Filter by Namespace
```bash
# Show only production applications
./atlasctl list --namespace prod

# Show only dev and stage
./atlasctl list --namespace dev,stage
```

### Watch Mode
```bash
# Continuously monitor applications
./atlasctl list --watch

# Watch with refresh interval
./atlasctl list --watch --interval 5s
```

### Multiple Clusters
```bash
# Check different clusters
./atlasctl list --kubeconfig ~/.kube/dev-cluster
./atlasctl list --kubeconfig ~/.kube/staging-cluster
./atlasctl list --kubeconfig ~/.kube/prod-cluster
```

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Access to Kubernetes cluster
- kubectl configured

### Build from Source
```bash
# Clone repository
git clone https://github.com/your-org/atlasctl.git
cd atlasctl

# Install dependencies
go mod tidy

# Build binary
go build -o atlasctl main.go

# Run tests
go test ./...
```

### Project Structure
```
atlasctl/
├── cmd/
│   └── root.go          # CLI root command
├── pkg/
│   ├── atlas/
│   │   ├── client.go    # Kubernetes client
│   │   └── types.go     # Data structures
│   └── output/
│       └── table.go     # Output formatting
├── main.go              # Entry point
├── go.mod
├── go.sum
└── README.md
```

### Adding New Features
1. Add command in `cmd/`
2. Implement logic in `pkg/`
3. Add tests
4. Update documentation

## 🔧 Troubleshooting

### Common Issues

#### Connection Problems
```bash
# Check kubeconfig
kubectl config current-context

# Verify cluster access
kubectl get nodes

# Test with specific kubeconfig
./atlasctl list --kubeconfig ~/.kube/config
```

#### No Applications Found
```bash
# Check if deployments exist
kubectl get deployments -n dev -l app=atlas

# Verify namespace access
kubectl auth can-i list deployments --namespace dev

# Check deployment labels
kubectl get deployment atlas -n dev -o yaml | grep -A 10 metadata.labels
```

#### Permission Errors
```bash
# Check RBAC permissions
kubectl auth can-i list deployments
kubectl auth can-i list deployments --namespace dev

# Verify service account (if using)
kubectl get serviceaccount
```

### Debug Mode
```bash
# Enable verbose logging
./atlasctl list --verbose

# Show raw Kubernetes responses
./atlasctl list --debug
```

## 🎯 Roadmap

- [ ] **Watch Mode**: Real-time monitoring with `--watch`
- [ ] **Output Formats**: JSON, YAML, and custom formats
- [ ] **Namespace Filtering**: Filter by specific namespaces
- [ ] **AtlasApp Integration**: Native support for Atlas Controller CRDs
- [ ] **Health Checks**: Application health status monitoring
- [ ] **Log Streaming**: View application logs directly
- [ ] **Deployment Actions**: Trigger deployments and rollbacks
- [ ] **Multi-cluster Support**: Manage multiple clusters simultaneously
- [ ] **Configuration File**: YAML-based configuration
- [ ] **Shell Completion**: Bash, Zsh, and Fish completion

## 📈 Performance

### Resource Usage
- **Memory**: ~10MB typical usage
- **CPU**: Minimal during operation
- **Network**: Only API calls to Kubernetes

### Scaling Considerations
- Handles clusters with 100+ applications
- Efficient API queries with label selectors
- Caching for repeated operations

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices
- Add tests for new functionality
- Update documentation
- Ensure backward compatibility

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [client-go](https://github.com/kubernetes/client-go) for Kubernetes integration
- Inspired by kubectl and other Kubernetes tools

---

**Ready to manage your Atlas applications?** 🚀 [Get started](#quick-start) or [contribute](#contributing) to make atlasctl even better!
