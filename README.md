# SimpleBank

SimpleBank is a backend training project built in Go around a small but realistic banking domain: users, accounts, authenticated access, and transactional money transfers. The goal is not just to make endpoints respond, but to practice the parts that usually break first in backend systems: database consistency, access control, repeatable setup, and test coverage.

## Highlights

- JWT-based authentication and ownership-based authorization
- PostgreSQL-backed persistence with SQLC-generated queries
- Deadlock-safe transfer transaction logic
- API unit tests, middleware tests, and database integration tests
- CI/CD through GitHub Actions
- Local automation through Bash scripts
- Containerized runtime with Docker and Docker Compose

## Tech Stack

- Go `1.25.0`
- Gin `v1.12.0`
- PostgreSQL `17.9-alpine`
- SQLC 
- golang-migrate 
- GoMock `v0.6.0`
- Testify `v1.11.1`
- JWT `v5 v5.3.1`
## Deployment
This project is configured and maintained with GitHub Actions workflows under .github/. The pipeline is used to validate and deploy the application as part of the training delivery flow.
## Authorization Rules

Protected APIs require `Authorization: Bearer <token>`.

Current business rules:

- `POST /accounts`: user can only create an account for their own username
- `GET /accounts/:id`: user can read only their own account
- `GET /accounts`: user sees only accounts they own
- `DELETE /accounts/:id`: user can delete only their own account
- `POST /transfers`: user can transfer money only from an account they own
- `GET /users/:username`: self-only access (token username must match path username)

## Environment Configuration

Create/update `app.env`:

```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=your-32-char-min-secret-key-here
ACCESS_TOKEN_DURATION=15m
```


## Quick Start

The fastest path is Docker. It avoids installing most local tooling and uses the same application image every time.

1. Build and start the stack:

   ```bash
   docker compose up --build
   ```

2. The API will be available at:

   ```text
   http://localhost:8080
   ```

3. Stop the stack:

   ```bash
   docker compose down
   ```

`docker compose up --build` waits for PostgreSQL, runs migrations, and only then starts the API server.

## Manual Run

If you want to run the app directly on your machine instead of inside Docker:

1. Install requirements:

   ```bash
   ./scripts/install_requirements.sh
   ```

2. Prepare local services and database:

   ```bash
   ./scripts/dev_up.sh
   ```

3. Start the server:

   ```bash
   DB_HOST=localhost make server
   ```

The manual path is mainly for development work where you want direct access to tools like `sqlc`, `mockgen`, or database test runs.

## Configuration

Application configuration lives in [app.env](/app.env).

Important values:

- `DB_SOURCE` defines the PostgreSQL connection string
- `TOKEN_SYMMETRIC_KEY` must be at least 32 characters
- `ACCESS_TOKEN_DURATION` controls JWT lifetime
- local runs use the `DB_SOURCE` value from `app.env`
- Docker Compose overrides `DB_SOURCE` so the API container can reach the `postgres` service

## What It Does

- creates users
- logs users in and returns JWT access tokens
- protects account and transfer routes with middleware
- allows each user to manage only their own accounts
- allows each user to fetch only their own profile
- performs transfers with validation for:
  - account existence
  - owner authorization
  - currency match
  - insufficient balance
  - self-transfer rejection

## Authorization Rules

Protected endpoints require:

```text
Authorization: Bearer <token>
```

Current rules:

- `POST /accounts`: authenticated user can create an account only for themselves
- `GET /accounts/:id`: authenticated user can read only their own account
- `GET /accounts`: authenticated user sees only their own accounts
- `DELETE /accounts/:id`: authenticated user can delete only their own account
- `POST /transfers`: authenticated user can transfer only from an account they own
- `GET /users/:username`: authenticated user can fetch only their own user record

## API Examples

Create a user:

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

Login:

```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username":"alice",
    "password":"secret123"
  }'
```

Call a protected endpoint:

```bash
curl -X GET "http://localhost:8080/accounts?page_id=1&page_size=5" \
  -H "Authorization: Bearer <access_token>"
```

## Tests

Run the full suite:

```bash
make test
```

Test layers:

- `api`: handler and middleware tests with mocks
- `db/sqlc`: integration-style tests against PostgreSQL
- `token`: JWT maker and payload tests
- `util`: utility-level tests such as password hashing

## Automation Scripts

- [install_requirements.sh](/scripts/install_requirements.sh)
  Installs or checks local requirements for macOS and Linux.
- [dev_up.sh](/scripts/dev_up.sh)
  Starts PostgreSQL, creates the database if needed, runs migrations, and regenerates code artifacts.
- [dev_down.sh](/scripts/dev_down.sh)
  Stops the local PostgreSQL container.

## Deployment

This project is deployed through GitHub Actions. The workflow is used to validate changes and keep delivery repeatable.

## License

This project is released under the MIT License. See [LICENSE](/LICENSE).
