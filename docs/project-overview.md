# Ministry Scheduler: Project Overview & Conventions

## Current Project Status

**Note**: The current codebase is a boilerplate Go application demonstrating feature-first architecture. For Phase 1 implementation, we will **completely rewrite** the application following the new distributed system design outlined in the Phase 1 plan.

## Current Codebase Analysis

### Existing Structure

The current project demonstrates a simple feature-first architecture:

```
ministry-scheduler/
├── cmd/app/main.go              # Simple HTTP server bootstrap
├── internal/
│   ├── features/users/          # User feature (example implementation)  
│   │   ├── user.go              # Domain entities & validation
│   │   ├── service.go           # Business logic  
│   │   ├── repository.go        # Data access (SQLite)
│   │   └── api.go               # HTTP handlers
│   └── shared/                  # Shared utilities
│       ├── database/            # SQLite setup
│       ├── middleware/          # Basic HTTP middleware  
│       └── env/                 # Environment helpers
└── test/features/users/         # Feature tests
```

### Current Technology Stack

- **Language**: Go 1.24.5
- **Database**: SQLite (single file database)
- **HTTP**: Go standard library
- **Architecture**: Feature-first monolith
- **Dependencies**: Minimal (only SQLite driver)

### What We'll Keep vs. Rewrite

#### Keep (Principles & Patterns)

- Feature-first organization philosophy
- Clean separation of concerns (domain, service, repository, API)
- Comprehensive testing approach
- Simple, pragmatic code style

#### Rewrite (Everything Else)

- **Architecture**: Monolith → Distributed app with background workers
- **Database**: SQLite → PostgreSQL cluster with replication + SQLc
- **Structure**: Single process → Web app + background workers
- **Message Queue**: None → Redis for job queues and caching
- **Authentication**: None → JWT + OAuth2 preparation
- **Communication**: Direct calls → HTTP + Redis job queues + PostgreSQL LISTEN/NOTIFY
- **Deployment**: Simple binary → Docker containers with load balancer

## New Project Structure & Conventions

### Directory Structure

```
ministry-scheduler/
├── cmd/                          # Application entrypoints
│   ├── web/main.go              # Web application server
│   └── worker/main.go           # Background worker process
├── internal/
│   ├── app/                     # Application core
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── auth/           # Authentication handlers
│   │   │   ├── users/          # User management handlers
│   │   │   ├── events/         # Event management handlers
│   │   │   └── schedules/      # Schedule management handlers
│   │   ├── services/           # Business logic services
│   │   │   ├── user/           # User service
│   │   │   ├── event/          # Event service
│   │   │   ├── schedule/       # Schedule service
│   │   │   └── notification/   # Notification service
│   │   └── workers/            # Background job workers
│   │       ├── email/          # Email sending worker
│   │       ├── notification/   # Push notification worker
│   │       └── scheduler/      # Recurring event generation worker
│   ├── domain/                 # Domain entities and business logic
│   │   ├── user/              # User domain
│   │   ├── event/             # Event domain
│   │   ├── schedule/          # Schedule domain
│   │   └── position/          # Position domain
│   ├── infrastructure/         # External system integrations
│   │   ├── database/          # Database connections, SQLc generated code, migrations
│   │   ├── redis/             # Redis client, job queue, caching
│   │   ├── auth/              # OAuth2 providers, JWT handling
│   │   ├── email/             # Email service integration
│   │   └── logging/           # Structured logging setup
│   └── shared/                # Cross-application utilities
│       ├── config/            # Configuration management
│       ├── errors/            # Common error types & handling
│       ├── middleware/        # HTTP middleware (auth, logging, etc.)
│       ├── jobs/              # Job queue interfaces and types
│       └── utils/             # Helper functions
├── sql/                       # Database schema & queries
│   ├── schema/               # PostgreSQL schemas (.sql files)
│   ├── queries/              # SQLc queries (.sql files)
│   └── migrations/           # Database migrations
├── docs/                      # Documentation
├── implemented/               # Implementation changelogs
├── deployments/              # Docker & deployment configurations
│   ├── docker-compose.yml    # Local development
│   ├── docker-compose.prod.yml # Production setup
│   └── nginx/               # Load balancer config
├── scripts/                  # Build & development scripts
└── tests/                   # Integration & end-to-end tests
```

## Code Organization Conventions

### 1. Application Structure

The application consists of two main processes that share the same codebase:

#### Web Application (cmd/web)

```
internal/app/handlers/users/
├── handler.go      # HTTP handlers & routing
├── middleware.go   # Handler-specific middleware
└── validation.go   # Request/response validation
```

#### Background Workers (cmd/worker)

```
internal/app/workers/email/
├── worker.go       # Job processing logic
├── templates.go    # Email templates
└── config.go       # Worker-specific configuration
```

### 2. Domain Layer

Domain packages contain pure business logic:

```
internal/domain/user/
├── user.go         # Entity definition & validation
├── repository.go   # Repository interface
├── service.go      # Domain service interface
└── errors.go       # Domain-specific errors
```

### 3. Database Layer

SQLc generates type-safe database code for PostgreSQL:

```
sql/
├── schema/001_users.sql      # PostgreSQL table definitions
├── queries/users.sql         # SQL queries with SQLc annotations
└── migrations/001_init.sql   # Database migrations (golang-migrate)
```

Database replication setup:

```
# Master database (writes)
postgresql://user:pass@db-master:5432/ministry_scheduler

# Read replicas (reads) 
postgresql://user:pass@db-replica-1:5432/ministry_scheduler
postgresql://user:pass@db-replica-2:5432/ministry_scheduler
```

### 4. Naming Conventions

#### File Naming

- **Go files**: `snake_case.go` (e.g., `user_service.go`)
- **SQL files**: `001_descriptive_name.sql` (numbered for ordering)
- **Config files**: `kebab-case.yaml` (e.g., `docker-compose.yaml`)

#### Go Code Conventions

- **Packages**: `lowercase` (e.g., `package user`)
- **Types**: `PascalCase` (e.g., `type UserService struct{}`)
- **Functions**: `PascalCase` for public, `camelCase` for private
- **Constants**: `PascalCase` for public, `camelCase` for private
- **Interfaces**: Often end with `-er` (e.g., `UserRepository`)

#### Database Conventions

- **Tables**: `snake_case` (e.g., `user_oauth_accounts`)
- **Columns**: `snake_case` (e.g., `created_at`, `user_id`)
- **Indexes**: `idx_table_column` (e.g., `idx_users_email`)
- **Foreign Keys**: `fk_table_reftable` (e.g., `fk_assignments_users`)

### 5. Error Handling Patterns

#### Domain Errors

```go
// Domain-specific errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email format")
)

// Error wrapping for context
func (s *userService) GetUser(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}
```

#### HTTP Error Responses

```go
// Standardized error response
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}

// Error handling middleware
func (h *Handler) handleError(w http.ResponseWriter, err error) {
    var statusCode int
    var errorCode string

    switch {
    case errors.Is(err, domain.ErrUserNotFound):
        statusCode = http.StatusNotFound
        errorCode = "USER_NOT_FOUND"
    default:
        statusCode = http.StatusInternalServerError
        errorCode = "INTERNAL_ERROR"
    }

    writeJSONError(w, statusCode, errorCode, err.Error())
}
```

### 6. Testing Conventions

#### Test File Naming

- Unit tests: `*_test.go` in same package
- Integration tests: `tests/integration/*_test.go`
- Mock files: `mocks/*_mock.go`

#### Test Structure

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        setup   func(*mocks.MockUserRepository)
        want    *User
        wantErr bool
    }{
        {
            name: "successful user creation",
            input: CreateUserRequest{Name: "John", Email: "john@test.com"},
            setup: func(repo *mocks.MockUserRepository) {
                repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&User{ID: 1}, nil)
            },
            want: &User{ID: 1},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Development Workflow

### 1. Database Changes

1. Create migration file in `sql/migrations/`
2. Update schema in `sql/schema/`
3. Add queries in `sql/queries/`
4. Run `sqlc generate` to update generated code
5. Update domain models if needed

### 2. New Feature Development

1. Define domain entities in `internal/domain/`
2. Create service implementation in `internal/services/`
3. Add database queries using SQLc
4. Implement HTTP handlers
5. Write comprehensive tests
6. Update API documentation

### 3. Job Queue Communication

Background jobs are processed via Redis queues:

```go
// Job queue example
type EmailJobPayload struct {
    UserID    int64  `json:"user_id"`
    Subject   string `json:"subject"`
    Body      string `json:"body"`
    Template  string `json:"template"`
}

// Enqueue job from web handler
func (h *Handler) SendWelcomeEmail(userID int64) error {
    payload := EmailJobPayload{
        UserID:   userID,
        Subject:  "Welcome to Ministry Scheduler",
        Template: "welcome",
    }
    
    task := asynq.NewTask("email:send", payload)
    return h.jobClient.Enqueue(task)
}

// Process job in worker
func (w *EmailWorker) ProcessEmailTask(ctx context.Context, t *asynq.Task) error {
    var payload EmailJobPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }
    
    // Send email logic
    return w.emailService.Send(ctx, payload)
}
```

## Configuration Management

### Environment-Based Configuration

```go
type Config struct {
    Port        string `env:"PORT" default:"8080"`
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    JWTSecret   string `env:"JWT_SECRET" required:"true"`
    LogLevel    string `env:"LOG_LEVEL" default:"info"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    return &cfg, envconfig.Process("", &cfg)
}
```

### Database and Redis Configuration

Configuration for distributed components:

```bash
# Database configuration
DATABASE_MASTER_URL=postgresql://user:pass@db-master:5432/ministry_scheduler
DATABASE_REPLICA_URLS=postgresql://user:pass@db-replica-1:5432/ministry_scheduler,postgresql://user:pass@db-replica-2:5432/ministry_scheduler

# Redis configuration
REDIS_URL=redis://redis:6379/0
REDIS_JOB_QUEUE_DB=1
REDIS_CACHE_DB=2
REDIS_SESSION_DB=3

# Load balancer
LOAD_BALANCER_URL=http://nginx:80
```

## Distributed System Patterns

### 1. Database Read/Write Splitting

```go
type DatabaseManager struct {
    master   *sql.DB
    replicas []*sql.DB
    rr       int // round-robin counter
}

func (dm *DatabaseManager) GetWriteDB() *sql.DB {
    return dm.master
}

func (dm *DatabaseManager) GetReadDB() *sql.DB {
    dm.rr = (dm.rr + 1) % len(dm.replicas)
    return dm.replicas[dm.rr]
}
```

### 2. Job Queue with Retry Logic

```go
type JobProcessor struct {
    client *asynq.Client
}

func (jp *JobProcessor) EnqueueWithRetry(taskType string, payload interface{}) error {
    task := asynq.NewTask(taskType, payload)
    
    // Configure retry policy
    opts := []asynq.Option{
        asynq.MaxRetry(3),
        asynq.Timeout(30 * time.Second),
        asynq.ProcessIn(10 * time.Second), // delay
    }
    
    return jp.client.Enqueue(task, opts...)
}
```

### 3. Redis Caching Pattern

```go
type CacheService struct {
    redis *redis.Client
    ttl   time.Duration
}

func (cs *CacheService) GetOrSet(ctx context.Context, key string, fetchFn func() (interface{}, error)) (interface{}, error) {
    // Try cache first
    cached, err := cs.redis.Get(ctx, key).Result()
    if err == nil {
        var result interface{}
        json.Unmarshal([]byte(cached), &result)
        return result, nil
    }
    
    // Cache miss - fetch and store
    data, err := fetchFn()
    if err != nil {
        return nil, err
    }
    
    serialized, _ := json.Marshal(data)
    cs.redis.Set(ctx, key, serialized, cs.ttl)
    
    return data, nil
}
```

### 4. Real-time Updates with PostgreSQL LISTEN/NOTIFY

```go
// Publisher (in web handler)
func (h *Handler) NotifyScheduleChange(ctx context.Context, scheduleID int64) error {
    payload := map[string]interface{}{
        "schedule_id": scheduleID,
        "action":     "updated",
        "timestamp":  time.Now(),
    }
    
    data, _ := json.Marshal(payload)
    _, err := h.db.ExecContext(ctx, "NOTIFY schedule_changes, $1", string(data))
    return err
}

// Subscriber (in WebSocket handler)
func (ws *WebSocketHandler) ListenForUpdates(ctx context.Context) {
    listener := pq.NewListener(ws.dbURL, 10*time.Second, time.Minute, nil)
    listener.Listen("schedule_changes")
    
    for {
        select {
        case n := <-listener.Notify:
            // Broadcast to connected WebSocket clients
            ws.broadcastToClients(n.Extra)
        case <-ctx.Done():
            return
        }
    }
}
```

### 5. Session Management with Redis

```go
type SessionManager struct {
    redis *redis.Client
    ttl   time.Duration
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID int64) (string, error) {
    sessionID := generateSessionID()
    sessionData := map[string]interface{}{
        "user_id":    userID,
        "created_at": time.Now(),
    }
    
    data, _ := json.Marshal(sessionData)
    err := sm.redis.Set(ctx, "session:"+sessionID, data, sm.ttl).Err()
    
    return sessionID, err
}

func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
    data, err := sm.redis.Get(ctx, "session:"+sessionID).Result()
    if err != nil {
        return nil, err
    }
    
    var session SessionData
    json.Unmarshal([]byte(data), &session)
    return &session, nil
}
```

This overview provides the foundation for understanding how we'll structure and organize the new distributed system with background workers, database replication, and Redis-based job queues while maintaining simplicity and avoiding over-engineering.
