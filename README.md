# ArusKita

**Platform Ekosistem Logistik Kebencanaan Terpadu** is a Go REST API backend for disaster-relief logistics, built for GarudaHacks 7.0. ArusKita connects Admin Posko, Donatur, Toko Mitra, Relawan Kurir, and Penyintas in a closed, auditable flow from incoming donations to verified last-mile delivery.

The system focuses on transparent disaster response: public map visibility, donation allocation, verified distribution proofs, QR/PIN-based custody handshakes, GPS-backed delivery verification, store disbursement tracking, and donor rewards. See [docs/PRD.md](docs/PRD.md) for the full product requirements document.

The backend follows a 3-layer architecture (`handler -> service -> repository`) with database access, authentication, middleware, configuration, payment integration, file upload support, and standardized JSON responses.

## Built With (Required)

A set of tags this project was developed with. This list covers the main languages, frameworks, databases, libraries, APIs, services, and deployment tools used in the project.

- **Language:** Go 1.25
- **Backend Framework:** Gin
- **Database:** MariaDB / MySQL
- **ORM:** GORM with `gorm.io/driver/mysql`
- **Authentication:** JWT, bcrypt, token revocation table for logout
- **Payment API:** Midtrans Core API / Snap integration
- **Storage API:** Supabase Storage
- **Email:** SMTP via Go `net/smtp`
- **Image Processing:** WebP support via `github.com/chai2010/webp`
- **Configuration:** `.env` loading with `godotenv`
- **Deployment:** Docker, Docker Compose, Alpine Linux runtime
- **Architecture:** REST API, repository-service-handler layering

## Tools Used

List of all technologies used to build the project.

| Tool / Technology | Usage |
| --- | --- |
| Go | Main backend programming language |
| Gin | HTTP routing, REST handlers, middleware integration |
| GORM | Database ORM and AutoMigrate |
| MariaDB 11.4 / MySQL | Relational database for users, reports, requests, payments, orders, custody logs, and disbursements |
| Docker | Containerized application build and runtime |
| Docker Compose | Local/app deployment with API service and MariaDB service |
| JWT (`github.com/golang-jwt/jwt/v5`) | Access token creation and validation |
| bcrypt (`golang.org/x/crypto`) | Password hashing |
| Midtrans Go SDK | Donation payment charge creation and notification handling |
| Supabase Storage Go SDK | Uploaded media/object storage support |
| SMTP / Gmail-compatible SMTP | OTP and verification email delivery |
| `github.com/google/uuid` | UUID primary keys and entity identifiers |
| `github.com/joho/godotenv` | Environment variable loading for local development |
| `github.com/chai2010/webp` | WebP image handling |
| GitHub Container Registry | Container image target used by `docker-compose.yml` |
| Alpine Linux | Minimal production image runtime |
| Makefile | Local developer shortcuts |

## Copyright Materials (Required)

**What is it?**

Please declare any third-party assets or copyrighted materials used in your project, such as icons, illustrations, images, datasets, music, etc.

This repository is primarily backend source code created by the team. No copyrighted images, music, illustrations, icons, or external datasets are bundled in the repository. Dummy image URLs used in seed/demo SQL point to generated placeholder images from `dummyimage.com` and are used only for development/testing examples. Third-party open-source Go packages are listed in [go.mod](go.mod) and governed by their respective licenses.

## Key Features

- Public disaster map with posko coordinates, disaster type, funding target, funded amount, and urgency level.
- Public transparency dashboard with total collected donations, verified disbursements, monthly disbursement breakdown, disaster allocation, and custody ledger.
- Donor registration, login, profile, donation payments, transaction history, and points/rewards.
- Admin posko onboarding, disaster event/report creation, profile metrics, and post handoff verification.
- Store order management, order acceptance, readiness status, disbursement dashboard, and goodness trail.
- Courier task claiming, location updates, pickup/delivery milestones, and custody handoff flow.
- Chain-of-custody logs with QR/PIN handshake support and audit hashes.
- JWT logout using server-side token revocation.

---

## Domain Model

The data model captures the full donation-to-delivery lifecycle described in the PRD (see `entity/`):

- **Users & Roles** — `User`, `Role`. Roles are seeded on startup (see [pkg/database/mariadb/seed.go](pkg/database/mariadb/seed.go)): `admin`, `donor`, `store`, and `relawan` (courier).
- **Profiles & Onboarding** — `AdminProfile` (NIK, affiliation), `DonorProfile` (phone number, preferences, consent), `RegistrationSession` (email/OTP-based signup flow for both admin and donor), `OtpCode`.
- **Posts & Disasters** — `Post` (posko/disaster site, with geofence radius), `DisasterEvent`, `DisasterReport` (field reports tied to a post/event).
- **Funding** — `Requests` (logistics needs with funding target/funded/reserved amounts), `Items` (structured needs list), `Wallets`, `WalletTransactions`, `Donations`, `PaymentTransactions` (Midtrans charge/notification records: QR/VA details, status, raw payloads).
- **Rewards & Gamification** — `PointAccount` (active/earned/redeemed totals per donor), `PointTransaction` (earn/redeem/adjustment ledger), `Reward`, `RewardClaim` (redemption of points for pulsa/voucher/donation rewards).
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
│   ├── user.go, role.go, otp.go, registration_session.go, revoked_token.go
│   ├── admin_profile.go, donor_profile.go
│   ├── post.go, disaster_events.go, disaster_report.go
│   ├── requests.go, items.go, wallets.go, wallet_transactions.go, donations.go
│   ├── payment_transaction.go, point.go, custody_handshake_tokens.go
│   ├── stores.go, orders.go, order_items.go
│   ├── custody_logs.go, delivery_verifications.go, disbursements.go
├── internal/                     # Core application code (3-layer architecture)
│   ├── handler/
│   │   └── rest/
│   │       ├── rest.go             # HTTP layer: route registration, server bootstrap
│   │       ├── auth.go             # Login, logout, registration + OTP flow
│   │       ├── admin_profile.go, admin_dashboard.go, admin_event.go
│   │       ├── donor_profile.go, donor_dashboard.go, donor_transaction.go
│   │       ├── donation_payment.go # Donation checkout (Midtrans) + webhook
│   │       ├── point.go            # Point dashboard, history, rewards, claims
│   │       ├── store_custody.go, store_disbursement.go, store_goodness.go
│   │       ├── courier_task.go, courier_goodness.go
│   │       ├── admin_custody.go
│   │       └── public_dashboard.go # Public dashboard endpoints
│   ├── repository/                # Data access layer: database operations via GORM
│   │   ├── repository.go          # Aggregates all repositories
│   │   ├── user.go, role.go, registration.go, otp.go
│   │   ├── post.go, disaster_event.go, disaster_report.go
│   │   ├── request.go, item.go, donation.go, wallet.go, wallet_transaction.go
│   │   ├── payment_transaction.go, point.go, revoked_token.go
│   │   ├── admin_profile.go, admin_dashboard.go, donor_profile.go
│   │   ├── order.go, order_item.go, disbursement.go
│   │   ├── delivery_verification.go, custody_log.go, custody_handshake_token.go
│   └── service/                   # Business logic layer
│       ├── service.go             # Aggregates all services
│       ├── auth.go, otp.go, user.go
│       ├── admin_profile.go, admin_dashboard.go, admin_event.go
│       ├── donor_profile.go, donor_dashboard.go, donor_transaction.go
│       ├── donation_payment.go     # Midtrans integration (charge + webhook handling)
│       ├── point.go                # Point accrual, history, rewards, claims
│       ├── store_custody.go, store_disbursement.go, store_goodness.go
│       ├── courier_task.go, courier_goodness.go, admin_custody.go
│       └── public_dashboard.go     # Public transparency/map/summary aggregation
├── model/                        # Request/response DTOs, one file per domain area
│   ├── auth.go, otp.go, user.go
│   ├── admin_profile.go, admin_dashboard.go, admin_event.go
│   ├── donor_profile.go, donor_dashboard.go, donor_transaction.go
│   ├── donation_payment.go, payment.go, point.go
│   ├── store_custody.go, store_disbursement.go, store_goodness.go
│   ├── courier_task.go, courier_goodness.go
│   └── public_dashboard.go
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
Owns all direct database interaction. Each method corresponds to a specific query or mutation. Receives a `*gorm.DB` instance and is the only layer allowed to call GORM methods. Aggregated in `repository.go` as `Repository` (including user/role/registration, public dashboard, disaster reports, requests/items, wallet/donation/payment, points/rewards, store/courier/order, disbursement, delivery verification, custody logs/tokens, and revoked token repositories).

**`internal/service/`**
Contains business logic. Depends on the repository for data access and on utilities such as `pkg/bcrypt`, `pkg/jwt`, `pkg/mail`, `pkg/hash`, Supabase, and Midtrans for cross-cutting concerns. Aggregated in `service.go` as `Service` across auth, public dashboard, admin, donor, store, courier, payment, points, custody, and disbursement domains.

**`internal/handler/rest/`**
The outermost layer. Binds HTTP routes via Gin, parses request bodies, calls the service, and writes back JSON responses using the shared `pkg/response` formatter.

---

## API Endpoints

All routes are namespaced under `/api/v1`. Currently implemented:

### Public

| Method | Path                       | Description                                                                                |
| ------ | -------------------------- | ------------------------------------------------------------------------------------------ |
| GET    | `/dashboard/summary`       | Aggregated funding summary per posko (target/funded/urgency)                               |
| GET    | `/dashboard/map`           | Posko locations + funding status for the interactive disaster map                          |
| GET    | `/dashboard/distributions` | Verified delivery proofs (photo, GPS, custody hash) per order                              |
| GET    | `/dashboard/transparency`  | Public transparency page: totals, monthly disbursements, ledger, verified fulfillment rate |
| POST   | `/payments/webhook`        | Midtrans payment notification webhook (updates donation/payment/point status)              |

### Auth & Registration

| Method | Path                               | Description                                                   |
| ------ | ---------------------------------- | ------------------------------------------------------------- |
| POST   | `/auth/login`                      | Login with email/password, returns JWT                        |
| POST   | `/auth/logout`                     | Logout current JWT by revoking the token server-side          |
| POST   | `/auth/register/request-otp`       | Start registration, send OTP to email                         |
| POST   | `/auth/register/verify-otp`        | Verify OTP for a pending registration session                 |
| POST   | `/auth/register/password`          | Set password for a verified registration session              |
| POST   | `/auth/register/admin/request-otp` | Start admin registration, send OTP to email                   |
| POST   | `/auth/register/admin/verify-otp`  | Verify OTP for a pending admin registration                   |
| POST   | `/auth/register/admin/password`    | Set password for a verified admin registration                |
| POST   | `/auth/register/admin/profile`     | Complete admin registration (NIK, affiliation, posko details) |
| POST   | `/auth/register/donor/profile`     | Complete donor registration (phone number, preferences)       |
| POST   | `/auth/register/store/profile`     | Complete store registration and store profile setup           |

### Admin (JWT + `admin` role required)

| Method | Path               | Description                                                 |
| ------ | ------------------ | ----------------------------------------------------------- |
| GET    | `/admin/dashboard` | Admin home dashboard: metrics/summary for the admin's posko |
| GET    | `/admin/profile`   | Admin profile (NIK, affiliation, aggregated report metrics) |
| POST   | `/admin/events`    | Create a disaster event/report tied to the admin's posko    |
| POST   | `/admin/custody/post-handoff` | Verify courier-to-posko handoff with QR/PIN custody data |

### Donor (JWT + `donor` role required)

| Method | Path                                         | Description                                                        |
| ------ | -------------------------------------------- | ------------------------------------------------------------------ |
| GET    | `/donor/profile`                             | Donor profile: verification status, level, lifetime donation stats |
| GET    | `/donor/dashboard/map`                       | Posko map view scoped to the donor experience                      |
| GET    | `/donor/dashboard/posts/:post_id`            | Detail of a single posko/request for donors                        |
| GET    | `/donor/donations/transactions`              | Donor's donation transaction history                               |
| GET    | `/donor/donations/transactions/:donation_id` | Detail of a single donation transaction                            |
| POST   | `/donor/donations/payments`                  | Create a donation payment (Midtrans charge: QR/VA)                 |
| GET    | `/donor/points`                              | Donor point dashboard (active/earned/redeemed totals)              |
| GET    | `/donor/points/history`                      | Point transaction history (earn/redeem/adjustment ledger)          |
| GET    | `/donor/points/rewards`                      | Browse claimable rewards (pulsa, voucher, donation)                |
| POST   | `/donor/points/rewards/claim`                | Claim a reward using accumulated points                            |

### Store (JWT + `store` role required)

| Method | Path                            | Description                                      |
| ------ | ------------------------------- | ------------------------------------------------ |
| GET    | `/store/profile`                | Store profile and verification data             |
| GET    | `/store/orders`                 | Store order list                                 |
| GET    | `/store/orders/:order_id`       | Detail of a store order                          |
| POST   | `/store/orders/:order_id/accept` | Accept an available order                       |
| POST   | `/store/orders/:order_id/ready` | Mark order ready for pickup                      |
| POST   | `/store/orders/:order_id/handoff-token` | Generate handoff token for courier pickup |
| GET    | `/store/disbursements/dashboard` | Store disbursement dashboard                    |
| GET    | `/store/goodness`               | Store goodness/contribution trail                |

### Courier (JWT + `relawan` role required)

| Method | Path                                  | Description                                      |
| ------ | ------------------------------------- | ------------------------------------------------ |
| GET    | `/courier/tasks`                      | Courier task list                                |
| GET    | `/courier/tasks/:order_id`            | Courier task detail                              |
| POST   | `/courier/tasks/:order_id/claim`      | Claim a delivery task                            |
| POST   | `/courier/tasks/:order_id/location`   | Update courier GPS location                      |
| POST   | `/courier/tasks/:order_id/arrived`    | Mark courier arrived at store                    |
| POST   | `/courier/tasks/:order_id/arrived-post` | Mark courier arrived at posko                  |
| POST   | `/courier/tasks/:order_id/handoff-token` | Generate handoff token for posko delivery    |
| POST   | `/courier/custody/store-handoff`      | Submit store-to-courier handoff proof            |
| GET    | `/courier/goodness`                   | Courier goodness/contribution trail              |

Authentication (JWT + bcrypt), server-side token revocation for logout, OTP-based email verification (`pkg/mail`), Midtrans payment integration, Supabase storage, CORS, and role-based auth middleware are wired in `cmd/app/main.go`.

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
  ├── supabase.Init()                  # Object storage client
  ├── config.LoadMidtransConfig()       # Payment gateway config
  ├── hash.Init()                      # NIK hashing util
  ├── service.NewService(...)          # Business logic
  ├── middleware.Init(service, jwt)    # Middleware chain
  └── rest.NewRest(service, middleware)
        ├── rest.MountEndpoint()       # Register routes
        └── rest.Run()                 # Start server
```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values before running.

| Variable | Description | Example |
| --- | --- | --- |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_NAME` | Database name | `garudahacks` |
| `DB_USER` | Database user | `garudahacks_user` |
| `DB_PASSWORD` | Database password | `secret` |
| `ADDRESS` | Server bind address | `localhost` |
| `PORT` | Server port | `8080` |
| `TIME_OUT_LIMIT` | Request timeout in seconds | `10` |
| `JWT_SECRET_KEY` | Secret key for signing JWTs | `a-string-secret-at-least-256-bits-long` |
| `JWT_EXP_TIME` | JWT expiration in hours | `1` |
| `SMTP_HOST` | SMTP server host for outgoing email | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USERNAME` | SMTP account used to send email | `youremail@gmail.com` |
| `SMTP_PASSWORD` | SMTP account password/app password | `yourpassword` |
| `SUPABASE_URL` | Supabase project URL | `https://yoururl.supabase.co` |
| `SUPABASE_TOKEN` | Supabase service/access token | `your-token` |
| `SUPABASE_BUCKET` | Supabase Storage bucket name | `your-bucket` |
| `MIDTRANS_CLIENT_KEY` | Midtrans client key | `your-client-key` |
| `MIDTRANS_SERVER_KEY` | Midtrans server key | `your-server-key` |
| `MIDTRANS_ENVIRONMENT` | Midtrans environment | `sandbox` |
| `NIK_HASH_SECRET` | Secret used for hashing sensitive identity data | `your-secret` |

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
   # Edit .env with your database, JWT, SMTP, Supabase, and Midtrans credentials
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

   To run with Docker Compose:

   ```bash
   docker compose up -d
   ```

5. **Continue building**
   - Define new domain models in `entity/` and register them in `pkg/database/mariadb/migrate.go`
   - Add request/response structs in `model/`
   - Implement repository methods in `internal/repository/`, then wire them into `Repository` in `repository.go`
   - Implement business logic in `internal/service/`, then wire it into `Service` in `service.go`
   - Register routes and handlers in `internal/handler/rest/`
