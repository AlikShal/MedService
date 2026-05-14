# Quick Deployment Guide - After Image Fixes

## 1️⃣ Build Images

```bash
docker compose build
```

This will create:
- medsync-auth-service:latest
- medsync-patient-service:latest
- medsync-doctor-service:latest
- medsync-appointment-service:latest
- medsync-chat-service:latest
- medsync-notification-service:latest
- medsync-frontend:latest

## 2️⃣ Deploy to Kubernetes

```bash
kubectl apply -f k8s/
```

## 3️⃣ Wait for Database

```bash
kubectl wait --for=condition=ready pod -l app=postgres -n medsync --timeout=300s
```

## 4️⃣ Check Status

```bash
kubectl get pods -n medsync
```

**Expected:** All pods in `Running` status with `1/1` Ready

## 5️⃣ Verify No Image Errors

```bash
# Should NOT show any ImagePullBackOff or ErrImagePull
kubectl get pods -n medsync | grep -iE "imagepull|errim" || echo "✓ No image errors"

# Verify all pods are Running
kubectl get pods -n medsync | grep Running | wc -l
# Should show: 10 pods
```

## 🔍 Troubleshooting

### Check specific pod

```bash
kubectl describe pod <pod-name> -n medsync
```

### View pod logs

```bash
kubectl logs <pod-name> -n medsync
```

### Verify image names

```bash
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].image}{"\n"}{end}'
```

### Verify imagePullPolicy

```bash
kubectl get pods -n medsync -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].imagePullPolicy}{"\n"}{end}'
```

All should show: `Never`

## 🎯 Expected Output

### After Step 1 (build)
```bash
$ docker images | grep medsync
medsync-auth-service              latest
medsync-patient-service           latest
medsync-doctor-service            latest
medsync-appointment-service       latest
medsync-chat-service              latest
medsync-notification-service      latest
medsync-frontend                  latest
```

### After Step 4 (check pods)
```bash
$ kubectl get pods -n medsync
NAME                                    READY   STATUS    RESTARTS   AGE
appointment-service-5d7c6f9d5d-2k5p4   1/1     Running   0          2m
appointment-service-5d7c6f9d5d-9x1z2   1/1     Running   0          2m
auth-service-7d4c8f9d5d-4a7b3c         1/1     Running   0          2m
auth-service-7d4c8f9d5d-8d9e5f         1/1     Running   0          2m
chat-service-2c3d4e5f-1g2h3i           1/1     Running   0          2m
chat-service-2c3d4e5f-9j8k7l           1/1     Running   0          2m
doctor-service-3d4e5f6g-2m3n4o         1/1     Running   0          2m
doctor-service-3d4e5f6g-8p9q0r         1/1     Running   0          2m
frontend-4e5f6g7h-3s4t5u              1/1     Running   0          2m
frontend-4e5f6g7h-9v8w0x              1/1     Running   0          2m
grafana-5f6g7h8i-4y5z1a               1/1     Running   0          2m
notification-service-6g7h8i9j-2b3c4d  1/1     Running   0          2m
notification-service-6g7h8i9j-8e9f0g  1/1     Running   0          2m
patient-service-7h8i9j0k-3h4i5j       1/1     Running   0          2m
patient-service-7h8i9j0k-9k0l1m       1/1     Running   0          2m
postgres-8i9j0k1l-4n5o6p             1/1     Running   0          2m
prometheus-9j0k1l2m-5q6r7s            1/1     Running   0          2m
```

## 📋 What Was Fixed

| Issue | Fix |
|-------|-----|
| ImagePullBackOff | Changed imagePullPolicy to `Never` |
| ErrImagePull | Added explicit image names to docker-compose |
| Local images not available | docker-compose now builds medsync-* names |
| Kubernetes couldn't find images | Image names now match between docker-compose and k8s |

## 📚 Related Documentation

- [k8s/IMAGE_FIX.md](IMAGE_FIX.md) - Detailed explanation of fixes
- [k8s/DEPLOYMENT.md](DEPLOYMENT.md) - Complete deployment guide
- [k8s/QUICK_START.md](QUICK_START.md) - Additional reference commands
