# Login Screen Application

A comprehensive Go web application demonstrating authentication best practices, following enterprise-grade architecture patterns inspired by the Perkeep project.

## Features

- **Secure Authentication**: JWT-based authentication with password hashing
- **Session Management**: Secure session handling with proper token validation
- **User Management**: Registration, login, and profile management
- **Enterprise Architecture**: Clean separation of concerns following Go best practices
- **Web Interface**: Modern HTML/CSS/JavaScript frontend
- **API Endpoints**: RESTful API for all authentication operations
- **Middleware**: Authentication and logging middleware
- **Configuration**: Environment-based configuration management

## Project Structure

```
login-app/
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
├── main.go                 # Application entry point
├── cmd/                    # Application commands
│   └── server/
│       └── main.go
├── internal/               # Private application code
│   ├── auth/              # Authentication logic
│   │   ├── handler.go     # HTTP handlers
│   │   ├── middleware.go  # Auth middleware
│   │   ├── service.go     # Business logic
│   │   └── types.go       # Auth-related types
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── storage/           # Data storage layer
│   │   ├── memory.go      # In-memory storage
│   │   └── user.go        # User storage interface
│   └── server/            # HTTP server setup
│       ├── handler.go     # Main server handler
│       ├── middleware.go  # Server middleware
│       └── routes.go      # Route definitions
├── web/                   # Web assets
│   ├── static/           # Static files (CSS, JS, images)
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/        # HTML templates
│       ├── login.html
│       ├── register.html
│       ├── dashboard.html
│       └── base.html
├── api/                   # API documentation
│   └── openapi.yaml
├── configs/               # Configuration files
│   ├── development.yaml
│   └── production.yaml
├── scripts/              # Build and deployment scripts
│   ├── build.sh
│   └── run.sh
└── docs/                 # Documentation
    └── architecture.md
```

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd login-app
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go

LOG_LEVEL=debug JWT_SECRET=dev-secret ./login-app -env=development
```

4. Open your browser and navigate to `http://localhost:8080`

### Configuration

The application uses environment variables for configuration:

- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT signing (required in production)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `ENVIRONMENT`: Application environment (development, production)

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `GET /api/auth/profile` - Get user profile (requires auth)

### Web Pages

- `GET /` - Landing page
- `GET /login` - Login page
- `GET /register` - Registration page
- `GET /dashboard` - User dashboard (requires auth)

## Architecture

This application follows enterprise Go architecture patterns:

### Clean Architecture Layers

1. **Handler Layer** (`internal/auth/handler.go`): HTTP request/response handling
2. **Service Layer** (`internal/auth/service.go`): Business logic and validation
3. **Storage Layer** (`internal/storage/`): Data persistence abstraction

### Key Design Patterns

- **Dependency Injection**: Services are injected into handlers
- **Interface Segregation**: Small, focused interfaces
- **Middleware Pattern**: Reusable cross-cutting concerns
- **Configuration Management**: Environment-based configuration
- **Error Handling**: Structured error responses

### Security Features

- **Password Hashing**: bcrypt for secure password storage
- **JWT Tokens**: Stateless authentication with configurable expiration
- **Input Validation**: Comprehensive request validation
- **CSRF Protection**: Cross-site request forgery protection
- **Secure Headers**: Security-focused HTTP headers

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
./scripts/build.sh
```

### Code Quality

This project follows Go best practices:

- Standard project layout
- Comprehensive error handling
- Proper logging
- Unit test coverage
- Code documentation
- Linting compliance (`golangci-lint`)

## Learning Objectives

This application demonstrates:

1. **Go Web Development**: Building robust web applications
2. **Authentication Patterns**: Industry-standard auth implementation
3. **Enterprise Architecture**: Scalable Go application structure
4. **Security Best Practices**: Secure web application development
5. **Testing Strategies**: Unit and integration testing
6. **Deployment Practices**: Production-ready Go applications

## Contributing

This is a reference implementation for learning purposes. Feel free to extend and modify for your specific needs.

## License

This project is part of the UdemyGolangApps learning repository.
