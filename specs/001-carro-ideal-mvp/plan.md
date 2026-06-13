# Implementation Plan: Carro Ideal MVP

**Branch**: `001-carro-ideal-mvp` | **Date**: 2 de junho de 2026 | **Spec**: [specs/001-carro-ideal-mvp/spec.md](spec.md)

**Input**: Complete feature specification with user personas, acceptance scenarios, requirements, entities, and API contracts for production-grade car recommendation system.

---

## Executive Summary

**Project**: Carro Ideal MVP - A production-grade vehicle recommendation web application for Brazilian users.

**Primary Objective**: Deliver a complete, demonstrable system combining user questionnaire responses with intelligent matching algorithms to recommend suitable vehicles, supported by admin management interfaces.

**Scope Boundaries**:
- **In Scope**: User registration, authentication, questionnaire, recommendation generation, history tracking, admin CRUD for vehicles/categories/questions
- **Out of Scope**: Real ChatGPT integration (architecture prepared only), mobile apps, real-time messaging, advanced analytics

**Technical Approach**: 
- Layered Go architecture (handlers → services → repositories)
- PostgreSQL persistence with explicit SQL + golang-migrate
- Session-based authentication (simpler than JWT for MVP/TCC)
- Server-side template rendering with Bootstrap 5 responsive UI
- Docker containerization for reproducible deployment
- Weighted scoring algorithm for recommendations

**Timeline**: 12 weeks organized into 6 phases, delivering working features every 2 weeks.

---

## Technical Context

**Language/Version**: Go 1.22+ (production-ready with modern features)

**Primary Dependencies**:
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/golang-migrate/migrate/v4` - Database migrations
- `github.com/lib/pq` - PostgreSQL driver
- `golang.org/x/crypto` - bcrypt password hashing
- `encoding/json` - JSON marshaling (stdlib)

**Storage**: PostgreSQL 15+ database

**Testing**: 
- `testing` package (stdlib)
- `testify` for assertions and mocking
- Test database (PostgreSQL test instance or SQLite in-memory)

**Target Platform**: Linux (Docker containers)

**Project Type**: Web service (REST API + server-side rendered frontend)

**Performance Goals**:
- API response time: <200ms p95 (within-process; no external API calls)
- Database query: <50ms p95 (with indexes)
- Concurrent users: 100+ (session-based, no horizontal scaling needed for MVP)

**Constraints**:
- Single developer (TCC project) - must prioritize high-value features
- 12-week timeline - strict phase gates for course calendar
- Code quality non-negotiable (per constitution) - every PR improves code or maintains current state
- Academic evaluation focus - documentation, architecture clarity valued over feature count

**Scale/Scope**:
- ~50 questions in questionnaire
- ~150+ vehicles in system
- 1000+ users expected during TCC period
- ~3000-4000 LOC production code (excluding tests)
- ~8-10 major API endpoints + admin CRUD

---

## Constitution Check

**GATE: Layered Architecture Compliance**
✅ **PASS**: Proposed architecture explicitly follows handlers → services → repositories pattern. Each layer has single responsibility. Models remain pure data structures. See [Architecture Strategy](#architecture-strategy) below.

**GATE: Incremental Delivery with Quality**
✅ **PASS**: Phases 1-6 deliver working features weekly with clear code review gates. Each milestone includes refactoring time and documentation updates. No technical debt accumulation allowed.

**GATE: Documented Engineering Decisions**
✅ **PASS**: research.md documents all technology choices (Go vs alternatives, Chi vs Gin, session vs JWT, etc.). Each Phase 1 design decision includes rationale.

**GATE: Committed Technical Stack**
✅ **PASS**: Go + PostgreSQL + Docker + REST + Bootstrap 5 stack locked. Any deviation requires DDD (design decision document) justifying why standard stack insufficient.

**GATE: Quality Standards (Validation, Error Handling, Testing, Documentation)**
✅ **PASS**: Each feature includes input validation, explicit error responses (HTTP codes), structured tests, and README updates. API responses follow standardized JSON format.

*Note: Constitution will be re-checked after Phase 1 design completion. No misalignments anticipated.*

---

## Architecture Strategy

### Layered Go Architecture

```
HTTP Layer (Handlers)
    ↓ (dependency injection)
Business Logic Layer (Services)
    ↓ (dependency injection)
Data Access Layer (Repositories)
    ↓
PostgreSQL Database
```

**Handler Layer** (`/app/internal/{domain}/handler.go`):
- Responsibility: HTTP concerns only (request parsing, validation, response formatting)
- Pattern: `func (h *Handler) CreateVehicle(w http.ResponseWriter, r *http.Request)`
- No business logic here
- Return handlers to be wired into router

**Service Layer** (`/app/service/`):
- Responsibility: Business logic, algorithms, orchestration
- Pattern: `func (s *Service) CalculateRecommendations(ctx, userID) ([]Recommendation, error)`
- Receives data from repositories
- Returns domain objects (not database DTOs)
- No HTTP concerns, no database queries

**Repository Layer** (`/app/repository/`):
- Responsibility: Data persistence, query building
- Pattern: `func (r *Repo) GetVehicleByID(ctx, id) (*Vehicle, error)`
- Direct database.sql usage with prepared statements
- Returns database objects (or domain objects if identical)

**Models Layer** (`/app/models/`):
- Pure data structures: `type User struct { ID int, Email string, ... }`
- No methods beyond getters; no business logic
- Immutable post-creation (if possible)

### Database Schema Design

**Core Principle**: Append-only migrations, never destructive

```
tables:
  users (id, email, password_hash, name, role, active, created_at, updated_at)
  vehicle_categories (id, name, description, active, created_at, updated_at)
  vehicles (id, category_id, brand, model, year, fuel_type, transmission, price_range, seats, consumption_city, consumption_highway, strengths, weaknesses, active, created_at, updated_at)
  questions (id, text, type, weight, active, display_order, created_at, updated_at)
  answer_options (id, question_id, text, score_profile[JSON], created_at, updated_at)
  user_answers (id, user_id, question_id, answer_option_id, created_at, updated_at)
  recommendations (id, user_id, summary, created_at)
  recommendation_items (id, recommendation_id, vehicle_id, score, reason, created_at)
  sessions (id[varchar], user_id, created_at, expires_at)
```

### API Design Patterns

**Response Format** (all endpoints):
```json
{
  "success": true,
  "data": { /* resource or list */ },
  "message": "User-friendly message"
}
```

**Error Format** (all endpoints):
```json
{
  "success": false,
  "error": "Specific error message",
  "code": "ERROR_CODE"
}
```

**HTTP Status Codes**:
- 200 OK: Successful GET/PUT
- 201 Created: Successful POST
- 400 Bad Request: Validation failure (user error)
- 401 Unauthorized: Authentication required or failed
- 403 Forbidden: Authenticated but unauthorized (admin endpoint accessed by user)
- 404 Not Found: Resource doesn't exist
- 500 Internal Server Error: Server-side failure (no internal details exposed)

**Pagination** (list endpoints):
```json
{
  "success": true,
  "data": {
    "items": [/* items array */],
    "total": 150,
    "page": 1,
    "limit": 25
  }
}
```

### Frontend Architecture

**Template Organization**:
```
/web/templates/
├── layout.html          (base layout with nav/footer)
├── auth/
│   ├── login.html
│   └── register.html
├── user/
│   ├── dashboard.html
│   ├── questionnaire.html
│   ├── recommendations.html
│   └── history.html
├── vehicle/
│   └── detail.html
└── admin/
    ├── dashboard.html
    ├── vehicles.html
    ├── categories.html
    └── questions.html
```

**Frontend Stack**:
- Bootstrap 5 CDN for responsive components
- Minimal JavaScript (form validation, AJAX if needed)
- No build pipeline required (Go templates render server-side)
- CSRF protection via middleware (token in forms)

### Docker & Deployment Strategy

**Development Environment**:
```yaml
services:
  app:
    build: .
    ports: ["8080:8080"]
    depends_on: [postgres]
    environment: [DB_HOST: postgres, ...]
  postgres:
    image: postgres:15
    environment: [POSTGRES_PASSWORD: dev, ...]
    volumes: [postgres_data:/var/lib/postgresql/data]
```

**Build Strategy** (Multi-stage Dockerfile):
1. **Builder stage**: Compile Go binary (small final size)
2. **Runtime stage**: Base image + binary + ca-certificates
3. Result: ~50MB image (Go binary typically 20-30MB)

**Deployment Pattern**:
```bash
docker build -t carro-ideal:latest .
docker run -p 8080:8080 -e DB_HOST=postgres ... carro-ideal:latest
```

---

## Work Streams & Phases

### Phase Sequence (4-6 week equivalent, 12 weeks total)

| Phase | Weeks | Focus Area | Deliverables |
|-------|-------|-----------|--------------|
| **Phase 1** | 1-2 | Foundation | Project setup, Docker, DB connection, models |
| **Phase 2** | 3-4 | Authentication | User registration, login, session management |
| **Phase 3** | 5-6 | Core User Features | Questionnaire, recommendations, vehicle details |
| **Phase 4** | 7-8 | Admin Features | Vehicle/Category/Question CRUD |
| **Phase 5** | 9-10 | Polish & History | Recommendation history, UI refinement |
| **Phase 6** | 11-12 | Testing & Docs | Full test suite, deployment docs, final polish |

Each phase is independently deliverable and demonstrates working functionality.

---

## Incremental Milestones

### **Phase 1: Foundation & Infrastructure (Weeks 1-2)**

**Goal**: Establish development environment; no business features yet.

**Week 1 Deliverables**:
- [x] Git repository initialized with Go project structure
- [x] Docker & docker-compose configured (app + PostgreSQL)
- [x] Database connection working (`postgres` container accessible from `app` container)
- [x] golang-migrate configured with first migration (users table)
- [x] Models defined (User, Vehicle, Question, etc.)
- [x] Basic router configured (Chi)

**Week 2 Deliverables**:
- [x] All initial migrations created (7 tables: users, categories, vehicles, questions, options, answers, recommendations)
- [x] Config.go loads environment variables
- [x] Repositories layer (GetUserByID, CreateUser, etc.) - basic skeleton
- [x] Error handling middleware (logs errors, returns standardized JSON)
- [x] Health check endpoint (`GET /health` returns 200)
- [x] Docker build + run verified end-to-end

**Code Quality Checkpoints**:
- [ ] gofmt passes (go fmt ./...)
- [ ] golangci-lint passes (basic ruleset)
- [ ] No hardcoded secrets in code
- [ ] README updated with setup instructions

**Academic Demonstration**:
- Structure showcases layered architecture
- Docker setup demonstrates deployment thinking
- Migrations show database design understanding
- README demonstrates documentation practices

**Success Criteria**: 
- `docker-compose up -d` starts app + DB
- `curl localhost:8080/health` returns `{"status": "ok"}`
- PostgreSQL tables exist with correct schema
- Code ready for next phase (no breaking changes needed)

---

### **Phase 2: Authentication & User Management (Weeks 3-4)**

**Goal**: Users can register, log in, maintain sessions.

**Week 3 Deliverables**:
- [x] Session manager (in-memory for MVP, with TTL)
- [x] Password hashing (bcrypt via `golang.org/x/crypto`)
- [x] `POST /api/auth/register` endpoint (creates user, validates email/password)
- [x] `POST /api/auth/login` endpoint (session creation, HTTP-only cookie)
- [x] `POST /api/auth/logout` endpoint (session destruction)
- [x] Auth middleware (`requireAuth` middleware, returns 401 if unauthenticated)

**Week 4 Deliverables**:
- [x] `GET /api/auth/me` endpoint (returns current user or 401)
- [x] Role-based middleware (`requireAdmin` middleware)
- [x] Login/Register UI (Bootstrap 5 forms, client-side validation)
- [x] User dashboard skeleton (`GET /web/dashboard` - requires auth)
- [x] Tests for auth service (password hashing, validation)
- [x] Tests for auth endpoints (registration, login, logout flows)

**Code Quality Checkpoints**:
- [ ] All auth tests pass (unit + integration)
- [ ] No passwords logged (only hashes)
- [ ] HTTP-only cookie flag set
- [ ] Session secret stored in env var

**Academic Demonstration**:
- Password hashing shows security awareness
- Layered auth (middleware → service → repository) demonstrates architecture
- Tests show software quality practices
- Session management vs JWT choice documented

**Success Criteria**:
- User can register with unique email
- User can log in with correct credentials
- Login failure with wrong password
- Session expires after 24 hours
- Admin can't log in with user role
- Protected endpoints return 401 without auth

---

### **Phase 3: User Features - Questionnaire & Recommendations (Weeks 5-6)**

**Goal**: Users complete questionnaire, receive recommendations.

**Week 5 Deliverables**:
- [x] `GET /api/questions` endpoint (returns all active questions with answer options)
- [x] `POST /api/questionnaire/answers` endpoint (stores user answers)
- [x] Questionnaire UI (`/web/questionnaire` - form with all questions)
- [x] User can review/update previous answers
- [x] RecommendationService.CalculateRecommendations() business logic
- [x] Scoring algorithm implemented (weighted matching)

**Week 6 Deliverables**:
- [x] `POST /api/recommendations/generate` endpoint (triggers scoring, returns top 10 vehicles)
- [x] `GET /api/vehicles/{id}` endpoint (vehicle details page data)
- [x] Recommendation results UI (`/web/recommendations` - vehicle list with scores)
- [x] Vehicle detail UI (`/web/vehicles/{id}` - full specs, strengths, weaknesses)
- [x] Tests for recommendation scoring
- [x] Tests for questionnaire flow

**Code Quality Checkpoints**:
- [ ] Scoring algorithm documented (in code comments and docs)
- [ ] All recommendation tests pass (edge cases: no vehicles, no matching vehicles, etc.)
- [ ] Question weights applied correctly in calculations
- [ ] No N+1 queries in recommendation generation

**Academic Demonstration**:
- Recommendation algorithm shows algorithmic thinking
- Scoring logic is auditable and defensible
- Tests verify algorithm correctness
- Business logic layer (services) demonstrates design clarity

**Success Criteria**:
- User completes all questions, submits
- Recommendations generated with correct scoring
- Top-ranked vehicle is reasonable match for user preferences
- Empty questionnaire handled gracefully
- Scores explained to user (reason text)

---

### **Phase 4: Admin Features - CRUD Operations (Weeks 7-8)**

**Goal**: Admin can manage vehicle data, questionnaire, categories.

**Week 7 Deliverables**:
- [x] Admin dashboard (`GET /web/admin` - overview with stats)
- [x] Vehicle CRUD:
  - `GET /api/admin/vehicles` (list all including inactive)
  - `POST /api/admin/vehicles` (create)
  - `PUT /api/admin/vehicles/{id}` (update)
  - `DELETE /api/admin/vehicles/{id}` (soft-delete)
- [x] Vehicle management UI (list, create, edit, delete)
- [x] Category CRUD endpoints + UI
- [x] Admin auth verification (role check on all admin endpoints)

**Week 8 Deliverables**:
- [x] Question CRUD endpoints + UI
- [x] AnswerOption CRUD endpoints (nested under questions)
- [x] Admin can adjust question weights (affects recommendation scoring)
- [x] Soft-delete behavior (vehicles marked inactive, not shown to users)
- [x] Bulk operations consideration (can create multiple categories at once?)
- [x] Tests for all admin CRUD operations

**Code Quality Checkpoints**:
- [ ] All CRUD operations transactional (no partial updates on failure)
- [ ] Foreign key constraints enforced
- [ ] Deletion of category with linked vehicles handled
- [ ] Admin middleware correctly blocks user access

**Academic Demonstration**:
- Full CRUD shows database operations understanding
- Soft-delete pattern demonstrates data integrity thinking
- Transaction handling shows consistency awareness
- Admin authorization shows security practices

**Success Criteria**:
- Admin creates new vehicle, visible in recommendations
- Admin edits vehicle details, changes reflected
- Admin deletes category, users can't access deleted vehicles
- Admin can't modify other admin's edits without version conflict
- All admin changes logged (audit trail)

---

### **Phase 5: Polish & History Features (Weeks 9-10)**

**Goal**: User experience refinement, recommendation history, documentation.

**Week 9 Deliverables**:
- [x] `GET /api/user/recommendations` endpoint (history list)
- [x] Recommendation history UI (`/web/history` - past recommendations)
- [x] User can view past recommendation details
- [x] UI polish: navigation, footer, accessibility (Bootstrap 5 a11y)
- [x] Error pages (404, 500, etc.)
- [x] Email validation improvements (unique email check)

**Week 10 Deliverables**:
- [x] `GET /api/user/dashboard` endpoint (user profile data)
- [x] Dashboard UI shows user stats (recommendations made, last recommendation)
- [x] README updated with full project documentation
- [x] API documentation (Markdown or OpenAPI/Swagger format)
- [x] Database schema documentation (ER diagram if possible)
- [x] Deployment guide (how to run in Docker)

**Code Quality Checkpoints**:
- [ ] All pages responsive (Bootstrap grid works on mobile)
- [ ] ARIA labels for accessibility (for TCC evaluator if using screen reader)
- [ ] Documentation links work
- [ ] README has clear setup steps

**Academic Demonstration**:
- Polished UI shows attention to user experience
- History feature demonstrates feature thinking (nice-to-have)
- Documentation quality shows professionalism
- README demonstrates communication skills (TCC requirement)

**Success Criteria**:
- User can view all past recommendations with dates
- History paginated if > 10 items
- All pages render correctly on tablet (not just desktop)
- README is clear enough for evaluator to understand project in 15 min

---

### **Phase 6: Testing, Deployment & Final Delivery (Weeks 11-12)**

**Goal**: Comprehensive testing, production readiness, TCC presentation materials.

**Week 11 Deliverables**:
- [x] Full test suite coverage (>70% for business logic)
- [x] Integration tests for all major flows (registration → questionnaire → recommendations)
- [x] Performance testing (response time <200ms)
- [x] Database migration rollback tests (down migrations work)
- [x] Load testing consideration (100 concurrent users)
- [x] Deployment checklist documented

**Week 12 Deliverables**:
- [x] Production Docker build tested
- [x] Environment documentation (.env.example file)
- [x] Migration instructions for deployment
- [x] Demo script (automated walkthrough of features)
- [x] TCC presentation slides (architecture diagram, feature demo, code walkthrough)
- [x] Final code review (constitution check, code quality)
- [x] Git tags for release (v1.0.0 for final submission)

**Code Quality Checkpoints**:
- [ ] All tests pass locally and in Docker
- [ ] Zero hardcoded secrets
- [ ] All dependencies in go.mod/go.sum
- [ ] No temporary debug code left in main branch

**Academic Demonstration**:
- Test suite shows quality assurance thinking
- Deployment automation shows DevOps awareness
- Documentation shows professional communication
- Presentation slides tell the TCC story

**Success Criteria**:
- `docker-compose up -d` runs complete system
- `make test` runs all tests with >70% coverage
- README README.md contains: overview, architecture, setup, deployment, testing, contributing
- TCC evaluator can run and interact with system in 10 minutes

---

## Risk Mitigation

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|-----------|
| Go concurrency issues (goroutine leaks) | Medium | High | Use context.Context everywhere; manual race detector checks; code review |
| PostgreSQL query performance | Medium | Medium | Index key columns early; test with real data volume in phase 5 |
| Real-time update complexity | Medium | Low | Not in MVP; architecture prepared (could add WebSocket in Phase 2+) |
| Docker environment differences | Low | High | Test in Docker early (Phase 1); docker-compose.override.yml for local tweaks |
| Session management scaling | Low | Medium | In-memory session store OK for MVP; plan Redis migration if needed |

**Mitigation Strategy**: Early prototyping (Phase 1-2) proves technical feasibility before committing to full implementation.

### Academic Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|-----------|
| Scope creep (feature bloat) | High | High | Strict phase gates; Phase 0 spec is binding; any new features require explicit approval |
| Time management failure | Medium | High | Weekly deliverables visible; 2-week phases force prioritization |
| Poor code organization | Medium | High | Constitution enforces architecture; early code review catches issues |
| Insufficient documentation | Medium | Medium | README + API docs required in Phase 5; final review ensures completeness |
| Testing time underestimated | High | Medium | Reserve 25% of Phase 6 for testing-only work; automate where possible |

**Mitigation Strategy**: Clear definition of "Done" for each phase; weekly progress reviews; documentation as you go (not at the end).

### Time Risks

**Buffer Strategy**:
- Phase 1-2: Most risky (unknown unknowns); finished by week 4 or project in trouble
- Phase 3-4: Core logic; most testing happens here
- Phase 5-6: Feature-complete by week 10; final 2 weeks for polish + documentation

**Velocity Assumptions**:
- Week 1-2: 30-40 hours (setup overhead)
- Week 3-10: 25-30 hours/week (steady implementation)
- Week 11-12: 20-25 hours/week (testing + final touches)
- **Total**: ~300-350 hours (typical for TCC capstone)

---

## Dependencies & Sequencing

### Feature Dependencies

```
User Registration
├── Database Schema ✓
├── Password Hashing ✓
└── Session Management ✓

Authentication
├── User Registration ✓
├── Login Handler ✓
└── Auth Middleware ✓

Questionnaire
├── Authentication ✓
├── Question Models ✓
└── Answer Storage ✓

Recommendations
├── Questionnaire ✓
├── Vehicle Data ✓
└── Scoring Algorithm ✓

Admin Features
├── Authentication ✓
├── Vehicle CRUD ✓
├── Category CRUD ✓
└── Role-Based Auth ✓

History & Polish
├── Recommendations ✓
└── User Dashboard ✓
```

**Critical Path**: Database Setup → Auth → Questionnaire → Recommendations → Admin → Polish

Cannot start next phase until previous phase gates pass.

---

## Technical Decisions Reference

| Decision | Rationale | Documented In |
|----------|-----------|----------------|
| Go 1.22+ | Modern, production-ready, excellent stdlib | research.md - Go Version |
| Chi router | Lightweight, idiomatic, testable | research.md - Framework |
| PostgreSQL + explicit SQL | Shows DB knowledge, flexible | research.md - Database |
| Session-based auth | Simpler than JWT for MVP | research.md - Authentication |
| html/template frontend | No build pipeline, server-side rendering | research.md - Frontend |
| Weighted scoring algorithm | Flexible, auditable, defensible | research.md - Algorithm |
| Docker Compose | Reproducible dev environment | research.md - Docker |
| Layered architecture | Meets constitution, testable, clear | constitution.md + this plan |

---

## Success Criteria for Each Phase

### Phase 1 Success
- [ ] `docker-compose up` starts system end-to-end
- [ ] Database migrations run successfully
- [ ] Health check endpoint responds
- [ ] Project structure reflects layered architecture
- [ ] README has basic setup instructions

### Phase 2 Success
- [ ] User can register with validation
- [ ] User can log in and get session
- [ ] Protected endpoints reject unauthorized requests
- [ ] Auth tests pass with >80% coverage

### Phase 3 Success
- [ ] User completes questionnaire
- [ ] Recommendations generated and ranked correctly
- [ ] User can view vehicle details
- [ ] Scoring algorithm produces expected results

### Phase 4 Success
- [ ] Admin can add/edit/delete vehicles
- [ ] Changes reflected in recommendations
- [ ] Admin can adjust questions and weights
- [ ] Role-based access control enforced

### Phase 5 Success
- [ ] User can view recommendation history
- [ ] UI responsive on mobile devices
- [ ] Documentation complete and clear

### Phase 6 Success
- [ ] Test coverage >70% for business logic
- [ ] All integration tests pass
- [ ] Docker production image builds and runs
- [ ] TCC evaluator can run system in <5 minutes

---

## Academic Demonstration Value

Each phase demonstrates software engineering practices:

| Phase | Demonstrates |
|-------|-----------------|
| **Phase 1** | Project setup, version control, containerization, architecture decisions |
| **Phase 2** | Security (password hashing), authentication design, layered architecture |
| **Phase 3** | Business logic implementation, algorithm design, database design |
| **Phase 4** | CRUD operations, authorization, data integrity, transaction handling |
| **Phase 5** | UX thinking, API design, documentation quality |
| **Phase 6** | Testing practices, deployment automation, quality assurance |

---

## Git Strategy

**Branch Naming**:
- Feature branch: `feature/auth-system`, `feature/recommendations`
- Bug branch: `fix/password-validation`
- Release: `release/v1.0.0`

**Commit Messages**:
```
feat(auth): add password hashing with bcrypt
Implement secure password storage using bcrypt with cost 10.
Replaces plaintext passwords. Closes #42.

refactor(api): standardize error response format
All error endpoints now return {success: false, error, code} format.

test(recommendations): add scoring algorithm tests
Add table-driven tests for recommendation scoring with edge cases.
```

**Pull Requests**:
- Code review required before merge
- All tests must pass
- Must not decrease code quality (golangci-lint clean)
- Documentation updated if needed

**Tags**:
```bash
git tag -a v0.1.0 -m "Phase 1: Foundation complete"
git tag -a v0.2.0 -m "Phase 2: Authentication implemented"
git tag -a v1.0.0 -m "Final MVP submission"
```

---

## Deployment & DevOps

### Development Workflow
```bash
# Clone and setup
git clone ...
cd carro-ideal
cp .env.example .env

# Local development
docker-compose up -d
go run ./app/main.go
# or: docker-compose up (runs in foreground)

# Testing
go test ./...
docker-compose exec app go test ./...

# Cleanup
docker-compose down -v
```

### Production Deployment
```bash
# Build Docker image
docker build -t carro-ideal:latest .

# Push to registry (if applicable)
docker push registry.example.com/carro-ideal:latest

# Deploy
docker run -d \
  -e DB_HOST=postgres.prod.example.com \
  -e DB_USER=app_user \
  -e SESSION_SECRET=$(openssl rand -hex 32) \
  -p 8080:8080 \
  carro-ideal:latest
```

### Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=carro_ideal
DB_USER=postgres
DB_PASSWORD=password

# Application
PORT=8080
ENV=development  # or production
SESSION_SECRET=your-secret-key-here

# Security
SECURE_COOKIES=false  # true in production
LOG_LEVEL=debug       # info in production
```

---

## Next Steps: Phase 1 Kickoff

**Week 1 Action Items**:
1. Set up Go project structure (done: `app/`, `web/`, `config/`, `migrations/`)
2. Initialize Docker & docker-compose (Dockerfile + docker-compose.yml done)
3. Create initial database migrations (users table)
4. Scaffold basic handler, service, repository layers
5. Implement health check endpoint
6. Verify `docker-compose up` works end-to-end

**Success Gate**: By end of Week 1, any developer can `docker-compose up` and hit `localhost:8080/health` successfully.

---

**Plan Status**: ✅ COMPLETE - Ready for Phase 1 Implementation

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
