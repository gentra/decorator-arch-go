# Task Management System - Technical Plan

## Project Overview

This document outlines the implementation plan for a Task Management System that serves as a comprehensive example of the Decorator Architecture pattern in Go. The system will demonstrate all key benefits of the architecture including flexibility, testability, separation of concerns, and easy microservice extraction.

## Architecture Goals Demonstration

| Architecture Goal | How Demo Achieves It |
|------------------|---------------------|
| **Monolith → Microservices** | Start with all domains in one service, then extract notification service |
| **Feature Flags** | Multiple algorithm implementations with strategy pattern |
| **Easy Testing** | Each decorator layer can be tested in isolation |
| **Separation of Concerns** | Each layer has single responsibility |
| **Composability** | Layers can be added/removed without affecting others |

## Domain Structure

### 1. User Domain (`internal/user/`)
**Interface**: User management, authentication, and preferences
```go
type Service interface {
    Register(ctx context.Context, data RegisterData) (*User, error)
    Login(ctx context.Context, email, password string) (*AuthResult, error)
    GetByID(ctx context.Context, id string) (*User, error)
    UpdateProfile(ctx context.Context, id string, data UpdateProfileData) (*User, error)
    GetPreferences(ctx context.Context, userID string) (*UserPreferences, error)
    UpdatePreferences(ctx context.Context, userID string, prefs UserPreferences) error
}
```

### 2. Task Domain (`internal/task/`)
**Interface**: Task CRUD operations, priority management, status updates
```go
type Service interface {
    Create(ctx context.Context, data CreateTaskData) (*Task, error)
    GetByID(ctx context.Context, id string) (*Task, error)
    List(ctx context.Context, filter TaskFilter) ([]*Task, error)
    Update(ctx context.Context, id string, data UpdateTaskData) (*Task, error)
    Delete(ctx context.Context, id string) error
    UpdateStatus(ctx context.Context, id string, status TaskStatus) (*Task, error)
    CalculatePriority(ctx context.Context, taskID string) (*PriorityResult, error)
}
```

### 3. Project Domain (`internal/project/`)
**Interface**: Workspace/project organization and team collaboration
```go
type Service interface {
    Create(ctx context.Context, data CreateProjectData) (*Project, error)
    GetByID(ctx context.Context, id string) (*Project, error)
    List(ctx context.Context, userID string) ([]*Project, error)
    AddMember(ctx context.Context, projectID, userID string, role Role) error
    RemoveMember(ctx context.Context, projectID, userID string) error
    GetMembers(ctx context.Context, projectID string) ([]*Member, error)
    Archive(ctx context.Context, id string) error
}
```

### 4. Notification Domain (`internal/notification/`)
**Interface**: Multi-channel notification system
```go
type Service interface {
    Send(ctx context.Context, notification Notification) error
    SendBulk(ctx context.Context, notifications []Notification) error
    GetHistory(ctx context.Context, userID string, filter HistoryFilter) ([]*NotificationHistory, error)
    UpdatePreferences(ctx context.Context, userID string, prefs NotificationPreferences) error
    MarkAsRead(ctx context.Context, notificationID string) error
}
```

## Decorator Layer Implementation

### Task Domain Example (Complete Decorator Chain)

```go
// Storage Layer (GORM)
taskGorm := gorm.NewTaskService(db)

// Caching Layer  
taskCache := redis.NewTaskService(taskGorm, redisClient, 5*time.Minute)

// Audit Layer
taskAudit := audit.NewTaskService(taskCache, auditLogger)

// Rate Limiting Layer
taskRateLimit := ratelimit.NewTaskService(taskAudit, limiter)

// Encryption Layer (for sensitive task data)
taskEncryption := encryption.NewTaskService(taskRateLimit, encryptionKey)

// Validation Layer
taskValidation := validation.NewTaskService(taskEncryption, validator)

// Business Logic Layer
taskCore := usecase.NewTaskService(taskValidation, notificationService, projectService)
```

### Decorator Responsibilities

| Layer | Responsibility | Example Implementation |
|-------|---------------|----------------------|
| **GORM** | Data persistence | CRUD operations, transactions, relationships |
| **Redis** | Caching | Cache-aside pattern, TTL management |
| **Audit** | Activity logging | Log all operations with metadata |
| **Rate Limit** | API protection | Per-user/IP rate limiting |
| **Encryption** | Data security | Encrypt/decrypt sensitive fields |
| **Validation** | Input validation | Business rule validation |
| **UseCase** | Business logic | Core domain logic, orchestration |

## Feature Flags & Strategy Pattern Implementation

### 1. Task Priority Algorithms

**Feature Flag**: `TASK_PRIORITY_ALGORITHM`

```go
// internal/task/prioritization/
├── simple.go          # Basic High/Medium/Low priority
├── eisenhower.go      # Eisenhower Matrix (Urgent/Important)
├── kanban.go          # Kanban-style prioritization
└── weighted.go        # Weighted scoring algorithm
```

**Factory Implementation**:
```go
func NewPriorityCalculator(algorithm string) PriorityCalculator {
    switch algorithm {
    case "eisenhower":
        return &EisenhowerCalculator{}
    case "kanban":
        return &KanbanCalculator{}
    case "weighted":
        return &WeightedCalculator{}
    default:
        return &SimpleCalculator{}
    }
}
```

### 2. Notification Channels

**Feature Flag**: `NOTIFICATION_CHANNELS`

```go
// internal/notification/channels/
├── email/             # Email notifications (SMTP)
├── push/              # Push notifications (Firebase)
├── sms/               # SMS notifications (Twilio)
├── slack/             # Slack integration
└── webhook/           # Generic webhook notifications
```

### 3. Authentication Methods

**Feature Flag**: `AUTH_METHOD`

```go
// internal/user/auth/
├── basic.go           # Username/password
├── oauth.go           # OAuth2 (Google, GitHub)
├── jwt.go             # JWT token-based
└── ldap.go            # LDAP integration
```

## Database Schema

### GORM Models

```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/datatypes"
)

// User represents a user in the system
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    FirstName    string    `gorm:"not null" json:"first_name"`
    LastName     string    `gorm:"not null" json:"last_name"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    
    // Relationships
    OwnedProjects []Project `gorm:"foreignKey:OwnerID" json:"owned_projects,omitempty"`
    AssignedTasks []Task    `gorm:"foreignKey:AssigneeID" json:"assigned_tasks,omitempty"`
    Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

// Project represents a project/workspace
type Project struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name        string    `gorm:"not null" json:"name"`
    Description string    `json:"description"`
    OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
    Status      string    `gorm:"default:active" json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relationships
    Owner User   `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
    Tasks []Task `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// Task represents a task within a project
type Task struct {
    ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Title       string     `gorm:"not null" json:"title"`
    Description string     `json:"description"`
    ProjectID   uuid.UUID  `gorm:"type:uuid;not null" json:"project_id"`
    AssigneeID  *uuid.UUID `gorm:"type:uuid" json:"assignee_id"`
    Status      string     `gorm:"default:todo" json:"status"`
    Priority    string     `gorm:"default:medium" json:"priority"`
    DueDate     *time.Time `json:"due_date"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    
    // Relationships
    Project  Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
    Assignee *User   `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
}

// Notification represents a notification sent to a user
type Notification struct {
    ID        uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID    uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
    Type      string           `gorm:"not null" json:"type"`
    Title     string           `gorm:"not null" json:"title"`
    Message   string           `gorm:"not null" json:"message"`
    Channel   string           `gorm:"not null" json:"channel"`
    Status    string           `gorm:"default:pending" json:"status"`
    Metadata  datatypes.JSON   `json:"metadata"`
    CreatedAt time.Time        `json:"created_at"`
    
    // Relationships
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return nil
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
    if n.ID == uuid.Nil {
        n.ID = uuid.New()
    }
    return nil
}
```

### Database Migrations using golang-migrate

#### Migration Setup
```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration directory
mkdir -p migrations

# Create initial migration
migrate create -ext sql -dir migrations -seq create_initial_tables
```

#### Migration Files
**migrations/000001_create_initial_tables.up.sql**
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    project_id UUID REFERENCES projects(id),
    assignee_id UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'todo',
    priority VARCHAR(50) DEFAULT 'medium',
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    channel VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_projects_owner_id ON projects(owner_id);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
```

**migrations/000001_create_initial_tables.down.sql**
```sql
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
```

## API Design

### REST Endpoints

```
# Authentication
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/logout
GET    /api/auth/me

# Users
GET    /api/users/profile
PUT    /api/users/profile
GET    /api/users/preferences
PUT    /api/users/preferences

# Projects
GET    /api/projects
POST   /api/projects
GET    /api/projects/:id
PUT    /api/projects/:id
DELETE /api/projects/:id
POST   /api/projects/:id/members
DELETE /api/projects/:id/members/:user_id

# Tasks
GET    /api/projects/:project_id/tasks
POST   /api/projects/:project_id/tasks
GET    /api/tasks/:id
PUT    /api/tasks/:id
DELETE /api/tasks/:id
PUT    /api/tasks/:id/status
POST   /api/tasks/:id/priority/calculate

# Notifications
GET    /api/notifications
PUT    /api/notifications/:id/read
POST   /api/notifications/test
GET    /api/notifications/preferences
PUT    /api/notifications/preferences
```

## Testing Strategy

### Unit Testing (Table Tests with Gherkin Syntax)
Each decorator layer will have comprehensive table-driven tests using Gherkin syntax for clarity.

#### Mock Generation
Use `testify/mock` to generate mocks for all service interfaces:

```bash
# Add testify to dependencies
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require

# Install mockery for generating testify mocks
go install github.com/vektra/mockery/v2@latest

# Generate mocks for all service interfaces using mockery config
mockery
```

#### Example Unit Test
```go
package redis_test

import (
    "context"
    "encoding/json"
    "errors"
    "testing"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "your-project/internal/task"
    "your-project/internal/task/mocks"
    taskRedis "your-project/internal/task/redis"
)

func TestTaskCacheService_GetByID(t *testing.T) {
    tests := []struct {
        name           string
        setupMocks     func(*mocks.Service, *redis.Client)
        taskID         string
        expectedTask   *task.Task
        expectedError  error
        expectedCalls  int
    }{
        {
            name: "Given task not in cache, When GetByID is called, Then should fetch from next service and cache result",
            setupMocks: func(mockNext *mocks.Service, redisClient *redis.Client) {
                taskData := &task.Task{ID: "task-1", Title: "Test Task"}
                mockNext.On("GetByID", mock.Anything, "task-1").Return(taskData, nil)
                redisClient.FlushAll(context.Background()) // Ensure cache is empty
            },
            taskID:        "task-1",
            expectedTask:  &task.Task{ID: "task-1", Title: "Test Task"},
            expectedError: nil,
            expectedCalls: 1,
        },
        {
            name: "Given task exists in cache, When GetByID is called, Then should return cached result without calling next service",
            setupMocks: func(mockNext *mocks.Service, redisClient *redis.Client) {
                taskData := &task.Task{ID: "task-2", Title: "Cached Task"}
                // Pre-populate cache
                cacheKey := "task:task-2"
                taskJSON, _ := json.Marshal(taskData)
                redisClient.Set(context.Background(), cacheKey, taskJSON, time.Minute)
            },
            taskID:        "task-2",
            expectedTask:  &task.Task{ID: "task-2", Title: "Cached Task"},
            expectedError: nil,
            expectedCalls: 0, // Should not call next service
        },
        {
            name: "Given next service returns error, When GetByID is called, Then should return error and not cache anything",
            setupMocks: func(mockNext *mocks.Service, redisClient *redis.Client) {
                mockNext.On("GetByID", mock.Anything, "task-3").Return(nil, errors.New("database error"))
                redisClient.FlushAll(context.Background())
            },
            taskID:        "task-3",
            expectedTask:  nil,
            expectedError: errors.New("database error"),
            expectedCalls: 1,
        },
        {
            name: "Given cache is down, When GetByID is called, Then should fallback to next service gracefully",
            setupMocks: func(mockNext *mocks.Service, redisClient *redis.Client) {
                taskData := &task.Task{ID: "task-4", Title: "Fallback Task"}
                mockNext.On("GetByID", mock.Anything, "task-4").Return(taskData, nil)
                // Simulate Redis down by closing connection
                redisClient.Close()
            },
            taskID:        "task-4",
            expectedTask:  &task.Task{ID: "task-4", Title: "Fallback Task"},
            expectedError: nil,
            expectedCalls: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockNext := mocks.NewService(t)
            redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
            cache := taskRedis.NewService(mockNext, redisClient, time.Minute)
            
            tt.setupMocks(mockNext, redisClient)
            
            // Act
            result, err := cache.GetByID(context.Background(), tt.taskID)
            
            // Assert
            if tt.expectedError != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedError.Error(), err.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedTask, result)
            }
            
            // Verify mock expectations
            mockNext.AssertNumberOfCalls(t, "GetByID", tt.expectedCalls)
            mockNext.AssertExpectations(t)
        })
    }
}
```

#### Mock Configuration (.mockery.yaml)
Create a `.mockery.yaml` file in the project root for consistent mock generation:

```yaml
with-expecter: true
all: true
dir: "{{.InterfaceDir}}/mocks"
filename: "{{.InterfaceName | snakecase}}.go"
mockname: "{{.InterfaceName}}"
outpkg: mocks
packages:
  your-project/internal/user:
    interfaces:
      Service:
  your-project/internal/task:
    interfaces:
      Service:
  your-project/internal/project:
    interfaces:
      Service:
  your-project/internal/notification:
    interfaces:
      Service:
```



### Performance Testing (Table Tests with Gherkin Syntax)
Benchmark each decorator layer with different scenarios:

```go
func BenchmarkTaskService_GetByID(b *testing.B) {
    benchmarks := []struct {
        name      string
        setupFunc func() task.Service
        taskID    string
    }{
        {
            name: "Given only GORM layer, When GetByID is called 1000 times, Then should measure baseline performance",
            setupFunc: func() task.Service {
                db := setupBenchmarkDB()
                return gorm.NewTaskService(db)
            },
            taskID: "benchmark-task-1",
        },
        {
            name: "Given GORM + cache layers, When GetByID is called 1000 times, Then should measure cache layer overhead",
            setupFunc: func() task.Service {
                db := setupBenchmarkDB()
                redisClient := setupBenchmarkRedis()
                gormService := gorm.NewTaskService(db)
                return redis.NewTaskService(gormService, redisClient, 5*time.Minute)
            },
            taskID: "benchmark-task-2",
        },
        {
            name: "Given full decorator chain, When GetByID is called 1000 times, Then should measure complete chain overhead",
            setupFunc: func() task.Service {
                db := setupBenchmarkDB()
                features := map[string]bool{
                    "ENABLE_CACHE":      true,
                    "ENABLE_AUDIT":      true,
                    "ENABLE_VALIDATION": true,
                    "ENABLE_ENCRYPTION": true,
                    "ENABLE_RATE_LIMIT": true,
                }
                return setupTaskServiceChain(db, features)
            },
            taskID: "benchmark-task-3",
        },
    }

    for _, bb := range benchmarks {
        b.Run(bb.name, func(b *testing.B) {
            // Arrange
            taskService := bb.setupFunc()
            ctx := context.Background()
            
            // Ensure task exists
            createTestTask(bb.taskID)
            
            // Act & Measure
            b.ResetTimer()
            b.ReportAllocs()
            
            for i := 0; i < b.N; i++ {
                _, err := taskService.GetByID(ctx, bb.taskID)
                if err != nil {
                    b.Fatalf("Benchmark failed: %v", err)
                }
            }
        })
    }
}
```

### Test Utilities
Helper functions to support the Gherkin-style table tests:

```go
package testutil

import (
    "database/sql"
    "time"

    "github.com/stretchr/testify/mock"
    
    "your-project/internal/task"
    "your-project/internal/task/mocks"
)

// Test helpers for verifying decorator effects
func isTaskInCache(taskID string) bool {
    // Implementation to check if task exists in Redis cache
}

func isInAuditLog(operation, entityID string) bool {
    // Implementation to check if operation was logged
}

func isValidTaskData(task *task.Task) bool {
    // Implementation to verify task data meets business rules
}

func isSensitiveDataEncrypted(taskID string) bool {
    // Implementation to verify sensitive fields are encrypted in DB
}

func getTaskCountInDB() int {
    // Implementation to count tasks in database
}

func setupTaskServiceChain(db *sql.DB, features map[string]bool) task.Service {
    // Implementation to build decorator chain based on feature flags
}

// Mock helper functions for common test setups
func SetupTaskServiceMock(t *testing.T) *mocks.Service {
    mockService := mocks.NewService(t)
    
    // Common mock expectations that most tests need
    mockService.EXPECT().GetByID(mock.Anything, mock.AnythingOfType("string")).
        Return(nil, errors.New("not found")).Maybe()
    
    return mockService
}

func SetupValidTaskMock(t *testing.T, taskID string) (*mocks.Service, *task.Task) {
    mockService := mocks.NewService(t)
    testTask := &task.Task{
        ID:          taskID,
        Title:       "Test Task",
        Description: "Test Description",
        Status:      "todo",
        Priority:    "medium",
        CreatedAt:   time.Now(),
    }
    
    mockService.EXPECT().GetByID(mock.Anything, taskID).
        Return(testTask, nil)
    
    return mockService, testTask
}
```

#### Makefile for Testing
Create a `Makefile` to automate testing workflows:

```makefile
.PHONY: test test-unit test-coverage mocks migrate migrate-up migrate-down

# Run all tests
test: mocks
	go test ./... -v

# Run only unit tests  
test-unit: mocks
	go test ./... -v -short

# Generate test coverage report
test-coverage: mocks
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Generate mocks using mockery
mocks:
	mockery

# Database migration commands
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)

# Clean generated files
clean:
	find . -name "mocks" -type d -exec rm -rf {} +
	rm -f coverage.out coverage.html
```

## Configuration Management

### Environment Variables
```bash
# Database
DATABASE_URL=postgres://user:pass@localhost/taskdb
REDIS_URL=redis://localhost:6379

# Feature Flags
TASK_PRIORITY_ALGORITHM=eisenhower
AUTH_METHOD=jwt
NOTIFICATION_CHANNELS=email,push,slack

# Layer Configuration
CACHE_TTL=300s
RATE_LIMIT_PER_MINUTE=60
AUDIT_LOG_LEVEL=info
ENCRYPTION_KEY=base64-encoded-key
```

### Feature Flag Service Integration
```go
type FeatureFlagService interface {
    IsEnabled(ctx context.Context, flag string) bool
    GetStringValue(ctx context.Context, flag string) string
    GetIntValue(ctx context.Context, flag string) int
}

// Usage in factory
func NewTaskService(flags FeatureFlagService, deps Dependencies) task.Service {
    algorithm := flags.GetStringValue(ctx, "TASK_PRIORITY_ALGORITHM")
    priorityCalculator := NewPriorityCalculator(algorithm)
    
    // Build decorator chain based on flags
    var service task.Service = gorm.NewTaskService(deps.DB)
    
    if flags.IsEnabled(ctx, "ENABLE_CACHE") {
        service = redis.NewTaskService(service, deps.Redis, 5*time.Minute)
    }
    
    if flags.IsEnabled(ctx, "ENABLE_AUDIT") {
        service = audit.NewTaskService(service, deps.AuditLogger)
    }
    
    return usecase.NewTaskService(service, priorityCalculator)
}
```

## Success Metrics

### Architecture Demonstration Success
- [ ] Each decorator layer can be added/removed without code changes
- [ ] Feature flags can switch between implementations at runtime
- [ ] Complete test coverage with isolated layer testing
- [ ] Successful microservice extraction with minimal changes
- [ ] Performance benchmarks showing layer overhead < 5%

### Code Quality Metrics
- [ ] Test coverage > 90%
- [ ] Cyclomatic complexity < 10 per function
- [ ] No circular dependencies
- [ ] All layers implement single responsibility principle
- [ ] Documentation coverage for all public interfaces

## Future Enhancements

### Advanced Patterns
- [ ] Event sourcing layer for audit trail
- [ ] CQRS implementation with separate read/write models
- [ ] Saga pattern for distributed transactions
- [ ] Circuit breaker pattern for external service calls

### Operational Features
- [ ] Health check endpoints for each layer
- [ ] Graceful shutdown with proper cleanup
- [ ] Configuration hot-reloading
- [ ] A/B testing framework integration

## Conclusion

This Task Management System will serve as a comprehensive example of the Decorator Architecture pattern, demonstrating all the key benefits outlined in the main README. The implementation will showcase how to build flexible, testable, and maintainable systems that can evolve from monoliths to microservices with minimal friction.
