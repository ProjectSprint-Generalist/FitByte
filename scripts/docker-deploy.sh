#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Check if .env file exists
if [ ! -f .env ]; then
    print_warning ".env file not found. Creating from template..."
    cat > .env << EOF
# Application Configuration
ENVIRONMENT=development
PORT=8080

# Database Configuration
DATABASE_URL=postgres://fitbyte:postgres@postgres:5432/fitbyte_db?sslmode=disable
DB_HOST=postgres
DB_PORT=5432
DB_USER=fitbyte
DB_PASSWORD=postgres
DB_NAME=fitbyte_db

# JWT Configuration
JWT_SECRET=generalist-production

# MinIO Configuration
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin123
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=fitbyte-files

# CORS Configuration
CORS_ALLOWED_ORIGINS=*

# Redis Configuration (Optional)
REDIS_URL=redis://redis:6379

# Logging Configuration
LOG_LEVEL=info
EOF
    print_warning "Created .env file with default values. You can edit it if needed."
fi

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Function to deploy development environment
deploy_dev() {
    print_status "Deploying FitByte in development mode..."
    
    # Stop existing containers
    docker-compose down --remove-orphans
    
    # Build and start services
    docker-compose up --build -d
    
    # Wait for services to be ready
    print_status "Waiting for services to start..."
    sleep 10
    
    # Check if services are running
    if docker-compose ps | grep -q "Up"; then
        print_status "✅ FitByte deployed successfully!"
        print_status "🌐 API: http://localhost:8080"
        print_status "🗄️  PostgreSQL: localhost:5432"
        print_status "📦 MinIO Console: http://localhost:9001"
        print_status "🔧 MinIO API: http://localhost:9000"
        print_status "📊 Redis: localhost:6379"
    else
        print_error "❌ Deployment failed. Check logs with: docker-compose logs"
        exit 1
    fi
}

# Function to deploy production environment
deploy_prod() {
    print_status "Deploying FitByte in production mode..."
    
    # Stop existing containers
    docker-compose -f docker-compose.prod.yml down --remove-orphans
    
    # Build and start services
    docker-compose -f docker-compose.prod.yml up --build -d
    
    # Wait for services to be ready
    print_status "Waiting for services to start..."
    sleep 15
    
    # Check if services are running
    if docker-compose -f docker-compose.prod.yml ps | grep -q "Up"; then
        print_status "✅ FitByte deployed successfully in production!"
        print_status "🌐 API: http://localhost:8080"
    else
        print_error "❌ Production deployment failed. Check logs with: docker-compose -f docker-compose.prod.yml logs"
        exit 1
    fi
}

# Function to show logs
show_logs() {
    if [ "$1" = "prod" ]; then
        docker-compose -f docker-compose.prod.yml logs -f
    else
        docker-compose logs -f
    fi
}

# Function to stop services
stop_services() {
    if [ "$1" = "prod" ]; then
        docker-compose -f docker-compose.prod.yml down
    else
        docker-compose down
    fi
    print_status "Services stopped."
}

# Function to clean up
cleanup() {
    print_status "Cleaning up Docker resources..."
    docker-compose down --volumes --remove-orphans
    docker system prune -f
    print_status "Cleanup completed."
}

# Main script logic
case "$1" in
    "dev")
        deploy_dev
        ;;
    "prod")
        deploy_prod
        ;;
    "logs")
        show_logs "$2"
        ;;
    "stop")
        stop_services "$2"
        ;;
    "cleanup")
        cleanup
        ;;
    *)
        echo "Usage: $0 {dev|prod|logs|stop|cleanup}"
        echo ""
        echo "Commands:"
        echo "  dev     - Deploy in development mode with all services"
        echo "  prod    - Deploy in production mode"
        echo "  logs    - Show logs (add 'prod' for production logs)"
        echo "  stop    - Stop services (add 'prod' for production)"
        echo "  cleanup - Clean up Docker resources"
        echo ""
        echo "Examples:"
        echo "  $0 dev              # Deploy development environment"
        echo "  $0 prod             # Deploy production environment"
        echo "  $0 logs             # Show development logs"
        echo "  $0 logs prod        # Show production logs"
        echo "  $0 stop             # Stop development services"
        echo "  $0 cleanup          # Clean up everything"
        exit 1
        ;;
esac
