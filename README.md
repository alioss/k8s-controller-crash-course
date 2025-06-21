# k8s-controller-tutorial

A starter template for building Kubernetes controllers or CLI tools in Go using [Cobra CLI](https://github.com/spf13/cobra-cli).

## Features

✅ **Cobra CLI Framework** - Professional command-line interface
✅ **Version Flag** - `--version` / `-v` to show application version
✅ **Multiple Commands** - `go-basic` and `status` commands
✅ **Go Structs & Methods** - Learn Go fundamentals
✅ **Unit Tests** - Test coverage for struct methods

## Prerequisites

- [Go](https://golang.org/dl/) 1.20 or newer
- Basic understanding of command-line tools

## Installation & Usage

### 1. Clone and Setup
```bash
git clone https://github.com/yourusername/k8s-controller-tutorial.git
cd k8s-controller-tutorial
go mod tidy
```

### 2. Run Commands

**Show help:**
```bash
go run main.go --help
```

**Show version:**
```bash
go run main.go --version
go run main.go -v
```

**Run Go basics example:**
```bash
go run main.go go-basic
```

**Show cluster status:**
```bash
go run main.go status
```

### 3. Build Binary (Optional)
```bash
go build -o controller .
./controller --help
```

## Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `help` | Show help information | `go run main.go --help` |
| `version` | Show version information | `go run main.go -v` |
| `go-basic` | Demonstrate Go structs and methods | `go run main.go go-basic` |
| `status` | Display cluster status information | `go run main.go status` |

## Project Structure

```
├── main.go              # Entry point
├── go.mod               # Go module dependencies
├── cmd/                 # Command implementations
│   ├── root.go         # Root command + version flag
│   ├── go_basic.go     # Go basics demonstration
│   ├── go_basic_test.go # Unit tests
│   └── status.go       # Status command
└── README.md           # This file
```

## Go Concepts Demonstrated

### Kubernetes Struct
```go
type Kubernetes struct {
    Name       string     `json:"name"`
    Version    string     `json:"version"`
    Users      []string   `json:"users,omitempty"`
    NodeNumber func() int `json:"-"`
}
```

### Methods with Different Receivers
```go
// Value receiver (read-only)
func (k8s Kubernetes) GetUsers() { ... }

// Pointer receiver (can modify)
func (k8s *Kubernetes) AddNewUser(user string) { ... }
```

### Cobra Command Structure
```go
var myCmd = &cobra.Command{
    Use:   "command-name",
    Short: "Brief description",
    Long:  "Detailed description",
    Run:   func(cmd *cobra.Command, args []string) { ... },
}
```

## Testing

Run all tests:
```bash
go test ./cmd
```

Run tests with coverage:
```bash
go test -cover ./cmd
```

## What's Next?

This is **Step 1** in the Kubernetes Controller tutorial series:

- ✅ **Step 1:** Golang CLI Application using Cobra *(current)*
- 🔄 **Step 2:** Zerolog for Log Levels
- 🔄 **Step 3:** pflag for Log Level Flags
- 🔄 **Step 4:** FastHTTP Server Command
- 🔄 **Step 5:** Makefile, Dockerfile, and GitHub Workflow

## License

MIT License. See [LICENSE](LICENSE) for details.
#