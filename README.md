# entrepreneur-pastoral
This is a Catholic Marketplace Platform built as a modular monolith that brings together Catholic entrepreneurs to share products, services, and job opportunities within the Catholic community.

## Architecture & Tech Stack
### Backend:
- **Language**: Go 1.24
- **Framework**: Chi router for HTTP routing
- **Architecture** Pattern: DDD (Domain-Driven Design) with modular monolith structure
- **Database**: PostgreSQL with sqlx and Squirrel query builder
- **Caching**: Redis for session management and rate limiting
- **Queue**: RabbitMQ (planned for notifications)
- **Authentication**: JWT tokens with bcrypt password hashing
- **Logging**: Uber Zap structured logging

### DevOps:
- **Containerization**: Docker & Docker Compose
- **Database** Migrations: golang-migrate

### Project Structure

```
├── cmd/api/                   # Application entry point
│   ├── main.go                # Server initialization & graceful shutdown
│   ├── orchestrator/          # Dependency injection & composition
│   └── router/                # HTTP routing & middleware setup
├── internal/                  # Business logic (modules)
│   ├── user/                  # User module (currently implemented)
│   ├── admin/                 # Admin module (placeholder)
│   ├── entrepreneur/          # Entrepreneur module (placeholder)
│   └── marketplace/           # Marketplace module (placeholder)
├── pkg/                       # Shared utilities
│   ├── config/                # Configuration management
│   ├── database/              # DB connection & migrations
│   ├── helper/                # Auth, constants, env helpers
│   ├── logger/                # Structured logging
│   └── storage/               # Redis cache abstraction
```

### Design Patterns Used
- **Repository Pattern** - Abstraction over data access
- **Unit of Work** - Transaction management
- **Dependency Injection** - Via Orchestrator
- **DTO Pattern** - Separation of API contracts from domain
- **Builder Pattern** - Squirrel query builder
- **Middleware Pattern** - Chi middleware stack

## Modules
The system will have the following modules:
- **User**: Users data, Job profile, Roles and Notifications management.
- **Admin**: Configuration related stuff (users role, business config/approval, permissions/role management and pastoral info).
- **Entrepreneur**: Business data and contact info, products and services offered and jobs available.
- **Marketplace**: A meeting point for all users and business, where business can find suppliers and workers and users can find services/products or jobs.
- **Notification**: Handles all stuff related to notifications, based on user config.

