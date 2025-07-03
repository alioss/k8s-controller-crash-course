# Atlas Controller Examples

## 1. Setup Namespaces
```bash
kubectl create namespace dev
kubectl create namespace stage
kubectl create namespace prod
kubectl create namespace atlas-system
```

## 2. Install CRD and Controller
```bash
# Install CRD
kubectl apply -f config/crd/atlasapp.yaml

# Create RBAC
kubectl apply -f config/rbac/role.yaml

# Deploy controller
kubectl apply -f config/manager/manager.yaml
```

## 3. Deploy to Dev Environment
```bash
kubectl apply -f examples/dev-deployment.yaml
```

## 4. Check Status
```bash
# Watch AtlasApp resources
kubectl get atlasapp -A -w

# Check deployment status
kubectl get deployment atlas -n dev

# Check controller logs
kubectl logs -n atlas-system deployment/atlas-controller
```

## 5. Auto-promotion to Stage
The controller will automatically promote to stage when dev is ready.

## 6. Manual Promotion to Prod
```bash
# Create prod deployment manually (requires approval)
kubectl apply -f examples/prod-deployment.yaml

# Or approve stage promotion
kubectl patch atlasapp atlas-stage -n stage --type='json' -p='[{"op": "replace", "path": "/spec/requireApproval", "value": false}]'
```
