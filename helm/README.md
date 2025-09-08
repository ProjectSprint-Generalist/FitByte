# FitByte Helm Chart

This Helm chart deploys the FitByte fitness tracking application with all its dependencies including PostgreSQL, MinIO, Prometheus, and Grafana.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- kubectl configured to connect to your Kubernetes cluster

## Quick Start

### 1. Build and Push Docker Image

First, build and push your FitByte application Docker image:

```bash
# Build the Docker image
docker build -t fitbyte:latest .

# Tag for your registry (replace with your registry)
docker tag fitbyte:latest your-registry.com/fitbyte:latest

# Push to registry
docker push your-registry.com/fitbyte:latest
```

### 2. Deploy with Default Values

```bash
# Deploy to default namespace
helm install fitbyte ./fitbyte

# Deploy to specific namespace
helm install fitbyte ./fitbyte -n fitbyte --create-namespace
```

### 3. Deploy with Custom Values

```bash
# Deploy with custom values file
helm install fitbyte ./fitbyte -f custom-values.yaml

# Deploy with inline values
helm install fitbyte ./fitbyte --set app.replicaCount=3 --set app.image.tag=v1.0.0
```

### 4. Using the Deployment Script

```bash
# Deploy with default settings
./deploy.sh

# Deploy with custom release name and namespace
./deploy.sh --release my-fitbyte --namespace my-namespace

# Deploy with custom values file
./deploy.sh --values custom-values.yaml

# Dry run to see what would be deployed
./deploy.sh --dry-run

# Upgrade existing deployment
./deploy.sh --upgrade
```

## Configuration

### Values File Structure

The main configuration is in `values.yaml`. Key sections:

#### Application Configuration
```yaml
app:
  replicaCount: 3
  image:
    repository: fitbyte
    tag: "latest"
  service:
    type: ClusterIP
    port: 8080
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi
```

#### Database Configuration
```yaml
database:
  postgresql:
    enabled: true
    auth:
      postgresPassword: "postgres_password"
      database: "fitbyte"
      username: "fitbyte"
    primary:
      persistence:
        enabled: true
        size: 10Gi
```

#### MinIO Configuration
```yaml
minio:
  enabled: true
  auth:
    rootUser: "minioadmin"
    rootPassword: "minioadmin"
  defaultBuckets: "fitbyte-uploads"
  persistence:
    enabled: true
    size: 50Gi
```

#### Monitoring Configuration
```yaml
prometheus:
  enabled: true
  server:
    persistentVolume:
      enabled: true
      size: 10Gi

grafana:
  enabled: true
  adminPassword: "admin"
  persistence:
    enabled: true
    size: 5Gi
```

## Environment Variables

The application uses the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Auto-generated |
| `JWT_SECRET` | JWT signing secret | From secret |
| `JWT_EXPIRY` | JWT token expiry | 24h |
| `MINIO_ENDPOINT` | MinIO server endpoint | Auto-generated |
| `MINIO_ACCESS_KEY` | MinIO access key | minioadmin |
| `MINIO_SECRET_KEY` | MinIO secret key | From secret |
| `MINIO_BUCKET` | MinIO bucket name | fitbyte-uploads |
| `PROMETHEUS_ENABLED` | Enable Prometheus metrics | true |
| `METRICS_PORT` | Metrics port | 9090 |

## Services and Ports

| Service | Port | Description |
|---------|------|-------------|
| FitByte API | 8080 | Main application API |
| FitByte Metrics | 9090 | Prometheus metrics endpoint |
| PostgreSQL | 5432 | Database |
| MinIO API | 9000 | Object storage API |
| MinIO Console | 9001 | Object storage web UI |
| Prometheus | 9090 | Monitoring |
| Grafana | 3000 | Dashboards |

## Ingress

The chart includes optional ingress configuration:

```yaml
app:
  ingress:
    enabled: true
    className: "nginx"
    hosts:
      - host: fitbyte.local
        paths:
          - path: /api
            pathType: Prefix
          - path: /metrics
            pathType: Prefix
```

## Monitoring

### Prometheus Metrics

The application exposes Prometheus metrics on port 9090:

- `fitbyte_http_requests_total` - HTTP request counter
- `fitbyte_http_request_duration_seconds` - HTTP request duration histogram
- `fitbyte_database_operations_total` - Database operation counter
- `fitbyte_upload_requests_total` - File upload counter
- `fitbyte_memory_usage_bytes` - Memory usage gauge

### Grafana Dashboards

Grafana is automatically configured with:
- Prometheus as data source
- Pre-built dashboards for FitByte monitoring
- Alerting rules for common issues

Access Grafana:
```bash
kubectl port-forward svc/fitbyte-grafana 3000:3000 -n fitbyte
# Open http://localhost:3000
# Username: admin, Password: admin
```

## Scaling

### Horizontal Pod Autoscaling

To enable HPA, add to your values:

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80
```

### Vertical Scaling

Adjust resource limits in values.yaml:

```yaml
app:
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 200m
      memory: 256Mi
```

## Security

### Secrets Management

Sensitive data is stored in Kubernetes secrets:

```yaml
secrets:
  create: true
  data:
    jwt-secret: "your-jwt-secret-key-change-in-production"
    db-password: "postgres_password"
    minio-secret: "minioadmin"
```

### Network Policies

Enable network policies for additional security:

```yaml
networkPolicy:
  enabled: true
  ingress: []
  egress: []
```

## Troubleshooting

### Common Issues

1. **Pod not starting**
   ```bash
   kubectl describe pod <pod-name> -n fitbyte
   kubectl logs <pod-name> -n fitbyte
   ```

2. **Database connection issues**
   ```bash
   kubectl logs deployment/fitbyte -n fitbyte | grep -i database
   ```

3. **MinIO connection issues**
   ```bash
   kubectl logs deployment/fitbyte -n fitbyte | grep -i minio
   ```

### Useful Commands

```bash
# Check all resources
kubectl get all -n fitbyte

# View logs
kubectl logs -f deployment/fitbyte -n fitbyte

# Port forward for testing
kubectl port-forward svc/fitbyte 8080:8080 -n fitbyte

# Check ingress
kubectl get ingress -n fitbyte

# Check persistent volumes
kubectl get pv,pvc -n fitbyte

# Check secrets
kubectl get secrets -n fitbyte

# Check configmaps
kubectl get configmaps -n fitbyte
```

## Upgrading

### Upgrade Application

```bash
# Update image tag
helm upgrade fitbyte ./fitbyte --set app.image.tag=v1.1.0

# Or with values file
helm upgrade fitbyte ./fitbyte -f new-values.yaml
```

### Upgrade Dependencies

```bash
# Update Helm dependencies
helm dependency update ./fitbyte

# Upgrade with updated dependencies
helm upgrade fitbyte ./fitbyte
```

## Uninstalling

```bash
# Uninstall release
helm uninstall fitbyte -n fitbyte

# Or using the script
./deploy.sh --uninstall
```

## Development

### Local Development

For local development, you can use the chart with local values:

```yaml
# local-values.yaml
app:
  image:
    repository: fitbyte
    tag: "dev"
  replicaCount: 1

database:
  postgresql:
    enabled: true
    auth:
      postgresPassword: "dev_password"

minio:
  enabled: true
  auth:
    rootPassword: "dev_password"

prometheus:
  enabled: false

grafana:
  enabled: false
```

Deploy with:
```bash
helm install fitbyte-dev ./fitbyte -f local-values.yaml
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the chart
5. Submit a pull request

## License

This project is licensed under the MIT License.
