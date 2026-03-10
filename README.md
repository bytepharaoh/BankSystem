# SimpleBank (Go + PostgreSQL)

SimpleBank is a backend training project focused on building a clean, testable banking API in Go.
It covers account management, user authentication, money transfers, transactional integrity, and ownership-based authorization.

## Project Highlights

- Deployed with CI/CD workflow using GitHub Actions
- Automated developer setup with Bash scripts (`scripts/`)
- Containerized application runtime with Docker and Docker Compose
- Strong test coverage across API and database layers

## Tech Stack

- Go (Gin for HTTP server)
- PostgreSQL
- SQLC (type-safe query generation)
- golang-migrate (schema migrations)
- JWT authentication
- GoMock + Testify (unit tests)

## Deployment

This project is configured and maintained with GitHub Actions workflows under `.github/`.
The pipeline is used to validate and deploy the application as part of the training delivery flow.

## What This Project Implements

- User registration and login
- JWT access token issuance
- Authentication middleware for protected routes
- Authorization rules based on token identity
- Account CRUD-style operations (create/get/list/delete)
- Transfer transaction with:
  - account existence checks
  - currency validation
  - insufficient balance protection
  - deadlock-safe transfer logic

## Authorization Rules

Protected APIs require `Authorization: Bearer <token>`.

Current business rules:

- `POST /accounts`: user can only create an account for their own username
- `GET /accounts/:id`: user can read only their own account
- `GET /accounts`: user sees only accounts they own
- `DELETE /accounts/:id`: user can delete only their own account
- `POST /transfers`: user can transfer money only from an account they own
- `GET /users/:username`: self-only access (token username must match path username)

## Project Structure

```text
api/           HTTP handlers, middleware, request/response models, API tests
db/migration/  SQL migrations
db/query/      SQL query files used by SQLC
db/sqlc/       generated SQLC code + transactional store + DB tests
db/mock/       generated GoMock store
util/          config, random helpers, password hashing, currency validation
toeken/        token interfaces + JWT maker + payload/tests
```

## Prerequisites

- Go 1.25+
- Docker
- PostgreSQL client tools (optional but useful)
- `migrate` CLI
- `sqlc`
- `mockgen` (from `go.uber.org/mock`)

## Environment Configuration

Create/update `app.env`:

```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=your-32-char-min-secret-key-here
ACCESS_TOKEN_DURATION=15m
```

Important:
- `TOKEN_SYMMETRIC_KEY` must be at least 32 characters.
- Do not use training secrets in real systems.

## Quick Start

1. Start PostgreSQL in Docker

```bash
make postgres
```

2. Create database

```bash
docker exec -it postgres17 createdb --username=root --owner=root simple_bank
```

3. Run migrations

```bash
make migrateup
```

4. Generate DB code and mocks (if needed)

```bash
make sqlc
make mock
```

5. Start API server

```bash
make server
```

Server listens on `SERVER_ADDRESS` (default: `0.0.0.0:8080`).

## Run Tests

Run all tests with coverage:

```bash
make test
```

Notes:
- `api` tests are unit tests with mocks.
- `db/sqlc` tests are integration-style tests and require a reachable PostgreSQL instance.

## API Examples

### Create user

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username":"alice",
    "password":"secret123",
    "full_name":"Alice Doe",
    "email":"alice@example.com"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username":"alice",
    "password":"secret123"
  }'
```

Use `access_token` from login response for protected endpoints:

```bash
curl -X GET "http://localhost:8080/accounts?page_id=1&page_size=5" \
  -H "Authorization: Bearer <access_token>"
```

## Common Dev Workflow

1. Add/modify SQL in `db/query/*.sql`
2. Run `make sqlc`
3. If store interface changed, run `make mock`
4. Update handlers/tests
5. Run `make test`

## Known Training Constraints

- The token package folder name is `toeken/` (intentionally kept as-is in this project).
- This project is training-focused; prioritize learning patterns and correctness over production polish.

## License

Training/educational project.

## Automation Scripts

Use these helper scripts to reduce setup friction:

- `scripts/install_requirements.sh`
  - Installs required local tooling (`docker`, `migrate`, `sqlc`, `mockgen`)
- `scripts/dev_up.sh`
  - Starts PostgreSQL container
  - Ensures DB exists
  - Runs migrations
  - Regenerates sqlc/mocks
- `scripts/dev_up.sh --run`
  - Does all setup and starts API server
- `scripts/dev_down.sh`
  - Stops local PostgreSQL container

Examples:

```bash
./scripts/install_requirements.sh
./scripts/dev_up.sh
./scripts/dev_up.sh --run
./scripts/dev_down.sh
```

## Run with Docker

Build and run API + PostgreSQL:

```bash
docker compose up --build
```

Server will be available at:

```text
http://localhost:8080
```

Stop and remove containers:

```bash
docker compose down
```
