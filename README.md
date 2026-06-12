# anyx-user-api

A production-style RESTful API built with Go. Manages users with `name` and `dob` (date of birth) and dynamically calculates each user's `age` on every read — age is never stored in the database.

## Tech Stack

| Layer | Technology |
|---|---|
| HTTP Framework | [GoFiber v2](https://gofiber.io/) |
| Database | PostgreSQL |
| DB Access Layer | [SQLC](https://sqlc.dev/) |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Config | [godotenv](https://github.com/joho/godotenv) |
| Container | Docker + Docker Compose |

## Project Structure

```
.
├── cmd/server/main.go          # Application entry point
├── config/config.go            # Environment variable loader
├── db/
│   ├── migrations/             # SQL migration files
│   ├── queries/users.sql       # SQLC-annotated queries
│   └── sqlc/                   # Auto-generated DB layer
├── internal/
│   ├── handler/                # HTTP handlers (parse → call service → respond)
│   ├── logger/                 # Uber Zap wrapper
│   ├── middleware/             # RequestID + Logger middleware
│   ├── models/                 # Request/response structs
│   ├── repository/             # DB access abstraction
│   ├── routes/                 # Route registration
│   └── service/                # Business logic + age calculation
├── .env.example                # Copy to .env and fill in values
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── sqlc.yaml
```

---

## Prerequisites

Make sure the following are installed:

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL 14+](https://www.postgresql.org/download/) (or use Docker)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [SQLC CLI](https://docs.sqlc.dev/en/latest/overview/install.html)
- (Optional) [Docker Desktop](https://www.docker.com/products/docker-desktop/)

---

## Option A: Run Locally (without Docker)

### Step 1 — Clone and set up environment

```bash
git clone https://github.com/YOUR_USERNAME/anyx-user-api.git
cd anyx-user-api
cp .env.example .env
```

Open `.env` and fill in your PostgreSQL credentials:

```env
PORT=3000
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/anyx_db?sslmode=disable
ENVIRONMENT=development
```

> ⚠️ **Human action required**: Create a PostgreSQL database named `anyx_db` first:
> ```sql
> CREATE DATABASE anyx_db;
> ```

### Step 2 — Run database migrations

```bash
make migrate-up
```

This creates the `users` table.

### Step 3 — Download Go dependencies

```bash
make tidy
```

### Step 4 — Start the server

```bash
make run
```

The API is now available at `http://localhost:3000`.

---

## Option B: Run with Docker Compose

> ⚠️ **Human action required**: Make sure Docker Desktop is running.

```bash
# Build and start both PostgreSQL and the API
make docker-up

# Run migrations inside Docker
docker exec anyx_api sh -c "migrate -path /app/db/migrations -database 'postgres://postgres:yourpassword@postgres:5432/anyx_db?sslmode=disable' up"

# Tail logs
make docker-logs
```

To stop everything:

```bash
make docker-down
```

---

## API Endpoints

Base URL: `http://localhost:3000`

### `POST /users` — Create a user

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "dob": "1990-05-10"}'
```

Response `201 Created`:
```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10" }
```

---

### `GET /users/:id` — Get user by ID (includes age)

```bash
curl http://localhost:3000/users/1
```

Response `200 OK`:
```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
```

---

### `PUT /users/:id` — Update a user

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated", "dob": "1991-03-15"}'
```

Response `200 OK`:
```json
{ "id": 1, "name": "Alice Updated", "dob": "1991-03-15" }
```

---

### `DELETE /users/:id` — Delete a user

```bash
curl -X DELETE http://localhost:3000/users/1
```

Response `204 No Content`

---

### `GET /users` — List all users (paginated)

```bash
curl "http://localhost:3000/users?page=1&limit=20"
```

Response `200 OK`:
```json
{
  "data": [
    { "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

---

### `GET /health` — Health check

```bash
curl http://localhost:3000/health
```

Response: `{ "status": "ok", "service": "anyx-user-api" }`

---

## Running Tests

```bash
make test
```

The unit tests cover the `CalculateAge` function including edge cases like:
- Birthday already passed this year
- Birthday is today
- Birthday is tomorrow
- Leap year birthdays

---

## SQLC Regeneration

If you modify `db/queries/users.sql` or the migration schema, regenerate the DB layer with:

```bash
make sqlc
```

---

## Response Headers

Every response includes:

| Header | Value |
|---|---|
| `X-Request-ID` | Unique UUID per request (reused from client if provided) |

---

## Human Interaction Checklist

Before running the project for the first time, you need to:

- [ ] Install Go 1.22+ and add it to PATH
- [ ] Install `golang-migrate` CLI
- [ ] Install `sqlc` CLI
- [ ] Create a PostgreSQL database (`anyx_db`)
- [ ] Copy `.env.example` → `.env` and fill in `DATABASE_URL`
- [ ] Run `make migrate-up` to create the `users` table
- [ ] Run `make tidy` to download dependencies
- [ ] Run `make run` to start the server
