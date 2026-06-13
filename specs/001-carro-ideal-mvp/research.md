# Research & Technical Clarifications: Carro Ideal MVP

**Date**: 2 de junho de 2026
**Status**: Complete - All clarifications resolved
**Branch**: `001-carro-ideal-mvp`

---

## Research Phase: Key Decisions & Rationale

### 1. Go Version & Dependencies

**Decision**: Go 1.22+ (latest stable)

**Rationale**:
- Go 1.22 introduced range-over-int and improved iterators, beneficial for querystring parsing
- Mature ecosystem for web development (Chi router, GORM, Postgres driver)
- Strong type system reduces runtime errors in TCC evaluation
- Fast compilation and single binary deployment fit Docker/TCC timeline

**Alternatives Considered**:
- Go 1.19: Too old, missing recent optimizations
- Go 1.21: Still viable but 1.22 has better features

**Migration Path**: None (greenfield project)

---

### 2. Web Framework & Router

**Decision**: Chi for HTTP routing + stdlib http.Handler interface

**Rationale**:
- Chi is lightweight, idiomatic Go, follows stdlib patterns
- Excellent middleware support (logging, auth, error handling)
- Easy to test (handlers are plain functions accepting http.ResponseWriter)
- No hidden magic; clear control flow for TCC defense
- Battle-tested in production (Hashicorp, Stripe use similar patterns)

**Alternatives Considered**:
- Gin: Popular but adds abstraction layer; heavier than Chi
- Echo: Similar to Gin; good but requires learning custom context
- stdlib only: Possible but more boilerplate for middleware

**Primary Dependencies**:
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/golang-migrate/migrate/v4` - Database migrations
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/jmoiron/sqlc` - OR `database/sql` manually (need clarity: sqlc recommended for type safety)
- `golang.org/x/crypto` - bcrypt password hashing
- `github.com/urfave/cli/v2` - CLI tool for admin operations if needed

---

### 3. Database & ORM Strategy

**Decision**: PostgreSQL 15+ with golang-migrate for migrations + database/sql with prepared statements (no ORM initially)

**Rationale**:
- Explicit SQL is easier to defend academically (shows DB knowledge)
- Avoids ORM lock-in; teaches database design principles
- PostgreSQL JSONB for flexible score_profile schema (AnswerOption.score_profile)
- golang-migrate provides up/down SQL files (meets TCC documentation requirements)
- Future: Can add GORM if complexity increases, but MVP should be explicit

**Alternatives Considered**:
- GORM: Fast development but less transparent for academic evaluation
- sqlc: Type-safe but adds build step; consider for Phase 2
- Raw sql.Rows: Too verbose; prepared statements + structured scanning better

**Schema Approach**:
- 7 tables: users, vehicle_categories, vehicles, questions, answer_options, user_answers, recommendations, recommendation_items
- All tables include created_at, updated_at (auto-timestamped)
- Soft-delete via `active` boolean field (no cascade deletes)
- Foreign keys with ON DELETE RESTRICT (prevent orphaning)

---

### 4. Session Management & Authentication

**Decision**: Session-based auth with HTTP-only cookies (not JWT)

**Rationale**:
- Simpler for MVP: store session ID server-side, no token expiration complexity
- HTTP-only cookies prevent XSS attacks (secure by default)
- Aligns with TCC timeline: less infrastructure (no token refresh logic)
- Easier to debug for professor (visible in browser DevTools)
- Standard golang.org/x/sessions or custom session store (Redis or in-memory for MVP)

**Alternatives Considered**:
- JWT: Stateless but adds token refresh complexity; overkill for MVP
- OAuth2: Unnecessary for internal TCC demo

**Implementation Details**:
- Session ID generated via crypto/rand
- Store in map (in-memory) for MVP; upgrade to Redis if needed
- 24-hour session expiry
- Secure flag on cookie (HTTPS in production)

---

### 5. Recommendation Algorithm

**Decision**: Weighted scoring based on user answers → vehicle attribute matches

**Rationale**:
- Each AnswerOption has score_profile JSON (e.g., `{"budget_score": 0.7, "efficiency": 0.3}`)
- Vehicle attributes (price_range, fuel_type, consumption) map to score dimensions
- Calculate composite score: sum of (answer_score_weight * vehicle_attribute_match)
- Ranking: sort vehicles by score descending, return top 10
- Demonstrates algorithmic thinking for TCC

**Score Calculation Example**:
```
User answers:
  - Budget: "R$ 50k-80k" → score_profile: {"budget": 0.4}
  - Fuel priority: "Efficiency" → score_profile: {"efficiency": 0.6}

Vehicle: Toyota Corolla
  - Price: 75k → budget match = 0.8 (in range) * 0.4 = 0.32
  - Consumption: 12 km/L → efficiency match = 0.7 * 0.6 = 0.42
  - Final score: 0.32 + 0.42 = 0.74 → 74%
```

**Alternatives Considered**:
- Machine learning: Overkill for MVP; requires data history
- Rule-based: Too rigid; score-based more flexible
- Future ChatGPT: Architecture prepared; implementation in Phase 2+

---

### 6. Frontend Architecture

**Decision**: Server-side template rendering (html/template) + Bootstrap 5 + minimal AJAX

**Rationale**:
- Go's html/template prevents XSS by default
- Bootstrap 5 CDN reduces build complexity
- Server-side rendering simpler for one-person TCC project
- Minimal AJAX for questionnaire (form submit → server → redirect)
- No Node.js build pipeline needed
- Easier to deploy: single binary includes templates and static files

**Alternatives Considered**:
- Single Page App (React/Vue): Overkill; adds Node build complexity
- Pure HTML: Boring for demonstration; Bootstrap improves UX

**Structure**:
- `/web/templates/` - Go template files organized by domain (auth, user, admin, web)
- `/web/static/` - CSS (Bootstrap CDN links), JS (minimal, maybe form validation)
- Template layout pattern: base.html + domain-specific partials

---

### 7. Docker & Deployment

**Decision**: Multi-stage Dockerfile + docker-compose for local dev + environment variables for config

**Rationale**:
- Build stage: compile Go binary (small final image)
- Runtime stage: base OS + binary + ca-certificates
- docker-compose for: app container + PostgreSQL container + volume for DB persistence
- Environment variables: DB connection, session secret, port
- Meets TCC requirement: `docker-compose up` should work immediately

**Development Workflow**:
```bash
docker-compose up -d       # Start app + PostgreSQL
go run ./app/main.go       # Local development (or docker-compose up in foreground)
docker-compose down -v     # Clean up including volumes
```

---

### 8. Testing Strategy

**Decision**: Unit tests for services; integration tests for API endpoints; repository tests with test database

**Rationale**:
- Unit tests for business logic (recommendation scoring, validation)
- Integration tests for API (handler → service → repository flow)
- Use database/sql with TestDB for repository tests (real queries, rollback after each test)
- TCC evaluation values: high code coverage demonstrates rigor

**Tools**:
- `testing` package (stdlib)
- `testify` for assertions if needed
- Test database: separate PostgreSQL instance or SQLite in-memory
- Table-driven tests for comprehensive scenarios

---

### 9. Error Handling & Logging

**Decision**: Structured logging with slog (Go 1.21+) + custom error types + HTTP status code mapping

**Rationale**:
- slog provides structured logging (JSON output for production, text for dev)
- Custom error types (e.g., ValidationError, NotFoundError) enable proper HTTP status mapping
- Middleware layer converts errors to appropriate HTTP responses
- Prevents information leakage (no internal stack traces to client)

**Error Response Format**:
```json
{
  "success": false,
  "error": "Email already registered",
  "code": "EMAIL_EXISTS"
}
```

---

### 10. Database Migrations

**Decision**: golang-migrate with numbered SQL files (up/down pattern)

**Rationale**:
- Up files (*.up.sql): Schema changes (CREATE, ALTER)
- Down files (*.down.sql): Rollback (DROP, RESTORE)
- Numbered sequentially: 0001, 0002, etc.
- Migrations versioned in git; repeatable and auditable
- TCC requirement: schema design must be documented

**Migration Approach**:
- One migration per logical change (e.g., "create_users_table", "add_soft_delete_to_vehicles")
- Down migrations are pure reversals (no data loss during dev)
- In production: apply migrations before app startup

---

### 11. API Versioning Strategy

**Decision**: No versioning prefix for MVP (all endpoints under `/api/...`); if needed, add `/api/v1/` later

**Rationale**:
- TCC deliverable unlikely to have breaking changes after launch
- If changes needed: support both v1 and v2 endpoints during transition
- Semantic Versioning for code releases (tags: v0.1.0, v0.2.0, etc.)

---

### 12. Configuration Management

**Decision**: Environment variables for secrets + config.go for structured config

**Rationale**:
- Environment variables for: DB connection, session secret, port, environment (dev/prod)
- Struct in config.go loaded from env at startup
- No hardcoded secrets in code
- .env file for local dev (excluded from git)
- Matches 12-factor app principles

**Key Config Values**:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `SESSION_SECRET`
- `PORT` (default 8080)
- `ENV` (dev/prod; controls logging, CORS, etc.)

---

## Technical Context Summary

| Aspect | Decision |
|--------|----------|
| **Language** | Go 1.22+ |
| **Framework** | Chi v5 (HTTP router) |
| **Database** | PostgreSQL 15+ |
| **ORM** | None (database/sql + prepared statements) |
| **Migrations** | golang-migrate |
| **Authentication** | Session-based HTTP-only cookies |
| **Frontend** | Server-side templates (html/template) + Bootstrap 5 |
| **Containerization** | Docker + docker-compose |
| **Testing** | stdlib testing + testify |
| **Logging** | slog (structured) |
| **API Format** | REST + JSON |
| **Error Handling** | Custom error types + middleware translation |
| **Configuration** | Environment variables |

---

## Resolved Clarifications

✅ **Go version clarified**: 1.22+ chosen (production-ready, modern features)
✅ **Framework selected**: Chi for routing (lightweight, testable, idiomatic)
✅ **Database approach determined**: PostgreSQL 15+ with explicit SQL + golang-migrate
✅ **Authentication method settled**: Session-based (simpler than JWT for MVP)
✅ **Frontend decided**: Server-side templates + Bootstrap 5 (no complex build pipeline)
✅ **Recommendation algorithm**: Weighted scoring with JSONB score_profile
✅ **Testing strategy**: Unit + Integration tests with structured approach
✅ **Error handling standardized**: Custom error types + HTTP status mapping
✅ **Deployment packaged**: Docker + docker-compose for reproducible environments
✅ **Configuration externalized**: Environment variables + config.go

---

## Next Phase: Design & Data Models

With these research decisions locked in, Phase 1 will produce:
1. **data-model.md** - Detailed PostgreSQL schema with relationships
2. **contracts/api.md** - Complete REST API endpoint documentation
3. **quickstart.md** - Instructions to build, run, and deploy locally

