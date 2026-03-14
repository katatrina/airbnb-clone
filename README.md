# Airbnb Clone

> A simplified Airbnb platform connecting people who have rooms/apartments for short-term rental (Hosts) with tourists who need accommodation (Guests). Built as a microservices architecture in Go.

## Tech Stack

| Category       | Technology                                              |
|----------------|---------------------------------------------------------|
| Language       | Go 1.25.5                                               |
| Framework      | [Gin](https://github.com/gin-gonic/gin)                |
| Database       | PostgreSQL 16                                           |
| Cache          | Redis 7                                                 |
| Authentication | JWT (HS256) via [golang-jwt](https://github.com/golang-jwt/jwt) |
| Config         | [Viper](https://github.com/spf13/viper)                |
| Validation     | [go-playground/validator](https://github.com/go-playground/validator) |
| DB Driver      | [pgx](https://github.com/jackc/pgx) v5                 |
| Migrations     | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Linting        | [golangci-lint](https://golangci-lint.run/)             |
| DevOps         | Docker, Docker Compose                                  |

## Architecture Overview

The project follows a **microservices architecture** with three independent services sharing a single PostgreSQL database. Each service is organized using **Clean Architecture** (Handler → Service → Repository).

```
airbnb-clone/
├── docker-compose.yml          # PostgreSQL + Redis
├── Makefile                    # Infrastructure commands
├── go.work                     # Go workspace (links all services)
├── pkg/                        # Shared packages
│   ├── middleware/              #   Auth middleware (JWT validation)
│   ├── token/                  #   JWT creation & verification
│   ├── request/                #   Validation & pagination helpers
│   └── response/               #   Standardized API responses
└── services/
    ├── user/                   # Port 8081 — Auth & profiles
    ├── listing/                # Port 8082 — Property listings & locations
    └── booking/                # Port 8083 — Booking management
```

Each service follows the same internal structure:

```
service/
├── cmd/api/main.go             # Entry point
├── config/                     # Configuration loading (.env)
├── internal/
│   ├── handler/                # HTTP handlers (routes)
│   ├── service/                # Business logic
│   ├── repository/             # Database access (pgx)
│   └── model/                  # Domain models
└── migrations/                 # SQL migration files
```

### Service Communication

```
┌──────────┐     ┌───────────┐     ┌───────────┐
│   User   │     │  Listing  │     │  Booking  │
│  :8081   │     │   :8082   │◄────│   :8083   │
└────┬─────┘     └─────┬─────┘     └─────┬─────┘
     │                 │                  │
     └────────────┬────┴──────────────────┘
                  ▼
            ┌──────────┐
            │ PostgreSQL│
            │  :5432    │
            └──────────┘
```

The **Booking** service calls the **Listing** service over HTTP to verify listing existence and retrieve pricing information when creating a booking.

## Getting Started

### Prerequisites

- Go 1.25.5+
- Docker & Docker Compose
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### 1. Clone the repository

```bash
git clone https://github.com/katatrina/airbnb-clone.git
cd airbnb-clone
```

### 2. Start infrastructure

```bash
make docker-up    # Starts PostgreSQL and Redis containers
```

### 3. Configure environment

Each service requires a `.env` file. Copy the examples:

```bash
cp services/user/.env.example services/user/.env
cp services/listing/.env.example services/listing/.env
cp services/booking/.env.example services/booking/.env
```

### 4. Run database migrations

```bash
cd services/user && make migrate-up && cd ../..
cd services/listing && make migrate-up && cd ../..
cd services/booking && make migrate-up && cd ../..
```

### 5. Load location data (optional)

The listing service uses Vietnamese administrative divisions (provinces, districts, wards). Load them with:

```bash
cd services/listing && go run ./cmd/import-locations && cd ../..
```

### 6. Run services

Run each service in a separate terminal:

```bash
# Terminal 1 — User Service (port 8081)
cd services/user && make server

# Terminal 2 — Listing Service (port 8082)
cd services/listing && make server

# Terminal 3 — Booking Service (port 8083)
cd services/booking && make server
```

### Useful commands

| Command              | Description                          |
|----------------------|--------------------------------------|
| `make docker-up`     | Start PostgreSQL and Redis           |
| `make docker-down`   | Stop all containers                  |
| `make migrate-up`    | Apply pending migrations (per service) |
| `make migrate-down`  | Rollback migrations (per service)    |
| `make server`        | Run the service (per service)        |
| `make lint`          | Run golangci-lint (per service)      |
| `make service-test`  | Run unit tests (user service)        |

## API Endpoints

All endpoints return a standardized JSON response:

```json
{
  "success": true,
  "code": "OK",
  "message": "...",
  "data": { },
  "meta": {
    "requestId": "uuid",
    "timestamp": 1234567890,
    "pagination": { "page": 1, "pageSize": 10, "total": 100, "totalPages": 10 }
  }
}
```

Protected endpoints require a `Authorization: Bearer <token>` header.

### Health Check

| Method | Endpoint  | Description          |
|--------|-----------|----------------------|
| GET    | `/health` | Available on all services |

### User Service `:8081`

| Method | Endpoint                | Auth | Description              |
|--------|-------------------------|------|--------------------------|
| POST   | `/api/v1/auth/register` | No   | Register a new user      |
| POST   | `/api/v1/auth/login`    | No   | Login and receive JWT    |
| GET    | `/api/v1/me/profile`    | Yes  | Get authenticated user's profile |

### Listing Service `:8082`

**Public**

| Method | Endpoint                              | Description                     |
|--------|---------------------------------------|---------------------------------|
| GET    | `/api/v1/listings`                    | List all active listings (paginated) |
| GET    | `/api/v1/listings/:id`                | Get a single listing            |
| GET    | `/api/v1/provinces`                   | List all provinces              |
| GET    | `/api/v1/provinces/:code/districts`   | List districts by province code |
| GET    | `/api/v1/districts/:code/wards`       | List wards by district code     |

**Protected (Host)**

| Method | Endpoint                                    | Description                |
|--------|---------------------------------------------|----------------------------|
| POST   | `/api/v1/me/listings`                       | Create a new listing       |
| GET    | `/api/v1/me/listings`                       | List host's own listings   |
| GET    | `/api/v1/me/listings/:id`                   | Get host's listing details |
| PATCH  | `/api/v1/me/listings/:id/basic-info`        | Update title, description, price |
| PATCH  | `/api/v1/me/listings/:id/address`           | Update listing address     |
| DELETE | `/api/v1/me/listings/:id`                   | Soft-delete a listing      |
| POST   | `/api/v1/me/listings/:id/publish`           | Publish listing (draft → active) |
| POST   | `/api/v1/me/listings/:id/deactivate`        | Deactivate listing         |
| POST   | `/api/v1/me/listings/:id/reactivate`        | Reactivate listing         |

### Booking Service `:8083`

All booking endpoints require authentication.

**Guest**

| Method | Endpoint                             | Description              |
|--------|--------------------------------------|--------------------------|
| POST   | `/api/v1/me/bookings`                | Create a booking         |
| GET    | `/api/v1/me/bookings`                | List guest's bookings    |
| GET    | `/api/v1/me/bookings/:id`            | Get booking details      |
| POST   | `/api/v1/me/bookings/:id/cancel`     | Cancel a booking         |

**Host**

| Method | Endpoint                                    | Description              |
|--------|---------------------------------------------|--------------------------|
| GET    | `/api/v1/me/hosting/bookings`               | List host's bookings     |
| POST   | `/api/v1/me/hosting/bookings/:id/confirm`   | Confirm a booking        |
| POST   | `/api/v1/me/hosting/bookings/:id/reject`    | Reject a booking         |
