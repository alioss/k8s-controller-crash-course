# Kubernetes Controller Development Course & Atlas Platform

🎓 **Complete journey from learning Kubernetes controllers to building a production-ready platform**

This repository contains both the **step-by-step course materials** for learning Kubernetes controller development and the **final implementation** - Atlas Platform, a complete deployment automation solution.

---

## 📚 Course: Learning Kubernetes Controllers

> **🎯 Educational Content**: Step-by-step tutorial on building Kubernetes controllers

### 📖 [Complete Course Documentation →](./README_course.md)

The course demonstrates how to build a Kubernetes controller from scratch, progressing through:
- ✅ Basic CLI application with Cobra
- ✅ Kubernetes API integration  
- ✅ Informers and reconcilers
- ✅ Production features (metrics, leader election)
- ✅ Platform-as-a-Service with Port.io integration

### 🔗 Reference Repository
📚 **Complete Tutorial Reference**: https://github.com/den-vasyliev/k8s-controller-tutorial-ref

---

## 🚀 Final Result: Atlas Platform

> **🎯 Production Solution**: Complete Kubernetes-native deployment automation platform

### 📖 [Complete Atlas Platform Documentation →](./README_atlas.md)

Based on the knowledge gained from the course, we created **Atlas Platform** - a comprehensive solution for automated application deployment and management across multiple environments.

### 🌟 Platform Components

| Component | Description | Documentation |
|-----------|-------------|---------------|
| **🎯 Atlas Controller** | Kubernetes controller with automated promotion workflows | [📖 Controller Docs](./atlas-controller/README.md) |
| **🔍 atlasctl CLI** | Command-line tool for multi-environment monitoring | [📖 CLI Docs](./atlasctl/README.md) |

### ✨ Key Features
- **Automated Promotion**: dev → stage → prod with approval gates
- **Health Monitoring**: Continuous application health checks
- **Production Safety**: Manual approval required for production deployments
- **Multi-environment View**: Monitor deployments across all clusters

---

## 🎯 Course Journey → Platform Result

| **Course Stage** | **What We Learned** | **Applied in Atlas Platform** |
|------------------|---------------------|-------------------------------|
| **Steps 1-3**: CLI Foundation | Cobra CLI, logging, flags | `atlasctl` command structure |
| **Steps 4-5**: HTTP & Build | FastHTTP server, Makefile, Docker | Controller HTTP endpoints, build system |
| **Steps 6-7**: Kubernetes Integration | Client-go, Informers | Real-time monitoring in `atlasctl` |
| **Steps 8-9**: Controller Runtime | API handlers, reconcilers | `AtlasApp` controller logic |
| **Step 10**: Production Features | Leader election, metrics | High-availability controller |
| **Steps 11-14**: Advanced Features | CRDs, Platform API, Auth | Complete Atlas Platform |

---

## 🚀 Quick Start

### 📚 For Learning (Course)
```bash
# Read the course documentation
open README_course.md

# Follow step-by-step branches
git checkout feature/step1-cobra-cli
git checkout feature/step2-zerolog-logging
# ... continue through all steps
```

### 🎯 For Production (Atlas Platform)
```bash
# Read platform documentation
open README_atlas.md

# Deploy Atlas Controller
cd atlas-controller
kubectl apply -f config/crd/atlasapp.yaml
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/manager/manager.yaml

# Use atlasctl for monitoring
cd ../atlasctl
go build -o atlasctl main.go
./atlasctl list --kubeconfig ~/.kube/config
```

---

## 📁 Repository Structure

```
kubernetes-controller-course/
├── README.md                    # 👈 This overview
├── README_course.md             # 📚 Course documentation
├── README_atlas.md              # 🎯 Atlas Platform documentation  
├── LICENSE                      # 📄 Apache 2.0 License
│
├── atlas-controller/            # 🎯 Production Controller
│   ├── README.md               # Controller documentation
│   ├── api/v1/                 # CRD definitions
│   ├── internal/controller/    # Controller logic
│   ├── config/                 # Kubernetes manifests
│   └── examples/               # Usage examples
│
└── atlasctl/                   # 🔍 CLI Tool
    ├── README.md               # CLI documentation
    ├── cmd/                    # Commands
    ├── pkg/                    # Core logic
    └── examples/               # Usage examples
```

---

## 🎓 Learning Path

### Phase 1: Master the Course
1. 📖 **Study**: [Course Documentation](./README_course.md)
2. 🔄 **Follow**: Step-by-step branches (`feature/step1-*` through `feature/step14-*`)
3. 🏗️ **Build**: Each component progressively
4. 🧪 **Test**: Your understanding at each step

### Phase 2: Deploy the Platform
1. 🎯 **Explore**: [Atlas Platform Documentation](./README_atlas.md)
2. 🚀 **Deploy**: [Atlas Controller](./atlas-controller/README.md)
3. 🔍 **Monitor**: [atlasctl CLI](./atlasctl/README.md)
4. 📊 **Operate**: Applications with automated promotion

---

## 🎯 Success Metrics

### ✅ Course Completion
- [x] Built CLI application with Cobra
- [x] Integrated Kubernetes client-go
- [x] Implemented informers and reconcilers
- [x] Added production features (metrics, leader election)
- [x] Created complete Platform-as-a-Service

### ✅ Platform Implementation
- [x] **Atlas Controller**: Automated deployment workflows
- [x] **atlasctl**: Multi-environment monitoring
- [x] **Production Tested**: Successfully deployed and operational
- [x] **Documentation**: Complete guides for all components

---

## 📖 Documentation Links

| Resource | Description | Link |
|----------|-------------|------|
| **📚 Course Guide** | Complete learning tutorial | [README_course.md](./README_course.md) |
| **🎯 Atlas Platform** | Production platform overview | [README_atlas.md](./README_atlas.md) |
| **🎯 Atlas Controller** | Kubernetes controller docs | [atlas-controller/README.md](./atlas-controller/README.md) |
| **🔍 atlasctl CLI** | Command-line tool docs | [atlasctl/README.md](./atlasctl/README.md) |
| **📝 Examples** | Usage examples | [atlas-controller/examples/](./atlas-controller/examples/) |

---

## 🤝 Contributing

### Course Improvements
- Submit issues for unclear course steps
- Suggest additional learning materials
- Share your learning experience

### Platform Enhancements
- Report bugs in Atlas Controller or atlasctl
- Propose new features for the platform
- Contribute code improvements

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

**Course Foundation**: Kubernetes community and controller development best practices  
**Platform Development**: Built with Kubebuilder, Controller Runtime, and Cobra  
**Inspiration**: GitOps principles and production Kubernetes patterns

---

**Ready to start your journey?** 🚀

**📚 Learn**: [Course Guide](./README_course.md) | **🎯 Deploy**: [Atlas Platform](./README_atlas.md)
