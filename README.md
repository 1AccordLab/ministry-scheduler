# Ministry Scheduler - Clean Architecture Go Example

This project demonstrates a robust Clean Architecture implementation in Go using only the standard library. It's a user management system that follows Clean Architecture principles with proper separation of concerns, dependency injection, and comprehensive testing.

## 🏗️ Architecture

```
ministry-scheduler/
├── cmd/app/                # Application entry point
├── internal/
│   ├── domain/             # Business entities and rules
│   ├── usecase/            # Business logic
│   ├── infra/              # Infrastructure (database, external services)
│   └── handler/            # HTTP handlers (presentation layer)
└── test/                   # Unit tests
    ├── domain/
    └── usecase/
```

## ✨ Features

- **Clean Architecture**: Proper separation of concerns with clear dependency rules
- **Standard Library Only**: No external frameworks, just pure Go
- **Comprehensive Testing**: Unit tests with mocks for all layers
- **REST API**: Full CRUD operations for users
- **SQLite Database**: Embedded database for simplicity
- **Graceful Shutdown**: Proper HTTP server lifecycle management
- **Context Support**: Request-scoped context throughout the application
- **Error Handling**: Proper error types and HTTP status codes
- **Input Validation**: Email format validation and length constraints
- **Request Logging**: HTTP request logging middleware
- **Rate Limiting**: Built-in limits for list operations
- **Duplicate Prevention**: Email uniqueness validation

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

## 🏛️ Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)

- Contains business entities (`User`)
- Defines business rules and validation
- Contains domain errors
- No dependencies on other layers

### 2. Use Case Layer (`internal/usecase/`)

- Contains application business logic
- Uses interfaces defined in the domain layer
- Orchestrates data flow between entities
- Depends only on the domain layer

### 3. Infrastructure Layer (`internal/infra/`)

- Implements external concerns (database, APIs)
- Implements interfaces defined in the domain layer
- Contains SQLite repository implementation
- Depends only on the domain layer

### 4. Handler Layer (`internal/handler/`)

- HTTP request/response handling
- Input validation and serialization
- Maps HTTP concerns to use case calls
- Depends on use case and domain layers

### 5. Main Layer (`cmd/app/`)

- Application assembly and dependency injection
- Configuration and environment setup
- HTTP server setup and lifecycle
- Dependencies wiring

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

## 🤔 Why This Architecture?

For beginners entering the Go ecosystem, this example provides:

1. **Pure Go Experience**: Learn Go without framework magic
2. **Industry Standards**: Clean Architecture is widely adopted
3. **Practical Patterns**: Real-world dependency injection and error handling
4. **Testing Culture**: Comprehensive test coverage with mocks
5. **Scalability**: Structure that grows with your application

## 🎓 Learning Path

1. Study the `domain` layer to understand business entities
2. Examine the `usecase` layer for business logic patterns
3. Look at the `infra` layer for data persistence patterns
4. Review the `handler` layer for HTTP API patterns
5. Understand dependency injection in `main.go`
6. Study the test files to learn Go testing patterns

This project serves as a foundation for building robust, maintainable Go applications using Clean Architecture principles with the standard library.
