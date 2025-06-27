# 🚀 Port.io Integration for Kubernetes Controller

Complete guide for integrating our Kubernetes FrontendPage Controller with Port.io Internal Developer Portal.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Setup Instructions](#setup-instructions)
- [Configuration](#configuration)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Advanced Features](#advanced-features)

## 🏗️ Overview

This integration transforms our Kubernetes Controller into a full-featured **Platform as a Service (PaaS)** by connecting it with Port.io Internal Developer Portal.

### What it provides:

✅ **Beautiful Self-Service UI** - No more YAML hell  
✅ **Real-time Status Sync** - Live updates from Kubernetes  
✅ **Software Catalog** - Complete visibility of all resources  
✅ **Developer Self-Service** - One-click deployments  
✅ **Automated Workflows** - GitOps-ready integration  

## 🏛️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Port.io UI    │───▶│  Our Controller  │───▶│   Kubernetes    │
│  Self-Service   │    │    + REST API    │    │    Cluster      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
        ▲                        │                        
        │                        ▼                        
        └──────── Port.io API ◀──┘                        
                 (Status sync)
```

### Data Flow:

1. **Developer** uses Port.io UI to create FrontendPage
2. **Port.io** sends webhook to our Controller API
3. **Controller** creates Kubernetes resources (Deployment + ConfigMap)
4. **Controller** syncs status back to Port.io in real-time
5. **Port.io** displays live status, metrics, and URLs

## 📚 Prerequisites

- Kubernetes cluster with kubectl access
- Go 1.21+ 
- Port.io account (free tier available)
- ngrok or similar tunneling service (for local development)

## 🛠️ Setup Instructions

### Step 1: Port.io Account Setup

1. **Create account** at [port.io](https://www.port.io)
2. **Get API credentials**:
   - Go to your Port app → Settings → Credentials
   - Copy `Client ID` and `Client Secret`
   - Note your API URL (`https://api.getport.io` or `https://api.us.getport.io`)

### Step 2: Create Blueprint in Port.io

1. **Navigate to:** Builder → Data model → New blueprint
2. **Use JSON mode** and paste:

```json
{
  "identifier": "frontendpage",
  "title": "Frontend Page",
  "icon": "Microservice",
  "description": "A custom web page deployed to Kubernetes",
  "schema": {
    "properties": {
      "contents": {
        "title": "HTML Contents",
        "type": "string",
        "description": "The HTML content of the page"
      },
      "image": {
        "title": "Container Image",
        "type": "string",
        "default": "nginx:latest",
        "description": "The container image to use"
      },
      "replicas": {
        "title": "Replicas",
        "type": "number",
        "minimum": 1,
        "maximum": 10,
        "default": 2,
        "description": "Number of replicas to run"
      },
      "namespace": {
        "title": "Kubernetes Namespace",
        "type": "string",
        "default": "default",
        "description": "The Kubernetes namespace"
      },
      "status": {
        "title": "Deployment Status",
        "type": "string",
        "enum": ["Creating", "Ready", "Failed"],
        "enumColors": {
          "Creating": "yellow",
          "Ready": "green",
          "Failed": "red"
        },
        "description": "Current deployment status"
      },
      "url": {
        "title": "Page URL",
        "type": "string",
        "format": "url",
        "description": "URL to access the deployed page"
      },
      "createdAt": {
        "title": "Created At",
        "type": "string",
        "format": "date-time",
        "description": "When the page was created"
      },
      "updatedAt": {
        "title": "Updated At", 
        "type": "string",
        "format": "date-time",
        "description": "When the page was last updated"
      }
    },
    "required": ["contents", "image", "replicas"]
  },
  "mirrorProperties": {},
  "calculationProperties": {},
  "relations": {}
}
```

3. **Save Blueprint**

### Step 3: Create Self-Service Action

1. **Navigate to:** Self-service → Actions → New Action
2. **Use JSON mode** and paste:

```json
{
  "identifier": "create_frontendpage",
  "title": "🚀 Create Frontend Page",
  "description": "Deploy a new frontend page to Kubernetes",
  "trigger": {
    "type": "self-service",
    "operation": "CREATE",
    "blueprintIdentifier": "frontendpage",
    "userInputs": {
      "properties": {
        "name": {
          "title": "Page Name",
          "type": "string",
          "pattern": "^[a-z0-9-]+$",
          "description": "Name for your frontend page (lowercase, numbers, hyphens only)"
        },
        "contents": {
          "title": "HTML Contents",
          "type": "string",
          "default": "<h1>Welcome!</h1><p>This is my awesome page created via Port!</p>",
          "description": "The HTML content of your page"
        },
        "image": {
          "title": "Container Image",
          "type": "string",
          "enum": ["nginx:latest", "nginx:alpine", "nginx:1.21"],
          "default": "nginx:latest",
          "description": "Choose the web server image"
        },
        "replicas": {
          "title": "Number of Replicas",
          "type": "number",
          "minimum": 1,
          "maximum": 5,
          "default": 2,
          "description": "How many instances to run"
        },
        "namespace": {
          "title": "Kubernetes Namespace",
          "type": "string",
          "enum": ["default", "development", "staging"],
          "default": "default",
          "description": "Target namespace for deployment"
        }
      },
      "required": ["name", "contents"]
    }
  },
  "invocationMethod": {
    "type": "WEBHOOK",
    "url": "YOUR_NGROK_URL/api/frontendpages",
    "method": "POST",
    "synchronized": false,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "metadata": {
        "name": "{{ .inputs.name }}",
        "namespace": "{{ .inputs.namespace }}"
      },
      "spec": {
        "contents": "{{ .inputs.contents }}",
        "image": "{{ .inputs.image }}",
        "replicas": "{{ .inputs.replicas | tonumber }}"
      }
    }
  }
}
```

3. **Important:** Replace `YOUR_NGROK_URL` with your actual ngrok URL
4. **Save Action**

### Step 4: Setup ngrok (for local development)

```bash
# Install ngrok
brew install ngrok

# Start tunnel
ngrok http 8080

# Copy the HTTPS URL (e.g., https://abc123.ngrok.io)
# Update the webhook URL in Port.io Action
```

### Step 5: Start Controller with Port Integration

```bash
# Set environment variables (optional)
export PORT_CLIENT_ID="your_client_id"
export PORT_CLIENT_SECRET="your_client_secret"

# Start controller with Port integration
go run main.go server \
  --port 8080 \
  --kubeconfig ~/.kube/config \
  --enable-leader-election=false \
  --port-base-url "https://api.getport.io" \
  --port-client-id "$PORT_CLIENT_ID" \
  --port-client-secret "$PORT_CLIENT_SECRET"
```

**Expected output:**
```
{"level":"info","message":"Setting up Port.io integration..."}
{"level":"info","message":"Successfully authenticated with Port.io"}
{"level":"info","message":"Successfully connected to Port.io"}
{"level":"info","message":"Starting FastHTTP server on :8080"}
```

## ⚙️ Configuration

### Command Line Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--port-base-url` | Port.io API URL | `https://api.getport.io` |
| `--port-client-id` | Port.io Client ID | `port_xxx` |
| `--port-client-secret` | Port.io Client Secret | `xxx` |

### Environment Variables

```bash
export PORT_BASE_URL="https://api.getport.io"
export PORT_CLIENT_ID="your_client_id" 
export PORT_CLIENT_SECRET="your_client_secret"
```

### Swagger API Documentation

The controller includes Swagger UI for API documentation:

- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **Swagger JSON:** `http://localhost:8080/swagger/doc.json`

#### Swagger Setup (Fixed Issues)

We fixed several Swagger integration issues:

1. **Added proper imports:**
```go
_ "github.com/yourusername/k8s-controller-tutorial/docs" // for swagger
```

2. **Added Swagger endpoints:**
```go
// Swagger JSON endpoint
router.GET("/swagger/doc.json", func(ctx *fasthttp.RequestCtx) {
    ctx.SetContentType("application/json")
    swaggerJSON := `{"swagger":"2.0",...}`
    ctx.SetBodyString(swaggerJSON)
})

// Swagger UI endpoint
router.GET("/swagger/index.html", func(ctx *fasthttp.RequestCtx) {
    ctx.SetContentType("text/html")
    swaggerHTML := `<!DOCTYPE html>...`
    ctx.SetBodyString(swaggerHTML)
})
```

3. **Fixed CDN imports:**
```html
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
<script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
```

## 🧪 Testing

### Test 1: Basic Functionality

```bash
# Test API directly
curl -X POST http://localhost:8080/api/frontendpages \
  -H 'Content-Type: application/json' \
  -d '{
    "metadata": {"name": "test-page", "namespace": "default"},
    "spec": {"contents": "<h1>Test</h1>", "image": "nginx:latest", "replicas": 1}
  }'

# Verify in Kubernetes
kubectl get frontendpage test-page
kubectl get deployment test-page
```

### Test 2: Port.io Integration

1. **Go to Port.io Software Catalog**
2. **Click "🚀 Create Frontend Page"**
3. **Fill form:**
   - Name: `port-test`
   - Contents: `<h1>Created from Port!</h1>`
   - Replicas: `2`
4. **Execute**

**Expected results:**
- ✅ Entity appears in Port.io Software Catalog
- ✅ Status changes from "Creating" to "Ready"
- ✅ Kubernetes resources created
- ✅ Real-time sync of replica count

### Test 3: Auto-sync

```bash
# Scale deployment
kubectl scale deployment port-test --replicas=4

# Check Port.io - should show 4 replicas within ~10 seconds

# Delete resource
kubectl delete frontendpage port-test

# Check Port.io - entity should disappear
```

## 🔧 Troubleshooting

### Common Issues

#### 1. "Destination is a private IP" Error

**Problem:** Port.io cannot reach localhost  
**Solution:** Use ngrok or public endpoint

```bash
ngrok http 8080
# Update webhook URL in Port.io Action
```

#### 2. Authentication Failed

**Problem:** Invalid Port.io credentials  
**Solution:** Verify credentials in Port.io settings

```bash
# Check logs for authentication errors
{"level":"warn","message":"Failed to authenticate with Port.io"}
```

#### 3. Swagger UI Not Loading

**Problem:** Missing swagger dependencies  
**Solution:** Ensure proper imports and CDN URLs

```go
_ "github.com/yourusername/k8s-controller-tutorial/docs"
```

#### 4. Port.io Sync Issues

**Problem:** Entities not updating in Port.io  
**Solution:** Check controller logs for Port API errors

```bash
# Look for Port sync errors
{"level":"warn","message":"Failed to sync to Port.io"}
```

### Debug Mode

```bash
# Run with debug logging
go run main.go server --log-level debug \
  --port 8080 \
  --kubeconfig ~/.kube/config \
  --port-base-url "https://api.getport.io" \
  --port-client-id "$PORT_CLIENT_ID" \
  --port-client-secret "$PORT_CLIENT_SECRET"
```

## 🚀 Advanced Features

### Port.io Scorecards

Add quality gates to your FrontendPages:

```yaml
# In Port.io Blueprint
"scorecards": [
  {
    "identifier": "security",
    "title": "Security Score",
    "rules": [
      {
        "property": "image",
        "operator": "contains",
        "value": "alpine"
      }
    ]
  }
]
```

### Port.io Automations

Trigger actions based on events:

```yaml
# Auto-notify on deployment failures
trigger:
  type: entity-updated
  blueprintIdentifier: frontendpage
  condition:
    property: status
    operator: "="
    value: "Failed"
action:
  type: webhook
  url: https://hooks.slack.com/your-webhook
```

### Production Deployment

```bash
# Build binary
go build -o frontendpage-controller main.go

# Deploy to Kubernetes
kubectl apply -f deployment.yaml

# Use ingress instead of ngrok
kubectl apply -f ingress.yaml
```

### Custom Metrics

Extend the integration with custom metrics:

```go
// Add to Port entity
entity.Properties["cpu_usage"] = getCPUUsage(deployment)
entity.Properties["memory_usage"] = getMemoryUsage(deployment)
entity.Properties["uptime"] = getUptime(deployment)
```

## 📊 Benefits

### Before Integration:
```
Developer → kubectl → YAML → Kubernetes
```

### After Integration:
```
Developer → Port.io UI → Beautiful Form → API → Kubernetes
                ↓
        Real-time status + Metrics + Self-Service
```

## 🎉 Result

**Complete Platform as a Service with:**

✅ **Self-Service Portal** - Developers create resources via UI  
✅ **Real-time Visibility** - Live status and metrics  
✅ **Automated Workflows** - No manual kubectl commands  
✅ **Quality Gates** - Scorecards and compliance  
✅ **Full Observability** - Complete audit trail  

**Your Kubernetes Controller is now a production-ready Internal Developer Platform!** 🚀

## 📝 Contributing

1. Fork the repository
2. Create feature branch
3. Add tests for Port.io integration
4. Submit pull request

## 📄 License

MIT License - see LICENSE file for details.
