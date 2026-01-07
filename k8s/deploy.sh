#!/bin/bash

# Ethos-Go Kubernetes Deployment Script
# This script deploys the complete Ethos-Go application to Kubernetes

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="ethos-go"
KUBECTL="kubectl"
K8S_DIR="$(dirname "$0")"

# Functions
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    print_info "Checking prerequisites..."

    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl not found. Please install kubectl."
        exit 1
    fi

    # Check cluster connection
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi

    print_info "Prerequisites check passed ✓"
}

create_namespace() {
    print_info "Creating namespace..."
    $KUBECTL apply -f "$K8S_DIR/01-namespace.yaml"
}

create_config() {
    print_info "Creating ConfigMap..."
    $KUBECTL apply -f "$K8S_DIR/02-configmap.yaml"
}

create_secrets() {
    print_info "Creating secrets..."

    # Check if secrets already exist
    if $KUBECTL get secret ethos-go-secret -n $NAMESPACE &> /dev/null; then
        print_warn "Secrets already exist. Skipping..."
        return
    fi

    print_warn "Please create secrets. You can either:"
    echo ""
    echo "Option 1: Apply the secret template (edit it first):"
    echo "  kubectl apply -f $K8S_DIR/03-secret.yaml"
    echo ""
    echo "Option 2: Create secret from command line:"
    echo "kubectl create secret generic ethos-go-secret \\"
    echo "  --from-literal=DB_PASSWORD=YOUR_DB_PASSWORD \\"
    echo "  --from-literal=REDIS_PASSWORD=YOUR_REDIS_PASSWORD \\"
    echo "  --from-literal=AUTH_JWT_SECRET=\$(openssl rand -base64 32) \\"
    echo "  --from-literal=SMTP_PASSWORD=YOUR_SMTP_PASSWORD \\"
    echo "  --from-literal=VAPID_PRIVATE_KEY=YOUR_VAPID_KEY \\"
    echo "  --namespace=$NAMESPACE"
    echo ""
    read -p "Press enter when secrets are created..."
}

deploy_postgresql() {
    print_info "Deploying PostgreSQL..."
    $KUBECTL apply -f "$K8S_DIR/04-postgres.yaml"

    print_info "Waiting for PostgreSQL to be ready..."
    $KUBECTL wait --for=condition=ready pod -l app=ethos-go-postgres --timeout=300s -n $NAMESPACE || {
        print_error "PostgreSQL failed to start"
        exit 1
    }
    print_info "PostgreSQL is ready ✓"
}

deploy_redis() {
    print_info "Deploying Redis..."
    $KUBECTL apply -f "$K8S_DIR/05-redis.yaml"

    print_info "Waiting for Redis to be ready..."
    $KUBECTL wait --for=condition=ready pod -l app=ethos-go-redis --timeout=300s -n $NAMESPACE || {
        print_error "Redis failed to start"
        exit 1
    }
    print_info "Redis is ready ✓"
}

deploy_app() {
    print_info "Deploying application (migrations will run automatically on startup)..."
    $KUBECTL apply -f "$K8S_DIR/06-app-deployment.yaml"

    print_info "Waiting for application to be ready..."
    $KUBECTL wait --for=condition=available deployment/ethos-go-app --timeout=300s -n $NAMESPACE || {
        print_error "Application failed to start"
        exit 1
    }
    print_info "Application is ready ✓"
}

deploy_worker() {
    print_info "Deploying worker..."
    $KUBECTL apply -f "$K8S_DIR/07-worker-deployment.yaml"

    print_info "Waiting for worker to be ready..."
    $KUBECTL wait --for=condition=available deployment/ethos-go-worker --timeout=300s -n $NAMESPACE || {
        print_error "Worker failed to start"
        exit 1
    }
    print_info "Worker is ready ✓"
}

deploy_services() {
    print_info "Deploying services..."
    $KUBECTL apply -f "$K8S_DIR/08-service.yaml"
}

deploy_gateway() {
    print_info "Deploying Gateway API resources..."
    print_warn "Make sure to update the hostname in 09-gateway.yaml before deploying!"
    print_warn "Prerequisites: Gateway API CRDs and a Gateway Controller (Envoy, Cilium, Istio, etc.)"
    read -p "Continue with Gateway deployment? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        $KUBECTL apply -f "$K8S_DIR/09-gateway.yaml"
    else
        print_warn "Skipping Gateway deployment"
    fi
}

deploy_autoscaling() {
    print_info "Deploying autoscaling (HPA & PDB)..."
    $KUBECTL apply -f "$K8S_DIR/10-hpa.yaml"
    $KUBECTL apply -f "$K8S_DIR/11-pdb.yaml"
}

deploy_network_policies() {
    print_info "Deploying network policies..."
    print_warn "Network policies require a CNI that supports them (Calico, Cilium, etc.)"
    read -p "Deploy network policies? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        $KUBECTL apply -f "$K8S_DIR/12-networkpolicy.yaml"
    else
        print_warn "Skipping network policies"
    fi
}

verify_deployment() {
    print_info "Verifying deployment..."

    echo ""
    echo "=== Pods ==="
    $KUBECTL get pods -n $NAMESPACE

    echo ""
    echo "=== Services ==="
    $KUBECTL get svc -n $NAMESPACE

    echo ""
    echo "=== Gateway ==="
    $KUBECTL get gateway,httproute -n $NAMESPACE 2>/dev/null || echo "No Gateway configured"

    echo ""
    echo "=== HPA ==="
    $KUBECTL get hpa -n $NAMESPACE 2>/dev/null || echo "No HPA configured"
}

show_access_info() {
    print_info "Deployment completed successfully! 🎉"
    echo ""
    echo "=== Access Information ==="

    # Get ingress URL
    GATEWAY_HOST=$($KUBECTL get gateway ethos-go-gateway -n $NAMESPACE -o jsonpath='{.spec.listeners[1].hostname}' 2>/dev/null || echo "")
    if [ -n "$GATEWAY_HOST" ]; then
        echo "Gateway URL: https://$GATEWAY_HOST"
    fi

    # Port forward command
    echo ""
    echo "For local development, use port forwarding:"
    echo "kubectl port-forward svc/ethos-go-app 8080:80 -n $NAMESPACE"
    echo "Then access: http://localhost:8080"

    echo ""
    echo "=== Useful Commands ==="
    echo "View logs:    kubectl logs -f deployment/ethos-go-app -n $NAMESPACE"
    echo "Get pods:     kubectl get pods -n $NAMESPACE"
    echo "Describe pod: kubectl describe pod <pod-name> -n $NAMESPACE"
    echo "Shell access: kubectl exec -it deployment/ethos-go-app -n $NAMESPACE -- sh"
}

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --dev       Deploy for development (skip ingress, HPA, network policies)"
    echo "  --prod      Deploy for production (full deployment)"
    echo "  --help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --dev    # Deploy for local development"
    echo "  $0 --prod   # Deploy for production"
}

# Main deployment flow
main() {
    local MODE="${1:---dev}"

    echo "======================================"
    echo "  Ethos-Go Kubernetes Deployment"
    echo "======================================"
    echo ""

    check_prerequisites
    create_namespace
    create_config
    create_secrets
    deploy_postgresql
    deploy_redis
    deploy_app
    deploy_worker
    deploy_services

    if [ "$MODE" == "--prod" ]; then
        deploy_gateway
        deploy_autoscaling
        deploy_network_policies
    else
        print_info "Skipping production components (gateway, HPA, network policies)"
        print_info "Use --prod flag for full production deployment"
    fi

    verify_deployment
    show_access_info
}

# Parse arguments
case "${1:-}" in
    --help|-h)
        usage
        exit 0
        ;;
    --prod|--dev|"")
        main "${1:---dev}"
        ;;
    *)
        print_error "Unknown option: $1"
        usage
        exit 1
        ;;
esac
