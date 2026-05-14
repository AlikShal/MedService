# MedSync Kubernetes Architecture

## Directory Structure

```
k8s/
├── namespace.yaml                 # Kubernetes namespace definition
│
├── configmaps/                    # Configuration data
│   ├── app-config.yaml           # Application settings (GIN_MODE, DB config, service URLs)
│   ├── database-config.yaml      # Database connection settings
│   ├── grafana-config.yaml       # Grafana configuration
│   ├── grafana-dashboards.yaml   # Grafana dashboard definitions
│   ├── grafana-provisioning.yaml # Grafana datasource configuration
│   ├── postgres-init.yaml        # PostgreSQL initialization scripts
│   └── prometheus-config.yaml    # Prometheus scrape configuration
│
├── secrets/                       # Sensitive data (base64 encoded)
│   └── app-secrets.yaml          # Credentials: DB password, JWT secret, admin account
│
├── deployments/                   # Kubernetes Deployments
│   ├── auth-service.yaml         # Authentication service (2 replicas)
│   ├── patient-service.yaml      # Patient management (2 replicas)
│   ├── doctor-service.yaml       # Doctor directory (2 replicas)
│   ├── appointment-service.yaml  # Appointment scheduling (2 replicas)
│   ├── chat-service.yaml         # Chat/messaging (2 replicas)
│   ├── notification-service.yaml # Notifications (2 replicas)
│   ├── frontend.yaml             # Frontend UI (2 replicas)
│   ├── postgres.yaml             # PostgreSQL database (1 replica)
│   ├── prometheus.yaml           # Prometheus monitoring (1 replica)
│   └── grafana.yaml              # Grafana dashboards (1 replica)
│
├── services/                      # Kubernetes Services
│   ├── auth-service.yaml         # ClusterIP: 8082
│   ├── patient-service.yaml      # ClusterIP: 8083
│   ├── doctor-service.yaml       # ClusterIP: 8080
│   ├── appointment-service.yaml  # ClusterIP: 8081
│   ├── chat-service.yaml         # ClusterIP: 8084
│   ├── notification-service.yaml # ClusterIP: 8085
│   ├── frontend.yaml             # NodePort: 30080 → 80
│   ├── postgres.yaml             # ClusterIP: 5432
│   ├── prometheus.yaml           # ClusterIP: 9090
│   └── grafana.yaml              # ClusterIP: 3000
│
├── hpa/                           # Horizontal Pod Autoscalers
│   ├── appointment-service-hpa.yaml  # Min: 2, Max: 5 replicas
│   └── notification-service-hpa.yaml # Min: 2, Max: 4 replicas
│
├── DEPLOYMENT.md                  # Comprehensive deployment guide
├── QUICK_START.md                 # Quick reference commands
├── PROBE_VERIFICATION.md          # Health probe verification guide
└── ARCHITECTURE.md                # This file
```

## Service Architecture

```
┌─────────────────────────────────────────────────────┐
│              Frontend (NodePort: 30080)             │
│                  (2 replicas)                       │
└────────────────────┬────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │   Auth   │  │ Patient  │  │  Doctor  │
  │ Service  │  │ Service  │  │ Service  │
  │   :8082  │  │  :8083   │  │  :8080   │
  │ (2x)     │  │  (2x)    │  │  (2x)    │
  └─────┬────┘  └─────┬────┘  └─────┬────┘
        │             │             │
        └─────────────┼─────────────┘
                      │
              ┌───────┴────────┐
              │                │
              ▼                ▼
        ┌──────────────┐  ┌──────────────┐
        │ Appointment  │  │    Chat      │
        │  Service     │  │  Service     │
        │   :8081      │  │   :8084      │
        │  (HPA: 2-5)  │  │   (2x)       │
        └──────┬───────┘  └──────┬───────┘
               │                 │
        ┌──────┴──────┬──────────┴──────┐
        │             │                 │
        ▼             ▼                 ▼
   ┌──────────┐ ┌──────────┐  ┌─────────────────┐
   │Notif.    │ │PostgreSQL│  │  Monitoring     │
   │Service   │ │Database  │  │ (Prometheus +   │
   │:8085     │ │ :5432    │  │  Grafana)       │
   │(HPA:2-4) │ │ (1x)     │  │                 │
   └──────────┘ └──────────┘  └─────────────────┘
```

## Deployment Strategy

### Phase 1: Infrastructure
1. Namespace creation
2. ConfigMaps & Secrets deployment
3. PostgreSQL deployment
4. Wait for DB readiness

### Phase 2: Core Services
5. Auth-Service (dependency for all)
6. Patient-Service, Doctor-Service
7. Appointment-Service (depends on doctor & patient)
8. Chat-Service (depends on appointment & patient)
9. Notification-Service (stateless, independent)

### Phase 3: Presentation & Monitoring
10. Frontend deployment
11. Prometheus deployment
12. Grafana deployment

### Phase 4: Scaling
13. HPA for appointment-service (critical path)
14. HPA for notification-service (async workload)

## Resource Allocation

### Backend Services
| Service | Memory Request | Memory Limit | CPU Request | CPU Limit |
|---------|---|---|---|---|
| auth-service | 128Mi | 256Mi | 100m | 500m |
| patient-service | 128Mi | 256Mi | 100m | 500m |
| doctor-service | 128Mi | 256Mi | 100m | 500m |
| appointment-service | 256Mi | 512Mi | 200m | 750m |
| chat-service | 128Mi | 256Mi | 100m | 500m |
| notification-service | 64Mi | 128Mi | 50m | 250m |

### Infrastructure Services
| Service | Memory Request | Memory Limit | CPU Request | CPU Limit |
|---------|---|---|---|---|
| frontend | 64Mi | 128Mi | 50m | 200m |
| postgres | 256Mi | 512Mi | 100m | 500m |
| prometheus | 256Mi | 512Mi | 100m | 500m |
| grafana | 128Mi | 256Mi | 100m | 500m |

**Total Minimum Allocation**: ~2.5 GB RAM, 1.5 CPUs

## Health Checks

### Readiness Probes
All services check `/health` endpoint:
- **Initial Delay**: 5s (time for service to start)
- **Period**: 10s (check every 10 seconds)
- **Timeout**: 5s (fail if no response in 5s)
- **Failure Threshold**: 3 (fail after 3 consecutive failures)

**Purpose**: Determines if pod should receive traffic from Services

### Liveness Probes
All services check `/health` endpoint:
- **Initial Delay**: 10s (time for service to start)
- **Period**: 30s (check every 30 seconds)
- **Timeout**: 10s (fail if no response in 10s)
- **Failure Threshold**: 3 (fail after 3 consecutive failures)

**Purpose**: Restarts pod if service becomes unresponsive

## Service Discovery

### DNS Names (ClusterIP Services)
```
<service-name>.<namespace>.svc.cluster.local:<port>

# Examples:
auth-service.medsync.svc.cluster.local:8082
postgres.medsync.svc.cluster.local:5432
prometheus.medsync.svc.cluster.local:9090

# Short form (within same namespace):
auth-service:8082
postgres:5432
prometheus:9090
```

### Service Communication
- **Internal**: Services communicate via ClusterIP (no external exposure)
- **Frontend Access**: NodePort on port 30080
- **Monitoring**: Port-forward for Prometheus (9090) and Grafana (3000)

## Data Persistence

### Current Setup (emptyDir)
- PostgreSQL data: `emptyDir` (ephemeral, lost on pod restart)
- Prometheus data: `emptyDir` (time-series data lost on restart)
- Grafana data: `emptyDir` (dashboards lost on restart)

### Production Recommendation
Use PersistentVolumes:
```yaml
volumeMounts:
- name: postgres-data
  mountPath: /var/lib/postgresql/data
volumes:
- name: postgres-data
  persistentVolumeClaim:
    claimName: postgres-pvc
```

## Horizontal Pod Autoscaling

### Appointment-Service HPA
```
Min Replicas: 2
Max Replicas: 5
Scale-Up: When CPU > 70% or Memory > 80%
Scale-Down: After 5 minutes below threshold, reduce by 50%
```

### Notification-Service HPA
```
Min Replicas: 2
Max Replicas: 4
Scale-Up: When CPU > 75% or Memory > 85%
Scale-Down: After 5 minutes below threshold, reduce by 50%
```

### Other Services
Static 2 replicas for high availability (no HPA).

## ConfigMap Usage

### app-config
- Service environment variables (GIN_MODE, ports)
- Inter-service URLs for HTTP calls
- JWT configuration
- Log levels

### database-config
- Database connection parameters (host, port, name, user)
- SSL mode settings

### grafana-config
- Admin credentials (username)
- Datasource configuration

### prometheus-config
- Scrape targets and intervals
- MedSync service endpoints
- Kubernetes API discovery

## Secret Usage

### app-secrets
- Database password
- JWT secret (encryption key)
- Admin credentials (email, password, full name)
- Grafana admin password

**Note**: Stored as base64-encoded strings in etcd. Use proper RBAC and encryption at rest in production.

## Environment Configuration Flow

```
┌─────────────────────────────────────────┐
│  ConfigMaps (Non-sensitive data)        │
│  - app-config                           │
│  - database-config                      │
│  - grafana-config                       │
└──────────────────┬──────────────────────┘
                   │
                   ├─► valueFrom.configMapKeyRef
                   │
┌──────────────────▼──────────────────────┐
│  Pod Environment Variables              │
└──────────────────┬──────────────────────┘
                   │
                   │
┌──────────────────▼──────────────────────┐
│  Secrets (Sensitive data)               │
│  - app-secrets                          │
└──────────────────┬──────────────────────┘
                   │
                   ├─► valueFrom.secretKeyRef
                   │
┌──────────────────▼──────────────────────┐
│  Pod Environment Variables              │
└─────────────────────────────────────────┘
```

## Deployment Validation

### Prerequisites Check
```bash
✓ Kubernetes cluster running
✓ kubectl configured
✓ Docker images built and available
```

### Deployment Steps
```bash
1. kubectl apply -f k8s/namespace.yaml
2. kubectl apply -f k8s/configmaps/
3. kubectl apply -f k8s/secrets/
4. kubectl apply -f k8s/deployments/postgres.yaml
5. kubectl wait --for=condition=ready pod -l app=postgres
6. kubectl apply -f k8s/deployments/
7. kubectl apply -f k8s/services/
8. kubectl apply -f k8s/hpa/
```

### Verification Checklist
```bash
✓ All pods in Running status
✓ All pods Ready condition = True
✓ No pod restarts (restartCount = 0)
✓ Services have endpoints
✓ Frontend accessible on NodePort:30080
✓ Prometheus scraping targets
✓ Grafana accessible
✓ Inter-service communication working
```

## Troubleshooting Decision Tree

```
Pod Not Running?
├─ CrashLoopBackOff
│  ├─ Check logs: kubectl logs <pod> -n medsync
│  ├─ Check events: kubectl describe pod <pod> -n medsync
│  └─ Check resource limits: is pod OOMKilled?
├─ Pending
│  ├─ Not enough resources: kubectl top nodes
│  ├─ Image not found: check image tag
│  └─ Node selector mismatch: check node labels
└─ ImagePullBackOff
   └─ Image not available in registry

Service Not Accessible?
├─ Check pod Ready status
├─ Check service endpoints: kubectl get endpoints
├─ Check port-forward or NodePort
└─ Test DNS from pod: kubectl exec pod -- nslookup service

Database Issues?
├─ Check postgres pod logs
├─ Verify credentials in secrets
├─ Test connection: kubectl exec postgres-pod -- psql -U postgres
└─ Check readiness probes

Probe Failures?
├─ Service crashing: check logs
├─ Probe timeout too short: increase timeoutSeconds
├─ Health endpoint unreachable: check service startup
└─ Container resource exhaustion: check limits
```

## Comparison: Docker Compose vs Kubernetes

| Aspect | Docker Compose | Kubernetes |
|--------|---|---|
| **Deployment** | `docker compose up` | `kubectl apply -f` |
| **Networking** | Custom bridge | ClusterIP Services |
| **Scaling** | Manual | HPA (automatic) |
| **Health Checks** | healthcheck | readinessProbe + livenessProbe |
| **Resource Limits** | deploy.resources | resources.requests/limits |
| **Data Persistence** | Named volumes | PersistentVolumes |
| **Configuration** | .env files | ConfigMaps + Secrets |
| **Debugging** | docker logs | kubectl logs |
| **High Availability** | Single machine | Multi-node clusters |
| **Use Case** | Local development | Production orchestration |

## Next Steps for Production

1. **Persistent Storage**: Add PersistentVolumes for PostgreSQL
2. **Ingress Controller**: Replace NodePort with Ingress
3. **SSL/TLS**: Add certificates for encrypted communication
4. **Network Policies**: Restrict inter-pod communication
5. **RBAC**: Implement role-based access control
6. **Monitoring**: Configure Prometheus scrape jobs and Grafana dashboards
7. **Logging**: Add centralized logging (ELK, Loki, etc.)
8. **Service Mesh**: Consider Istio for advanced traffic management
9. **Multi-Cluster**: Deploy across regions for disaster recovery
10. **GitOps**: Use ArgoCD or Flux for declarative deployments
