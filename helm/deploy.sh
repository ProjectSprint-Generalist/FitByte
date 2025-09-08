#!/bin/bash

# FitByte Helm Deployment Script
# This script deploys the FitByte application using Helm

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[HEADER]${NC} $1"
}

# Default values
RELEASE_NAME="fitbyte"
NAMESPACE="fitbyte"
CHART_PATH="./fitbyte"
VALUES_FILE=""
DRY_RUN=false
UPGRADE=false
UNINSTALL=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -r|--release)
            RELEASE_NAME="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -f|--values)
            VALUES_FILE="$2"
            shift 2
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -u|--upgrade)
            UPGRADE=true
            shift
            ;;
        --uninstall)
            UNINSTALL=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  -r, --release NAME     Release name (default: fitbyte)"
            echo "  -n, --namespace NAME   Namespace (default: fitbyte)"
            echo "  -f, --values FILE      Values file to use"
            echo "  -d, --dry-run          Dry run mode"
            echo "  -u, --upgrade          Upgrade existing release"
            echo "  --uninstall            Uninstall release"
            echo "  -h, --help             Show this help"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

print_header "FitByte Helm Deployment Script"
echo "======================================"

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    print_error "Helm is not installed. Please install Helm first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if cluster is accessible
if ! kubectl cluster-info &> /dev/null; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

print_status "Connected to Kubernetes cluster: $(kubectl config current-context)"

# Create namespace if it doesn't exist
print_status "Creating namespace '$NAMESPACE' if it doesn't exist..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Handle uninstall
if [ "$UNINSTALL" = true ]; then
    print_status "Uninstalling release '$RELEASE_NAME'..."
    helm uninstall $RELEASE_NAME -n $NAMESPACE
    print_status "Release uninstalled successfully!"
    exit 0
fi

# Build values file path
if [ -z "$VALUES_FILE" ]; then
    VALUES_FILE="$CHART_PATH/values.yaml"
fi

# Check if values file exists
if [ ! -f "$VALUES_FILE" ]; then
    print_error "Values file not found: $VALUES_FILE"
    exit 1
fi

# Prepare helm command
HELM_CMD="helm"
if [ "$UPGRADE" = true ]; then
    HELM_CMD="$HELM_CMD upgrade --install"
else
    HELM_CMD="$HELM_CMD install"
fi

HELM_CMD="$HELM_CMD $RELEASE_NAME $CHART_PATH"
HELM_CMD="$HELM_CMD -n $NAMESPACE"
HELM_CMD="$HELM_CMD -f $VALUES_FILE"

if [ "$DRY_RUN" = true ]; then
    HELM_CMD="$HELM_CMD --dry-run --debug"
fi

# Execute helm command
print_status "Executing Helm command..."
print_status "Command: $HELM_CMD"
echo ""

if [ "$DRY_RUN" = true ]; then
    print_warning "DRY RUN MODE - No changes will be made"
    eval $HELM_CMD
else
    eval $HELM_CMD
    
    if [ $? -eq 0 ]; then
        print_status "Deployment completed successfully!"
        
        # Wait for deployment to be ready
        print_status "Waiting for deployment to be ready..."
        kubectl wait --for=condition=available --timeout=300s deployment/$RELEASE_NAME -n $NAMESPACE
        
        # Show status
        print_status "Deployment status:"
        kubectl get pods -n $NAMESPACE -l "app.kubernetes.io/instance=$RELEASE_NAME"
        
        # Show services
        print_status "Services:"
        kubectl get services -n $NAMESPACE -l "app.kubernetes.io/instance=$RELEASE_NAME"
        
        # Show ingress if enabled
        if kubectl get ingress -n $NAMESPACE -l "app.kubernetes.io/instance=$RELEASE_NAME" &> /dev/null; then
            print_status "Ingress:"
            kubectl get ingress -n $NAMESPACE -l "app.kubernetes.io/instance=$RELEASE_NAME"
        fi
        
        echo ""
        print_status "To view logs: kubectl logs -f deployment/$RELEASE_NAME -n $NAMESPACE"
        print_status "To port forward: kubectl port-forward service/$RELEASE_NAME 8080:8080 -n $NAMESPACE"
        
    else
        print_error "Deployment failed!"
        exit 1
    fi
fi
