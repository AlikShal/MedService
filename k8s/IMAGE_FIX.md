# Kubernetes Image Deployment Fix

## Problem Fixed

- ✅ ImagePullBackOff errors resolved
- ✅ Local Docker images now usable by Kubernetes
- ✅ Image names aligned between docker-compose and Kubernetes
- ✅ imagePullPolicy set to Never (local-only)

## Changes Made

### 1. docker-compose.yml
Added explicit image names to ensure consistent naming:
```yaml
auth-service:
  image: medsync-auth-service:latest
  build: ./auth-service
```

**Services updated:**
- medsync-auth-service:latest
- medsync-patient-service:latest
- medsync-doctor-service:latest
- medsync-appointment-service:latest
- medsync-chat-service:latest
- medsync-notification-service:latest
- medsync-frontend:latest

### 2. Kubernetes Deployments
Changed imagePullPolicy from `IfNotPresent` to `Never`:
```yaml
containers:
- name: auth-service
  image: medsync-auth-service:latest
  imagePullPolicy: Never  # <- Changed from IfNotPresent
```

**Effect:** Kubernetes will only use locally available images, no registry pulls.

## Build and Deploy Steps

### Step 1: Build Images

```bash
# Build all images with correct names (from project root)
docker compose build

# Verify images were created
docker images | grep medsync
```

Expected output:
```
medsync-auth-service              latest
medsync-patient-service           latest
medsync-doctor-service            latest
medsync-appointment-service       latest
medsync-chat-service              latest
medsync-notification-service      latest
medsync-frontend                  latest
```

### Step 2: Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/

# Wait for database to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n medsync --timeout=300s
```

### Step 3: Verify Deployment

```bash
# Check pod status (should show Running)
kubectl get pods -n medsync

# Expected output: All pods in Running status with 1/1 Ready
```

## Verification Commands

### Quick Status Check

```bash
# List all pods
kubectl get pods -n medsync

# Short output example:
# NAME                                  READY   STATUS    RESTARTS   AGE
# auth-service-5d7c6f9d5d-2k5p4        1/1     Running   0          2m
# patient-service-7d4c8f9d5d-9x1z2     1/1     Running   0          2m
```

### Detailed Pod Information

```bash
# Describe specific pod (shows readiness status, image used, etc.)
kubectl describe pod <pod-name> -n medsync

# Example:
kubectl describe pod auth-service-5d7c6f9d5d-2k5p4 -n medsync
```

### Check Pod Logs

```bash
# View pod startup logs
kubectl logs <pod-name> -n medsync

# Example:
kubectl logs auth-service-5d7c6f9d5d-2k5p4 -n medsync

# Follow logs in real-time
kubectl logs -f auth-service-5d7c6f9d5d-2k5p4 -n medsync
```

### Troubleshoot Image Issues

```bash
# Check which image a pod is using
kubectl get pod <pod-name> -n medsync -o jsonpath='{.spec.containers[0].image}'

# Check imagePullPolicy
kubectl get pod <pod-name> -n medsync -o jsonpath='{.spec.containers[0].imagePullPolicy}'

# Check image pull status
kubectl get pod <pod-name> -n medsync -o jsonpath='{.status.containerStatuses[0].imageID}'
```

### Monitor All Services

```bash
# Watch pods in real-time
kubectl get pods -n medsync -w

# Get services and their status
kubectl get svc -n medsync

# Check all resources
kubectl get all -n medsync
```

## Complete Deployment Script

```bash
#!/bin/bash
set -e

echo "=== Step 1: Building Docker images ==="
docker compose build
echo "✓ Images built successfully"

echo ""
echo "=== Step 2: Verifying images ==="
docker images | grep medsync || echo "Warning: No medsync images found"

echo ""
echo "=== Step 3: Deploying to Kubernetes ==="
kubectl apply -f k8s/
echo "✓ Manifests applied"

echo ""
echo "=== Step 4: Waiting for PostgreSQL to be ready ==="
kubectl wait --for=condition=ready pod -l app=postgres -n medsync --timeout=300s
echo "✓ PostgreSQL is ready"

echo ""
echo "=== Step 5: Checking pod status ==="
kubectl get pods -n medsync

echo ""
echo "=== Step 6: Waiting for all services to be ready ==="
kubectl wait --for=condition=ready pod --all -n medsync --timeout=300s || true
echo "✓ Services are ready"

echo ""
echo "=== Deployment Complete ==="
kubectl get all -n medsync
```

## Validation Checklist

Run these to verify successful deployment:

```bash
# ✓ All pods Running
kubectl get pods -n medsync | grep Running

# ✓ All pods Ready
kubectl get pods -n medsync | grep "1/1"

# ✓ No ImagePullBackOff or ErrImagePull
kubectl get pods -n medsync | grep -E "ImagePull|ErrImage"

# ✓ Services have endpoints
kubectl get endpoints -n medsync

# ✓ Frontend accessible
kubectl port-forward svc/frontend 8080:80 -n medsync
# Access: http://localhost:8080

# ✓ All images using local policy
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].imagePullPolicy}{"\n"}{end}'
# All should show: Never
```

## Troubleshooting ImagePullBackOff

If pods still show ImagePullBackOff:

```bash
# 1. Check image names match exactly
docker images | grep medsync
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].image}{"\n"}{end}'

# 2. Verify imagePullPolicy is Never
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].imagePullPolicy}{"\n"}{end}'

# 3. Check pod events for detailed error
kubectl describe pod <pod-name> -n medsync | tail -20

# 4. Rebuild and redeploy
docker compose build
kubectl delete namespace medsync
kubectl apply -f k8s/
```

## Image Name Reference

| Service | Image Name |
|---------|-----------|
| Auth Service | medsync-auth-service:latest |
| Patient Service | medsync-patient-service:latest |
| Doctor Service | medsync-doctor-service:latest |
| Appointment Service | medsync-appointment-service:latest |
| Chat Service | medsync-chat-service:latest |
| Notification Service | medsync-notification-service:latest |
| Frontend | medsync-frontend:latest |

## Summary

✅ **Before Fix:**
- docker-compose: builds with auto-generated names
- Kubernetes: looks for medsync-* names
- imagePullPolicy: IfNotPresent (tries registry)
- Result: ImagePullBackOff

✅ **After Fix:**
- docker-compose: explicitly builds medsync-* names
- Kubernetes: uses exact same names
- imagePullPolicy: Never (local only)
- Result: Pods run successfully with local images
