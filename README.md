# FitByte API

A RESTful API built with Go and Gin framework, featuring JWT authentication, PostgreSQL integration, and clean architecture principles.

## Features

- 🚀 **Gin Framework** - Fast HTTP web framework
- 🔐 **JWT Authentication** - Secure user authentication with JWT tokens
- 🗄️ **PostgreSQL Integration** - Database persistence with GORM
- 🔒 **Password Security** - bcrypt password hashing
- 📝 **Structured Logging** - Using zerolog for efficient logging
- 🛡️ **CORS Support** - Cross-origin resource sharing
- 🏗️ **Clean Architecture** - Organized project structure
- 📊 **Health Checks** - Built-in health and readiness endpoints
- 🔧 **Environment Configuration** - Easy configuration management
- 📋 **Standard REST API** - Following REST conventions
- 🔄 **Database Migrations** - Automatic schema management

## Project Structure

```
fitbyte/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── .env                    # Environment variables
├── README.md              # This file
└── internal/              # Private application code
    ├── config/            # Configuration management
    │   └── config.go
    ├── database/          # Database connection and migrations
    │   └── database.go
    ├── handlers/          # HTTP request handlers
    │   ├── auth.go        # Authentication handlers
    │   ├── health.go      # Health check handlers
    │   └── user.go        # User management handlers
    ├── middleware/        # HTTP middleware
    │   ├── auth.go        # JWT authentication middleware
    │   ├── cors.go        # CORS middleware
    │   ├── logger.go      # Logging middleware
    │   └── recovery.go    # Panic recovery middleware
    ├── models/            # Data models
    │   ├── response.go    # API response models
    │   └── user.go        # User models
    ├── routes/            # Route definitions
    │   └── routes.go
    └── services/          # Business logic services
        ├── jwt_service.go # JWT token management
        └── user_service.go # User business logic
```

## Getting Started

### Prerequisites

- Go 1.25.0 or higher
- PostgreSQL 14 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd fitbyte
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**
   ```sql
   CREATE DATABASE fitbyte_db;
   CREATE USER fitbyte_user WITH PASSWORD 'fitbyte_password';
   GRANT ALL PRIVILEGES ON DATABASE fitbyte_db TO fitbyte_user;
   ```

4. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Public Endpoints

#### Health Check
- `GET /api/v1/health/` - Health status
- `GET /api/v1/health/ready` - Readiness check

#### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Protected Endpoints (Require JWT Token)

#### User Profile
- `GET /api/v1/profile/` - Get current user profile

#### User Management
- `GET /api/v1/users/` - Get all users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users/` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Root
- `GET /` - API information

## Authentication

### Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Using JWT Token
```bash
curl -X GET http://localhost:8080/api/v1/profile/ \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Data Models

### User Model
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "preference": "metric",
  "weightUnit": "kg",
  "heightUnit": "cm", 
  "weight": 75.5,
  "height": 180.0,
  "imageUri": "https://example.com/image.jpg",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### Authentication Response
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `JWT_SECRET` | JWT signing secret | Required |

### Example .env file
```env
ENVIRONMENT=development
PORT=8080
DATABASE_URL=postgres://fitbyte_user:fitbyte_password@localhost:5432/fitbyte_db?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

## Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Tokens**: 24-hour expiration with refresh capability
- **Protected Routes**: Middleware-based authentication
- **Input Validation**: Request validation with Gin binding
- **CORS Support**: Configurable cross-origin resource sharing

## Development

### Adding New Endpoints

1. **Create a new handler** in `internal/handlers/`
2. **Define models** in `internal/models/` if needed
3. **Add business logic** in `internal/services/` if needed
4. **Add routes** in `internal/routes/routes.go`
5. **Update main.go** to initialize the new handler

### Example: Adding a Product Handler

```go
// internal/handlers/product.go
type ProductHandler struct {
    productService *services.ProductService
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
    // Implementation
}
```

```go
// internal/routes/routes.go
products := protected.Group("/products")
{
    products.GET("/", productHandler.GetProducts)
}
```

## Database

The application uses PostgreSQL with GORM for database operations:

- **Automatic Migrations**: Database schema is created automatically on startup
- **Connection Pooling**: Managed by GORM
- **Transaction Support**: Built-in transaction management
- **Query Logging**: SQL queries are logged in development mode

## Building for Production

```bash
# Build the application
go build -o fitbyte main.go

# Run the binary
./fitbyte
```

## Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o fitbyte main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/fitbyte .
CMD ["./fitbyte"]
```

## Testing

### Manual Testing

1. **Start the application**
   ```bash
   go run main.go
   ```

2. **Test registration**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
   ```

3. **Test login**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

4. **Test protected endpoint**
   ```bash
   curl -X GET http://localhost:8080/api/v1/profile/ \
     -H "Authorization: Bearer <token-from-login>"
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Next Steps

- [ ] Add input validation middleware
- [ ] Add rate limiting
- [ ] Add API documentation (Swagger)
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Add Docker Compose setup
- [ ] Add CI/CD pipeline
- [ ] Add email verification
- [ ] Add password reset functionality
- [ ] Add user roles and permissions