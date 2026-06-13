# Quickstart Guide: Carro Ideal MVP

**Date**: 2 de junho de 2026  
**Status**: Phase 1 Complete - Ready for Implementation

---

## Prerequisites

- **Go**: 1.22+ ([download](https://golang.org/dl/))
- **Docker**: 20.10+ ([download](https://www.docker.com/products/docker-desktop))
- **Docker Compose**: 2.0+ (included with Docker Desktop)
- **PostgreSQL Client** (optional, for direct DB access): `brew install postgresql` (macOS) or similar
- **Git**: For version control

---

## Project Setup (5 minutes)

### 1. Clone Repository

```bash
git clone https://github.com/kleberfranco/carro-ideal.git
cd carro-ideal
```

### 2. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

**Default `.env` contents** (for local development):
```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=carro_ideal
DB_USER=postgres
DB_PASSWORD=dev_password

# Application
PORT=8080
ENV=development
SESSION_SECRET=your-secret-key-change-in-production
SECURE_COOKIES=false
LOG_LEVEL=debug
```

### 3. Start Docker Services

Start PostgreSQL and the application in Docker:

```bash
docker-compose up -d
```

Verify both services are running:

```bash
docker-compose ps
# Should show:
# NAME           STATUS
# carro-ideal    Up X seconds (exit 0)
# postgres       Up X seconds
```

### 4. Verify Setup

Check health endpoint:

```bash
curl http://localhost:8080/health
# Should return: {"status":"healthy","timestamp":"2024-06-01T15:00:00Z"}
```

---

## Running the Application

### Option A: In Docker (Recommended for TCC Submission)

**Start Everything**:
```bash
docker-compose up -d
# Application available at http://localhost:8080
```

**View Logs**:
```bash
docker-compose logs -f app
# Press Ctrl+C to exit
```

**Stop Services**:
```bash
docker-compose down
# Add -v to also remove database volumes (WARNING: deletes data)
```

### Option B: Local Development (Go Running Natively)

**Prerequisites**: PostgreSQL must be running

```bash
# Start PostgreSQL in Docker
docker-compose up -d postgres

# Install Go dependencies
go mod download

# Run database migrations (first time only)
go run ./cmd/migrate/main.go up

# Start application
go run ./app/main.go
# Application available at http://localhost:8080
```

**Hot Reload** (optional, for faster development):
```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Run with air
air
```

---

## First Steps: Explore the Application

### 1. Register as User

**Via Browser**:
1. Navigate to http://localhost:8080
2. Click "Sign Up"
3. Fill in:
   - Name: "Test User"
   - Email: "test@example.com"
   - Password: "TestPass123"
4. Click "Register"

**Via cURL**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "TestPass123"
  }'

# Response (201):
# {"success":true,"data":{"user_id":1,"email":"test@example.com"},"message":"Registration successful"}
```

### 2. Login

**Via cURL**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123"
  }'

# Response (200):
# {"success":true,"data":{"user_id":1,"email":"test@example.com","role":"USER"},"message":"Login successful"}
```

The `-c cookies.txt` flag saves the session cookie for subsequent requests.

### 3. Get Questionnaire Questions

```bash
curl -X GET http://localhost:8080/api/questions \
  -b cookies.txt

# Response: List of all active questions with answer options
```

### 4. Answer Questionnaire

```bash
curl -X POST http://localhost:8080/api/questionnaire/answers \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"question_id": 1, "answer_option_id": 2},
      {"question_id": 2, "answer_option_id": 5}
    ]
  }'
```

### 5. Generate Recommendations

```bash
curl -X POST http://localhost:8080/api/recommendations/generate \
  -b cookies.txt \
  -H "Content-Type: application/json"

# Response: Top 10 recommended vehicles with scores and reasons
```

---

## Admin Access

### Create Admin User

Admin users must be created directly in the database (no self-registration for admin):

```bash
docker-compose exec postgres psql -U postgres -d carro_ideal -c "
UPDATE users SET role = 'ADMIN' WHERE email = 'test@example.com';
"
```

### Admin Dashboard

**Via Browser**:
1. Login as admin user
2. Navigate to http://localhost:8080/web/admin
3. Manage:
   - Vehicles (Create, Edit, Delete)
   - Categories (Create, Edit, Delete)
   - Questions (Create, Edit, Delete)
   - Answer Options (Create, Edit, Delete)

**Via API**:
```bash
# List all vehicles (including inactive)
curl -X GET http://localhost:8080/api/admin/vehicles \
  -b cookies.txt

# Create vehicle
curl -X POST http://localhost:8080/api/admin/vehicles \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "brand": "Toyota",
    "model": "Corolla",
    "year": 2023,
    "fuel_type": "Gasoline",
    "transmission": "Automatic",
    "price_range": "80000-120000",
    "seats": 5,
    "trunk_capacity": 300,
    "consumption_city": 9.5,
    "consumption_highway": 13.2,
    "strengths": "Reliable",
    "weaknesses": "Limited cargo"
  }'
```

---

## Testing the Application

### Run Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Integration Tests

```bash
# Start test database
docker-compose up -d postgres

# Run integration tests (marked with // +build integration)
go test -tags=integration ./...
```

### Manual Testing Script

```bash
# Create test data and verify core flows
./scripts/test-flow.sh
```

---

## Database Management

### Access PostgreSQL Directly

```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d carro_ideal

# Common queries:
# List users
SELECT id, email, name, role FROM users;

# List vehicles
SELECT id, brand, model, year FROM vehicles;

# List recommendations for a user
SELECT r.id, r.created_at, COUNT(ri.id) as vehicle_count
FROM recommendations r
LEFT JOIN recommendation_items ri ON ri.recommendation_id = r.id
WHERE r.user_id = 1
GROUP BY r.id, r.created_at;

# Exit psql
\q
```

### Reset Database

```bash
# Delete all data but keep schema
docker-compose exec postgres psql -U postgres -d carro_ideal -c "
TRUNCATE TABLE recommendation_items CASCADE;
TRUNCATE TABLE recommendations CASCADE;
TRUNCATE TABLE user_answers CASCADE;
TRUNCATE TABLE answer_options CASCADE;
TRUNCATE TABLE questions CASCADE;
TRUNCATE TABLE vehicles CASCADE;
TRUNCATE TABLE vehicle_categories CASCADE;
TRUNCATE TABLE sessions CASCADE;
TRUNCATE TABLE users CASCADE;
"

# Or completely reset with volumes
docker-compose down -v
docker-compose up -d  # Will re-run migrations and seed data
```

### View Database Schema

```bash
docker-compose exec postgres psql -U postgres -d carro_ideal -c "
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
"

# View specific table schema
docker-compose exec postgres psql -U postgres -d carro_ideal -c "\d vehicles"
```

---

## Troubleshooting

### Application Won't Start

**Error**: Connection refused at `postgres:5432`

**Solution**:
```bash
# Ensure PostgreSQL is running
docker-compose ps

# If postgres not running, start it
docker-compose up -d postgres

# Wait 10 seconds for database to be ready
sleep 10

# Then start app
docker-compose up -d app
```

### Migrations Failed

**Error**: "No migrations found" or schema issues

**Solution**:
```bash
# Check migration status
docker-compose exec app go run ./cmd/migrate/main.go version

# Force re-run migrations
docker-compose down -v
docker-compose up -d
```

### Port Already in Use

**Error**: "Address already in use: :8080"

**Solution**:
```bash
# Change port in .env
echo "PORT=8081" >> .env

# Or kill existing process
lsof -ti :8080 | xargs kill -9
```

### Database Corrupted

**Solution**:
```bash
# Full reset
docker-compose down -v
rm -rf postgres_data  # If persisted locally
docker-compose up -d
```

---

## Development Workflow

### Making Changes

1. **Edit code** in `app/` or `web/` directories
2. **For local dev with hot reload**:
   ```bash
   air  # Auto-restarts on file changes
   ```
3. **For Docker**:
   ```bash
   docker-compose up --build  # Rebuild image with changes
   ```

### Adding Database Migrations

1. Create new migration files:
   ```bash
   mkdir -p migrations
   touch migrations/000X_description.up.sql
   touch migrations/000X_description.down.sql
   ```

2. Write SQL (see [data-model.md](data-model.md) for schema)

3. Run migrations:
   ```bash
   docker-compose exec app go run ./cmd/migrate/main.go up
   ```

### Adding New API Endpoint

1. Create handler in `app/internal/{domain}/handler.go`
2. Create service in `app/service/`
3. Create repository in `app/repository/` if DB access needed
4. Wire handler into router in `app/internal/{domain}/routes.go`
5. Document in [contracts/api.md](contracts/api.md)
6. Add tests in `*_test.go` files

Example:
```go
// app/internal/vehicle/handler.go
func (h *Handler) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
  id := chi.URLParam(r, "id")
  vehicle, err := h.service.GetVehicleByID(r.Context(), id)
  if err != nil {
    h.error(w, err)
    return
  }
  h.json(w, http.StatusOK, vehicle)
}
```

---

## Code Quality & Standards

### Run Linter

```bash
# Install golangci-lint
brew install golangci-lint

# Run linter
golangci-lint run ./...

# Fix issues automatically
golangci-lint run --fix ./...
```

### Format Code

```bash
go fmt ./...
```

### Check Dependencies

```bash
go mod tidy
go mod verify
```

---

## Deployment Checklist

Before final submission:

- [ ] All tests pass: `go test ./...`
- [ ] No lint errors: `golangci-lint run ./...`
- [ ] Code formatted: `go fmt ./...`
- [ ] Docker builds: `docker build -t carro-ideal .`
- [ ] Docker runs: `docker-compose up -d && curl localhost:8080/health`
- [ ] README complete and clear
- [ ] API documentation up to date
- [ ] Database schema documented
- [ ] No hardcoded secrets in code
- [ ] Environment variables documented in `.env.example`
- [ ] Git history clean with meaningful commits
- [ ] Tags created: `git tag v1.0.0`

---

## Production Deployment (Future)

### Build Production Docker Image

```bash
docker build -f Dockerfile.prod -t carro-ideal:latest .
```

### Environment Variables (Production)

```bash
DB_HOST=prod-postgres.example.com
DB_PORT=5432
DB_NAME=carro_ideal
DB_USER=app_user
DB_PASSWORD=strong_password_here
PORT=8080
ENV=production
SESSION_SECRET=$(openssl rand -hex 32)
SECURE_COOKIES=true
LOG_LEVEL=info
```

### Run Production Container

```bash
docker run -d \
  --name carro-ideal \
  -p 80:8080 \
  -e DB_HOST=postgres.prod \
  -e DB_USER=app_user \
  -e DB_PASSWORD=strong_password \
  -e SESSION_SECRET=your-secret-key \
  -e ENV=production \
  carro-ideal:latest
```

---

## Additional Resources

- **Go Documentation**: https://golang.org/doc
- **PostgreSQL Documentation**: https://www.postgresql.org/docs/
- **Docker Documentation**: https://docs.docker.com/
- **Bootstrap 5 Documentation**: https://getbootstrap.com/docs/5.0/
- **Project Specification**: [spec.md](spec.md)
- **Implementation Plan**: [plan.md](plan.md)
- **API Contract**: [contracts/api.md](contracts/api.md)
- **Data Model**: [data-model.md](data-model.md)

---

## Support & Issues

For issues or questions:

1. Check this Quickstart guide
2. Review [data-model.md](data-model.md) for schema questions
3. Review [contracts/api.md](contracts/api.md) for API questions
4. Check application logs: `docker-compose logs app`
5. Check database logs: `docker-compose logs postgres`

---

**Quickstart Status**: ✅ COMPLETE - Ready for Implementation