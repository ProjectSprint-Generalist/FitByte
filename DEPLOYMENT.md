# FitByte Kubernetes + Helm Deployment Guide

This guide shows you how to deploy your FitByte application using Kubernetes and Helm.

## 🚀 Quick Start

### 1. Prerequisites

- Kubernetes cluster (1.19+)
- Helm 3.0+
- kubectl configured
- Docker (for building images)

### 2. Build and Push Docker Image

```bash
# Build the image
docker build -t fitbyte:latest .

# Tag for your registry
docker tag fitbyte:latest your-registry.com/fitbyte:latest

# Push to registry
docker push your-registry.com/fitbyte:latest
```

### 3. Deploy with Helm

```bash
# Navigate to helm directory
cd helm

# Deploy with default values
./deploy.sh

# Or deploy with custom values
helm install fitbyte ./fitbyte -f custom-values.yaml
```

## 📋 What Gets Deployed

### Core Application
- **FitByte API** - Your Go application (3 replicas)
- **PostgreSQL** - Database with persistent storage
- **MinIO** - Object storage for file uploads

### Monitoring Stack
- **Prometheus** - Metrics collection and alerting
- **Grafana** - Dashboards and visualization
- **ServiceMonitor** - Automatic metrics scraping

### Networking
- **Ingress** - External access with SSL/TLS
- **Services** - Internal service discovery
- **ConfigMaps & Secrets** - Configuration management

## 🔧 Configuration

### Environment Variables

The application automatically configures:
- `DATABASE_URL` - PostgreSQL connection
- `JWT_SECRET` - JWT signing key
- `MINIO_ENDPOINT` - Object storage endpoint
- `PROMETHEUS_ENABLED` - Metrics collection

### Customization

Edit `helm/fitbyte/values.yaml` to customize:
- Resource limits
- Replica count
- Storage sizes
- Ingress configuration
- Monitoring settings

## 📊 Monitoring

### Access Dashboards

```bash
# Grafana
kubectl port-forward svc/fitbyte-grafana 3000:3000 -n fitbyte
# Open http://localhost:3000 (admin/admin)

# Prometheus
kubectl port-forward svc/fitbyte-prometheus-server 9090:9090 -n fitbyte
# Open http://localhost:9090
```

### Metrics Available

- HTTP request rates and latencies
- Database operation metrics
- File upload statistics
- Memory and CPU usage
- Custom application metrics

## 🔒 Security

### Secrets Management

Sensitive data is stored in Kubernetes secrets:
- JWT signing keys
- Database passwords
- MinIO credentials

### Network Policies

Optional network policies restrict traffic between components.

## 📈 Scaling

### Horizontal Scaling

```yaml
# In values.yaml
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70
```

### Vertical Scaling

```yaml
# In values.yaml
app:
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 200m
      memory: 256Mi
```

## 🛠️ Troubleshooting

### Common Commands

```bash
# Check deployment status
kubectl get pods -n fitbyte

# View logs
kubectl logs -f deployment/fitbyte -n fitbyte

# Check services
kubectl get svc -n fitbyte

# Check ingress
kubectl get ingress -n fitbyte

# Port forward for testing
kubectl port-forward svc/fitbyte 8080:8080 -n fitbyte
```

### Health Checks

```bash
# API health
curl http://localhost:8080/api/v1/health

# Metrics endpoint
curl http://localhost:8080/metrics
```

## 🔄 Updates

### Upgrade Application

```bash
# Update image tag
helm upgrade fitbyte ./fitbyte --set app.image.tag=v1.1.0

# Or with values file
helm upgrade fitbyte ./fitbyte -f new-values.yaml
```

### Rollback

```bash
# List releases
helm list -n fitbyte

# Rollback to previous version
helm rollback fitbyte 1 -n fitbyte
```

## 🗑️ Cleanup

```bash
# Uninstall everything
helm uninstall fitbyte -n fitbyte

# Or using the script
./deploy.sh --uninstall
```

## 📚 Production Considerations

### Security
- Change all default passwords
- Use proper SSL certificates
- Enable network policies
- Regular security updates

### Performance
- Configure resource limits appropriately
- Use fast storage classes
- Enable horizontal pod autoscaling
- Monitor and tune database settings

### Backup
- Regular database backups
- MinIO data replication
- Configuration backup
- Disaster recovery plan

## 🆘 Support

For issues and questions:
1. Check the logs: `kubectl logs -f deployment/fitbyte -n fitbyte`
2. Verify configuration: `kubectl describe deployment fitbyte -n fitbyte`
3. Check resource usage: `kubectl top pods -n fitbyte`
4. Review the Helm chart documentation in `helm/README.md`

## 📖 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
