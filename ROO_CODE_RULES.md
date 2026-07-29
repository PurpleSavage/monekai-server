# Project Architecture

This project follows a **Hexagonal Architecture (Ports & Adapters)** combined with **Domain-Driven Design (DDD)** principles.

The main goal of this architecture is to keep the business logic isolated from external frameworks, databases, third-party services, and delivery mechanisms, making the codebase highly maintainable, scalable, testable, and easy to evolve.

---

# Project Structure

```text
/cmd
│── /entrypoint                # Application entry point (main.go)
│── /scripts                   # Database migration scripts

/configurations
│── /migrations                # Database migrations
│── /persistence               # GORM models
│── /persistenceconnections    # Database connections (currently PostgreSQL)

/docs
│── /assets                    # Documentation assets
│── *.md                       # Technical documentation

/modules
│── /community
│── /notifications
│── /payments
│── /sampler
│── /shared
```

---

# Module Architecture

Every module follows the exact same internal architecture.

For simplicity, the **Community** module is shown below, but **Notifications**, **Payments**, and **Sampler** all follow the same folder structure and architectural principles.

```text
/community
│
├── application
│   ├── dtos
│   │   ├── requests
│   │   └── responses
│   │
│   ├── ports
│   └── usecases
│
├── domain
│   ├── entities
│   ├── enums
│   └── valueobjects
│
├── infrastructure
│   ├── controllers
│   ├── middlewares
│   ├── in-adapters
│   ├── out-adapters
│   ├── mappers
│   └── raws
│
└── community_module.go
```

---

# Layer Responsibilities

## Application

The Application layer orchestrates business operations.

### `dtos`

Contains all Data Transfer Objects.

- `requests/` → DTOs received from external clients.
- `responses/` → DTOs returned by the application.

### `ports`

Defines interfaces (Ports) that abstract infrastructure implementations.

These interfaces describe *what* the application needs without knowing *how* it is implemented.

### `usecases`

Contains the application's use cases.

A use case coordinates domain objects and ports to execute a business operation.

---

## Domain

The Domain layer contains the pure business logic.

It must never depend on frameworks, databases, HTTP, or external SDKs.

### `entities`

Business entities.

### `valueobjects`

Immutable domain objects with their own validation rules.

Value Objects are responsible for guaranteeing domain correctness.

### `enums`

Strongly typed domain enumerations.

---

## Infrastructure

Infrastructure contains all implementation details.

Anything related to frameworks, HTTP, databases, SDKs, external APIs, or persistence belongs here.

### `controllers`

HTTP controllers.

### `middlewares`

Module-specific middlewares.

### `in-adapters`

Internal adapters.

Examples:

- Observer implementations
- Wrapper services around libraries
- Internal event dispatchers
- Utility services

### `out-adapters`

Outbound adapters.

These communicate with external systems such as:

- Databases
- Cloud services
- External APIs
- Third-party SDKs

### `mappers`

Responsible for converting objects between different layers.

Examples:

- DTO → Entity
- Entity → DTO
- Raw → Entity
- Raw → DTO

### `raws`

Contains raw database projections.

Useful when SQL joins return flattened structures that later need to be mapped into richer domain models.

---

## Module Bootstrap

Each module contains a bootstrap file.

Example:

```
community_module.go
```

The bootstrap is responsible for:

- Initializing repositories
- Initializing services
- Initializing adapters
- Initializing use cases
- Initializing controllers
- Wiring dependency injection

---

# Shared Module

The `shared` module contains components used across the entire application.

```text
/shared
│── auth
│── common
```

## Auth

Contains the complete authentication system.

Examples:

- JWT
- Authentication middleware
- Token services
- Authorization utilities

## Common

Contains reusable components shared across multiple modules.

Examples:

- Shared services
- Utilities
- Event dispatcher
- Common domain logic

---

# Dependency Rules

## Dependency Injection

Dependency Injection is mandatory for all services.

Repositories, use cases, adapters, strategies, factories, services, and similar components should always be injected instead of instantiated directly.

Example:

```go
type CommunityRepository struct {
	db *gorm.DB
}

func NewCommunityRepository(
	db *gorm.DB,
) communityports.CommunityPersistencePort {
	return &CommunityRepository{
		db: db,
	}
}
```

Always return the interface (Port), not the concrete implementation.

---

## Module Independence

Modules must remain independent.

Allowed dependencies:

```text
Community
      │
      ▼
Shared
```

Not allowed:

```text
Shared
      │
      ▼
Community
```

In other words:

- Feature modules may depend on `shared`.
- `shared` must never depend on feature modules.

---

## Cross-Module Communication

Feature modules should not directly call each other whenever possible.

Instead, use events through the shared Observer pattern.

The shared event dispatcher is located at:

```text
/shared/common/infrastructure/in-adapters/event_dispatcher.go
```

This keeps modules loosely coupled and easier to evolve independently.

---

# Middleware Rules

Module-specific middlewares:

- Initialized inside the module bootstrap.
- Injected only into that module.

Shared middlewares:

- Initialized in the application entrypoint.
- Used globally across the application.

The same rule applies to shared services.

---

# Error Handling

There are two different error packages.

## Infrastructure Errors

```text
/shared/common/infrastructure/errors
```

Used inside infrastructure implementations.

Examples:

- Repository errors
- External SDK failures
- Database-related errors
- Adapter errors

---

## Domain Errors

```text
/shared/common/domain/errors
```

Used by the Domain layer.

Examples:

- Value Object validation
- Business rule violations
- Domain service validations

Domain errors should never depend on infrastructure.

---

# Domain Services

Some services belong exclusively to the Domain layer.

Characteristics:

- Pure business logic
- No ports
- No databases
- No HTTP
- No external SDKs
- No infrastructure dependencies

They should only use the Go standard library and domain objects.

---

# Design Principles

The architecture is intentionally flexible.

Whenever appropriate, additional design patterns may be introduced, including:

- Strategy
- Factory
- Builder
- Observer
- Decorator
- Chain of Responsibility

The primary objective is always:

- Loose coupling
- High cohesion
- Easy testing
- Scalability
- Maintainability
- Easy replacement of implementations

---

# Documentation Rules

Whenever a feature introduces significant business logic or architecture, a dedicated markdown document should be added under `/docs`.

Example:

```text
/docs/payment-flow.md
/docs/sampler-generation.md
/docs/webhooks.md
```

Documentation should explain:

- Business flow
- Design decisions
- Architecture
- Sequence diagrams (when useful)
- Important implementation details

---

# Concurrency

Whenever it provides a measurable benefit, Go routines should be used to improve performance.

Concurrency should only be introduced when it improves:

- Throughput
- Latency
- User experience

Avoid unnecessary goroutines that increase complexity without performance gains.

---

# General Guidelines

- Always respect the architecture.
- Keep business logic inside the Domain layer.
- Keep orchestration inside the Application layer.
- Keep implementation details inside Infrastructure.
- Prefer interfaces over concrete implementations.
- Keep modules independent.
- Favor composition over inheritance.
- Use dependency injection consistently.
- Document complex features.
- Keep the codebase clean, modular, and easy to extend.

---

# Technologies

- Go
- PostgreSQL
- GORM
- Chi Router
- Docker
- Docker Compose
- Cloudflare R2
- Paddle
- Replicate
- Server-Sent Events (SSE)