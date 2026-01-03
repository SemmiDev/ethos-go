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
    $KUBECTL apply -f base/namespace.yaml
}

create_secrets() {
    print_info "Creating secrets..."

    # Check if secrets already exist
    if $KUBECTL get secret app-secrets -n $NAMESPACE &> /dev/null; then
        print_warn "Secrets already exist. Skipping..."
        return
    fi

    print_warn "Please create secrets manually:"
    echo ""
    echo "kubectl create secret generic app-secrets \\"
    echo "  --from-literal=DB_PASSWORD=YOUR_DB_PASSWORD \\"
    echo "  --from-literal=REDIS_PASSWORD=YOUR_REDIS_PASSWORD \\"
    echo "  --from-literal=AUTH_JWT_SECRET=\$(openssl rand -base64 32) \\"
    echo "  --namespace=$NAMESPACE"
    echo ""
    read -p "Press enter when secrets are created..."
}

deploy_configmap() {
    print_info "Deploying ConfigMap..."
    $KUBECTL apply -f base/configmap.yaml
}

deploy_postgresql() {
    print_info "Deploying PostgreSQL..."
    $KUBECTL apply -f base/postgresql.yaml

    print_info "Waiting for PostgreSQL to be ready..."
    $KUBECTL wait --for=condition=ready pod -l app=postgresql --timeout=300s -n $NAMESPACE || {
        print_error "PostgreSQL failed to start"
        exit 1
    }
    print_info "PostgreSQL is ready ✓"
}

deploy_redis() {
    print_info "Deploying Redis..."
    $KUBECTL apply -f base/redis.yaml

    print_info "Waiting for Redis to be ready..."
    $KUBECTL wait --for=condition=ready pod -l app=redis --timeout=300s -n $NAMESPACE || {
        print_error "Redis failed to start"
        exit 1
    }
    print_info "Redis is ready ✓"
}

run_migrations() {
    print_info "Running database migrations..."
    $KUBECTL apply -f base/migration-job.yaml

    print_info "Waiting for migration to complete..."
    $KUBECTL wait --for=condition=complete job/migration --timeout=300s -n $NAMESPACE || {
        print_error "Migration failed"
        $KUBECTL logs job/migration -n $NAMESPACE
        exit 1
    }
    print_info "Migrations completed ✓"
}

deploy_app() {
    print_info "Deploying application..."
    $KUBECTL apply -f base/app.yaml

    print_info "Waiting for application to be ready..."
    $KUBECTL wait --for=condition=available deployment/ethos-go-app --timeout=300s -n $NAMESPACE || {
        print_error "Application failed to start"
        exit 1
    }
    print_info "Application is ready ✓"
}

deploy_services() {
    print_info "Deploying services..."
    $KUBECTL apply -f base/service.yaml
}

deploy_ingress() {
    print_info "Deploying ingress..."
    print_warn "Make sure to update the domain in ingress.yaml before deploying!"
    read -p "Continue? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        $KUBECTL apply -f base/ingress.yaml
    else
        print_warn "Skipping ingress deployment"
    fi
}

deploy_autoscaling() {
    print_info "Deploying autoscaling..."
    $KUBECTL apply -f base/hpa.yaml
    $KUBECTL apply -f base/pdb.yaml
}

deploy_network_policies() {
    print_info "Deploying network policies..."
    $KUBECTL apply -f base/networkpolicy.yaml
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
    echo "=== Ingress ==="
    $KUBECTL get ingress -n $NAMESPACE

    echo ""
    echo "=== HPA ==="
    $KUBECTL get hpa -n $NAMESPACE
}

show_access_info() {
    print_info "Deployment completed successfully! 🎉"
    echo ""
    echo "=== Access Information ==="

    # Get ingress URL
    INGRESS_HOST=$($KUBECTL get ingress ethos-go-ingress -n $NAMESPACE -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo "Not configured")
    echo "API URL: https://$INGRESS_HOST"

    # Port forward command
    echo ""
    echo "Or use port forwarding:"
    echo "kubectl port-forward svc/ethos-go-app 8080:80 -n $NAMESPACE"
    echo "Then access: http://localhost:8080"

    echo ""
    echo "=== Useful Commands ==="
    echo "View logs:    kubectl logs -f deployment/ethos-go-app -n $NAMESPACE"
    echo "Get pods:     kubectl get pods -n $NAMESPACE"
    echo "Describe pod: kubectl describe pod <pod-name> -n $NAMESPACE"
    echo "Shell access: kubectl exec -it deployment/ethos-go-app -n $NAMESPACE -- sh"
}

# Main deployment flow
main() {
    echo "======================================"
    echo "  Ethos-Go Kubernetes Deployment"
    echo "======================================"
    echo ""

    check_prerequisites
    create_namespace
    create_secrets
    deploy_configmap
    deploy_postgresql
    deploy_redis
    run_migrations
    deploy_app
    deploy_services
    deploy_ingress
    deploy_autoscaling
    deploy_network_policies
    verify_deployment
    show_access_info
}

# Run main function
main
