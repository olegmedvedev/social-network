# Social Network API (Go + GraphQL + PostgreSQL)

Educational project: a basic social network with registration, login, and profile viewing.

## Stack
- Go (gqlgen, database/sql)
- PostgreSQL (Docker)
- Goose (migrations)
- GraphQL API
- godotenv (for local development)

## Quick Start

### 1. Clone the repository and go to the project folder

```sh
git clone ...
cd social-network
```

### 2. Set up environment variables

Create a `.env` file based on `.env.example` and fill in your values:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=socialuser
DB_PASSWORD=socialpass
DB_NAME=socialdb
JWT_SECRET=supersecretkey
```

### 3. Start PostgreSQL with Docker

```sh
docker compose up -d
```

### 4. Run migrations

```sh
goose -dir ./migrations postgres "host=localhost port=5432 user=socialuser password=socialpass dbname=socialdb sslmode=disable" up
```

### 5. Start the server

```sh
go run cmd/gateway/main.go
```

GraphQL Playground will be available at http://localhost:8080/

---

## Features
- User registration (name, email, password)
- Login (email, password)
- Get profile by id and current user

---

## Best practices for AWS deployment
- All secrets and parameters (passwords, keys) are set via environment variables or AWS Secrets Manager.
- In production, godotenv is not used — environment variables are set by the cloud platform (ECS, Lambda, EC2, etc.).
- Do not commit `.env` or real secrets to the repository!

---

## Migrations
- [Goose](https://github.com/pressly/goose) is used for migrations.
- All migrations are in the `migrations/` folder.

---
