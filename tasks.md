# Tasks: Carro Ideal MVP - Car Recommendation System

**Status**: Ready for Implementation  
**Total Tasks**: 72 tasks across 6 phases  
**Estimated Duration**: 110-120 hours (single developer)  
**MVP Checkpoint**: After Phase 3 (complete working recommendations engine)

---

## Quick Navigation

- [Phase 1: Foundation & Infrastructure](#phase-1-foundation--infrastructure-weeks-1-2)
- [Phase 2: Authentication System](#phase-2-authentication-system-weeks-3-4)
- [Phase 3: Core Recommendation Engine](#phase-3-core-recommendation-engine-weeks-5-6)
- [Phase 4: Admin Management Panel](#phase-4-admin-management-panel-weeks-7-8)
- [Phase 5: Polish & Deployment Prep](#phase-5-polish--deployment-prep-weeks-9-10)
- [Phase 6: Testing, QA & Documentation](#phase-6-testing-qa--documentation-weeks-11-12)
- [Dependency Graph](#dependency-graph)
- [Parallel Execution Strategy](#parallel-execution-strategy)

---

## Format Reference

Each task uses the format: `- [ ] [ID] [P?] [Story?] Description`

- **[P]**: Task can be parallelized (independent files, no dependencies)
- **[Story]**: User Story ID (US1, US2, etc.) - only for feature phases
- **ID**: Sequential task identifier for dependency tracking
- Each task includes effort estimate in parentheses: XS (15min), S (30min), M (1h), L (2h), XL (3h)

---

## Phase 1: Foundation & Infrastructure (Weeks 1-2)

**Purpose**: Project setup, Docker configuration, database connection, and basic health check endpoint

**Goal**: Establish development environment with working database connection and deployment pipeline

**Checkpoint**: After Phase 1, developers can run `docker-compose up` and access a healthy API

### Phase 1 Dependencies

None - this is the blocking foundation for all subsequent work.

### Phase 1 Tasks

- [X] T001 Create GitHub repository with initial Go project structure (app/, config/, migrations/, web/) **[XS - 15min]**
  - **Acceptance**: Repository initialized with .gitignore for Go, README.md with project vision, initial commit
  
- [X] T002 Initialize Go module and add core dependencies to go.mod **[S - 30min]**
  - **Dependencies**: T001
  - **Acceptance**: `go mod tidy` succeeds; dependencies include chi/v5, golang-migrate/migrate/v4, lib/pq, golang.org/x/crypto
  
- [X] T003 [P] Create .env.example and environment configuration loader in config/config.go **[M - 1h]**
  - **Acceptance**: Config loads from .env; includes DATABASE_URL, PORT, ENVIRONMENT variables; panics with helpful message if required vars missing
  
- [X] T004 [P] Create Dockerfile for multi-stage Go binary build **[M - 1h]**
  - **Acceptance**: Builds single binary; no source code in final image; respects .dockerignore
  
- [X] T005 [P] Create docker-compose.yml with PostgreSQL 15 and Go app services **[M - 1h]**
  - **Acceptance**: `docker-compose up` starts both services; PostgreSQL initializes with empty database; networking configured
  
- [X] T006 Create app/main.go with basic HTTP server listening on configured port **[S - 30min]**
  - **Dependencies**: T001, T002, T003
  - **Acceptance**: Server starts; logs port on startup; graceful shutdown on SIGTERM
  
- [X] T007 [P] Create migrations directory structure with golang-migrate setup **[M - 1h]**
  - **Acceptance**: migrations/ contains numbered .up.sql and .down.sql files; migration runner integrated into app startup
  
- [X] T008 [P] Create database connection pool in app/db/connection.go using lib/pq **[L - 2h]**
  - **Acceptance**: Connection pool initialized with sensible defaults (max 25 connections); connection tested before server starts; error returns HTTP 503 if DB unavailable
  
- [X] T009 Create database health check handler (GET /health endpoint) **[M - 1h]**
  - **Dependencies**: T006, T008
  - **Acceptance**: Returns 200 OK with JSON {success: true, status: "healthy"} when DB accessible; returns 503 if DB down
  
- [X] T010 [P] Create 0001_create_users_table.up.sql migration with users table schema **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), email (varchar unique), password_hash (varchar), name (varchar), role (varchar default 'USER'), active (boolean default true), created_at, updated_at; UNIQUE index on email
  
- [X] T011 [P] Create 0001_create_users_table.down.sql migration **[S - 30min]**
  - **Acceptance**: DROP TABLE users CASCADE; no errors on revert
  
- [X] T012 Run migrations on Docker startup and verify first user table created **[S - 30min]**
  - **Dependencies**: T007, T008, T010, T011
  - **Acceptance**: `docker-compose up` runs migrations automatically; SELECT * FROM users returns empty result
  
- [X] T013 [P] Initialize Chi router in app/internal/routes.go with middleware structure **[L - 2h]**
  - **Acceptance**: Router created; middleware chain defined (logging, panic recovery, JSON responses); routes organized by domain (auth, user, admin, web)
  
- [X] T014 [P] Create basic error handler middleware responding with standardized JSON error format **[M - 1h]**
  - **Acceptance**: All errors return {success: false, error: "message", code: "ERROR_CODE"}; HTTP status codes accurate (400 user error, 500 server error)
  
- [X] T015 [P] Create JSON response wrapper middleware ensuring all responses {success: true, data: {...}} **[M - 1h]**
  - **Acceptance**: All handlers return wrapped responses; consistency across all endpoints
  
- [X] T016 Create comprehensive README.md with project vision, architecture diagram, setup instructions **[L - 2h]**
  - **Dependencies**: T006, T012, T013
  - **Acceptance**: README includes: project description, layered architecture diagram, local setup steps, Docker setup steps, available endpoints, how to run migrations, database diagram, deployment notes

### Phase 1 Summary

- **Total Tasks**: 16
- **Estimated Effort**: 20 hours
- **Parallel Tasks**: T003, T004, T005, T007, T008, T010, T011, T013, T014, T015 (can run simultaneously)
- **Critical Path**: T001 → T002 → T006 → T009 → T012
- **Output**: Working Docker development environment with database, health check endpoint, and structured routing

---

## Phase 2: Authentication System (Weeks 3-4)

**Purpose**: User registration, login, logout, session management, and authentication middleware

**Goal**: Complete user authentication system allowing users to create accounts, log in securely, and maintain sessions

**MVP User Stories**: US1 (Registration) & US2 (Session Management)  
**Blocks**: Phase 3 (all protected endpoints require authentication)

**Independent Test**: (1) Register new user → (2) Verify in database → (3) Login with credentials → (4) Access protected endpoint → (5) Logout → (6) Verify redirect to login

### Phase 2 Dependencies

- Blocked by: T012 (users table must exist)

### Phase 2 Tasks

- [ ] T017 [P] Create User model in app/models/user.go with fields matching database schema **[S - 30min]**
  - **Acceptance**: Struct includes ID, Email, Name, PasswordHash, Role, Active, CreatedAt, UpdatedAt; JSON tags for API responses
  
- [ ] T018 [P] Create UserRepository in app/repository/user_repo.go with CRUD methods **[L - 2h]**
  - **Acceptance**: Methods: GetByID, GetByEmail, Create, Update, Delete (soft); uses prepared statements; no SQL injection vulnerabilities
  
- [ ] T019 [P] Implement bcrypt password hashing in auth service (app/service/auth_service.go) **[M - 1h]**
  - **Acceptance**: HashPassword uses bcrypt cost=10; VerifyPassword compares hash correctly; never stores plaintext passwords
  
- [ ] T020 Create session management service in app/service/auth_service.go **[L - 2h]**
  - **Dependencies**: T017, T018, T019
  - **Acceptance**: GenerateSessionID creates cryptographically secure 32-byte random ID; StoreSession persists to in-memory map with 24-hour expiry; GetSession retrieves and validates expiry; DestroySession removes from map
  
- [ ] T021 Create authentication middleware (app/internal/middleware.go) checking session cookie **[L - 2h]**
  - **Dependencies**: T020
  - **Acceptance**: Validates session cookie presence; returns 401 if missing or expired; injects user into request context; skips auth for public endpoints (/health, /api/auth/login, /api/auth/register, /)
  
- [ ] T022 [P] Create user registration handler (POST /api/auth/register) **[L - 2h]**
  - **Acceptance**: Validates email format, password minimum 8 chars; checks email uniqueness; returns 400 with specific error if validation fails; creates user with hashed password; sets session cookie; returns 201 with user data
  
- [ ] T023 [P] Create user login handler (POST /api/auth/login) **[L - 2h]**
  - **Acceptance**: Validates email exists; verifies password hash; returns 401 with generic message if failed; creates session; sets cookie; returns 200 with user data and role
  
- [ ] T024 [P] Create logout handler (POST /api/auth/logout) **[M - 1h]**
  - **Acceptance**: Destroys session; clears session cookie; returns 200 OK; user cannot access protected endpoints after logout
  
- [ ] T025 [P] Create GET /api/auth/me endpoint returning current user **[M - 1h]**
  - **Acceptance**: Requires authentication; returns current user data; returns 401 if unauthenticated
  
- [ ] T026 Create session middleware with CSRF protection and cookie security flags **[M - 1h]**
  - **Dependencies**: T021
  - **Acceptance**: Cookie set with HttpOnly, Secure, SameSite=Lax flags; session stored server-side; 24-hour expiry enforced
  
- [ ] T027 Create registration form template (web/templates/register.html) **[M - 1h]**
  - **Acceptance**: HTML form with name, email, password, confirm password fields; Bootstrap 5 styling; client-side validation; submit POSTs to /api/auth/register
  
- [ ] T028 Create login form template (web/templates/login.html) **[M - 1h]**
  - **Acceptance**: HTML form with email, password fields; Bootstrap 5 styling; remember me checkbox (optional); submit POSTs to /api/auth/login
  
- [ ] T029 Create user dashboard template (web/templates/user/dashboard.html) **[M - 1h]**
  - **Dependencies**: T028
  - **Acceptance**: Displays logged-in user name; shows welcome message; links to questionnaire; displays recommendation history summary; logout button
  
- [ ] T030 Create frontend session verification script (web/static/js/session.js) **[M - 1h]**
  - **Acceptance**: On page load, calls GET /api/auth/me; if 401, redirects to login; stores user data in localStorage; updates navbar with user name
  
- [ ] T031 Create authentication routes handler (app/internal/auth/routes.go) **[S - 30min]**
  - **Dependencies**: T022, T023, T024, T025
  - **Acceptance**: Routes POST /api/auth/register, POST /api/auth/login, POST /api/auth/logout, GET /api/auth/me to handlers
  
- [ ] T032 Create web landing page (web/templates/index.html) with links to login/register **[M - 1h]**
  - **Acceptance**: Landing page displays project description; buttons for login/register visible; responsive Bootstrap 5 design

### Phase 2 Tests (Optional but recommended for TCC evaluation)

- [ ] T033 [P] [US1] Unit test UserRepository.Create validates unique email constraint **[M - 1h]**
  - **Acceptance**: Test in tests/unit/repository/user_repo_test.go; covers: valid user creation, duplicate email error, password hash stored not plaintext
  
- [ ] T034 [P] [US2] Unit test authentication middleware rejects expired sessions **[M - 1h]**
  - **Acceptance**: Test in tests/unit/middleware/auth_test.go; covers: valid session accepted, expired session rejected, missing session rejected
  
- [ ] T035 [P] [US1] Integration test registration flow end-to-end **[L - 2h]**
  - **Acceptance**: Test in tests/integration/auth_test.go; POST /api/auth/register → user created in DB → login succeeds

### Phase 2 Summary

- **Total Tasks**: 19 (16 implementation + 3 tests)
- **Estimated Effort**: 27 hours (24 implementation + 3 tests)
- **Parallel Tasks**: T017, T018, T019, T022, T023, T024, T025, T027, T028, T033, T034, T035 (can run in parallel)
- **Critical Path**: T017 → T018 → T020 → T021 → T022, T023, T024, T025
- **Output**: Complete authentication system ready for protected endpoints

---

## Phase 3: Core Recommendation Engine (Weeks 5-6)

**Purpose**: Questionnaire system, recommendation algorithm, and recommendation history

**Goal**: Complete working recommendation engine that generates vehicle recommendations based on user preferences

**MVP User Stories**: US3 (Questionnaire), US4 (Recommendations), US5 (Vehicle Details), US6 (History)  
**🎯 CHECKPOINT**: After Phase 3, MVP is functionally complete and demonstrable

**Independent Test**: (1) Login → (2) Answer questionnaire → (3) Generate recommendations → (4) View results with scoring → (5) View vehicle details → (6) Access history → (7) Compare recommendations

### Phase 3 Dependencies

- Blocked by: T032 (Phase 2 must complete first)

### Phase 3 Database Migrations

- [ ] T036 [P] Create 0002_create_questions_table.up.sql migration **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), question_text (text), question_type (varchar), weight (numeric), display_order (int), active (boolean), created_at, updated_at; INDEX on active and display_order
  
- [ ] T037 [P] Create 0002_create_questions_table.down.sql migration **[S - 30min]**
  
- [ ] T038 [P] Create 0003_create_answer_options_table.up.sql migration **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), question_id (fk to questions), option_text (text), score_profile (jsonb), display_order (int), active (boolean), created_at, updated_at
  
- [ ] T039 [P] Create 0003_create_answer_options_table.down.sql migration **[S - 30min]**
  
- [ ] T040 [P] Create 0004_create_user_answers_table.up.sql migration **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), user_id (fk), question_id (fk), answer_option_id (fk), created_at, updated_at; INDEX on user_id and created_at
  
- [ ] T041 [P] Create 0004_create_user_answers_table.down.sql migration **[S - 30min]**
  
- [ ] T042 [P] Create 0005_create_vehicles_table.up.sql migration **[L - 2h]**
  - **Acceptance**: Table includes: id, brand, model, version, year, fuel_type, transmission, seats, trunk_capacity, consumption_city, consumption_highway, price_range, category_id (fk), strengths (text), weaknesses (text), active (boolean), created_at, updated_at; INDEX on active and category_id
  
- [ ] T043 [P] Create 0005_create_vehicles_table.down.sql migration **[S - 30min]**
  
- [ ] T044 [P] Create 0006_create_vehicle_categories_table.up.sql migration **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), name (varchar unique), description (text), created_at, updated_at
  
- [ ] T045 [P] Create 0006_create_vehicle_categories_table.down.sql migration **[S - 30min]**
  
- [ ] T046 [P] Create 0007_create_recommendations_table.up.sql migration **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), user_id (fk), created_at; stores reference to recommendation set, user answers snapshot for reproducibility
  
- [ ] T047 [P] Create 0007_create_recommendations_table.down.sql migration **[S - 30min]**
  
- [ ] T048 [P] Create 0008_create_recommendation_items_table.up.sql migration **[M - 1h]**
  - **Acceptance**: Table includes: id (serial PK), recommendation_id (fk), vehicle_id (fk), rank (int), score (numeric), reason (text), created_at; stores individual vehicle recommendations with scores
  
- [ ] T049 [P] Create 0008_create_recommendation_items_table.down.sql migration **[S - 30min]**

### Phase 3 Models and Repositories

- [ ] T050 [P] Create Question, AnswerOption, UserAnswer models in app/models/ **[L - 2h]**
  - **Acceptance**: All models include JSON tags; AnswerOption.ScoreProfile is map[string]float64 (JSONB); no business logic in models
  
- [ ] T051 [P] Create Vehicle and VehicleCategory models in app/models/ **[M - 1h]**
  - **Acceptance**: Vehicle includes all attributes; Category has name and description; all JSON tagged
  
- [ ] T052 [P] Create Recommendation and RecommendationItem models **[M - 1h]**
  - **Acceptance**: Models include ID, UserID, Rank, Score, Reason, CreatedAt; ready for API responses
  
- [ ] T053 [P] Create QuestionRepository with GetActive, GetAll, GetByID methods **[L - 2h]**
  - **Acceptance**: Uses prepared statements; handles JSONB parsing for score profiles; efficient query ordering by display_order
  
- [ ] T054 [P] Create VehicleRepository with GetActive, GetByCategory, GetAll, GetByID methods **[L - 2h]**
  - **Acceptance**: GetActive returns only active vehicles; efficient filtering; join with categories for denormalized response
  
- [ ] T055 [P] Create RecommendationRepository with Create, GetByUser, GetByID, GetUserHistory methods **[L - 2h]**
  - **Acceptance**: Create stores recommendation with all recommendation_items in transaction; GetUserHistory paginated; efficient joins with vehicles
  
- [ ] T056 [P] Create UserAnswerRepository with Create, GetUserLatest, DeleteByRecommendation methods **[M - 1h]**
  - **Acceptance**: Create stores all answers in transaction; GetUserLatest retrieves most recent questionnaire answers
  
- [ ] T057 [P] Create seed data loader (app/db/seed.go) with sample questions and vehicles **[L - 2h]**
  - **Acceptance**: Creates ~20 sample vehicles with category, ~15 sample questions with answer options, all marked active

### Phase 3 Business Logic Services

- [ ] T058 Create QuestionnaireService in app/service/ **[M - 1h]**
  - **Acceptance**: GetActiveQuestions returns ordered questions with answer options; validates all questions answered; no database queries in service
  
- [ ] T059 Create RecommendationService with scoring algorithm in app/service/ **[L - 2h]**
  - **Acceptance**: GenerateRecommendations takes user answers, calculates weighted scores for each vehicle, returns top 10 sorted by score; algorithm explained in code comments with example
  
- [ ] T060 [P] Create VehicleService in app/service/ **[S - 30min]**
  - **Acceptance**: GetActiveVehicles, GetByCategory, GetByID; handles repository calls and error mapping
  
- [ ] T061 [P] Create RecommendationHistoryService **[M - 1h]**
  - **Acceptance**: GetUserHistory with pagination; GetRecommendationDetails; proper error handling for not found

### Phase 3 API Handlers

- [ ] T062 [P] Create GET /api/questions handler returning active questions with options **[M - 1h]**
  - **Acceptance**: Returns 200 with questions in order; empty array if none active; returns 401 if not authenticated
  
- [ ] T063 [P] Create POST /api/recommendations/generate handler **[L - 2h]**
  - **Acceptance**: Accepts user answers, validates all questions answered, generates recommendations, stores in DB, returns 201 with scored results
  
- [ ] T064 [P] Create GET /api/recommendations endpoint returning user history **[M - 1h]**
  - **Acceptance**: Returns paginated list of past recommendations; accepts page and limit query params; returns 200 with total count
  
- [ ] T065 [P] Create GET /api/recommendations/{id} endpoint returning recommendation details **[M - 1h]**
  - **Acceptance**: Returns recommendation with all recommendation_items sorted by rank; includes vehicle details and scores; returns 404 if not found or unauthorized
  
- [ ] T066 [P] Create GET /api/vehicles endpoint returning all active vehicles **[M - 1h]**
  - **Acceptance**: Supports category_id query param for filtering; returns with category info; paginated; returns 200 with metadata
  
- [ ] T067 [P] Create GET /api/vehicles/{id} endpoint returning vehicle details **[S - 30min]**
  - **Acceptance**: Returns complete vehicle with category; returns 404 if not found
  
- [ ] T068 Create recommendation routes handler (app/internal/api/routes.go) **[S - 30min]**
  - **Dependencies**: T062, T063, T064, T065, T066, T067
  - **Acceptance**: Routes all recommendation and vehicle endpoints to handlers

### Phase 3 Frontend Templates

- [ ] T069 Create questionnaire form template (web/templates/user/questionnaire.html) **[L - 2h]**
  - **Acceptance**: Displays all active questions; different input types for different question types (radio, checkbox, range slider); submit button POSTs to /api/recommendations/generate; responsive Bootstrap 5 design
  
- [ ] T070 [P] Create recommendations results page (web/templates/user/recommendations.html) **[L - 2h]**
  - **Acceptance**: Lists ranked vehicles with scores; color-coded score visualization; click to see vehicle details; back button returns to questionnaire
  
- [ ] T071 [P] Create vehicle detail modal/page (web/templates/user/vehicle_detail.html) **[L - 2h]**
  - **Acceptance**: Shows complete vehicle info; specifications organized in sections; reason explanation displayed; comparison link to similar vehicles
  
- [ ] T072 [P] Create recommendation history page (web/templates/user/history.html) **[M - 1h]**
  - **Acceptance**: Lists all user's past recommendations with dates; click to view details; comparison option available; shows vehicle count per recommendation

### Phase 3 Tests (Recommended)

- [ ] T073 [P] [US4] Unit test RecommendationService scoring algorithm **[L - 2h]**
  - **Acceptance**: Test in tests/unit/service/recommendation_test.go; verify score calculation with known inputs; edge cases (no matches, perfect match, tie-breaking)
  
- [ ] T074 [P] [US3] Integration test questionnaire submission flow **[M - 1h]**
  - **Acceptance**: Test in tests/integration/questionnaire_test.go; POST answers → stored in DB → retrievable in history
  
- [ ] T075 [P] [US4] Integration test recommendation generation full flow **[L - 2h]**
  - **Acceptance**: Complete flow from questionnaire to recommendation storage to retrieval; verify scores calculated correctly

### Phase 3 Summary

- **Total Tasks**: 40 (27 implementation + 3 tests + 10 database migrations)
- **Estimated Effort**: 55 hours
- **Parallel Tasks**: T036-T049 (migrations), T050-T057 (models/repos), T058-T061 (services), T062-T067 (handlers), T069-T072 (templates), T073-T075 (tests)
- **Critical Path**: T036-T049 → T050-T052 → T053-T056 → T058-T059 → T062-T063
- **Output**: Complete working MVP with questionnaire, recommendations, and history

---

## Phase 4: Admin Management Panel (Weeks 7-8)

**Purpose**: Admin CRUD operations for vehicles, categories, questions, and answers

**Goal**: Complete admin interface for system data management

**MVP User Stories**: US7 (Vehicle CRUD), US8 (Category CRUD), US9 (Question CRUD)  
**Blocks**: Phase 5 (admin features needed for demo)

**Independent Test**: (1) Login as admin → (2) Create vehicle → (3) Edit category → (4) Update question → (5) Verify reflected in user questionnaire → (6) Delete answer option

### Phase 4 Dependencies

- Blocked by: T072 (user features must be complete first)

### Phase 4 Models and Repositories

- [ ] T076 [P] Extend QuestionRepository with Create, Update, Delete methods **[M - 1h]**
  - **Acceptance**: Create validates question_text non-empty; Update supports partial updates; Delete marks inactive (soft delete)
  
- [ ] T077 [P] Extend AnswerOptionRepository with Create, Update, Delete methods **[M - 1h]**
  - **Acceptance**: Create validates score_profile JSONB format; Update allows score_profile modification; Delete marks inactive
  
- [ ] T078 [P] Extend VehicleRepository with Create, Update, Delete methods **[L - 2h]**
  - **Acceptance**: Create validates all required fields; Update supports partial updates; Delete marks inactive; proper category foreign key validation
  
- [ ] T079 [P] Extend VehicleCategory Repository with Create, Update, Delete methods **[M - 1h]**
  - **Acceptance**: Create validates unique category name; Update allows name and description changes; Delete soft delete if vehicles linked

### Phase 4 Business Logic Services

- [ ] T080 Create AdminVehicleService in app/service/ **[M - 1h]**
  - **Acceptance**: CreateVehicle validates all fields, creates with category link; UpdateVehicle partial updates; DeleteVehicle soft delete; GetAdminVehicles returns with stats
  
- [ ] T081 [P] Create AdminQuestionService **[M - 1h]**
  - **Acceptance**: Create/Update/Delete questions; GetAdminQuestions with metadata; ordering support
  
- [ ] T082 [P] Create AdminCategoryService **[S - 30min]**
  - **Acceptance**: CRUD operations for categories; GetAdminCategories with vehicle count; soft delete handling
  
- [ ] T083 [P] Create AdminAnswerOptionService **[M - 1h]**
  - **Acceptance**: CRUD for answer options; validates score_profile; linked question validation

### Phase 4 Admin Authorization

- [ ] T084 [P] Create admin role middleware in app/internal/middleware.go **[M - 1h]**
  - **Acceptance**: Checks user role = 'ADMIN'; returns 403 Forbidden if not admin; applied to all admin routes
  
- [ ] T085 [P] Create admin dashboard authorization check **[S - 30min]**
  - **Acceptance**: Only admins can access /admin/* routes; anonymous users redirected to login; regular users receive 403

### Phase 4 API Handlers

- [ ] T086 [P] Create POST /api/admin/vehicles handler for vehicle creation **[M - 1h]**
  - **Acceptance**: Validates all required fields; returns 201 with created vehicle; returns 400 with validation errors; admin only
  
- [ ] T087 [P] Create PUT /api/admin/vehicles/{id} handler for vehicle updates **[M - 1h]**
  - **Acceptance**: Supports partial updates; validates field types; returns 200 with updated vehicle; returns 404 if not found; admin only
  
- [ ] T088 [P] Create DELETE /api/admin/vehicles/{id} handler for soft delete **[S - 30min]**
  - **Acceptance**: Marks vehicle inactive; returns 200 OK; vehicle no longer appears in user recommendations; admin only
  
- [ ] T089 [P] Create GET /api/admin/vehicles endpoint listing all vehicles with stats **[M - 1h]**
  - **Acceptance**: Returns paginated list; includes recommendation count per vehicle; shows active/inactive status; admin only
  
- [ ] T090 [P] Create POST /api/admin/categories handler for category creation **[S - 30min]**
  - **Acceptance**: Validates unique name; returns 201; admin only
  
- [ ] T091 [P] Create PUT /api/admin/categories/{id} handler **[S - 30min]**
  - **Acceptance**: Updates name and description; returns 200; admin only
  
- [ ] T092 [P] Create DELETE /api/admin/categories/{id} handler **[S - 30min]**
  - **Acceptance**: Soft delete; handles linked vehicles; returns 200; admin only
  
- [ ] T093 [P] Create GET /api/admin/categories endpoint **[S - 30min]**
  - **Acceptance**: Lists all categories with vehicle count; includes inactive; admin only
  
- [ ] T094 [P] Create POST /api/admin/questions handler **[M - 1h]**
  - **Acceptance**: Validates question_text and weight; creates with active=true; returns 201; admin only
  
- [ ] T095 [P] Create PUT /api/admin/questions/{id} handler **[M - 1h]**
  - **Acceptance**: Updates all fields; returns 200; admin only
  
- [ ] T096 [P] Create DELETE /api/admin/questions/{id} handler **[S - 30min]**
  - **Acceptance**: Soft delete; removes from user questionnaire; returns 200; admin only
  
- [ ] T097 [P] Create GET /api/admin/questions endpoint **[M - 1h]**
  - **Acceptance**: Lists all questions with answer option count; includes inactive; admin only
  
- [ ] T098 [P] Create POST /api/admin/questions/{id}/options handler **[M - 1h]**
  - **Acceptance**: Creates answer option with score_profile JSONB; validates format; returns 201; admin only
  
- [ ] T099 [P] Create PUT /api/admin/questions/{id}/options/{optionId} handler **[M - 1h]**
  - **Acceptance**: Updates option_text and score_profile; returns 200; admin only
  
- [ ] T100 [P] Create DELETE /api/admin/questions/{id}/options/{optionId} handler **[S - 30min]**
  - **Acceptance**: Soft delete answer option; returns 200; admin only
  
- [ ] T101 Create admin routes handler (app/internal/admin/routes.go) **[S - 30min]**
  - **Dependencies**: T086-T100
  - **Acceptance**: Routes all admin CRUD endpoints with authorization middleware applied

### Phase 4 Frontend Admin Templates

- [ ] T102 Create admin dashboard template (web/templates/admin/dashboard.html) **[L - 2h]**
  - **Acceptance**: Displays key metrics (total users, vehicles, recommendations); links to all admin functions; responsive layout; Bootstrap 5
  
- [ ] T103 [P] Create vehicle management page (web/templates/admin/vehicles.html) **[L - 2h]**
  - **Acceptance**: Lists all vehicles in table; search/filter by brand; create/edit/delete buttons; inline stats; pagination
  
- [ ] T104 [P] Create vehicle form template (web/templates/admin/vehicle_form.html) **[L - 2h]**
  - **Acceptance**: Form with all vehicle fields; category dropdown; supports create and edit modes; form validation; submit POSTs to appropriate endpoint
  
- [ ] T105 [P] Create category management page (web/templates/admin/categories.html) **[M - 1h]**
  - **Acceptance**: Lists categories with vehicle count; create/edit/delete buttons; simple table layout
  
- [ ] T106 [P] Create category form template (web/templates/admin/category_form.html) **[M - 1h]**
  - **Acceptance**: Form with name and description fields; create/edit modes; validation
  
- [ ] T107 [P] Create question management page (web/templates/admin/questions.html) **[L - 2h]**
  - **Acceptance**: Lists questions with answer count; search capability; create/edit/delete buttons; display order editing
  
- [ ] T108 [P] Create question form template (web/templates/admin/question_form.html) **[L - 2h]**
  - **Acceptance**: Form with question_text, type, weight, display_order; nested form for answer options; create/edit modes; dynamic option field addition
  
- [ ] T109 Create admin authentication verification **[M - 1h]**
  - **Acceptance**: Frontend checks user role on page load; redirects non-admins to dashboard; admin nav menu only shows for admins

### Phase 4 Summary

- **Total Tasks**: 34 (26 implementation + 8 frontend)
- **Estimated Effort**: 42 hours
- **Parallel Tasks**: T076-T079 (repos), T080-T083 (services), T084-T085 (auth), T086-T100 (handlers), T102-T109 (templates)
- **Critical Path**: T076 → T080 → T086 → T102
- **Output**: Complete admin panel for data management

---

## Phase 5: Polish & Deployment Prep (Weeks 9-10)

**Purpose**: Refinement, optimization, and deployment readiness

**Goal**: Production-ready application with error handling, logging, and optimization

**User Story**: US10 (Admin Dashboard with stats) + cross-cutting concerns

### Phase 5 Dependencies

- Blocked by: T101 (admin routes complete)

### Phase 5 Infrastructure Tasks

- [ ] T110 [P] Implement structured logging with log levels (DEBUG, INFO, WARN, ERROR) **[M - 1h]**
  - **Acceptance**: All handlers log requests with method, path, status code, duration; sensitive data never logged; logs to stdout for Docker; different log levels for dev vs prod
  
- [ ] T111 [P] Implement global error recovery middleware (panic catcher) **[M - 1h]**
  - **Acceptance**: Catches all panics; logs stack trace; returns 500 with generic error message to client; no internal details exposed
  
- [ ] T112 [P] Create request ID middleware for tracing across logs **[M - 1h]**
  - **Acceptance**: Generates unique request ID; includes in all logs and responses; helps troubleshooting in production
  
- [ ] T113 [P] Implement CORS middleware with configurable origins **[M - 1h]**
  - **Acceptance**: Allows configured frontend origin; preflight requests handled; credentials included; secure by default
  
- [ ] T114 [P] Create rate limiting middleware (basic token bucket) **[L - 2h]**
  - **Acceptance**: Limits requests per IP; returns 429 Too Many Requests when exceeded; configurable limits; stores state in memory
  
- [ ] T115 [P] Implement input sanitization middleware preventing XSS **[M - 1h]**
  - **Acceptance**: Trims whitespace; HTML escapes user input; validates input types; applied to all POST/PUT handlers

### Phase 5 Performance Optimization

- [ ] T116 [P] Add database query caching for frequently accessed data **[L - 2h]**
  - **Acceptance**: Cache active vehicles, questions, categories with TTL; invalidate on updates; reduces DB hits by 50%+
  
- [ ] T117 [P] Implement pagination optimizations (offset/limit best practices) **[M - 1h]**
  - **Acceptance**: Recommendations list paginated with page size limit (max 100); index used for efficient queries
  
- [ ] T118 [P] Create database connection pool tuning documentation **[S - 30min]**
  - **Acceptance**: Document pool size rationale; benchmark concurrent connections; include tuning recommendations
  
- [ ] T119 [P] Optimize vehicle recommendation scoring algorithm for large datasets **[M - 1h]**
  - **Acceptance**: Algorithmic review; potential pre-calculation of certain scores; performance benchmark included

### Phase 5 Feature Completeness

- [ ] T120 [P] Create admin dashboard with system statistics (GET /api/admin/dashboard) **[L - 2h]**
  - **Acceptance**: Returns: total users, total vehicles, total recommendations, active users this week, new users this week
  
- [ ] T121 [P] Implement comparison view for multiple recommendations **[M - 1h]**
  - **Acceptance**: Side-by-side vehicle comparison; shows score differences; highlights key differences
  
- [ ] T122 [P] Create vehicle recommendation reason explanations with NLP-ready structure **[M - 1h]**
  - **Acceptance**: Reason text includes matching criteria; structure allows future AI enhancement; clear explanation of why vehicle matched
  
- [ ] T123 [P] Implement recommendation update flow (user can re-answer questionnaire) **[M - 1h]**
  - **Acceptance**: User can resubmit questionnaire; creates new recommendation entry; history preserved
  
- [ ] T124 Create seed data initialization on first app startup **[M - 1h]**
  - **Dependencies**: T057
  - **Acceptance**: Checks if vehicles table empty; if yes, runs seed script; idempotent (safe to run multiple times)

### Phase 5 Frontend Polish

- [ ] T125 [P] Implement responsive Bootstrap 5 design for mobile (max-width: 768px) **[L - 2h]**
  - **Acceptance**: All pages render correctly on mobile; touch-friendly buttons; readable text; no horizontal scrolling
  
- [ ] T126 [P] Create comprehensive form validation UI feedback **[M - 1h]**
  - **Acceptance**: Error messages appear below fields; success messages after submission; loading indicators during requests
  
- [ ] T127 [P] Add loading skeletons for slow API calls **[M - 1h]**
  - **Acceptance**: Skeleton screens displayed while recommendations loading; improves perceived performance
  
- [ ] T128 [P] Implement dark mode toggle (optional luxury feature) **[M - 1h]**
  - **Acceptance**: Dark mode CSS created; toggle saves preference; localStorage persists choice
  
- [ ] T129 Create frontend error boundary and graceful error display **[M - 1h]**
  - **Acceptance**: API errors display friendly messages; 404 shows "not found" page; 500 shows "try again" message

### Phase 5 Documentation

- [ ] T130 Create comprehensive API documentation with examples **[M - 1h]**
  - **Acceptance**: All endpoints documented with request/response examples; error codes explained; authentication requirements clear
  
- [ ] T131 [P] Create database schema documentation with ERD **[M - 1h]**
  - **Acceptance**: ERD diagram included; all tables documented; relationships explained; indexes noted
  
- [ ] T132 [P] Create recommendation algorithm documentation **[M - 1h]**
  - **Acceptance**: Algorithm explanation with examples; score calculation walkthrough; edge cases noted
  
- [ ] T133 [P] Create deployment guide (Docker, environment config, migrations) **[L - 2h]**
  - **Acceptance**: Step-by-step deployment instructions; production checklist included; troubleshooting guide
  
- [ ] T134 Create developer setup guide for local development **[M - 1h]**
  - **Acceptance**: Prerequisites listed; setup commands provided; common issues and solutions; testing instructions

### Phase 5 Security Hardening

- [ ] T135 [P] Implement HTTPS/TLS readiness (certificate configuration) **[S - 30min]**
  - **Acceptance**: Application supports TLS certificates; documented for production deployment
  
- [ ] T136 [P] Add SQL injection prevention audit **[M - 1h]**
  - **Acceptance**: All queries use prepared statements; no string concatenation in SQL; code review confirms safety
  
- [ ] T137 [P] Add XSS prevention audit (template escaping review) **[M - 1h]**
  - **Acceptance**: All user input in templates escaped; HTML template engine used properly; no raw HTML injection possible
  
- [ ] T138 [P] Implement CSRF token protection for state-changing requests **[M - 1h]**
  - **Acceptance**: CSRF tokens generated and validated; forms include token; returns 403 if token invalid
  
- [ ] T139 Create security documentation and best practices guide **[M - 1h]**
  - **Acceptance**: Security considerations documented; password requirements explained; session security explained

### Phase 5 Summary

- **Total Tasks**: 30 (infrastructure, optimization, features, frontend, docs, security)
- **Estimated Effort**: 35 hours
- **Parallel Tasks**: T110-T115 (middleware), T116-T119 (optimization), T120-T124 (features), T125-T129 (frontend), T130-T134 (docs), T135-T139 (security)
- **Critical Path**: Most tasks independent; T124 depends on T057
- **Output**: Production-ready application with documentation and security

---

## Phase 6: Testing, QA & Documentation (Weeks 11-12)

**Purpose**: Comprehensive testing, manual QA, and final documentation

**Goal**: Production-grade quality assurance and academic documentation

### Phase 6 Dependencies

- Blocked by: T134 (Phase 5 complete)

### Phase 6 Testing

- [ ] T140 [P] Create comprehensive unit tests for all service layer classes **[L - 2h]**
  - **Acceptance**: Tests in tests/unit/service/; >80% code coverage; mocked repositories; edge cases covered
  
- [ ] T141 [P] Create integration tests for database repository layer **[L - 2h]**
  - **Acceptance**: Tests in tests/integration/repository/; uses test database; verifies SQL correctness; transaction handling tested
  
- [ ] T142 [P] Create API contract tests for all public endpoints **[L - 2h]**
  - **Acceptance**: Tests in tests/integration/api/; verifies request/response schema; error responses correct; status codes accurate
  
- [ ] T143 [P] Create end-to-end user flow tests **[L - 2h]**
  - **Acceptance**: Tests: register → login → questionnaire → recommendations → history → logout; covers main happy path; all steps verified
  
- [ ] T144 [P] Create admin workflow tests **[M - 1h]**
  - **Acceptance**: Tests: admin login → vehicle CRUD → category management → question updates; authorization verified

### Phase 6 Manual QA

- [ ] T145 Create comprehensive manual testing checklist **[M - 1h]**
  - **Acceptance**: Checklist includes: registration flow, authentication, questionnaire, recommendations, history, admin CRUD, error scenarios, mobile responsiveness, cross-browser testing
  
- [ ] T146 [P] Execute manual QA on registration and authentication flows **[M - 1h]**
  - **Acceptance**: Document results; any bugs filed as issues; edge cases tested (duplicate email, weak password, expired session)
  
- [ ] T147 [P] Execute manual QA on questionnaire and recommendation flows **[M - 1h]**
  - **Acceptance**: Complete questionnaire multiple times; verify recommendations differ; verify history stores correctly; test with different answer combinations
  
- [ ] T148 [P] Execute manual QA on admin operations **[M - 1h]**
  - **Acceptance**: Create/update/delete all entity types; verify reflected in user interface; test authorization (non-admin access denied)
  
- [ ] T149 [P] Cross-browser testing (Chrome, Firefox, Safari) **[M - 1h]**
  - **Acceptance**: Application works on all major browsers; responsive design works on mobile/tablet/desktop; no console errors
  
- [ ] T150 [P] Performance testing and load testing **[L - 2h]**
  - **Acceptance**: Load 100 concurrent users; measure response times; identify bottlenecks; document results; optimize if needed

### Phase 6 Security Testing

- [ ] T151 [P] Manual security testing for XSS vulnerabilities **[M - 1h]**
  - **Acceptance**: Test input fields with <script> tags; verify output is escaped; verify no JavaScript execution
  
- [ ] T152 [P] Manual security testing for SQL injection **[M - 1h]**
  - **Acceptance**: Test with SQL commands in input fields; verify queries use prepared statements; no SQL execution
  
- [ ] T153 [P] Authentication security review **[S - 30min]**
  - **Acceptance**: Session tokens validated; expired sessions rejected; CSRF tokens checked; cookie flags secure

### Phase 6 Documentation Finalization

- [ ] T154 Create final comprehensive README.md for project submission **[L - 2h]**
  - **Acceptance**: README includes: project vision, features, technical stack, architecture diagram, setup instructions, API overview, database schema, deployment guide, testing results, TCC-specific documentation
  
- [ ] T155 [P] Create TCC project defense presentation outline **[L - 2h]**
  - **Acceptance**: Presentation deck includes: problem statement, solution architecture, key design decisions, implementation highlights, code examples, testing results, lessons learned
  
- [ ] T156 [P] Create project architecture documentation with diagrams **[M - 1h]**
  - **Acceptance**: Layered architecture diagram; database ERD; API flow diagram; sequence diagram for key user flows
  
- [ ] T157 [P] Create decision records for all major technical choices **[L - 2h]**
  - **Acceptance**: Document: Go choice, PostgreSQL choice, REST vs alternatives, session vs JWT, recommendation algorithm, architecture pattern
  
- [ ] T158 [P] Create performance analysis and optimization report **[M - 1h]**
  - **Acceptance**: Benchmark results for key operations; identify hotspots; document optimization strategies used
  
- [ ] T159 Create lessons learned document and future improvements list **[M - 1h]**
  - **Acceptance**: Reflect on challenges faced; document solutions; list potential Phase 2+ improvements (ChatGPT integration, analytics, etc.)

### Phase 6 Final Checks

- [ ] T160 [P] Verify all code formatting with gofmt **[S - 30min]**
  - **Acceptance**: No formatting issues; all .go files properly formatted
  
- [ ] T161 [P] Run linter (golangci-lint) and fix any warnings **[M - 1h]**
  - **Acceptance**: No lint warnings; documentation for any exceptions needed
  
- [ ] T162 [P] Clean up any dead code or commented-out code **[S - 30min]**
  - **Acceptance**: Codebase clean; no debugging artifacts; production-ready
  
- [ ] T163 [P] Verify all environment variables documented in .env.example **[S - 30min]**
  - **Acceptance**: .env.example complete; all required and optional vars explained
  
- [ ] T164 [P] Final Docker build and test in clean environment **[M - 1h]**
  - **Acceptance**: `docker-compose up` works; all migrations run; seed data loads; health check passes
  
- [ ] T165 Create final submission checklist and project completion summary **[M - 1h]**
  - **Acceptance**: Checklist covers all requirements; project demonstrated working; documentation complete; ready for defense

### Phase 6 Summary

- **Total Tasks**: 26 (5 unit tests + 7 QA + 3 security + 6 documentation + 5 final checks)
- **Estimated Effort**: 30 hours
- **Parallel Tasks**: T140-T144 (tests), T146-T150 (QA), T151-T153 (security), T154-T159 (docs), T160-T165 (final)
- **Critical Path**: Sequential in phases but parallel within phase
- **Output**: Fully tested, documented, production-ready TCC submission

---

## Summary Statistics

### Total Project Metrics

- **Total Tasks**: 72 (T001-T165)
- **Total Estimated Effort**: 110-120 hours
- **Number of Phases**: 6
- **Average Phase Duration**: 2 weeks
- **Average Tasks per Phase**: 12

### Effort Distribution by Phase

| Phase | Tasks | Effort | Focus |
|-------|-------|--------|-------|
| Phase 1 | 16 | 20h | Foundation |
| Phase 2 | 19 | 27h | Authentication |
| Phase 3 | 40 | 55h | MVP Features |
| Phase 4 | 34 | 42h | Admin Panel |
| Phase 5 | 30 | 35h | Polish |
| Phase 6 | 26 | 30h | Testing & Docs |
| **Total** | **72** | **110-120h** | **Complete** |

### Task Breakdown by Type

- **Infrastructure & Setup**: 16 tasks (22%)
- **Database & Migrations**: 17 tasks (24%)
- **Backend Logic & APIs**: 18 tasks (25%)
- **Frontend & Templates**: 12 tasks (17%)
- **Testing**: 5 tasks (7%)
- **Documentation**: 4 tasks (6%)

### Parallelization Opportunities

- **Phase 1**: 10 tasks can parallelize (T003-T005, T007-T008, T010-T011, T013-T015)
- **Phase 2**: 9 tasks can parallelize (T017-T019, T022-T025, T027-T028)
- **Phase 3**: 20+ tasks can parallelize (migrations, models, services, handlers, templates)
- **Phase 4**: 26 tasks can parallelize (most handlers, templates, services)
- **Phase 5**: 28 tasks can parallelize (infrastructure, optimization, frontend, docs, security)
- **Phase 6**: 20+ tasks can parallelize (unit tests, QA, docs)

**Single Developer Timeline**: 110-120 hours ≈ 6-7 weeks at 20h/week, fitting 12-week calendar perfectly with buffer

---

## Dependency Graph

### Phase-Level Dependencies

```
Phase 1 (Foundation)
    ↓
Phase 2 (Authentication)
    ↓
Phase 3 (Core MVP) ← This is MVP checkpoint
    ↓
Phase 4 (Admin) ← Runs parallel to Phase 5 possible
    ↓
Phase 5 (Polish)
    ↓
Phase 6 (Testing & QA)
```

### Critical Path (Blocking Sequence)

1. T001-T016 (Phase 1: Must complete before anything else)
2. T017-T032 (Phase 2: Auth required for protected endpoints)
3. T036-T068 (Phase 3: Migrations first, then models/services/handlers)
4. T076-T109 (Phase 4: Build on Phase 3)
5. T110-T139 (Phase 5: Polish and optimize)
6. T140-T165 (Phase 6: Final testing and docs)

### Independent Work Streams Within Phases

**Phase 3 Parallel Streams** (after T032 complete):
- Stream A: Migrations (T036-T049)
- Stream B: Models & Repos (T050-T056)
- Stream C: Services (T058-T061)
- Stream D: Handlers (T062-T067)
- Stream E: Templates (T069-T072)
- Stream F: Tests (T073-T075)

Streams can work in parallel with dependency: A → B → C → D (for implementation), E parallel to D, F parallel to D/E.

---

## Parallel Execution Strategy

### Week 1-2 (Phase 1) - Sequential, Single Developer

All tasks in Phase 1 are blocking. Recommend conservative sequential execution:

1. **Day 1-2**: T001-T002, T006 (repo, Go setup, main.go)
2. **Day 3-4**: T003-T005, T013-T015 (config, Docker, routing)
3. **Day 5**: T007-T012 (migrations and database)
4. **Day 6-7**: T009, T016 (health check, README)
5. **Day 8-10**: Buffer and testing

### Week 3-4 (Phase 2) - Parallel Possible

Recommend grouping by layer:

- **Developer**: T017-T021 (Models, Repos, Services)
- **Parallel**: T022-T025 (Handlers)
- **Parallel**: T027-T032 (Frontend templates)
- **Tests**: T033-T035 (unit and integration tests)

### Week 5-6 (Phase 3) - Maximum Parallelization Possible

**If two developers available**:
- Developer A: Migrations (T036-T049), Models (T050-T052)
- Developer B: Frontend (T069-T072)
- Both: Services (T058-T061), Handlers (T062-T067)

**Single developer approach**:
- Day 1-2: Migrations and Models (T036-T052)
- Day 3-4: Services (T058-T061)
- Day 5: Handlers (T062-T067)
- Day 6-7: Frontend (T069-T072)
- Day 8-10: Tests (T073-T075)

### Week 7-8 (Phase 4) - Same approach as Phase 3

Most tasks independent after T076. Recommend by entity type:

- **Vehicles**: T076-T080, T086-T089, T103-T104
- **Categories**: T079, T082, T090-T093, T105-T106
- **Questions**: T077-T078, T081, T094-T100, T107-T108

### Week 9-10 (Phase 5) - Maximum Parallelization

All categories independent:

- **Infrastructure**: T110-T115
- **Performance**: T116-T119
- **Features**: T120-T124
- **Frontend**: T125-T129
- **Docs**: T130-T134
- **Security**: T135-T139

Can be distributed across multiple developers or handled sequentially by single developer with high parallelization of reading/documentation work.

### Week 11-12 (Phase 6) - Testing & QA

**Recommended Single Developer Sequence**:
1. Week 11a: Unit tests (T140) + Integration tests (T141-T142)
2. Week 11b: E2E tests (T143-T144)
3. Week 12a: Manual QA (T145-T150)
4. Week 12b: Final docs and checks (T154-T165)

Tests can be written before implementation (TDD) for accelerated timeline.

---

## MVP Checkpoint (After Phase 3)

After completing Phase 3 (T001-T075), the system has:

- ✅ Complete user registration and authentication
- ✅ Questionnaire with multiple questions and answer options
- ✅ Recommendation generation with scoring algorithm
- ✅ Recommendation history tracking
- ✅ Vehicle detail pages
- ✅ Responsive Bootstrap 5 UI
- ✅ All working features demonstrated end-to-end
- ✅ Basic tests for critical flows

**Status**: Fully functional MVP ready for academic evaluation

**Time to MVP**: ~55 hours (approximately 3-4 weeks at 15-20h/week)

Phases 4-6 add:
- Admin management interface
- Performance optimizations
- Security hardening
- Comprehensive testing
- Complete documentation
- Deployment readiness

---

## Success Criteria by Task

Each task can be validated using the Acceptance Criteria provided in its definition. For example:

- **T001**: Check GitHub repo created with correct structure + .gitignore + initial README
- **T009**: Test `/health` endpoint returns 200 OK with JSON response when DB connected
- **T022**: Test POST /api/auth/register with valid data creates user and returns 201
- **T063**: Test POST /api/recommendations/generate returns scored results sorted by score

General success definition for all tasks:

✅ **Task Complete When**:
1. Code written and committed to feature branch
2. All Acceptance Criteria met and verified
3. No lint warnings (gofmt, golangci-lint pass)
4. Manual testing confirms behavior
5. Related tests passing (if applicable)
6. Code review approved before merge to main

---

## Notes for TCC Evaluation

This task plan demonstrates:

1. **Layered Architecture**: Every layer clearly separated (Handlers → Services → Repositories)
2. **Incremental Delivery**: Features added in prioritized phases with working code each cycle
3. **Quality Focus**: Testing, documentation, and code quality built-in throughout
4. **Academic Rigor**: Decision documentation, security considerations, performance analysis
5. **Professional Practices**: Git workflow, code review gates, CI/CD readiness
6. **Scalable Design**: Foundation supports future enhancements (ChatGPT integration, analytics)

Recommended presentation flow for TCC defense:
- Phase 1: Technical foundation (architecture, database, infrastructure)
- Phase 2: Security implementation (authentication, session management)
- Phase 3: Algorithm implementation (recommendation scoring)
- Phase 4: Admin capabilities (CRUD operations, authorization)
- Phase 5: Production readiness (performance, security, docs)
- Phase 6: Quality assurance (testing results, bug fixes)

---

**Document Created**: 2026-06-02  
**Version**: 1.0  
**Status**: Ready for Implementation  
**Total Duration**: 12 weeks (110-120 hours estimated)
