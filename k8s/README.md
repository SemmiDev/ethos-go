# Ethos-Go Kubernetes Deployment Guide

This directory contains the production-ready Kubernetes manifests for deploying the Ethos-Go stack.

## 📂 Architecture & Components

All manifests are located in the `k8s/` directory.

| Component      | Kind          | File                     | Description                                                      |
| -------------- | ------------- | ------------------------ | ---------------------------------------------------------------- |
| **Namespace**  | `Namespace`   | `namespace.yaml`         | Isolates resources in `ethos-go`.                                |
| **Config**     | `ConfigMap`   | `configmap.yaml`         | Application configuration (e.g., DB_HOST, LOG_LEVEL).            |
| **Security**   | `Secret`      | `secret.yaml`            | Sensitive data (DB_PASSWORD, JWT_SECRET).                        |
| **Database**   | `StatefulSet` | `postgres.yaml`          | PostgreSQL 17. PVC ensures data persists across pod restarts.    |
| **Cache**      | `Deployment`  | `redis.yaml`             | Redis 8.0. Stateless deployment using ephemeral storage.         |
| **API/Web**    | `Deployment`  | `app-deployment.yaml`    | Main Go binary serving API & Frontend. Exposed via LoadBalancer. |
| **Worker**     | `Deployment`  | `worker-deployment.yaml` | Same Go binary running in worker mode (Asynq background tasks).  |
| **Migrations** | `Job`         | `migration-job.yaml`     | Ephemeral job to run schema migrations via `golang-migrate`.     |

---

## 🚀 Quick Start Deployment

**Prerequisites:**

- Kubernetes Cluster (Docker Desktop, Kind, Minikube, or Cloud).
- `kubectl` configured.
- Docker for building images.

### 1. Build & Push Image

The manifests are configured to pull `sammidev/ethos-go:latest`.

```bash
# Build
docker build -t sammidev/ethos-go:latest .

# Push (Required for cluster to pull)
docker push sammidev/ethos-go:latest
```

### 2. Infrastructure Setup

Deploy the foundation first.

```bash
# Namespace
kubectl apply -f k8s/namespace.yaml

# Configuration
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml

# Data Services
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
```

**Wait** until Postgres and Redis are running:

```bash
kubectl get pods -n ethos-go -w
```

### 3. Database Migration

Initialize the database schema.

```bash
kubectl apply -f k8s/migration-job.yaml
```

_Verify success:_ `kubectl logs job/ethos-go-migration -n ethos-go` (Should say `1/u init_schema`).

### 4. Deploy Application

Launch the API and Worker.

```bash
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/worker-deployment.yaml
```

### 5. Access

Open **[http://localhost:8080](http://localhost:8080)**.

If using Minikube or cloud, get the external IP:

```bash
kubectl get svc ethos-go-app -n ethos-go
```

---

## 🔍 Advanced Debugging & Troubleshooting

### 1. Pod Health & States

| Status             | Likely Cause                                                   | Debug Action                                                |
| ------------------ | -------------------------------------------------------------- | ----------------------------------------------------------- |
| `CrashLoopBackOff` | Application panic, DB connection failure, or missing env vars. | `kubectl logs <pod> -n ethos-go` (Check for panic/error)    |
| `ImagePullBackOff` | Image name typo, tag missing, or private registry auth failed. | `kubectl describe pod <pod> -n ethos-go` (Look at 'Events') |
| `Pending`          | No nodes available, PVC not bound, or insufficient CPU/RAM.    | `kubectl describe pod <pod> -n ethos-go`                    |
| `Evicted`          | Node out of disk or memory.                                    | Check node resources: `kubectl top nodes`                   |

### 2. Log Analysis strategies

**View logs of a crashing pod (previous instance):**
If a pod crashes immediately, current logs might be empty. Check the _previous_ run:

```bash
kubectl logs <pod-name> -n ethos-go --previous
```

**Stream logs from all application pods:**

```bash
kubectl logs -f -l app=ethos-go -n ethos-go
```

**Stream worker logs:**

```bash
kubectl logs -f -l app=ethos-go-worker -n ethos-go
```

### 3. Networking & Connectivity

Since production images (Distroless/Alpine) are minimal, they might lack `curl` or `ping`.

**A. Validate Service Discovery**
Use a debug container to test DNS and connections inside the cluster.

```bash
# Launch a temporary debug shell in the namespace
kubectl run -i --tty --rm debug-shell --image=curlimages/curl --restart=Never -n ethos-go -- sh
```

**B. Inside the debug shell:**

```bash
# 1. Test DNS resolution
nslookup ethos-go-postgres
nslookup ethos-go-redis
nslookup ethos-go-app

# 2. Test Connection to App
curl -v http://ethos-go-app:8080/health

# 3. Test Connection to Postgres Port
curl -v telnet://ethos-go-postgres:5432
```

### 4. Storage (StatefulSet) Issues

If `ethos-go-postgres-0` is stuck in `Pending`:

1. **Check PVC Status:**

   ```bash
   kubectl get pvc -n ethos-go
   ```

   Status must be `Bound`.

2. **Check StorageClass:**
   If using Kind/Minikube, ensure a default storage class exists:

   ```bash
   kubectl get sc
   ```

3. **Reset Database (Data Wipe):**
   _Warning: This deletes all data._
   ```bash
   kubectl delete statefulset ethos-go-postgres -n ethos-go
   kubectl delete pvc postgres-data-ethos-go-postgres-0 -n ethos-go
   kubectl apply -f k8s/postgres.yaml
   ```

### 5. Configuration Verification

Ensure pods are receiving the correct environment variables.

**Inspect a running pod's environment:**

```bash
kubectl exec <pod-name> -n ethos-go -- env | grep DB_
```

**Decode Secrets:**
To see what password Kubernetes is actually injecting:

```bash
kubectl get secret ethos-go-secret -n ethos-go -o jsonpath="{.data.DB_PASSWORD}" | base64 --decode
```

---

## 🔄 Maintenance

### Updating the Application

When you change code:

1. Rebuild & Push:
   ```bash
   docker build -t sammidev/ethos-go:latest .
   docker push sammidev/ethos-go:latest
   ```
2. Rolling Restart:
   ```bash
   kubectl rollout restart deployment -n ethos-go
   ```

### Scaling

Scale workers to handle more background tasks:

```bash
kubectl scale deployment ethos-go-worker --replicas=3 -n ethos-go
```

---

## 🧹 Clean Up

To remove all resources created by this project:

### Option 1: Delete Namespace (Recommended)

This removes everything including data volumes.

```bash
kubectl delete namespace ethos-go
```

### Option 2: Delete Resources Only

Keeps the namespace but removes deployments and services.

```bash
kubectl delete -f k8s/
```
