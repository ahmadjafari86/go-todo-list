# 📝 Golang Gin Todo API

A backend **Todo List** service implemented with **Go (Golang)** and **Gin Gonic**.  
This project includes JWT authentication, CRUD management for todos, modular structure, and production-ready configuration.

---

## ✨ Features

- Modular and maintainable structure (`cmd/`, `internal/`, `configs/`)
- JWT-based authentication (each user can only access their own todos)
- Input validation with [go-playground/validator](https://github.com/go-playground/validator)
- Error handling compliant with [RFC7807 Problem+JSON](https://datatracker.ietf.org/doc/html/rfc7807)
- PostgreSQL with GORM (Connection Pool + Timeouts)
- Graceful Shutdown support
- Structured Logging with [zerolog](https://github.com/rs/zerolog)
- Integration testing with [Testcontainers-Go](https://github.com/testcontainers/testcontainers-go)
- API documentation with Swagger/OpenAPI

---

## 📂 Project Structure

```
.
├── cmd/server/main.go     # entrypoint
├── configs/               # configuration and .env
├── internal/
│   ├── handlers/          # gin handlers
│   ├── middleware/        # JWT and middlewares
│   ├── models/            # User, Todo, DTOs
│   ├── repository/        # data access (DB)
│   ├── service/           # business logic
│   └── validation/        # error handling
├── docs/                  # swagger generated files
├── tests/                 # integration tests (testcontainers)
├── go.mod / go.sum
└── Dockerfile
```

---

## 🚀 Getting Started

### 1) Run with Docker Compose (recommended)

```bash
docker-compose up --build
```

The service will run at `http://localhost:8282`.  
Postgres will run on port `5432`.

---

### 2) Run locally (without Docker)

#### Prerequisites:

- Go 1.22+
- PostgreSQL running

#### Steps:

1. Create `.env` file:

   ```env
   APP_PORT=8080
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=todo_db
   JWT_SECRET=supersecretjwtkey
   ```

2. Create database:

   ```bash
   createdb todo_db
   ```

3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

---

## 📖 API Documentation (Swagger)

After running the project, open:

👉 [http://localhost:8282/swagger/index.html](http://localhost:8282/swagger/index.html)

---

## 🔑 Authentication

All `/api/todos/*` routes require JWT authentication.

1. Register a new user:

   ```http
   POST /auth/register
   ```

2. Login:

   ```http
   POST /auth/login
   ```

   Response:

   ```json
   { "token": "jwt.token.here" }
   ```

3. Send the token in the request header:
   ```http
   Authorization: Bearer jwt.token.here
   ```

---

## 🧪 Tests (Testcontainers)

To run integration tests:

```bash
go test ./tests/... -v
```

These tests spin up a real PostgreSQL container with Testcontainers and validate CRUD & Auth functionality.

---

## 🛠️ Development Commands

### Generate Swagger docs

```bash
swag init -g cmd/server/main.go -o docs
```

### Run locally

```bash
go run cmd/server/main.go
```

---

## 📜 License

MIT License
