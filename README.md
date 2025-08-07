# Ministry Scheduler - Feature-First Go Architecture

This project demonstrates a pragmatic feature-first architecture in Go using only the standard library. It organizes code by features rather than technical layers, making it easier to find and modify related functionality while maintaining clean separation of concerns.

## 🏗️ Architecture

```
ministry-scheduler/
├── cmd/app/                 # Application entry point
├── internal/
│   ├── features/           # Feature-based organization
│   │   └── users/          # User management feature
│   │       ├── user.go     # Domain entities and validation
│   │       ├── service.go  # Business logic
│   │       ├── repository.go # Data access
│   │       └── api.go      # HTTP handlers
│   └── shared/             # Shared utilities
│       ├── database/       # Database connection and setup
│       ├── middleware/     # HTTP middleware
│       └── types/          # Common types and utilities
└── test/
    └── features/
        └── users/          # Feature-specific tests
```

## ✨ Features

- **Feature-First Organization**: Code organized by business features, not technical layers
- **Go Simplicity**: Embraces Go's philosophy of simplicity over enterprise complexity
- **Standard Library Only**: No external frameworks, just pure Go
- **Self-Contained Features**: Each feature contains all related code (domain, service, data, API)
- **Shared Utilities**: Common functionality organized in shared packages
- **Comprehensive Testing**: Feature-specific unit tests with mocks
- **REST API**: Full CRUD operations for users
- **SQLite Database**: Embedded database for simplicity
- **Graceful Shutdown**: Proper HTTP server lifecycle management
- **Context Support**: Request-scoped context throughout the application
- **Error Handling**: Proper error types and HTTP status codes
- **Input Validation**: Email format validation and length constraints
- **Request Logging**: HTTP request logging middleware
- **Easy to Navigate**: Find all user-related code in one place

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Installation

1. Clone the repository:

```bash
git clone <your-repo-url>
cd ministry-scheduler
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the application:

```bash
go run cmd/app/main.go
```

The server will start on port 8080 by default.

## 📋 API Endpoints

### Health Check

```bash
curl http://localhost:8080/health
```

### Users

**Create User**

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

**Get User**

```bash
curl http://localhost:8080/users/1
```

**Update User**

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com"}'
```

**List Users**

```bash
curl "http://localhost:8080/users?limit=10&offset=0"
```

**Delete User**

```bash
curl -X DELETE http://localhost:8080/users/1
```

## 🧪 Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## 🎯 Feature-First Architecture

### Features (`internal/features/`)

Each feature is self-contained with all related functionality:

#### User Feature (`internal/features/users/`)
- **`user.go`**: Domain entities, validation, and interfaces
- **`service.go`**: Business logic and use cases
- **`repository.go`**: Data access layer
- **`api.go`**: HTTP handlers and routing

### Shared Components (`internal/shared/`)

Common utilities shared across features:

#### Database (`internal/shared/database/`)
- Database connection and initialization
- Schema creation and migrations

#### Middleware (`internal/shared/middleware/`)
- HTTP middleware (logging, CORS, etc.)
- Cross-cutting concerns

#### Types (`internal/shared/types/`)
- Common types and utilities
- Environment variable helpers

### Main Application (`cmd/app/`)
- Application bootstrap and dependency injection
- Feature wiring and HTTP server setup

## 🎯 Key Benefits

1. **Testability**: Each layer can be tested in isolation using mocks
2. **Maintainability**: Clear separation makes code easy to understand and modify
3. **Independence**: Business logic is independent of frameworks and external tools
4. **Flexibility**: Easy to swap implementations (e.g., different databases)
5. **SOLID Principles**: Follows SOLID design principles throughout

## 🔧 Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database file path (default: users.db)

## 📊 Example Usage

```bash
# Start the server
go run cmd/app/main.go

# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Smith", "email": "alice@example.com"}'

# Response:
# {
#   "id": 1,
#   "name": "Alice Smith",
#   "email": "alice@example.com",
#   "created_at": "2024-01-15T10:30:00Z",
#   "updated_at": "2024-01-15T10:30:00Z"
# }

# Get the user
curl http://localhost:8080/users/1

# List users
curl http://localhost:8080/users
```

## 🤔 Why Feature-First Architecture?

This approach provides several advantages over traditional layered architecture:

1. **Easier Navigation**: All user-related code is in one place
2. **Faster Development**: No jumping between domain/usecase/handler folders
3. **Better Team Collaboration**: Different developers can work on different features
4. **Go Philosophy**: Embraces simplicity over enterprise complexity
5. **Natural Scaling**: Add new features without reorganizing existing code
6. **Reduced Cognitive Load**: Focus on business features, not technical layers

## 🎓 Learning Path

1. **Start with a Feature**: Explore `internal/features/users/` to see how everything is organized
2. **Understand the Domain**: Look at `user.go` for entities, validation, and business rules
3. **Follow the Flow**: See how `service.go` coordinates between repository and domain logic
4. **Study Data Access**: Review `repository.go` for database interaction patterns
5. **Examine the API**: Check `api.go` for HTTP handling and JSON responses
6. **Review Shared Code**: Look at `internal/shared/` for common utilities
7. **Understand Wiring**: See how `main.go` brings everything together
8. **Study Tests**: Check `test/features/users/` for testing patterns

This project serves as a foundation for building pragmatic, maintainable Go applications using a feature-first approach with the standard library.

## 🔄 Adding New Features

To add a new feature (e.g., `scheduling`):

1. Create `internal/features/scheduling/` directory
2. Add domain entities in `schedule.go`
3. Implement business logic in `service.go`  
4. Add data access in `repository.go`
5. Create HTTP handlers in `api.go`
6. Wire it up in `main.go`
7. Add tests in `test/features/scheduling/`

Each feature is independent and self-contained!
