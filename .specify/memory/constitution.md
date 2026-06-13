<!-- SYNC IMPACT REPORT (Initial Constitution)
Version: 1.0.0 (initial)
Status: Ratified
Modified Principles: All 5 principles added (Architecture, Incremental Delivery, Academic Excellence, Technical Stack, Quality Standards)
Added Sections: Technical Constraints, Development Workflow & Documentation Standards, Governance
Templates Requiring Updates:
  - plan-template.md: Constitution Check section will enforce layered architecture + academic requirements ✅
  - spec-template.md: User scenarios must align with incremental delivery principle ✅
  - tasks-template.md: Task categorization per principles (architecture/quality/docs) ✅
Follow-up TODOs: None (all placeholders filled with concrete guidance)
-->

# Carro Ideal Constitution

A software engineering capstone project delivering a complete, production-grade car recommendation system. This constitution balances academic excellence with professional development practices.

## Core Principles

### I. Layered Architecture with Clear Separation of Concerns

**Non-negotiable Rule**: Every feature MUST follow Go layered architecture patterns: HTTP handlers → business services → repository layer → database. No layer dependencies may bypass or short-circuit this hierarchy.

**Specific Requirements**:
- Handlers layer (`/app/internal/{domain}/handler.go`, `user_handler.go`, etc.) → only HTTP concern, request validation, error translation
- Service layer (`/app/service/`) → business logic, algorithm implementation, service composition, NO database queries here
- Repository layer (`/app/repository/`) → ONLY database queries, data persistence, query builders
- Models layer (`/app/models/`) → pure data structures with no business logic, immutable domain boundaries
- Each layer MUST be independently testable with mocked dependencies

**Why This Matters**: Demonstrates architectural maturity for TCC evaluation; enables feature testing without infrastructure; improves maintainability as requirements evolve.

---

### II. Incremental Delivery with Code Quality as Non-Negotiable

**Non-negotiable Rule**: Every merged feature MUST maintain current code quality or improve it. Code MUST be cleaner after the feature than before. Incremental does NOT mean compromising on standards.

**Specific Requirements**:
- Each feature delivery includes: functionality + refactoring debt paydown + documentation updates
- Code review gates MUST verify: naming clarity, DRY principle compliance, error handling patterns, no duplicate repository queries
- Go formatting (gofmt) and lint (golangci-lint) MUST pass before merge—automatic enforcement via CI
- Incremental commits with clear commit messages: "feat: [layer] [what] (why)" pattern

**Why This Matters**: Delivers working features weekly while maintaining a codebase that looks intentional and well-crafted; TCC evaluators will see continuous improvement rhythm, not technical debt accumulation.

---

### III. Academic Excellence Through Documented Engineering Decisions

**Non-negotiable Rule**: Every significant architectural decision, technology choice, or deviation from standard patterns MUST be documented with decision rationale in `/docs/decisions/` before implementation.

**Specific Requirements**:
- Design Decision Documents (DDDs) REQUIRED for: new service creation, database schema changes, API contract changes, third-party dependency adoption
- Each DDD MUST include: context (problem), alternatives considered, chosen solution, trade-offs, migration path if applicable
- Code comments MUST explain WHY, not WHAT (what is visible in code; why requires explanation)
- Repository patterns and optimization choices MUST be justified (e.g., "eager loading used here because N+1 pattern would exceed performance SLA")

**Why This Matters**: Demonstrates systematic engineering practices; provides TCC defense material showing rigorous decision-making; allows future developers to understand intent, not just implementation.

---

### IV. Committed Technical Stack—No Substitutions Without Justification

**Immutable Stack**: Go 1.22+ backend, PostgreSQL 15+ database, Docker containers, REST API only, Bootstrap 5 frontend.

**Why These Choices**:
- Go: Type-safe, fast compilation, excellent concurrency primitives, clear error handling culture
- PostgreSQL: Reliability (ACID), JSON support for flexible schemas, proven at scale
- Docker: Reproducible development environment, consistent deployment, TCC-friendly (runs anywhere)
- REST: Stateless, cacheable, standard HTTP semantics, easy to document and test
- Bootstrap 5: Production-ready components, accessibility built-in, reduces CSS burden

**Deviation Rule**: Any technology substitution (e.g., "use gRPC instead of REST") MUST be documented in a DDD explaining why the standard stack insufficient for this specific use case.

---

### V. Quality Standards: Validation, Error Handling, Testing Mindset, Documentation

**Non-negotiable Rule**: Every feature MUST include validated input handling, explicit error responses, and usable documentation. No exceptions.

**Specific Requirements**:
- Input Validation: All HTTP handlers MUST validate incoming data before passing to services (schema, ranges, business rules)
- Error Handling: MUST return HTTP status codes truthfully (400 for user error, 500 for server error, never confuse the client)
- Error Messages: MUST be specific and actionable (not "error", but "email already registered" or "car year must be 1990 or later")
- Testing Mindset: Unit tests for business logic (services), integration tests for database queries, handler tests with mock services
- Documentation: README updates for new features, API endpoint descriptions in code comments, database schema documented in migrations
- Browser DevTools Friendliness: Errors logged to console, API responses inspectable, no secrets in frontend code

**Why This Matters**: Professional systems fail gracefully; TCC evaluators will recognize production thinking; users (including test users) get useful feedback instead of confusion.

---

## Technical Constraints & Guardrails

**Database Immutability**: Schema migrations are append-only (`.up.sql`) and fully reversible (`.down.sql`). Never destructive. Every migration file numbered sequentially.

**Containerization**: Application MUST run identically in Docker and locally. No "works on my machine" exceptions. `docker-compose up` MUST work end-to-end without configuration.

**API Stability**: Existing endpoints MUST NOT break between versions. Version bumps (Semantic Versioning) reflect contract changes. Deprecated endpoints return 410 Gone with migration instructions.

**Configuration**: Environment-specific config (database URL, API keys, port) loaded from `.env` via `config/config.go`. No hardcoded secrets, no environment-specific code branches.

---

## Development Workflow & Standards

**Incremental Workflow**:
1. **Feature Branch**: Create from main via `###-feature-name` naming convention (e.g., `001-user-authentication`)
2. **Specification**: Document user scenarios, requirements, and success criteria BEFORE coding
3. **Implementation**: Follow layers; commit incrementally; run tests locally
4. **Code Review**: Verify architecture compliance, naming clarity, no security gaps
5. **Merge**: All CI checks pass; author provides deployment notes

**Testing Gates**:
- Unit tests for services and repositories (mock external dependencies)
- Integration tests for API handlers (use test database)
- No merge without test coverage for new business logic
- Manual testing checklist provided for UI changes

**Documentation Standards**:
- Feature README in `/docs/features/` with user flow diagrams
- API changes documented with example requests/responses
- Database schema changes explained in migration comments
- Complex algorithms documented with comments explaining intent

---

## Governance

**Constitution Supersedes All**: When constitution conflicts with convenience or shortcuts, constitution wins. Exceptions REQUIRE documented justification and explicit team approval.

**Amendments**: This constitution evolves based on lessons learned. Changes require:
1. Assessment of current principle effectiveness
2. Proposed new wording with rationale
3. Impact analysis (what breaks? what needs updating?)
4. Version bump (MAJOR for principle removal, MINOR for clarification/addition, PATCH for wording only)
5. All dependent templates updated simultaneously

**Compliance Review**: Every feature delivery (end of 2-week cycle or milestone) includes: principle checklist verified, no shortcuts accumulated, architectural integrity confirmed.

**Runtime Guidance**: For day-to-day decisions not covered by constitution, see `.github/copilot-instructions.md` and branch-specific documentation.

---

**Version**: 1.0.0 | **Ratified**: 2026-06-02 | **Last Amended**: 2026-06-02
