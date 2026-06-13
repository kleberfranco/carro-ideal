# Carro Ideal

Carro Ideal is a Go-based vehicle recommendation project designed to provide a production-ready foundation for a layered web application with PostgreSQL persistence, containerized deployment, and health checks.

## Phase 1 Completion

This repository now includes the Phase 1 foundation and infrastructure:

- Go project structure with `app/`, `config/`, `migrations/`, and `web/`
- Environment-based configuration via `.env` and `config/config.go`
- PostgreSQL connection pooling and health checks
- Automatic database migrations on startup
- Chi v5 router with API, admin, web, and health endpoints
- Dockerfile and `docker-compose.yml` for local development
- `.gitignore`, `.dockerignore`, and `.env.example`

## Local Setup

1. Copy the example environment file:

```sh
cp .env.example .env
```

2. Build and start the application with Docker Compose:

```sh
docker compose up --build
```

3. Confirm the application is healthy:

```sh
curl http://localhost:8080/health
```

## Available Endpoints

- `GET /health` - application and database health check
- `POST /api/auth/register` - user registration placeholder
- `POST /api/auth/login` - user login placeholder
- `GET /web` and `GET /` - frontend landing page placeholder

## Notes

- The application runs on port `8080` by default.
- Database migrations run automatically when the app starts.
- PostgreSQL service is configured via Docker Compose on port `5432`.
