# Ministry Scheduler: Phase 1 Implementation Plan

## Overview

This document outlines the Phase 1 implementation plan for the Ministry Scheduler system - a distributed, modern Go application that will serve as both a learning project for distributed systems and a production-ready church ministry scheduling solution.

## Project Philosophy

**KEEP THINGS SIMPLE. DO NOT OVER ENGINEER.**

While we're building a distributed system to learn distributed patterns, we'll prioritize:

- Simple, clear code over complex abstractions
- Proven patterns over experimental approaches
- Incremental complexity - start simple, add complexity only when needed
- Production readiness with learning opportunities

## Architecture Overview

### Distributed System Design Goals

We'll build a distributed system that demonstrates key concepts while remaining practical:

1. **Architecture**: Horizontally scalable application with background workers
2. **Data Layer**: PostgreSQL cluster with master-slave replication and SQLc for type-safe SQL
3. **Message Queue**: Redis for job queues and pub/sub patterns
4. **Caching**: Redis for session storage and application caching
5. **Communication**: HTTP for web requests, Redis for async job processing, PostgreSQL LISTEN/NOTIFY for real-time
6. **Authentication**: JWT + OAuth2 with Google/Line integration (prepared, implemented later)
7. **Observability**: Structured logging, metrics, distributed tracing
8. **Deployment**: Docker containers with load balancer

### High-Level System Architecture

```
                    ┌─────────────────┐
                    │  Load Balancer  │
                    │   (nginx/       │
                    │   traefik)      │
                    └─────────────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │  Web App    │ │  Web App    │ │  Web App    │
    │ Instance 1  │ │ Instance 2  │ │ Instance N  │
    └─────────────┘ └─────────────┘ └─────────────┘
              │              │              │
              └──────────────┼──────────────┘
                             ▼
                    ┌─────────────────┐
                    │     Redis       │
                    │  ┌─────────────┐ │
                    │  │ Job Queue   │ │
                    │  └─────────────┘ │
                    │  ┌─────────────┐ │
                    │  │   Cache     │ │
                    │  └─────────────┘ │
                    │  ┌─────────────┐ │
                    │  │  Sessions   │ │
                    │  └─────────────┘ │
                    └─────────────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │ Background  │ │ Background  │ │ Background  │
    │  Worker 1   │ │  Worker 2   │ │  Worker N   │
    └─────────────┘ └─────────────┘ └─────────────┘
              │              │              │
              └──────────────┼──────────────┘
                             ▼
                    ┌─────────────────┐
                    │  PostgreSQL     │
                    │    Cluster      │
                    │ ┌─────────────┐ │
                    │ │   Master    │ │
                    │ │  (Write)    │ │
                    │ └─────────────┘ │
                    │        │        │
                    │ ┌──────▼──────┐ │
                    │ │   Replica   │ │
                    │ │   (Read)    │ │
                    │ └─────────────┘ │
                    │ ┌─────────────┐ │
                    │ │   Replica   │ │
                    │ │   (Read)    │ │
                    │ └─────────────┘ │
                    └─────────────────┘
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Gateway   │    │    Auth     │    │   Users     │
│   Service   │────│   Service   │────│   Service   │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       │            ┌─────────────┐    ┌─────────────┐
       │────────────│   Events    │────│ Schedules   │
       │            │   Service   │    │   Service   │
       │            └─────────────┘    └─────────────┘
       │                   │                   │
┌─────────────┐    ┌─────────────────────────────────┐
│ PostgreSQL  │    │         Shared Libs             │
│  Database   │    │  (domain, utils, middleware)    │
└─────────────┘    └─────────────────────────────────┘
```

## Project Structure

### Distributed Application Architecture

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

### Code Organization Principles

1. **Single Application, Multiple Processes**: Web server + background workers sharing the same codebase
2. **Domain-First**: Business logic separated from infrastructure
3. **Async Processing**: Heavy operations moved to background jobs
4. **Shared Components**: Common functionality in internal packages
5. **Database-Centric**: PostgreSQL as the source of truth with proper replication
6. **Message-Driven**: Redis for job queues and real-time communication

## Technology Stack

### Core Technologies

- **Language**: Go 1.21+
- **Database**: PostgreSQL with master-slave replication
- **Connection Pooling**: PgBouncer for database connection management
- **Message Queue & Cache**: Redis for job queues, pub/sub, and caching
- **SQL Management**: SQLc for type-safe database operations
- **HTTP**: Standard library + gorilla/mux for routing
- **Authentication**: JWT + OAuth2 (Google, Line) - prepared for Phase 2
- **Configuration**: Environment variables + YAML for complex configs

### Development Tools

- **Database Migrations**: golang-migrate/migrate
- **Job Queue**: github.com/hibiken/asynq (Redis-based)
- **Caching**: github.com/go-redis/redis/v9
- **API Documentation**: OpenAPI 3.0 spec generation
- **Linting**: golangci-lint
- **Testing**: Standard testing package + testify for assertions

### Distributed System Components

- **Load Balancing**: nginx/traefik for request distribution
- **Database Replication**: PostgreSQL streaming replication
- **Caching Strategy**: Redis for application cache + session storage
- **Background Processing**: Redis job queues with worker pools
- **Real-time Communication**: PostgreSQL LISTEN/NOTIFY + WebSocket
- **Monitoring**: Prometheus metrics + structured logging
- **Tracing**: OpenTelemetry (basic implementation)

## Phase 1 Features

### 1. Web Application (cmd/web)

- **Purpose**: Main HTTP server handling user requests
- **Responsibilities**:
  - Serve HTTP API endpoints
  - Handle user authentication (JWT)
  - Process synchronous operations
  - Queue background jobs
  - Serve static files (if any)
- **Processes**: Multiple instances behind load balancer
- **Database**: Read/write to PostgreSQL cluster

### 2. Background Workers (cmd/worker)

- **Purpose**: Process asynchronous tasks
- **Responsibilities**:
  - Email notifications and reminders
  - Recurring event generation
  - Data cleanup and maintenance
  - Heavy report generation
  - Third-party API integrations
- **Processes**: Worker pool processing Redis job queues
- **Database**: Read/write to PostgreSQL cluster

### 3. Core Feature Areas

#### User Management

- **Web Handlers**: User registration, profile management, authentication
- **Background Jobs**: Welcome emails, profile image processing
- **Database**: Users table with OAuth2 preparation
- **Endpoints**: `/api/users/*`, `/api/auth/*`

#### Event Management

- **Web Handlers**: CRUD operations for ministry events
- **Background Jobs**: Recurring event creation, event reminders
- **Database**: Events table with recurrence rules
- **Endpoints**: `/api/events/*`

#### Position Management

- **Web Handlers**: Define ministry positions and requirements
- **Background Jobs**: Position availability notifications
- **Database**: Positions table
- **Endpoints**: `/api/positions/*`

#### Schedule Management

- **Web Handlers**: Create assignments, manage schedules
- **Background Jobs**: Schedule conflict detection, assignment reminders
- **Database**: Assignments table linking users, events, positions
- **Endpoints**: `/api/schedules/*`, `/api/assignments/*`

#### Real-time Features

- **Web Handlers**: WebSocket connections for live updates
- **Background Jobs**: PostgreSQL LISTEN/NOTIFY message broadcasting
- **Database**: PostgreSQL NOTIFY for live schedule changes
- **Endpoints**: `/ws` (WebSocket endpoint)

## Database Design with SQLc

### SQLc Setup

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sql/queries/"
    schema: "sql/schema/"
    gen:
      go:
        package: "db"
        out: "internal/infrastructure/database/generated"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
```

### Core Database Schema

#### Users Table

```sql
-- sql/schema/001_users.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('leader', 'member')),
    avatar_url TEXT,
    google_id VARCHAR(255) UNIQUE,
    line_id VARCHAR(255) UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
CREATE INDEX idx_users_line_id ON users(line_id) WHERE line_id IS NOT NULL;
```

#### Positions Table

```sql
-- sql/schema/002_positions.sql
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7), -- hex color code
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_positions_active ON positions(is_active);
```

#### Events Table

```sql
-- sql/schema/003_events.sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    is_recurring BOOLEAN DEFAULT false,
    recurrence_rule JSONB, -- RFC 5545 RRULE as JSON
    created_by INTEGER NOT NULL REFERENCES users(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_created_by ON events(created_by);
CREATE INDEX idx_events_active ON events(is_active);
```

#### Assignments Table

```sql
-- sql/schema/004_assignments.sql
CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    position_id INTEGER NOT NULL REFERENCES positions(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'assigned'
        CHECK (status IN ('assigned', 'confirmed', 'declined', 'substitute_needed')),
    notes TEXT,
    assigned_by INTEGER NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(event_id, position_id, user_id)
);

CREATE INDEX idx_assignments_event_id ON assignments(event_id);
CREATE INDEX idx_assignments_user_id ON assignments(user_id);
CREATE INDEX idx_assignments_status ON assignments(status);
```

### Sample SQLc Queries

```sql
-- sql/queries/users.sql

-- name: CreateUser :one
INSERT INTO users (email, name, role, google_id, line_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 AND is_active = true;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1 AND is_active = true;

-- name: UpdateUser :one
UPDATE users 
SET name = COALESCE($2, name),
    avatar_url = COALESCE($3, avatar_url),
    updated_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING *;
```

## Implementation Order

### Phase 1.1: Foundation & Database Setup (Week 1-2)

1. Set up project structure and build system
2. Configure PostgreSQL master-slave replication with Docker Compose
3. Set up Redis for job queues and caching
4. Configure SQLc and database migrations
5. Implement basic web application (single instance)
6. Implement user management with JWT authentication
7. Set up basic background worker infrastructure

### Phase 1.2: Core Features & Job Processing (Week 3-4)  

1. Implement position management
2. Implement event management with basic CRUD
3. Implement schedule/assignment management
4. Add background job processing (email notifications, etc.)
5. Implement database read/write splitting
6. Add Redis caching for frequently accessed data

### Phase 1.3: Distribution & Scaling (Week 5-6)

1. Set up load balancer (nginx) for multiple web instances
2. Implement PostgreSQL LISTEN/NOTIFY for real-time updates
3. Add WebSocket support for live schedule changes
4. Implement worker pools and job retry mechanisms
5. Add comprehensive monitoring and observability
6. Performance testing and optimization
7. Add docker-compose for full distributed setup

## Distributed System Learning Points

### Key Concepts to Implement and Learn

1. **Horizontal Scaling**
   - Multiple application instances behind load balancer
   - Stateless application design
   - Session management with Redis
   - Load balancing strategies (round-robin, least connections)

2. **Database Distribution**
   - Master-slave replication setup
   - Read/write splitting in application code
   - Connection pooling with PgBouncer
   - Database failover and recovery

3. **Asynchronous Processing**
   - Job queues with Redis and Asynq
   - Worker pool management
   - Job retry mechanisms and dead letter queues
   - Background task scheduling

4. **Caching Strategies**
   - Application-level caching with Redis
   - Cache invalidation patterns
   - Cache-aside vs write-through patterns
   - Session storage in Redis

5. **Real-time Communication**
   - WebSocket connections for live updates
   - PostgreSQL LISTEN/NOTIFY for pub/sub
   - Broadcasting messages across multiple app instances
   - Connection management and cleanup

6. **Data Consistency**
   - ACID transactions in PostgreSQL
   - Eventual consistency with background jobs
   - Optimistic locking for concurrent updates
   - Database constraints and data integrity

7. **Observability & Monitoring**
   - Structured logging across processes
   - Application metrics (request rate, error rate, response time)
   - Database performance monitoring
   - Queue depth and worker performance monitoring

8. **Resilience Patterns**
   - Graceful shutdown of web servers and workers
   - Health checks for all components
   - Circuit breaker for external service calls
   - Timeout handling and retry logic

### Learning Path

1. **Week 1-2**: Single instance with database replication
2. **Week 3-4**: Add background jobs and Redis caching
3. **Week 5-6**: Scale to multiple instances with load balancing
4. **Ongoing**: Monitor, optimize, and add resilience patterns

## OAuth2 Authentication Design (Prepared for Phase 2)

### Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │    Auth     │    │   Google/   │
│ (Frontend)  │    │   Service   │    │    Line     │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       │ 1. Login Request  │                   │
       │──────────────────→│                   │
       │                   │ 2. OAuth2 Flow    │
       │                   │──────────────────→│
       │                   │←──────────────────│
       │ 3. JWT Token      │ 3. User Info      │
       │←──────────────────│                   │
```

### Database Design for OAuth2

```sql
CREATE TABLE oauth_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- 'google', 'line'
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL, -- Encrypted
    auth_url VARCHAR(500) NOT NULL,
    token_url VARCHAR(500) NOT NULL,
    user_info_url VARCHAR(500) NOT NULL,
    scopes TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE user_oauth_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    provider_id INTEGER NOT NULL REFERENCES oauth_providers(id),
    provider_user_id VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255),
    provider_name VARCHAR(255),
    access_token TEXT, -- Encrypted
    refresh_token TEXT, -- Encrypted
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(provider_id, provider_user_id)
);
```

## Success Criteria

### Phase 1 Completion Indicators

1. **Distributed Application**: Web app + background workers running with load balancer
2. **Database Cluster**: PostgreSQL master-slave replication working
3. **Job Processing**: Redis job queues with worker pools processing tasks
4. **API Coverage**: All CRUD operations for core entities (users, events, positions, assignments)
5. **Basic Authentication**: JWT-based auth system
6. **Real-time Features**: WebSocket + PostgreSQL LISTEN/NOTIFY working
7. **Caching**: Redis caching for performance optimization
8. **Documentation**: API docs and deployment guides
9. **Testing**: Unit and integration test coverage >80%
10. **Deployment**: Docker-compose for full distributed setup

### Non-Functional Requirements

1. **Performance**: API response times <200ms for simple operations
2. **Scalability**: Multiple web instances handling increased load
3. **Reliability**: Application handles database failover and Redis outages gracefully
4. **Data Consistency**: Database transactions maintain data integrity
5. **Maintainability**: Clear code structure and documentation

## Next Phases Preview

- **Phase 2**: Leave management, shift swapping, OAuth2 integration
- **Phase 3**: Advanced scheduling rules, notifications
- **Phase 4**: External integrations (Line Bot, Google Calendar)
- **Phase 5**: Real-time features, analytics, reporting

## Getting Started

1. Review this plan and provide feedback
2. Set up development environment (Go, PostgreSQL, Docker)
3. Initialize project structure
4. Begin with Phase 1.1 implementation

This plan balances learning distributed systems concepts with delivering a practical ministry scheduling solution while maintaining simplicity and avoiding over-engineering.
