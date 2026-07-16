# ArusKita

**Platform Ekosistem Logistik Kebencanaan Terpadu** — a Go REST API backend for a disaster-relief logistics platform built for GarudaHacks 7.0. It connects five actors (Admin Posko, Donatur, Toko Mitra, Relawan Kurir, and Penyintas) in a closed, auditable loop from incoming donations to verified delivery, using real-time geospatial mapping, autonomous order matching, and a QR/PIN-based chain of custody. See [docs/PRD.md](docs/PRD.md) for the full product requirements document.

The backend follows a simple 3-layer architecture (handler → service → repository) with the common pieces needed for an HTTP service: database access, authentication utilities, middleware, configuration, and standardized JSON responses.

## Tech Stack

| Package                                                       | Purpose                                                 |
| ------------------------------------------------------------- | ------------------------------------------------------- |
| [Gin](https://github.com/gin-gonic/gin)                       | HTTP web framework                                      |
| [GORM](https://gorm.io)                                       | ORM for database operations                             |
| [gorm/driver/mysql](https://github.com/go-gorm/mysql)         | MariaDB/MySQL driver                                    |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt)           | JWT authentication                                      |
| [google/uuid](https://github.com/google/uuid)                 | UUID generation                                         |
| [joho/godotenv](https://github.com/joho/godotenv)             | `.env` file loading                                     |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | Bcrypt password hashing                                 |
| `net/smtp`                                                    | Transactional email (verification codes, notifications) |

---

## Domain Model

The data model captures the full donation-to-delivery lifecycle described in the PRD (see `entity/`):

- **Users & Roles** — `User`, `Role`. Roles are seeded on startup (see [pkg/database/mariadb/seed.go](pkg/database/mariadb/seed.go)): `admin`, `donor`, `store`, and `relawan` (courier).
- **Posts & Disasters** — `Post` (posko/disaster site, with geofence radius), `DisasterEvent`, `DisasterReport` (field reports tied to a post/event).
- **Funding** — `Requests` (logistics needs with funding target/funded/reserved amounts), `Items` (structured needs list), `Wallets`, `WalletTransactions`, `Donations`.
- **Fulfillment** — `Stores` (Toko Mitra), `Orders`, `OrderItems`, `CustodyLogs` (hashed chain-of-custody handshakes), `DeliveryVerification` (forced-camera proof with GPS), `Disbursements` (payouts to stores).

Run `mariadb.Migrate` auto-migrates all of the above; `mariadb.Seed` seeds the four system roles and a default admin account.

---

## Folder Structure

```
project-root/
├── cmd/
│   └── app/
│       └── main.go               # Entry point: wires all dependencies and starts the server
├── docs/
│   └── PRD.md                    # Product requirements document
├── entity/                       # GORM models (mapped to database tables)
│   ├── user.go, role.go
│   ├── post.go, disaster_events.go, disaster_report.go
│   ├── requests.go, items.go, wallets.go, wallet_transactions.go, donations.go
│   ├── stores.go, orders.go, order_items.go
│   ├── custody_logs.go, delivery_verifications.go, disbursements.go
├── internal/                     # Core application code (3-layer architecture)
│   ├── handler/
│   │   └── rest/
│   │       ├── rest.go            # HTTP layer: route registration, server bootstrap
│   │       └── public_dashboard.go # Public dashboard endpoints
│   ├── repository/                # Data access layer: database operations via GORM
│   │   ├── repository.go          # Aggregates all repositories
│   │   ├── user.go, post.go, disaster_event.go, disaster_report.go
│   │   ├── request.go, donation.go, disbursement.go
│   │   ├── delivery_verification.go, custody_log.go
│   └── service/                   # Business logic layer
│       ├── service.go             # Aggregates all services
│       ├── user.go
│       └── public_dashboard.go    # Public transparency/map/summary aggregation
├── model/
│   ├── user.go
│   └── public_dashboard.go        # DTOs for public dashboard request/response
├── pkg/                           # Shared utilities
│   ├── bcrypt/
│   │   └── bcrypt.go              # Password hashing (bcrypt, cost=10)
│   ├── config/
│   │   ├── config.go              # Loads .env file via godotenv
│   │   └── database.go            # Builds the DSN connection string
│   ├── constant/
│   │   └── role.go                # Application-wide constants (e.g. role UUIDs)
│   ├── database/
│   │   └── mariadb/
│   │       ├── mariadb.go         # Opens the GORM database connection
│   │       ├── migrate.go         # Runs GORM AutoMigrate on startup
│   │       └── seed.go            # Seeds system roles + default admin user
│   ├── errors/
│   │   └── errors.go              # Custom AppError type with HTTP status codes
│   ├── jwt/
│   │   └── jwt.go                 # JWT creation and validation
│   ├── mail/
│   │   ├── mail.go                # SMTP sending + verification code generation
│   │   ├── template.go            # Email templates
│   │   └── verification.go        # Verification code flow
│   ├── middleware/
│   │   ├── middleware.go          # Middleware container (auth guards, etc.)
│   │   ├── authentication.go      # JWT-based auth guard
│   │   └── cors.go                # CORS handling
│   └── response/
│       └── response.go            # Standardized JSON response envelope
├── .env.example                   # Environment variable template
├── go.mod
└── go.sum
```

---

## Architecture

The core application code lives in the `internal/` directory, which keeps a strict 3-layer separation of concerns.

```
HTTP Request
     │
     ▼
┌─────────────────────┐
│   Handler / REST    │  ← Receives requests, validates input, returns responses
│  internal/handler/  │
└────────┬────────────┘
         │ calls
         ▼
┌─────────────────────┐
│      Service        │  ← Business logic, orchestrates data flow
│  internal/service/  │
└────────┬────────────┘
         │ calls
         ▼
┌─────────────────────┐
│     Repository      │  ← Database operations (GORM queries)
│ internal/repository/│
└────────┬────────────┘
         │
         ▼
      Database
    (MariaDB / MySQL)
```

### Layer Responsibilities

**`internal/repository/`**
Owns all direct database interaction. Each method corresponds to a specific query or mutation. Receives a `*gorm.DB` instance and is the only layer allowed to call GORM methods. Aggregated in `repository.go` as `Repository` (currently: `UserRepository`, `PostRepository`, `DisasterReportRepository`, `DisasterEventRepository`, `RequestRepository`, `DeliveryVerificationRepository`, `DonationRepository`, `DisbursementRepository`, `CustodyLogRepository`).

**`internal/service/`**
Contains business logic. Depends on the repository for data access and on `pkg/bcrypt`, `pkg/jwt`, and `pkg/mail` for cross-cutting concerns. Never imports GORM directly. Aggregated in `service.go` as `Service` (currently: `UserService`, `PublicDashboardService`).

**`internal/handler/rest/`**
The outermost layer. Binds HTTP routes via Gin, parses request bodies, calls the service, and writes back JSON responses using the shared `pkg/response` formatter.

---

## API Endpoints

All routes are namespaced under `/api/v1`. Currently implemented:

| Method | Path                       | Description                                                                                |
| ------ | -------------------------- | ------------------------------------------------------------------------------------------ |
| GET    | `/dashboard/summary`       | Aggregated funding summary per posko (target/funded/urgency)                               |
| GET    | `/dashboard/map`           | Posko locations + funding status for the interactive disaster map                          |
| GET    | `/dashboard/distributions` | Verified delivery proofs (photo, GPS, custody hash) per order                              |
| GET    | `/dashboard/transparency`  | Public transparency page: totals, monthly disbursements, ledger, verified fulfillment rate |

Authentication (JWT + bcrypt), email verification (`pkg/mail`), and CORS/auth middleware are already wired in `cmd/app/main.go`, ready for additional route groups (donation, order matching, custody, disbursement, etc.) to be layered in as those modules are built out per the PRD roadmap.

---

## Dependency Injection

All dependencies are wired manually in `cmd/app/main.go` using constructor injection — no DI framework required.

```
main()
  ├── config.LoadEnvironment()         # Load .env
  ├── mariadb.ConnectDatabase()        # Open *gorm.DB
  ├── mariadb.Migrate()                # Auto-migrate tables
  ├── mariadb.Seed()                   # Seed system roles + default admin
  │
  ├── repository.NewRepository(db)     # Data layer
  ├── bcrypt.Init()                    # Password util
  ├── jwt.Init()                       # Auth util
  ├── service.NewService(repo, bcrypt, jwt)   # Business logic
  ├── middleware.Init(service, jwt)    # Middleware chain
  └── rest.NewRest(service, middleware)
        ├── rest.MountEndpoint()       # Register routes
        └── rest.Run()                 # Start server
```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values before running.

| Variable         | Description                               | Example               |
| ---------------- | ----------------------------------------- | --------------------- |
| `DB_HOST`        | Database host                             | `localhost`           |
| `DB_PORT`        | Database port                             | `3306`                |
| `DB_NAME`        | Database name                             | `myapp`               |
| `DB_USER`        | Database user                             | `root`                |
| `DB_PASSWORD`    | Database password                         | `secret`              |
| `ADDRESS`        | Server bind address                       | `localhost`           |
| `PORT`           | Server port                               | `8080`                |
| `TIME_OUT_LIMIT` | Request timeout (seconds)                 | `10`                  |
| `JWT_SECRET_KEY` | Secret key for signing JWTs (min 256-bit) | `a-very-long-secret`  |
| `JWT_EXP_TIME`   | JWT expiration in hours                   | `1`                   |
| `SMTP_HOST`      | SMTP server host for outgoing email       | `smtp.gmail.com`      |
| `SMTP_PORT`      | SMTP server port                          | `587`                 |
| `SMTP_USERNAME`  | SMTP account used to send email           | `youremail@gmail.com` |
| `SMTP_PASSWORD`  | SMTP account password/app password        | `yourpassword`        |

---

## Getting Started

1. **Clone the project**

   ```bash
   git clone <repository-url> <project-name>
   cd <project-name>
   ```

2. **Set up environment**

   ```bash
   cp .env.example .env
   # Edit .env with your database, JWT, and SMTP credentials
   ```

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

4. **Run**

   ```bash
   go run cmd/app/main.go
   ```

   ```bash
   air
   ```

   On startup, the app connects to MariaDB, runs `AutoMigrate` for all entities, and seeds the `admin`, `donor`, `store`, and `relawan` roles plus a default admin user (`admin@example.com` / `admin123` — change this before any real deployment).

5. **Continue building**
   - Define new domain models in `entity/` and register them in `pkg/database/mariadb/migrate.go`
   - Add request/response structs in `model/`
   - Implement repository methods in `internal/repository/`, then wire them into `Repository` in `repository.go`
   - Implement business logic in `internal/service/`, then wire it into `Service` in `service.go`
   - Register routes and handlers in `internal/handler/rest/`
