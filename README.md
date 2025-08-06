# 🚀 3rEco NextGen Platform

A full-stack enterprise management platform featuring:

- **Backend:** Go + Fiber REST API (authentication, MFA, RBAC, audit logging, organizations, users, roles)
- **Frontend:** Modern React (Vite, TanStack Router, Tailwind CSS) web app for user, organization, and role management

---

## 🛠️ Getting Started

### Prerequisites

- [Go](https://golang.org/) >= 1.24.5
- PostgreSQL database
- Node.js & Yarn/NPM (for frontend)
- PM2 (optional, for production deployment)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd nextgen-threereco
   ```
2. **Backend setup:**
   ```bash
   go mod download
   go run cmd/api/main.go
   ```
   The API runs at `http://localhost:6173`.
3. **Frontend setup:**
   ```bash
   cd frontend
   yarn install # or npm install
   yarn build   # or npm run build
   npx serve dist --port 5177
   ```
   The frontend runs at `http://localhost:5177` and connects to the API.
4. **Environment:**
   - Configure database/API keys in `env/env.go` (backend)
   - Set `VITE_API_URL` in `.env` (frontend) if needed

---

## 🏗️ Project Structure

```
nextgen-threereco/
├── cmd/api/                  # Backend entrypoint & HTTP routing
│   ├── main.go               # API server setup
│   └── http/                 # HTTP routers & middleware
│       ├── authentication/   # Auth endpoints (login, logout, MFA)
│       └── middleware/       # Auth/session middleware
├── env/                      # Environment config (env.go)
├── internal/
│   ├── constants/            # Error/status constants
│   ├── models/               # Data models (User, Organization, Role, AuditLog)
│   ├── routing/              # OpenAPI schemas, route definitions
│   ├── services/             # Business logic (users, roles, orgs)
│   ├── sessions/             # Session management
│   └── storage/              # Database connection/migrations
├── frontend/
│   ├── src/                  # React app source code
│   │   ├── components/ui/    # UI components (Table, Select, Sheet, etc.)
│   │   ├── hooks/            # Custom React hooks
│   │   ├── lib/              # Utilities (permissions, API client)
│   │   ├── providers/        # Context providers (auth, theme)
│   │   ├── routes/           # Route definitions (TanStack Router)
│   │   └── styles.css        # Tailwind & custom styles
│   ├── dist/                 # Production build output
│   ├── vite.config.ts        # Vite config
│   └── tsconfig.json         # TypeScript config
├── ecosystem.config.js       # PM2 deployment config
├── go.mod                    # Go module definition
└── go.sum                    # Go module checksums
```

---

## ✨ Features

### 🔐 Authentication & Security

- Email/password login (bcrypt)
- Multi-factor authentication (MFA, TOTP)
- Session management (PostgreSQL-backed)
- Role-based access control (RBAC)
- Microsoft OAuth SSO (enterprise)

### 👥 User Management

- Registration, profile, and password management
- Organization-based user grouping
- Self-referencing modification tracking
- Primary organization assignment

### 🏢 Organization Management

- Multi-tenant structure
- Domain-based identification
- Owner/user/role associations

### 📋 Role & Permission System

- Flexible, string-based permissions
- Organization-scoped roles
- Permission inheritance/checking
- Dynamic assignment

### 📊 Audit Logging

- Tracks all CRUD operations
- JSON data snapshots
- User attribution
- Automatic timestamps

### 📖 API Documentation

- OpenAPI 3.0 spec (auto-generated)
- Interactive docs via Scalar
- Real-time spec at `/api/api-spec`

---

## 📦 Key Dependencies

**Backend:**

- [Fiber v2](https://gofiber.io/) - Web framework
- [GORM](https://gorm.io/) - ORM
- [OpenAPI 3](https://github.com/getkin/kin-openapi) - API spec
- [go-json](https://github.com/goccy/go-json) - Fast JSON
- [UUID](https://github.com/google/uuid) - UUIDs
- [OTP](https://github.com/pquerna/otp) - MFA
- [bcrypt](https://golang.org/x/crypto/bcrypt) - Password hashing

**Frontend:**

- [React](https://react.dev/) + [Vite](https://vitejs.dev/) - SPA
- [TanStack Router](https://tanstack.com/router) - Routing
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [React Query](https://tanstack.com/query) - Data fetching

---

## 🌐 API Endpoints

### Authentication

- `POST /api/v2/authentication/login` — Login
- `POST /api/v2/authentication/logout` — Logout
- `GET /api/v2/authentication/check` — Check session
- `POST /api/v2/authentication/mfa/enable` — Enable MFA
- `POST /api/v2/authentication/mfa/verify` — Verify MFA

### System

- `GET /api/health` — Health check
- `GET /api/api-spec` — OpenAPI spec
- `GET /api/api-doc` — Interactive docs

---

## 🗃️ Data Models

### User

- Email/password (bcrypt)
- MFA secret & status
- Profile: name, phone, job title, image
- Roles (many-to-many)
- Organizations (many-to-many)
- Primary organization
- ModifiedBy (self-referencing)
- Created/updated timestamps

### Organization

- Name, domain (unique)
- Owner (user)
- Users, roles (many-to-many)
- ModifiedBy (user)
- Created/updated timestamps

### Role

- Name, description
- Permissions (string[])
- Users, organizations (many-to-many)
- ModifiedBy (user)
- Created/updated timestamps

### AuditLog

- Table name, operation type
- Object ID, data (JSON)
- User (who performed action)
- Created/updated timestamps

---

## 🚀 Deployment

### Development

- Backend: `go run cmd/api/main.go` (http://localhost:6173)
- Frontend: `npx serve dist --port 5177` (http://localhost:5177)

### Production

- Use `pm2 start ecosystem.config.js` to run both backend and serve frontend from `frontend/dist`.

### Environment Configuration

- Backend: `env/env.go` (DB, OAuth, cookie, mode)
- Frontend: `.env` (VITE_API_URL)

---

## 🖥️ Frontend

- Built with React, Vite, TanStack Router, Tailwind CSS
- UI components for tables, forms, modals, sheets, sidebar, etc.
- Authentication context/provider for user session
- Theme provider (light/dark/system)
- Permission checks via `lib/permissions.ts`
- API client auto-configured for backend URL
- Route-based code splitting and navigation
- Responsive/mobile support via custom hooks
- Entry: `frontend/src/main.tsx`, routes in `frontend/src/routes/`
- Production build: `frontend/dist/`

---

## 🔧 Development Workflow

### Backend

- Auto-connects to PostgreSQL, runs migrations, seeds initial data
- Add new routes: create handler in `cmd/api/http/`, define OpenAPI schema, register in router
- Add custom middleware in `cmd/api/http/middleware/`

### Frontend

- Add new pages/routes in `frontend/src/routes/`
- Add UI components in `frontend/src/components/ui/`
- Use context providers for authentication, theme, etc.
- Use hooks/utilities for permissions, API calls, mobile detection

---

## 📚 Learn More

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [React Documentation](https://react.dev/)
- [TanStack Router](https://tanstack.com/router)
- [Tailwind CSS](https://tailwindcss.com/)

---

**Built with ❤️ for enterprise-grade applications** 🎉
