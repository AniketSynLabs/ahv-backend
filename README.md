# AHV Worldwide — Backend

Go + Echo + PostgreSQL API server for the AHV Worldwide website.

## Stack

- **Language:** Go 1.22
- **Framework:** Echo v4
- **Database:** PostgreSQL via pgx/v5
- **Auth:** HTTP Basic Auth on all `/api/admin/*` routes

## Getting Started

### 1. Configure environment

Edit `.env` with your database credentials:

```
PORT=8090
ADMIN_USERNAME=admin
ADMIN_PASSWORD=ahv-admin-2024

DB_HOST=your-db-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=disable
```

### 2. Run the server

```bash
go run ./cmd/server
```

On first start it auto-runs the schema migration (creates the `ahv_worldwide` schema and tables).

### 3. API Endpoints

#### Public
| Method | Path | Description |
|--------|------|-------------|
| POST | /api/leads | Submit a contact lead |
| GET | /api/settings | Get site settings |
| GET | /health | Health check |

#### Admin (Basic Auth: username + password from .env)
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/admin/leads | List all leads (optional ?status=New) |
| PUT | /api/admin/leads/:id/status | Update lead status |
| DELETE | /api/admin/leads/:id | Delete a lead |
| GET | /api/admin/settings | Get site settings |
| PUT | /api/admin/settings | Update site settings |

### Lead status values

New → Contacted → Closed

## Build for production

```bash
go build -o ahv-backend ./cmd/server
./ahv-backend
```
