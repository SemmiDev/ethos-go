# Ethos-Go Kubernetes Deployment Guide

This directory contains production-ready Kubernetes manifests for deploying the Ethos-Go stack.

## 📂 Architecture & Components

All manifests are numbered by deployment order:

| #   | File                        | Kind            | Description                                                            |
| --- | --------------------------- | --------------- | ---------------------------------------------------------------------- |
| 01  | `01-namespace.yaml`         | `Namespace`     | Isolates resources in `ethos-go`.                                      |
| 02  | `02-configmap.yaml`         | `ConfigMap`     | Application configuration (DB_HOST, LOG_LEVEL, etc.).                  |
| 03  | `03-secret.yaml`            | `Secret`        | Sensitive data (DB_PASSWORD, JWT_SECRET).                              |
| 04  | `04-postgres.yaml`          | `StatefulSet`   | PostgreSQL 17. PVC ensures data persists across pod restarts.          |
| 05  | `05-redis.yaml`             | `Deployment`    | Redis 8.0. Stateless deployment using ephemeral storage.               |
| 06  | `06-app-deployment.yaml`    | `Deployment`    | Main Go binary serving API & Frontend. **Runs migrations on startup.** |
| 07  | `07-worker-deployment.yaml` | `Deployment`    | Same Go binary running in worker mode (Asynq background tasks).        |
| 08  | `08-service.yaml`           | `Service`       | ClusterIP service for internal access.                                 |
| 09  | `09-gateway.yaml`           | `Gateway`       | Gateway API with HTTPRoute for routing and TLS.                        |
| 10  | `10-hpa.yaml`               | `HPA`           | Horizontal Pod Autoscaler (CPU/Memory based).                          |
| 11  | `11-pdb.yaml`               | `PDB`           | Pod Disruption Budget for high availability.                           |
| 12  | `12-networkpolicy.yaml`     | `NetworkPolicy` | Network isolation between pods.                                        |

> **Note:** Database migrations are embedded in the application and run automatically on startup.

---

## 🚀 Quick Start Deployment

**Prerequisites:**

- Kubernetes Cluster (Docker Desktop, Kind, Minikube, or Cloud).
- `kubectl` configured.
- Docker for building images.

### Option 1: Using Deploy Script (Recommended)

```bash
# Development deployment (skips gateway, HPA, network policies)
./k8s/deploy.sh --dev

# Production deployment (full deployment)
./k8s/deploy.sh --prod
```

### Option 2: Manual Deployment

#### 1. Build & Push Image

```bash
docker build -t sammidev/ethos-go:latest .
docker push sammidev/ethos-go:latest
```

#### 2. Infrastructure Setup

```bash
kubectl apply -f k8s/01-namespace.yaml
kubectl apply -f k8s/02-configmap.yaml
kubectl apply -f k8s/03-secret.yaml
kubectl apply -f k8s/04-postgres.yaml
kubectl apply -f k8s/05-redis.yaml
```

Wait for pods:

```bash
kubectl get pods -n ethos-go -w
```

#### 3. Deploy Application

```bash
kubectl apply -f k8s/06-app-deployment.yaml
kubectl apply -f k8s/07-worker-deployment.yaml
kubectl apply -f k8s/08-service.yaml
```

#### 4. Access (Development)

```bash
kubectl port-forward svc/ethos-go-app 8080:80 -n ethos-go
```

Open **[http://localhost:8080](http://localhost:8080)**.

---

## 🏭 Production Deployment

### Gateway API (Modern Alternative to Ingress)

The project uses [Gateway API](https://gateway-api.sigs.k8s.io/) instead of legacy Ingress for production traffic routing.

**Prerequisites:**

```bash
# 1. Install Gateway API CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/standard-install.yaml

# 2. Install a Gateway Controller (choose one):
# - Envoy Gateway: https://gateway.envoyproxy.io/
# - Cilium: https://cilium.io/
# - Istio: https://istio.io/
# - NGINX Gateway Fabric: https://github.com/nginxinc/nginx-gateway-fabric
```

**Deploy Gateway:**

```bash
# Update hostname in 09-gateway.yaml first!
kubectl apply -f k8s/09-gateway.yaml
```

### Full Production Deployment

```bash
# Gateway API resources
kubectl apply -f k8s/09-gateway.yaml

# Auto-scaling
kubectl apply -f k8s/10-hpa.yaml
kubectl apply -f k8s/11-pdb.yaml

# Network policies (requires compatible CNI)
kubectl apply -f k8s/12-networkpolicy.yaml
```

### Production Requirements

| Component         | Requirement                                                 |
| ----------------- | ----------------------------------------------------------- |
| **Gateway**       | Gateway API CRDs + Gateway Controller (Envoy, Cilium, etc.) |
| **TLS**           | cert-manager with `letsencrypt-prod` cluster issuer         |
| **HPA**           | metrics-server installed                                    |
| **NetworkPolicy** | CNI that supports NetworkPolicy (Calico, Cilium)            |

---

## 🔍 Debugging & Troubleshooting

### Pod Health

| Status             | Likely Cause                                          | Debug Action                     |
| ------------------ | ----------------------------------------------------- | -------------------------------- |
| `CrashLoopBackOff` | Application panic, DB connection failure, missing env | `kubectl logs <pod> -n ethos-go` |
| `ImagePullBackOff` | Image name typo, tag missing, registry auth failed    | `kubectl describe pod <pod>`     |
| `Pending`          | No nodes, PVC not bound, insufficient resources       | `kubectl describe pod <pod>`     |

### Log Analysis

```bash
kubectl logs -f -l app=ethos-go -n ethos-go
kubectl logs -f -l app=ethos-go-worker -n ethos-go
```

### Gateway Troubleshooting

```bash
# Check Gateway status
kubectl get gateway,httproute -n ethos-go

# Describe Gateway for events
kubectl describe gateway ethos-go-gateway -n ethos-go
```

---

## 🔄 Maintenance

### Updating the Application

```bash
docker build -t sammidev/ethos-go:latest .
docker push sammidev/ethos-go:latest
kubectl rollout restart deployment -n ethos-go
```

### Scaling

```bash
kubectl scale deployment ethos-go-worker --replicas=3 -n ethos-go
```

---

## 🧹 Clean Up

```bash
kubectl delete namespace ethos-go
```
