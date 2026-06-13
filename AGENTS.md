# AGENTS.md

## Assignment Overview

Build a RESTful API in Go for managing users with `name` and `dob` fields. The API must store each user's date of birth in the database and calculate the user's age dynamically when returning user details.

This project is for an internship backend development assignment, so prioritize correctness, clear structure, readable code, and a professional README over unnecessary complexity.

## Core Objective

- Implement a Go backend API for CRUD operations on users.
- Persist users in a SQL database.
- Store `dob` as a database date value.
- Return `age` dynamically when fetching user details or listing users.
- Use the required libraries and project structure from the assignment.

## Required Tech Stack

- Language: Go
- HTTP framework: GoFiber
- Database: PostgreSQL or MySQL
- SQL access layer: SQLC
- Logging: Uber Zap
- Validation: `go-playground/validator`
- Date and age calculation: Go standard `time` package

Prefer PostgreSQL unless the user explicitly chooses MySQL.

## Required Project Structure

Follow this directory structure unless there is a strong reason to adapt it:

```text
/cmd/server/main.go
/config/
/db/migrations/
/db/sqlc/<generated>
/internal/
|-- handler/
|-- repository/
|-- service/
|-- routes/
|-- middleware/
|-- models/
`-- logger/
```

Suggested responsibilities:

- `cmd/server/main.go`: application entrypoint, config loading, logger setup, DB connection, Fiber app startup.
- `config/`: environment configuration and typed config structs.
- `db/migrations/`: SQL migration files for the database schema.
- `db/query/` or similar: SQLC query files if needed by the selected SQLC config.
- `db/sqlc/`: generated SQLC database access code.
- `internal/handler/`: HTTP handlers, request parsing, response formatting.
- `internal/repository/`: database-facing repository wrapper around SQLC generated methods.
- `internal/service/`: business logic such as age calculation and validation orchestration.
- `internal/routes/`: route registration.
- `internal/middleware/`: request ID, request duration logging, recovery, and related middleware.
- `internal/models/`: request/response DTOs and domain models.
- `internal/logger/`: Zap logger initialization and helpers.

## Required Database Table

Create a `users` table with this schema:

| Field | Type | Constraints |
| --- | --- | --- |
| `id` | `SERIAL` | `PRIMARY KEY` |
| `name` | `TEXT` | `NOT NULL` |
| `dob` | `DATE` | `NOT NULL` |

For PostgreSQL, `SERIAL PRIMARY KEY` is acceptable. For MySQL, use the equivalent auto-incrementing integer primary key.

## API Contract

All request and response bodies should be JSON.

### Create User

`POST /users`

Request:

```json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

Response:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

Expected status codes:

- `201 Created` on success.
- `400 Bad Request` for malformed JSON or validation errors.
- `500 Internal Server Error` for unexpected failures.

### Get User by ID

`GET /users/:id`

Response:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

Expected status codes:

- `200 OK` on success.
- `400 Bad Request` for invalid IDs.
- `404 Not Found` when the user does not exist.
- `500 Internal Server Error` for unexpected failures.

### Update User

`PUT /users/:id`

Request:

```json
{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

Response:

```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

Expected status codes:

- `200 OK` on success.
- `400 Bad Request` for invalid IDs, malformed JSON, or validation errors.
- `404 Not Found` when the user does not exist.
- `500 Internal Server Error` for unexpected failures.

### Delete User

`DELETE /users/:id`

Response:

- HTTP `204 No Content`

Expected status codes:

- `204 No Content` on success.
- `400 Bad Request` for invalid IDs.
- `404 Not Found` when the user does not exist.
- `500 Internal Server Error` for unexpected failures.

### List All Users

`GET /users`

Response:

```json
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 34
  }
]
```

Expected status codes:

- `200 OK` on success.
- `500 Internal Server Error` for unexpected failures.

Optional bonus: add pagination to `GET /users`, for example with `page` and `limit` query parameters. If pagination is added, document it in `README.md`.

## Age Calculation Rules

- Age must be calculated dynamically from `dob`; do not store `age` in the database.
- Use Go's `time` package.
- Correctly handle whether the user's birthday has occurred yet in the current year.
- Use the current date at request time.
- Treat `dob` as a calendar date in `YYYY-MM-DD` format.
- Reject future dates of birth.
- Add a focused unit test for the age calculation function if time allows; this is listed as a bonus but is high-value.

Suggested behavior:

```text
age = currentYear - birthYear
if current month/day is before birth month/day:
    age--
```

## Validation Rules

Use `go-playground/validator` for input validation.

Validate:

- `name` is required.
- `name` is not only whitespace.
- `dob` is required.
- `dob` matches `YYYY-MM-DD`.
- `dob` is a valid calendar date.
- `dob` is not in the future.

Return clean JSON error responses. Keep error messages helpful but concise.

Example error response:

```json
{
  "error": "validation_error",
  "message": "dob must be a valid date in YYYY-MM-DD format"
}
```

## Logging Requirements

Use Uber Zap for logging.

Log key actions such as:

- Server startup.
- Database connection success/failure.
- User creation.
- User retrieval.
- User update.
- User deletion.
- Request validation failures where useful.
- Unexpected internal errors.

Avoid logging sensitive data. This assignment has no sensitive fields, but keep logs professional and structured.

## Middleware Expectations

Required by the base assignment:

- Clean HTTP status codes.
- Clean error handling.

Bonus middleware:

- Inject a `requestId` header into responses.
- Log request duration.

Suggested headers:

- Request ID response header: `X-Request-Id`

Suggested request log fields:

- `request_id`
- `method`
- `path`
- `status`
- `duration_ms`

## SQLC Guidance

Use SQLC for the DB access layer.

Expected SQLC pieces:

- A `sqlc.yaml` configuration file.
- SQL migration for the `users` table.
- Query definitions for create, get by ID, update, delete, and list users.
- Generated Go code under `db/sqlc/` or another documented generated path.

Prefer simple explicit SQL queries over ORM-style abstractions.

Expected queries:

- `CreateUser`
- `GetUserByID`
- `UpdateUser`
- `DeleteUser`
- `ListUsers`

If pagination is implemented:

- `ListUsersPaginated`

## Error Handling Standards

Use consistent JSON responses for errors, except for `204 No Content`, which must have no body.

Recommended error response shape:

```json
{
  "error": "not_found",
  "message": "user not found"
}
```

Common error codes:

- `bad_request`
- `validation_error`
- `not_found`
- `internal_error`

Map database "no rows" errors to `404 Not Found`.

## README Requirements

The final submission must include a clear `README.md` with:

- Project overview.
- Tech stack.
- Prerequisites.
- Environment variables.
- Database setup instructions.
- Migration instructions.
- SQLC generation instructions.
- How to run the server.
- API endpoint examples.
- How to run tests.
- Docker instructions if Docker support is added.

The assignment submission requires pushing code to GitHub and sharing the repository link.

## Optional Bonus Features

Only add these after the core API is working:

- Docker support.
- Pagination for `GET /users`.
- Unit tests for the age calculation function.
- Middleware that injects `requestId` into responses.
- Middleware that logs request duration.

Recommended priority:

1. Unit test age calculation.
2. Request ID and request duration middleware.
3. Docker support.
4. Pagination.

## Implementation Priorities

Work in this order:

1. Initialize Go module and dependencies.
2. Create the required project structure.
3. Add config and logger setup.
4. Add database migration and SQLC config/query files.
5. Generate SQLC code.
6. Implement repository layer.
7. Implement service layer, including age calculation and validation.
8. Implement HTTP handlers.
9. Register routes.
10. Add middleware.
11. Add tests for age calculation.
12. Add README setup and API documentation.
13. Verify the project builds and tests pass.

## Verification Checklist

Before considering the assignment complete, verify:

- `go fmt` has been run.
- `go test ./...` passes.
- The server starts successfully.
- SQLC generated code is present and used.
- Migration creates the expected `users` table.
- `POST /users` creates a user.
- `GET /users/:id` returns user with dynamic `age`.
- `PUT /users/:id` updates `name` and `dob`.
- `DELETE /users/:id` returns `204 No Content`.
- `GET /users` returns a list with dynamic `age`.
- Invalid input returns `400`.
- Missing user returns `404`.
- Unexpected failures return `500`.
- README has setup and run instructions.

## Notes For Future Codex Work

- Do not store `age` in the database.
- Keep the project aligned with the required directory structure.
- Use SQLC rather than hand-written generic DB scanning in handlers.
- Keep handlers thin; put business logic in services.
- Prefer deterministic, testable age calculation by accepting a `today time.Time` argument in the core age function, then wrapping it for request-time use.
- Do not over-engineer authentication, authorization, roles, or unrelated features; they are outside the assignment scope.
- Keep responses close to the examples in the screenshots unless a documented improvement is needed.
