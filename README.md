# Todo API

## Overview

Todo API is a RESTful task-management backend written in Go. It provides user registration and login, JWT-based authentication, and user-scoped CRUD operations for todo items.

The API is designed as a small, maintainable backend service with PostgreSQL persistence. It can be run locally with Go or through Docker Compose, with a development container configured for hot reloading.

### Target users

Todo API is intended for:
- Developers building a todo or productivity client that needs a simple authenticated backend.
- Learners exploring REST API design, JWT authentication, PostgreSQL, and Go project structure.
- Small applications that need isolated task lists for multiple users.

### Key features
- 🔐 **User authentication** - Register accounts with bcrypt-hashed passwords and sign in to receive a JWT.
- ✅ **Authenticated task management** - Create, read, update, and delete tasks belonging to the signed-in user.
- 🗄️ **PostgreSQL persistence** - Store users and tasks in PostgreSQL with versioned SQL migrations.
- 🧱 **Layered project structure** - Separate configuration, database access, handlers, middleware, models, and repositories.
- 🐳 **Docker support** - Run the API and PostgreSQL together with Docker Compose.
- 🔄 **Development hot reload** - Use Air in the development image to rebuild the API when source files change.

### Tech Stack
- 🧩 **Language**: Go 1.26.5
- 🚀 **HTTP framework**: Gin v1.12.0
- 💾 **Database**: PostgreSQL 16
- 🔌 **Database driver**: pgx v5
- 🔑 **Authentication**: JWT with HS256
- 🔒 **Password hashing**: bcrypt
- 🐳 **Containerization**: Docker and Docker Compose
- 🛠️ **Database migrations**: golang-migrate
- 📦 **Configuration**: `.env` files and environment variables via `godotenv`

### API endpoints

The default base URL is `http://localhost:8000`.

| Method | Endpoint | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/` | No | Health response confirming the service and database are available. |
| `POST` | `/auth/signup` | No | Create a user account. |
| `POST` | `/auth/login` | No | Authenticate and receive a JWT. |
| `POST` | `/task` | Bearer token | Create a task. |
| `GET` | `/task/all` | Bearer token | List the authenticated user's tasks. |
| `GET` | `/task/:id` | Bearer token | Get one of the authenticated user's tasks. |
| `PUT` | `/task/:id` | Bearer token | Update one of the authenticated user's tasks. |
| `DELETE` | `/task/:id` | Bearer token | Delete one of the authenticated user's tasks. |

Protected endpoints require the header:

```http
Authorization: Bearer <jwt>
```

#### Example requests

Create an account:

```bash
curl -X POST http://localhost:8000/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","confirm_password":"password123"}'
```

Log in:

```bash
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Create a task using the token returned by login:

```bash
curl -X POST http://localhost:8000/task \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Read the API documentation","description":"Review the available endpoints","is_completed":false}'
```

## Product thinking

Todo API was built around a focused backend workflow:
- 👤 Keep each user's tasks isolated through authenticated, user-scoped queries.
- 🔒 Protect account credentials with bcrypt and protect API access with expiring JWTs.
- 🧱 Keep responsibilities separated so handlers, repositories, and middleware can evolve independently.
- 🛠️ Make local development repeatable with migrations and containerized PostgreSQL.

## Developer instructions

### Prerequisites

Choose one of the following setup options:
- Go 1.26.5 or later, PostgreSQL 16, and the `migrate` CLI for local development.
- Docker and Docker Compose for the containerized setup.

### Environment variables

Create a `.env` file in the project root:

```dotenv
DB_USER=superuser
DB_PASSWORD=your_password
DB_NAME=todo_app_db
DB_HOST=localhost
DB_PORT=5432
DATABASE_URL=postgres://superuser:your_password@localhost:5432/todo_app_db?sslmode=disable
APP_PORT=8000
SECRET_KEY=replace_with_a_long_random_secret
```

Keep `.env` out of version control and use a strong, unique `SECRET_KEY` outside local development.

### Local installation

1. Clone the repository and enter the project directory:

```bash
git clone <repository-url>
cd to-do-app
```

2. Install the migration CLI:

```bash
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

3. Start PostgreSQL. You can start only the database service with Docker Compose:

```bash
docker compose up -d db
```

4. Apply the database migrations:

```bash
./scripts/migrate.sh
```

5. Start the API:

```bash
go run ./cmd/api
```

The API is available at `http://localhost:8000` unless `APP_PORT` is changed.

### Docker installation guide

#### Prerequisites

Install Docker Desktop or Docker Engine with the Docker Compose plugin. Verify that both commands are available:

```bash
docker --version
docker compose version
```

#### Run with Docker Compose

1. Create a `.env` file in the project root. Use `db` as the database host because the API connects to PostgreSQL through the Compose network:

```dotenv
DB_USER=superuser
DB_PASSWORD=your_password
DB_NAME=todo_app_db
DB_HOST=db
DB_PORT=5432
DATABASE_URL=postgres://superuser:your_password@db:5432/todo_app_db?sslmode=disable
APP_PORT=8000
SECRET_KEY=replace_with_a_long_random_secret
```

2. Build the development image and start the API with PostgreSQL:

```bash
docker compose up --build
```

3. In a separate terminal, apply the database migrations:

```bash
./scripts/migrate.sh
```

4. Open the API at `http://localhost:8000`. The development image uses Air for hot reloading.

Stop the services and remove the containers with:

```bash
docker compose down
```

The PostgreSQL data is stored in the `postgres_data` named volume and remains available after stopping the services. To remove the volume as well, use:

```bash
docker compose down -v
```

#### Build the production image

Build the production image defined in `Dockerfile`:

```bash
docker build -t todo-api .
```

Run the image with environment variables and publish the API port:

```bash
docker run --env-file .env -p 8000:8000 todo-api
```

The production container still requires access to a running PostgreSQL instance. Set `DATABASE_URL` to a hostname reachable from the container before starting it.

### Project structure

```text
.
├── cmd/
│   └── api/              Application entrypoint and route registration
├── internal/
│   ├── config/           Environment and application configuration
│   ├── db/               PostgreSQL connection setup
│   ├── handlers/         HTTP request handlers
│   ├── middleware/       JWT authentication middleware
│   ├── models/           User and task data models
│   └── repository/       Database queries and persistence operations
├── migrations/           Versioned PostgreSQL schema migrations
└── scripts/              Development and migration scripts
```

## 🤝 Contributor expectations

For a bug fix or new feature, create a branch and open a pull request:

```bash
git checkout -b <name-of-your-branch>
```

Please keep changes focused, update migrations when the data model changes, and include tests for new behavior where practical. You can also open an issue to report a bug or suggest an improvement.

## 🙏 Request

If this project is useful to you, consider starring the repository and sharing feedback through an issue or pull request.