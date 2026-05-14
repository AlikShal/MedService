# Kubernetes Probe Verification Guide

This guide explains how to verify that readiness and liveness probes are working correctly in the MedSync Kubernetes deployment.

## Probe Configuration

All services are configured with health checks:

### Readiness Probe
- **Path**: `/health` endpoint
- **Initial Delay**: 5 seconds
- **Period**: 10 seconds
- **Timeout**: 5 seconds
- **Failure Threshold**: 3

### Liveness Probe
- **Path**: `/health` endpoint
- **Initial Delay**: 10 seconds
- **Period**: 30 seconds
- **Timeout**: 10 seconds
- **Failure Threshold**: 3

## Verify Probes are Configured

### 1. Check Probe Configuration in Pod

```bash
# View complete probe configuration
kubectl get pod <pod-name> -n medsync -o jsonpath='{.spec.containers[0].readinessProbe}'

# Example output:
# {"httpGet":{"path":"/health","port":8082},"initialDelaySeconds":5,"periodSeconds":10,"timeoutSeconds":5,"failureThreshold":3}
```

### 2. Check All Pods' Probe Status

```bash
# List all pods with probe info
kubectl describe pod -n medsync | grep -A 5 "Readiness\|Liveness"
```

## Verify Probes are Working

### 1. Check Pod Ready Status

```bash
# View pod ready status (should show 1/1 Ready)
kubectl get pods -n medsync

# Example output:
# NAME                                 READY   STATUS    RESTARTS   AGE
# auth-service-7d4c8f9d5d-2k5p4      1/1     Running   0          2m
```

### 2. Check Ready Condition

```bash
# Show detailed ready status for all pods
kubectl get pods -n medsync -o custom-columns=NAME:.metadata.name,READY:.status.conditions[?(@.type==\"Ready\")].status,REASON:.status.conditions[?(@.type==\"Ready\")].reason

# Example output:
# NAME                                   READY   REASON
# auth-service-7d4c8f9d5d-2k5p4        True    
# auth-service-7d4c8f9d5d-9x1z2        True    
```

### 3. Manually Test /health Endpoint

```bash
# Test health endpoint from within a pod
kubectl exec -it <pod-name> -n medsync -- /bin/sh

# Inside the pod, test the health endpoint
wget -q -O- http://localhost:8082/health
# Expected: JSON response indicating service is healthy

# Or using curl (if available)
curl -s http://localhost:8082/health
```

### 4. Check Container Readiness Status

```bash
# View detailed container readiness information
kubectl get pod <pod-name> -n medsync -o jsonpath='{.status.containerStatuses[0].ready}'

# Should return: true
```

### 5. Check Ready Condition Details

```bash
# Get full pod status including conditions
kubectl describe pod <pod-name> -n medsync | grep -A 20 "Status:"

# Look for:
# Type              Status
# ----              ------
# Initialized       True
# Ready             True          <-- Should be True
# ContainersReady   True
# PodScheduled      True
```

## Verify Liveness Probes

### 1. Check Restart Count

```bash
# View restart count (should be 0 or low)
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.containerStatuses[0].restartCount}{"\n"}{end}'

# High restart count indicates liveness probe failures
```

### 2. Check Pod Events

```bash
# View pod events (shows probe activity)
kubectl describe pod <pod-name> -n medsync | grep -A 10 "Events:"

# Look for:
# Liveness probe succeeded
# Readiness probe succeeded
# Successfully assigned to node
```

### 3. Monitor Probe Activity

```bash
# Watch events in real-time
kubectl get events -n medsync --sort-by='.lastTimestamp' -w

# Filter for probe events
kubectl get events -n medsync --field-selector type=Normal | grep -i probe
```

## Simulating Probe Failures

### 1. Stop Service to Trigger Probe Failure

```bash
# Exec into pod
kubectl exec -it <pod-name> -n medsync -- /bin/sh

# Kill the service process (if Go, look for the binary)
kill -9 $(pgrep -f 'auth-service|patient-service|doctor-service|appointment-service|chat-service|notification-service')

# Exit the pod
exit

# Wait ~30 seconds and watch the pod restart
kubectl get pods -n medsync -w
```

### 2. Observe Pod Recovery

```bash
# After killing the service:
# - Pod status changes to CrashLoopBackOff
# - Restart count increases
# - Pod automatically restarts

# Monitor restart
watch kubectl get pod <pod-name> -n medsync -o wide

# Check logs to see service starting again
kubectl logs -f <pod-name> -n medsync
```

## Probe Success Indicators

### What "Working Correctly" Looks Like

1. **Pod Status**: `1/1 Running`
2. **Ready Condition**: `True`
3. **Restart Count**: `0` (or very low)
4. **Recent Events**: Include successful probe messages
5. **No errors**: No `Liveness probe failed` or `Readiness probe failed` in events

### Verify All Services

```bash
# Check all pods at once
kubectl get pods -n medsync \
  -o custom-columns=NAME:.metadata.name,READY:.status.conditions[?(@.type==\"Ready\")].status,STATUS:.status.phase,RESTARTS:.status.containerStatuses[0].restartCount,AGE:.metadata.creationTimestamp

# All pods should show:
# READY: True
# STATUS: Running
# RESTARTS: 0
```

## Resource Requests/Limits Verification

Verify probes have adequate resources to run:

```bash
# Check resource allocation
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].resources}{"\n"}{end}' | grep -o '"requests":[^}]*' | head -5

# Verify services have CPU/memory allocated for probes
# Resources should be sufficient for probe execution
```

## Advanced Probe Troubleshooting

### 1. Check Probe Timeout Behavior

```bash
# Add logging to probe behavior
kubectl logs <pod-name> -n medsync | grep -i 'health\|ready\|liveness'
```

### 2. Increase Verbosity for Probe Debug

```bash
# View kubelet logs (if accessible)
# Note: kubelet runs on each node

# For Minikube:
minikube ssh
docker logs $(docker ps | grep kubelet | awk '{print $1}')
```

### 3. Simulate Slow Response

```bash
# If probe timeout needs adjustment, modify deployment:
kubectl edit deployment auth-service -n medsync

# Look for:
# readinessProbe:
#   httpGet:
#     path: /health
#     port: 8082
#   timeoutSeconds: 5  <-- Increase if needed
```

## Probe Configuration Best Practices

### Current Configuration
- **Readiness**: 5s initial delay, 10s period → fast feedback
- **Liveness**: 10s initial delay, 30s period → prevents unnecessary restarts
- **Failure threshold**: 3 → allows temporary failures

### When to Adjust

```bash
# If probes fail constantly:
# 1. Increase initialDelaySeconds (give app time to start)
# 2. Check application logs
# 3. Verify /health endpoint is responding

# If probes are too aggressive:
# 1. Increase periodSeconds
# 2. Increase failureThreshold
# 3. Increase timeoutSeconds

# Edit deployment to adjust:
kubectl edit deployment <deployment-name> -n medsync
```

## Summary Verification Command

```bash
# Single command to verify all probes are working:
echo "=== Pod Status ===" && \
kubectl get pods -n medsync && \
echo -e "\n=== Ready Status ===" && \
kubectl get pods -n medsync -o custom-columns=NAME:.metadata.name,READY:.status.conditions[?(@.type==\"Ready\")].status && \
echo -e "\n=== Restart Count ===" && \
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.containerStatuses[0].restartCount}{"\n"}{end}' && \
echo -e "\n=== Recent Probe Events ===" && \
kubectl get events -n medsync --field-selector type=Normal | grep -i probe | tail -5
```

If all outputs show:
- Pods in `Running` status
- `Ready` = `True`
- `Restart Count` = 0
- Recent probe success events

**Then probes are working correctly! ✓**
