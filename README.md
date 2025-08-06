# 🚀 3rEco NextGen API

A comprehensive enterprise management API built with Go and Fiber, featuring authentication, multi-factor authentication, role-based access control, and audit logging.

---

## 🛠️ Getting Started

### Prerequisites

- 🐹 [Go](https://golang.org/) (version 1.24.5 or later)
- 🐘 PostgreSQL database
- 🔧 PM2 (for production deployment)

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

3. Set up your PostgreSQL database and configure environment variables in `env/env.go`

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

---

## 🏗️ Project Structure

```
threereco-nextgen/
├── cmd/
│   └── api/
│       ├── main.go                     # Application entry point
│       └── http/
│           ├── http.go                 # HTTP router setup
│           ├── authentication/         # Authentication routes
│           └── middleware/             # HTTP middleware
├── env/
│   └── env.go                         # Environment configuration
├── internal/
│   ├── constants/                     # Application constants
│   ├── models/                        # Data models
│   │   ├── audit_log.go              # Audit logging model
│   │   ├── organization.go           # Organization model
│   │   ├── role.go                   # Role & permissions model
│   │   └── user.go                   # User model with MFA
│   ├── routing/                       # OpenAPI routing framework
│   │   ├── bodies/                   # Request body schemas
│   │   ├── properties/               # OpenAPI property definitions
│   │   └── schemas/                  # OpenAPI schemas
│   ├── services/                      # Business logic layer
│   │   ├── organizations/            # Organization service
│   │   ├── roles/                    # Role management service
│   │   └── users/                    # User management service
│   ├── sessions/                      # Session management
│   └── storage/                       # Database layer
├── ecosystem.config.js                # PM2 deployment configuration
├── go.mod                            # Go module definition
└── go.sum                            # Go module checksums
```

---

## � Features

### 🔐 Authentication & Security

- **Email/Password Authentication** with secure password hashing
- **Multi-Factor Authentication (MFA)** using TOTP
- **Session Management** with PostgreSQL-backed sessions
- **Role-Based Access Control** with granular permissions
- **Microsoft OAuth Integration** (configured for enterprise SSO)

### 👥 User Management

- User registration and profile management
- Organization-based user grouping
- Self-referencing user modification tracking
- Primary organization assignment

### 🏢 Organization Management

- Multi-tenant organization structure
- Domain-based organization identification
- Owner and user associations
- Organization-specific roles and permissions

### 📋 Role & Permission System

- Flexible role-based permissions
- Organization-scoped roles
- Permission inheritance and checking
- Dynamic role assignment

### 📊 Audit Logging

- Comprehensive audit trail for all data changes
- JSON-based change tracking
- User attribution for all modifications
- Automatic timestamping

### 📖 API Documentation

- **OpenAPI 3.0 Specification** with full schema definitions
- **Interactive API Documentation** via Scalar
- Auto-generated request/response schemas
- Real-time API specification at `/api/api-spec`

---

## 📦 Key Dependencies

- **[Fiber v2](https://gofiber.io/)** - Fast Express-inspired web framework
- **[GORM](https://gorm.io/)** - Go ORM with advanced features
- **[OpenAPI 3](https://github.com/getkin/kin-openapi)** - API specification and validation
- **[go-json](https://github.com/goccy/go-json)** - High-performance JSON processing
- **[UUID](https://github.com/google/uuid)** - UUID generation and parsing
- **[OTP](https://github.com/pquerna/otp)** - TOTP multi-factor authentication
- **[bcrypt](https://golang.org/x/crypto/bcrypt)** - Secure password hashing

---

## 🌐 API Endpoints

### Authentication

- `POST /api/v2/authentication/login` - User login
- `POST /api/v2/authentication/logout` - User logout
- `GET /api/v2/authentication/check` - Check authentication status
- `POST /api/v2/authentication/mfa/enable` - Enable MFA
- `POST /api/v2/authentication/mfa/verify` - Verify MFA token

### System

- `GET /api/health` - Health check endpoint
- `GET /api/api-spec` - OpenAPI specification
- `GET /api/api-doc` - Interactive API documentation

---

## 🗃️ Data Models

### User Model

- **Authentication**: Email/password with MFA support
- **Profile**: Name, phone, job title, profile image
- **Associations**: Multiple organizations, roles, and permissions
- **Security**: Encrypted password and MFA secret storage

### Organization Model

- **Identity**: Unique name and domain
- **Ownership**: Owner user with administrative rights
- **Members**: Associated users and their roles
- **Auditing**: Creation and modification tracking

### Role Model

- **Permissions**: Flexible string-based permission system
- **Associations**: Users and organizations
- **Validation**: Built-in permission checking methods

### Audit Log Model

- **Comprehensive Tracking**: All CRUD operations
- **Data Snapshots**: JSON representation of changes
- **Attribution**: User who performed the action
- **Timestamps**: Automatic creation and update times

---

## 🚀 Deployment

### Development

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:6173`

### Production with PM2

```bash
pm2 start ecosystem.config.js
```

### Environment Configuration

Configure the following in `env/env.go`:

- `POSTGRES_DSN` - PostgreSQL connection string
- `MICROSOFT_CLIENT_ID/SECRET` - OAuth credentials
- `COOKIE_DOMAIN` - Session cookie domain
- `MODE` - Application mode (development/production)

---

## 🔧 Development

### Database Setup

The application automatically:

1. Connects to PostgreSQL
2. Runs database migrations
3. Seeds initial data (admin user, roles, organization)

### Adding New Routes

1. Create route handlers in appropriate `cmd/api/http/` subdirectory
2. Define OpenAPI schemas in `internal/routing/schemas/`
3. Register routes in the router's `InitializeRoutes()` method

### Custom Middleware

Add middleware in `cmd/api/http/middleware/` and register in the router setup.

---

## 📚 Learn More

- 📖 [Fiber Documentation](https://docs.gofiber.io/)
- � [GORM Documentation](https://gorm.io/docs/)
- 🔐 [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- 🐹 [Go Documentation](https://golang.org/doc/)
- 🐘 [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

**Built with ❤️ for enterprise-grade applications** 🎉
