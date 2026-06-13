# Go Backend Users API

RESTful API built with GoFiber to manage users with `name` and `dob`. The API stores date of birth in PostgreSQL and calculates age dynamically when users are fetched.

## Tech Stack

- Go
- GoFiber
- PostgreSQL
- SQLC
- Uber Zap
- go-playground/validator

## Requirements

- Go 1.23+
- PostgreSQL
- SQLC
- Optional: Docker and Docker Compose

## Environment

Copy the sample environment file and update values if needed:

```bash
cp .env.example .env
```

Required variables:

```text
APP_ENV=development
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/users_api?sslmode=disable
```

On Windows PowerShell:

```powershell
$env:APP_ENV="development"
$env:PORT="8080"
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/users_api?sslmode=disable"
```

## Database Setup

Create the database:

```sql
CREATE DATABASE users_api;
```

Run migrations with golang-migrate:

```bash
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/users_api?sslmode=disable" up
```

Rollback:

```bash
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/users_api?sslmode=disable" down
```

## SQLC

The SQLC config and generated code are included. To regenerate after changing SQL:

```bash
sqlc generate
```

## Run

Install dependencies and run:

```bash
go mod tidy
go run ./cmd/server
```

Health check:

```bash
curl http://localhost:8080/healthz
```

## API Endpoints

### Create User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

Response:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID

```bash
curl http://localhost:8080/users/1
```

Response:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 36
}
```

### Update User

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```

Response:

```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/users/1 -i
```

Response:

```text
HTTP/1.1 204 No Content
```

### List Users

```bash
curl http://localhost:8080/users
```

Response:

```json
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 36
  }
]
```

## Error Response

Errors use a consistent JSON shape:

```json
{
  "error": "validation_error",
  "message": "dob must be a valid date in YYYY-MM-DD format"
}
```

## Tests

Run all tests:

```bash
go test ./...
```

## Docker

Start PostgreSQL:

```bash
docker compose up -d db
```

Run migrations against the Docker database:

```bash
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/users_api?sslmode=disable" up
```

Start the API:

```bash
docker compose up --build api
```

The API will be available at:

```text
http://localhost:8080
```
