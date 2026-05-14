# Kubernetes One-Liner Deployment & Verification

## Quick Start (Docker Desktop)

### Build Images
```bash
docker build -t medsync-auth-service:latest ./auth-service && \
docker build -t medsync-patient-service:latest ./patient-service && \
docker build -t medsync-doctor-service:latest ./doctor-service && \
docker build -t medsync-appointment-service:latest ./appointment-service && \
docker build -t medsync-chat-service:latest ./chat-service && \
docker build -t medsync-notification-service:latest ./notification-service && \
docker build -t medsync-frontend:latest ./frontend
```

### Deploy All Resources
```bash
kubectl apply -f k8s/ && kubectl wait --for=condition=ready pod -l app=postgres -n medsync --timeout=300s
```

### Verify Deployment
```bash
kubectl get pods -n medsync && kubectl get svc -n medsync
```

### Access Services
```bash
# Frontend
kubectl port-forward svc/frontend 8080:80 -n medsync

# Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n medsync

# Grafana
kubectl port-forward svc/grafana 3000:3000 -n medsync
```

## Verification Commands

```bash
# Get all pods status
kubectl get pods -n medsync

# Check deployments
kubectl get deployments -n medsync

# Check services
kubectl get services -n medsync

# Check HPA status
kubectl get hpa -n medsync

# View detailed pod info
kubectl describe pod -n medsync

# View pod logs
kubectl logs <pod-name> -n medsync

# Check resource usage (requires metrics-server)
kubectl top pods -n medsync

# Wait for all pods to be ready
kubectl wait --for=condition=ready pod --all -n medsync --timeout=300s

# Port-forward to frontend
kubectl port-forward svc/frontend 8080:80 -n medsync

# Port-forward to Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n medsync

# Port-forward to Grafana
kubectl port-forward svc/grafana 3000:3000 -n medsync

# Delete all resources
kubectl delete namespace medsync

# Watch pods in real-time
kubectl get pods -n medsync -w

# Check probe status
kubectl get pods -n medsync -o custom-columns=NAME:.metadata.name,READY:.status.conditions[?(@.type==\"Ready\")].status,RESTARTS:.status.containerStatuses[0].restartCount

# Describe a specific pod (e.g., auth-service)
kubectl describe pod -l app=auth-service -n medsync

# View all events
kubectl get events -n medsync

# Check HPA scaling
kubectl get hpa -n medsync -w

# Scale a deployment
kubectl scale deployment auth-service --replicas=3 -n medsync

# Check ConfigMaps
kubectl get configmaps -n medsync

# Check Secrets
kubectl get secrets -n medsync

# Restart a deployment
kubectl rollout restart deployment/auth-service -n medsync

# Check rollout status
kubectl rollout status deployment/auth-service -n medsync
```

## Monitoring

```bash
# Watch all pods
watch kubectl get pods -n medsync

# Monitor resource usage
kubectl top pods -n medsync

# Monitor HPA
watch kubectl get hpa -n medsync

# Check logs in real-time
kubectl logs -f <pod-name> -n medsync

# Monitor events
watch kubectl get events -n medsync
```

## Troubleshooting

```bash
# Why isn't pod starting?
kubectl describe pod <pod-name> -n medsync

# Check pod logs
kubectl logs <pod-name> -n medsync

# Check previous logs if crashed
kubectl logs <pod-name> --previous -n medsync

# Get detailed pod info
kubectl get pod <pod-name> -n medsync -o yaml

# Check service endpoints
kubectl get endpoints -n medsync

# Test DNS resolution
kubectl exec -it <pod-name> -n medsync -- nslookup auth-service

# Test service connectivity
kubectl exec -it <pod-name> -n medsync -- wget -q -O- http://auth-service:8082/health
```

## Cleanup

```bash
# Delete everything
kubectl delete namespace medsync

# Or delete specific resource
kubectl delete deployment auth-service -n medsync
```
