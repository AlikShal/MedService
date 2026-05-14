# Kubernetes Deployment Guide

## Overview

This minimal Kubernetes setup provides an academically complete deployment of the MedSync platform. It preserves the existing Docker Compose environment while adding clean Kubernetes orchestration.

## Prerequisites

### 1. Docker Images

Build and push service images with the correct tags:

```bash
# From project root, build images with proper tags
docker build -t medsync-auth-service:latest ./auth-service
docker build -t medsync-patient-service:latest ./patient-service
docker build -t medsync-doctor-service:latest ./doctor-service
docker build -t medsync-appointment-service:latest ./appointment-service
docker build -t medsync-chat-service:latest ./chat-service
docker build -t medsync-notification-service:latest ./notification-service
docker build -t medsync-frontend:latest ./frontend
```

### 2. Kubernetes Cluster

Choose one of these options:

#### Option A: Docker Desktop Kubernetes (Recommended for Mac/Windows)

1. Enable Kubernetes in Docker Desktop settings
2. Verify cluster:
   ```bash
   kubectl cluster-info
   kubectl version
   ```

#### Option B: Minikube (Linux/Mac/Windows)

1. Install minikube:
   ```bash
   # macOS
   brew install minikube

   # Linux
   curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
   sudo install minikube-linux-amd64 /usr/local/bin/minikube

   # Windows (PowerShell, with admin)
   choco install minikube
   ```

2. Start cluster:
   ```bash
   minikube start --memory=4096 --cpus=4
   minikube dashboard  # (optional) Opens web UI
   ```

3. Point Docker to Minikube:
   ```bash
   eval $(minikube docker-env)
   # Then rebuild images for Minikube
   ```

## Deployment

### 1. Load Images into Kubernetes (Minikube only)

```bash
# If using Minikube with IfNotPresent pull policy
eval $(minikube docker-env)
docker build -t medsync-auth-service:latest ./auth-service
# ... rebuild all services
```

### 2. Create Namespace and Secrets

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Verify namespace
kubectl get namespace
```

### 3. Deploy ConfigMaps and Secrets

```bash
# Create all ConfigMaps
kubectl apply -f k8s/configmaps/

# Create Secrets
kubectl apply -f k8s/secrets/

# Verify
kubectl get configmaps -n medsync
kubectl get secrets -n medsync
```

### 4. Deploy Core Infrastructure

```bash
# Deploy PostgreSQL first (services depend on it)
kubectl apply -f k8s/deployments/postgres.yaml
kubectl apply -f k8s/services/postgres.yaml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n medsync --timeout=300s
```

### 5. Deploy Backend Services

```bash
# Deploy all backend services and their services
kubectl apply -f k8s/deployments/auth-service.yaml
kubectl apply -f k8s/services/auth-service.yaml

kubectl apply -f k8s/deployments/patient-service.yaml
kubectl apply -f k8s/services/patient-service.yaml

kubectl apply -f k8s/deployments/doctor-service.yaml
kubectl apply -f k8s/services/doctor-service.yaml

kubectl apply -f k8s/deployments/appointment-service.yaml
kubectl apply -f k8s/services/appointment-service.yaml

kubectl apply -f k8s/deployments/chat-service.yaml
kubectl apply -f k8s/services/chat-service.yaml

kubectl apply -f k8s/deployments/notification-service.yaml
kubectl apply -f k8s/services/notification-service.yaml
```

### 6. Deploy Monitoring Stack

```bash
# Deploy Prometheus
kubectl apply -f k8s/deployments/prometheus.yaml
kubectl apply -f k8s/services/prometheus.yaml

# Deploy Grafana
kubectl apply -f k8s/deployments/grafana.yaml
kubectl apply -f k8s/services/grafana.yaml
```

### 7. Deploy Frontend

```bash
kubectl apply -f k8s/deployments/frontend.yaml
kubectl apply -f k8s/services/frontend.yaml
```

### 8. Deploy Horizontal Pod Autoscalers

```bash
# Deploy HPA for critical services
kubectl apply -f k8s/hpa/appointment-service-hpa.yaml
kubectl apply -f k8s/hpa/notification-service-hpa.yaml
```

### Complete Deployment (All-in-One)

```bash
# Deploy everything at once
kubectl apply -f k8s/

# Verify all resources created
kubectl get all -n medsync
```

## Verification and Validation

### 1. Check Pod Status

```bash
# List all pods
kubectl get pods -n medsync

# List pods with more details
kubectl get pods -n medsync -o wide

# Watch pods in real-time
kubectl get pods -n medsync -w

# Expected output: all pods in Running state
```

### 2. Check Services

```bash
# List all services
kubectl get svc -n medsync

# Get service details
kubectl describe svc frontend -n medsync

# Port-forward to frontend for testing
kubectl port-forward svc/frontend 8080:80 -n medsync
# Then access http://localhost:8080
```

### 3. Check Deployment Status

```bash
# Check deployment status
kubectl get deployments -n medsync

# Detailed deployment status
kubectl describe deployment auth-service -n medsync

# Rollout status
kubectl rollout status deployment/auth-service -n medsync
```

### 4. Check Pod Readiness and Liveness

```bash
# View probe status
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.containerStatuses[0].readinessProbes}{"\n"}{end}'

# Describe pod for detailed probe information
kubectl describe pod <pod-name> -n medsync
```

### 5. Check Container Logs

```bash
# View pod logs
kubectl logs <pod-name> -n medsync

# Follow logs in real-time
kubectl logs -f <pod-name> -n medsync

# View logs from all pods in a deployment
kubectl logs -l app=auth-service -n medsync

# View previous logs (if container restarted)
kubectl logs <pod-name> --previous -n medsync
```

### 6. Check Resource Usage

```bash
# View resource metrics (requires metrics-server)
kubectl top nodes
kubectl top pods -n medsync
kubectl top pod <pod-name> -n medsync
```

### 7. Check HPA Status

```bash
# View HPA status
kubectl get hpa -n medsync

# Detailed HPA information
kubectl describe hpa appointment-service-hpa -n medsync

# Watch HPA scaling in real-time
kubectl get hpa -n medsync -w
```

### 8. Check ConfigMaps and Secrets

```bash
# List ConfigMaps
kubectl get configmaps -n medsync

# View ConfigMap content
kubectl get configmap app-config -n medsync -o yaml

# List Secrets
kubectl get secrets -n medsync

# Note: Secrets are base64 encoded; decode with:
kubectl get secret app-secrets -n medsync -o jsonpath='{.data.DB_PASSWORD}' | base64 -d
```

### 9. Test Service-to-Service Communication

```bash
# Exec into a pod and test connectivity
kubectl exec -it <pod-name> -n medsync -- /bin/sh

# Inside pod, test connectivity to other services
wget -q -O- http://auth-service:8082/health
wget -q -O- http://patient-service:8083/health
wget -q -O- http://doctor-service:8080/health
```

### 10. Access Services

```bash
# Frontend (NodePort)
# For Docker Desktop: http://localhost:30080
# For Minikube: http://$(minikube ip):30080

# Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n medsync
# Access: http://localhost:9090

# Grafana
kubectl port-forward svc/grafana 3000:3000 -n medsync
# Access: http://localhost:3000
# Credentials: admin / admin
```

## Probe Verification

### Readiness Probes

Verify services are ready to accept traffic:

```bash
# Readiness probes check the /health endpoint
# Services respond with 200 OK when ready

# Test manually
kubectl exec <pod-name> -n medsync -- wget -q -O- http://localhost:8082/health

# Check probe status in pod description
kubectl describe pod <pod-name> -n medsync | grep -A 5 "Readiness"
```

### Liveness Probes

Verify pods are alive and restart if not:

```bash
# Liveness probes ensure pod restarts if service dies
# Configured with 10s initialDelaySeconds, 30s periodSeconds

# Monitor pod restarts
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.containerStatuses[0].restartCount}{"\n"}{end}'

# Check liveness status
kubectl describe pod <pod-name> -n medsync | grep -A 5 "Liveness"
```

## Scaling

### Manual Scaling

```bash
# Scale auth-service to 5 replicas
kubectl scale deployment auth-service --replicas=5 -n medsync

# Scale back down
kubectl scale deployment auth-service --replicas=2 -n medsync

# Verify scaling
kubectl get deployment auth-service -n medsync
```

### Automatic Scaling (HPA)

```bash
# Monitor HPA scaling
kubectl get hpa -n medsync -w

# When appointment-service CPU > 70% or memory > 80%, it scales up
# When below thresholds, it scales down (with 5-minute stabilization)
```

## Troubleshooting

### Pod Not Starting

```bash
# Check pod status and events
kubectl describe pod <pod-name> -n medsync

# Check logs
kubectl logs <pod-name> -n medsync

# Common issues:
# - Image not found: ensure image tags match Deployment
# - CrashLoopBackOff: check logs for application errors
# - Pending: check node resources with `kubectl top nodes`
```

### Service Connectivity Issues

```bash
# Check service endpoints
kubectl get endpoints -n medsync

# Test DNS resolution (from within a pod)
kubectl exec -it <pod-name> -n medsync -- nslookup auth-service

# Test connectivity
kubectl exec -it <pod-name> -n medsync -- wget -q -O- http://auth-service:8082/health
```

### Database Connection Issues

```bash
# Check PostgreSQL pod logs
kubectl logs -l app=postgres -n medsync

# Verify database credentials in secrets
kubectl get secret app-secrets -n medsync -o yaml

# Test PostgreSQL connectivity
kubectl exec -it <postgres-pod> -n medsync -- psql -U postgres -d medical_platform -c "SELECT VERSION();"
```

### Probe Failures

```bash
# Identify unhealthy pods
kubectl get pods -n medsync | grep -v Running

# Check probe configuration
kubectl get pod <pod-name> -n medsync -o jsonpath='{.spec.containers[*].readinessProbe}'

# Manually test probe endpoint
kubectl exec <pod-name> -n medsync -- wget -q -O- http://localhost:8082/health
```

## Common Commands

```bash
# All-in-one deployment
kubectl apply -f k8s/

# Full cleanup
kubectl delete namespace medsync

# Restart a service
kubectl rollout restart deployment/auth-service -n medsync

# View resource requests/limits
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].resources}{"\n"}{end}'

# Export configuration
kubectl get all -n medsync -o yaml > medsync-backup.yaml

# Monitor events
kubectl get events -n medsync --sort-by='.lastTimestamp'
```

## Architecture Notes

### Service Discovery

- Services communicate via DNS: `<service-name>:<port>`
- Example: `http://auth-service:8082`
- Automatic DNS resolution within the cluster

### Resource Limits

- Backend services: 128-256Mi memory, 100-500m CPU
- Infrastructure: 64-256Mi memory, 50-500m CPU
- Based on docker-compose resource limits

### Scaling Strategy

- **Appointment Service**: HPA min=2, max=5 (critical path)
- **Notification Service**: HPA min=2, max=4 (asynchronous)
- Other services: static 2 replicas for availability

### Data Persistence

- PostgreSQL uses `emptyDir` (data lost on pod restart)
- For persistent data: mount PersistentVolume to `/var/lib/postgresql/data`
- See ADVANCED.md for persistent storage setup

## Next Steps

1. **Persistence**: Add PersistentVolumes for PostgreSQL
2. **Ingress**: Add Ingress Controller for production routing
3. **Monitoring**: Configure Prometheus scrape jobs for metrics
4. **Networking**: Add NetworkPolicies for security
5. **RBAC**: Implement Role-Based Access Control
6. **Custom Dashboards**: Build Grafana dashboards for metrics
