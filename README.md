# 🚀 Thusa One Project

A Go-based application with data models and environment configuration for enterprise integrations.

---

## 🛠️ Getting Started

### Prerequisites

- 🐹 [Go](https://golang.org/) (version 1.24.5 or later)
- 🐘 PostgreSQL database

### Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd threereco-nextgen
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

---

## 🏗️ Project Structure

```
threereco-nextgen/
├── env/                    # Environment configuration
│   └── env.go              # Environment variables and constants
├── internal/               # Internal application code
│   └── models/             # Data models
│       ├── audit_log.go    # Audit logging model
│       ├── organization.go # Organization model
│       └── user.go         # User model
├── ecosystem.config.js     # PM2 process management config
├── go.mod                  # Go module definition
└── go.sum                  # Go module checksums
```

---

## 📦 Dependencies

The project uses the following key dependencies:

- **[GORM](https://gorm.io/)** - Go ORM for database operations
- **[go-json](https://github.com/goccy/go-json)** - Fast JSON encoding/decoding
- **[UUID](https://github.com/google/uuid)** - UUID generation and parsing

---

## 🗃️ Data Models

### User Model

The `User` model represents application users with the following features:

- Unique UUID identifier
- Email address (unique)
- Organization association
- Self-referencing creator tracking
- Automatic audit logging

### Organization Model

The `Organization` model represents business entities with:

- Unique UUID identifier
- Organization name and domain
- Creator tracking
- Automatic audit logging

### Audit Log Model

The `AuditLog` model tracks all data changes:

- Event descriptions
- JSON data snapshots
- Automatic timestamping

---

## ⚙️ Environment Configuration

Environment variables are defined in `env/env.go` and include:

- **Database connections** for PostgreSQL and various warehouse schemas
- **API integrations** for Microsoft, Autotask, CyberCNS, RocketCyber, and more
- **Authentication secrets** and API keys
- **Service configurations**

---

## 🏗️ Using the Models

### Basic Usage

```go
import "github.com/connor-davis/threereco-nextgen/internal/models"

// Create a new organization
org := models.Organization{
    Name:   "Example Corp",
    Domain: "example.com",
}

// Create a new user
user := models.User{
    Email:          "user@example.com",
    OrganizationId: org.Id,
}
```

### GORM Integration

All models include GORM hooks for automatic audit logging:

```go
// Models automatically create audit logs when:
// - AfterCreate: Record is created
// - AfterUpdate: Record is updated
// - AfterDelete: Record is deleted

// Each model has a ToJSON() method for serialization
jsonData, err := user.ToJSON()
```

## 📚 Learn More

Expand your knowledge with these resources:

- 📖 [GORM Documentation](https://gorm.io/docs/)
- 🐹 [Go Documentation](https://golang.org/doc/)
- 🐘 [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

Happy coding! 🎉
