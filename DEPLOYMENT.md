# FitByte Deployment Guide

This guide covers different deployment options for the FitByte API application.

## 🐳 Docker Deployment

### Prerequisites
- Docker and Docker Compose installed
- At least 4GB RAM available
- Ports 8080, 5432, 9000, 9001, 6379 available

### Quick Start

1. **Clone and setup**
   ```bash
   git clone <your-repo-url>
   cd FitByte
   # The deployment script will create .env automatically
   ```

2. **Deploy with Docker Compose**
   ```bash
   # Development deployment (includes all services)
   ./scripts/docker-deploy.sh dev
   
   # Production deployment (app only)
   ./scripts/docker-deploy.sh prod
   ```

3. **Access the application**
   - API: http://localhost:8080
   - MinIO Console: http://localhost:9001
   - PostgreSQL: localhost:5432

### Manual Docker Commands

```bash
# Build the image
docker build -t fitbyte:latest .

# Run with external database
docker run -d \
  --name fitbyte-api \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e JWT_SECRET="your-secret" \
  fitbyte:latest

# Run with docker-compose
docker-compose up -d
```

## ☸️ Kubernetes Deployment

### Prerequisites
- Kubernetes cluster (1.20+)
- Helm 3.x installed
- kubectl configured

### Deploy with Helm

1. **Add required repositories**
   ```bash
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
   helm repo update
   ```

2. **Deploy to development**
   ```bash
   # Install dependencies
   helm install postgresql bitnami/postgresql \
     --set auth.postgresPassword=postgres \
     --set auth.database=fitbyte
   
   helm install minio bitnami/minio \
     --set auth.rootUser=minioadmin \
     --set auth.rootPassword=minioadmin123
   
   # Deploy FitByte
   helm install fitbyte ./helm/fitbyte
   ```

3. **Deploy to production**
   ```bash
   # Create secrets first
   kubectl create secret generic fitbyte-secrets \
     --from-literal=jwt-secret="your-production-secret" \
     --from-literal=database-url="postgres://user:pass@external-db:5432/fitbyte" \
     --from-literal=minio-access-key="your-access-key" \
     --from-literal=minio-secret-key="your-secret-key"
   
   # Deploy with production values
   helm install fitbyte ./helm/fitbyte -f ./helm/fitbyte/values-production.yaml
   ```

### Verify Deployment

```bash
# Check pods
kubectl get pods -l app=fitbyte

# Check services
kubectl get svc

# Check ingress
kubectl get ingress

# View logs
kubectl logs -l app=fitbyte -f
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENVIRONMENT` | Application environment | `development` | No |
| `PORT` | Server port | `8080` | No |
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9000` | Yes |
| `MINIO_ACCESS_KEY` | MinIO access key | - | Yes |
| `MINIO_SECRET_KEY` | MinIO secret key | - | Yes |
| `MINIO_USE_SSL` | Use SSL for MinIO | `false` | No |
| `MINIO_BUCKET_NAME` | MinIO bucket name | `fitbyte-files` | No |

### Database Setup

1. **PostgreSQL**
   ```sql
   CREATE DATABASE fitbyte_db;
   CREATE USER fitbyte_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE fitbyte_db TO fitbyte_user;
   ```

2. **MinIO Bucket**
   ```bash
   # Create bucket using MinIO client
   mc mb minio/fitbyte-files
   mc policy set public minio/fitbyte-files
   ```

## 📊 Monitoring

### Health Checks

- **Liveness**: `GET /v1/health`
- **Readiness**: `GET /v1/health/ready`
- **Metrics**: `GET /metrics` (Prometheus format)

### Prometheus Metrics

The application exposes metrics on port 9090:
- HTTP request metrics
- Database connection metrics
- Custom business metrics

### Grafana Dashboard

Access Grafana at `http://grafana.fitbyte.local` (default credentials: admin/admin)

## 🔒 Security

### Production Security Checklist

- [ ] Change all default passwords
- [ ] Use strong JWT secrets (32+ characters)
- [ ] Enable SSL/TLS for all services
- [ ] Configure proper network policies
- [ ] Use secrets management (Vault, AWS Secrets Manager)
- [ ] Enable audit logging
- [ ] Regular security updates

### Secrets Management

```bash
# Kubernetes secrets
kubectl create secret generic fitbyte-secrets \
  --from-literal=jwt-secret="$(openssl rand -base64 32)" \
  --from-literal=database-url="postgres://..." \
  --from-literal=minio-access-key="..." \
  --from-literal=minio-secret-key="..."

# Docker secrets
echo "your-secret" | docker secret create jwt-secret -
```

## 🚀 CI/CD Pipeline

### GitHub Actions Example

```yaml
name: Deploy FitByte

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: docker build -t fitbyte:${{ github.sha }} .
      
      - name: Deploy to Kubernetes
        run: |
          helm upgrade --install fitbyte ./helm/fitbyte \
            --set app.image.tag=${{ github.sha }} \
            --set app.image.repository=your-registry.com/fitbyte
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check database connectivity
   kubectl exec -it <pod-name> -- nc -zv postgres-service 5432
   
   # Check database logs
   kubectl logs <postgres-pod>
   ```

2. **MinIO Connection Failed**
   ```bash
   # Check MinIO service
   kubectl get svc minio
   
   # Test MinIO connectivity
   kubectl exec -it <pod-name> -- nc -zv minio-service 9000
   ```

3. **Application Won't Start**
   ```bash
   # Check application logs
   kubectl logs <fitbyte-pod> -f
   
   # Check pod status
   kubectl describe pod <fitbyte-pod>
   ```

### Useful Commands

```bash
# View all resources
kubectl get all -l app=fitbyte

# Port forward for local testing
kubectl port-forward svc/fitbyte-service 8080:8080

# Scale application
kubectl scale deployment fitbyte --replicas=5

# Rolling update
kubectl rollout restart deployment/fitbyte

# Check rollout status
kubectl rollout status deployment/fitbyte
```

## 📈 Scaling

### Horizontal Scaling

```bash
# Scale using kubectl
kubectl scale deployment fitbyte --replicas=10

# Scale using Helm
helm upgrade fitbyte ./helm/fitbyte --set app.replicaCount=10

# Enable HPA
helm upgrade fitbyte ./helm/fitbyte --set autoscaling.enabled=true
```

### Vertical Scaling

```yaml
# Update resources in values.yaml
resources:
  limits:
    cpu: 2000m
    memory: 2Gi
  requests:
    cpu: 1000m
    memory: 1Gi
```

## 🔄 Updates and Maintenance

### Rolling Updates

```bash
# Update image
helm upgrade fitbyte ./helm/fitbyte --set app.image.tag=v1.1.0

# Rollback if needed
helm rollback fitbyte 1
```

### Database Migrations

The application automatically runs migrations on startup. For production:

1. Backup database before updates
2. Test migrations in staging
3. Use blue-green deployment for zero-downtime updates

## 📞 Support

For deployment issues:
1. Check the logs: `kubectl logs -l app=fitbyte`
2. Verify configuration: `kubectl describe pod <pod-name>`
3. Check resource usage: `kubectl top pods`
4. Review this documentation
5. Create an issue in the repository