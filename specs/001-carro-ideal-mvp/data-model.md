# Data Model & Database Schema: Carro Ideal MVP

**Date**: 2 de junho de 2026  
**Status**: Phase 1 Design Complete  
**Target**: PostgreSQL 15+

---

## Schema Overview

Carro Ideal uses a normalized relational schema with 8 core tables. All tables include `created_at` and `updated_at` timestamps managed by the application (not database triggers).

### Entity Relationship Diagram

```
┌─────────────┐
│   users     │ (id, email, password_hash, name, role, active)
└──────┬──────┘
       │
       ├──────────────────┬──────────────────┐
       ▼                  ▼                  ▼
   ┌───────────┐  ┌──────────────┐  ┌─────────────────┐
   │user_answers│ │recommendations│ │    sessions     │
   └─────┬─────┘  └──────┬───────┘  └─────────────────┘
         │               │
         ├──────┬────────┘
         ▼      ▼
    ┌─────────────┐     ┌──────────────────┐
    │ questions   │     │recommendation_   │
    │             │     │    items         │
    └──────┬──────┘     └────────┬─────────┘
           │                     │
           ▼                     ▼
    ┌─────────────────┐    ┌──────────┐
    │ answer_options  │    │ vehicles │
    └─────────────────┘    └─────┬────┘
                                 │
                                 ▼
                        ┌──────────────────┐
                        │vehicle_categories│
                        └──────────────────┘
```

---

## Table Specifications

### 1. `users` Table

**Purpose**: Stores user account information and role assignment.

**Columns**:
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,  -- bcrypt hash (60 chars)
  name VARCHAR(100) NOT NULL,
  role VARCHAR(20) NOT NULL,             -- 'USER' or 'ADMIN'
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id` (automatic)
- UNIQUE on `email` (enforces no duplicate registrations)

**Constraints**:
- `email` must be valid (application validates format; DB only checks unique)
- `password_hash` never null (application ensures before INSERT)
- `role` enforced as 'USER' or 'ADMIN' (application validation)
- `active` enables soft-delete (deactivate user without deleting records)

**Example Row**:
```json
{
  "id": 1,
  "email": "maria@email.com",
  "password_hash": "$2a$10$...(bcrypt hash)...",
  "name": "Maria Silva",
  "role": "USER",
  "active": true,
  "created_at": "2024-06-01T10:00:00Z",
  "updated_at": "2024-06-01T10:00:00Z"
}
```

---

### 2. `vehicle_categories` Table

**Purpose**: Organizes vehicles into semantic types (SUV, Sedan, etc.) for user understanding.

**Columns**:
```sql
CREATE TABLE vehicle_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- UNIQUE on `name` (prevent duplicate category names)

**Constraints**:
- `name` is unique and required
- `active` allows hiding categories without deleting

**Example Rows**:
```json
[
  { "id": 1, "name": "Sedan", "description": "Passenger cars with trunk" },
  { "id": 2, "name": "SUV", "description": "Sport Utility Vehicles" },
  { "id": 3, "name": "Hatchback", "description": "Compact cars with cargo flexibility" }
]
```

---

### 3. `vehicles` Table

**Purpose**: Core vehicle data used for recommendations and admin management.

**Columns**:
```sql
CREATE TABLE vehicles (
  id SERIAL PRIMARY KEY,
  category_id INTEGER NOT NULL REFERENCES vehicle_categories(id),
  brand VARCHAR(100) NOT NULL,              -- e.g., "Toyota"
  model VARCHAR(100) NOT NULL,              -- e.g., "Corolla"
  version VARCHAR(100),                     -- e.g., "GLi 1.6"
  year INTEGER NOT NULL,                    -- 1980-2050
  fuel_type VARCHAR(50) NOT NULL,           -- 'Gasoline', 'Diesel', 'Ethanol', 'Hybrid', 'Electric'
  transmission VARCHAR(50) NOT NULL,       -- 'Manual', 'Automatic', 'CVT'
  price_range VARCHAR(50),                  -- e.g., "80000-120000" (BRL)
  seats INTEGER,                            -- 2-8
  trunk_capacity INTEGER,                   -- liters
  consumption_city DECIMAL(5,2),            -- km/L
  consumption_highway DECIMAL(5,2),         -- km/L
  description TEXT,
  strengths TEXT,                           -- Bullet list or formatted text
  weaknesses TEXT,                          -- Bullet list or formatted text
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- FOREIGN KEY on `category_id` (link to categories table)
- INDEX on `active` (filter queries for active vehicles only)
- INDEX on `brand, model` (search queries)

**Constraints**:
- `category_id` must reference existing category
- `year` validated by application (1980-2050)
- `fuel_type` validated by application (enum)
- `transmission` validated by application (enum)
- `active` used for soft-delete

**Example Row**:
```json
{
  "id": 1,
  "category_id": 1,
  "brand": "Toyota",
  "model": "Corolla",
  "version": "GLi 1.6 2023",
  "year": 2023,
  "fuel_type": "Gasoline",
  "transmission": "Automatic",
  "price_range": "80000-120000",
  "seats": 5,
  "trunk_capacity": 300,
  "consumption_city": 9.5,
  "consumption_highway": 13.2,
  "description": "Reliable daily driver with excellent fuel economy",
  "strengths": "Fuel efficient, reliable, good resale value",
  "weaknesses": "Limited cargo space, not ideal for off-road",
  "active": true,
  "created_at": "2024-01-15T08:00:00Z",
  "updated_at": "2024-06-01T14:30:00Z"
}
```

---

### 4. `questions` Table

**Purpose**: Stores questionnaire questions that users answer to build preference profile.

**Columns**:
```sql
CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  text TEXT NOT NULL,                       -- "What is your budget?"
  type VARCHAR(50) NOT NULL,                -- 'SINGLE_CHOICE', 'MULTIPLE_CHOICE', 'SCALE'
  weight DECIMAL(3,2) DEFAULT 0.1,          -- 0.0-1.0, importance in scoring
  active BOOLEAN DEFAULT true,
  display_order INTEGER,                    -- Sort order in questionnaire
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- INDEX on `active` (filter for active questions)
- INDEX on `display_order` (sort questionnaire)

**Constraints**:
- `type` validated by application (enum: SINGLE_CHOICE, MULTIPLE_CHOICE, SCALE)
- `weight` between 0.0 and 1.0
- `display_order` nullable (admin can set custom order or leave NULL for append order)

**Example Rows**:
```json
[
  {
    "id": 1,
    "text": "What is your budget (BRL)?",
    "type": "SINGLE_CHOICE",
    "weight": 0.4,
    "display_order": 1,
    "active": true
  },
  {
    "id": 2,
    "text": "How important is fuel efficiency?",
    "type": "SINGLE_CHOICE",
    "weight": 0.3,
    "display_order": 2,
    "active": true
  }
]
```

---

### 5. `answer_options` Table

**Purpose**: Selectable answer choices for questions with scoring profiles.

**Columns**:
```sql
CREATE TABLE answer_options (
  id SERIAL PRIMARY KEY,
  question_id INTEGER NOT NULL REFERENCES questions(id),
  text VARCHAR(500) NOT NULL,               -- "R$ 50.000 - R$ 80.000"
  score_profile JSONB NOT NULL,             -- {"budget_score": 0.7, "efficiency_score": 0.3}
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- FOREIGN KEY on `question_id`

**Constraints**:
- `question_id` must reference existing question
- `score_profile` must be valid JSON with scoring weights

**Example Rows**:
```json
[
  {
    "id": 1,
    "question_id": 1,
    "text": "R$ 30.000 - R$ 50.000",
    "score_profile": {
      "budget_score": 0.3,
      "efficiency_score": 0.7,
      "comfort_score": 0.4
    }
  },
  {
    "id": 2,
    "question_id": 1,
    "text": "R$ 50.000 - R$ 80.000",
    "score_profile": {
      "budget_score": 0.7,
      "efficiency_score": 0.5,
      "comfort_score": 0.6
    }
  }
]
```

**Score Profile Explanation**:
- Keys represent scoring dimensions (budget_score, efficiency_score, comfort_score, etc.)
- Values are decimal weights (0.0-1.0) indicating how this answer contributes to each dimension
- Algorithm multiplies these weights by vehicle attribute matches to calculate final score

---

### 6. `user_answers` Table

**Purpose**: Records user responses to questionnaire, one row per answer.

**Columns**:
```sql
CREATE TABLE user_answers (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  question_id INTEGER NOT NULL REFERENCES questions(id),
  answer_option_id INTEGER NOT NULL REFERENCES answer_options(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- FOREIGN KEY on `user_id`, `question_id`, `answer_option_id`
- COMPOSITE INDEX on `(user_id, question_id)` (find user's answer to specific question)
- INDEX on `(user_id, created_at)` (find user's most recent answers)

**Constraints**:
- All foreign keys required (no null references)
- Implicit: only one answer per user per question (application enforces, or add UNIQUE constraint)

**Example Rows**:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "question_id": 1,
    "answer_option_id": 2,
    "created_at": "2024-06-01T11:00:00Z",
    "updated_at": "2024-06-01T11:00:00Z"
  },
  {
    "id": 2,
    "user_id": 1,
    "question_id": 2,
    "answer_option_id": 5,
    "created_at": "2024-06-01T11:05:00Z",
    "updated_at": "2024-06-01T11:05:00Z"
  }
]
```

---

### 7. `recommendations` Table

**Purpose**: Represents a single recommendation generation session (parent record for recommendation results).

**Columns**:
```sql
CREATE TABLE recommendations (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  summary TEXT,                             -- Overall recommendation summary
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- FOREIGN KEY on `user_id`
- INDEX on `(user_id, created_at)` (find user's recommendations in order)

**Constraints**:
- `user_id` must reference existing user

**Example Row**:
```json
{
  "id": 1,
  "user_id": 1,
  "summary": "Based on your budget and efficiency preferences, these vehicles are ideal",
  "created_at": "2024-06-01T12:00:00Z"
}
```

---

### 8. `recommendation_items` Table

**Purpose**: Individual vehicle recommendations within a recommendation session (child records with scores).

**Columns**:
```sql
CREATE TABLE recommendation_items (
  id SERIAL PRIMARY KEY,
  recommendation_id INTEGER NOT NULL REFERENCES recommendations(id),
  vehicle_id INTEGER NOT NULL REFERENCES vehicles(id),
  score DECIMAL(5,2) NOT NULL,              -- 0-100
  reason TEXT,                              -- "Good match on budget and fuel efficiency"
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- PRIMARY KEY on `id`
- FOREIGN KEY on `recommendation_id`, `vehicle_id`
- INDEX on `(recommendation_id, score DESC)` (get top vehicles for a recommendation)

**Constraints**:
- `score` between 0 and 100
- Both foreign keys required

**Example Rows**:
```json
[
  {
    "id": 1,
    "recommendation_id": 1,
    "vehicle_id": 1,
    "score": 92.5,
    "reason": "Excellent fuel efficiency match, budget range fits perfectly, reliable brand"
  },
  {
    "id": 2,
    "recommendation_id": 1,
    "vehicle_id": 3,
    "score": 85.0,
    "reason": "Good fuel efficiency, budget match, slightly higher price"
  }
]
```

---

### 9. `sessions` Table

**Purpose**: Server-side session storage for user authentication.

**Columns**:
```sql
CREATE TABLE sessions (
  id VARCHAR(128) PRIMARY KEY,              -- Random session token
  user_id INTEGER NOT NULL REFERENCES users(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL             -- 24 hours from creation
);
```

**Indexes**:
- PRIMARY KEY on `id` (session lookup)
- FOREIGN KEY on `user_id`
- INDEX on `expires_at` (cleanup query for expired sessions)

**Constraints**:
- `id` is unique (by design, generated via crypto/rand)
- `expires_at` must be in future (application validates)

**Example Row**:
```json
{
  "id": "a3f8d9e2c1b5f7a9d4e6c8b0f2a4d6e8",
  "user_id": 1,
  "created_at": "2024-06-01T13:00:00Z",
  "expires_at": "2024-06-02T13:00:00Z"
}
```

---

## Schema Migration Strategy

All schema changes via `golang-migrate` with numbered SQL files:

```
migrations/
├── 0001_create_users_table.up.sql
├── 0001_create_users_table.down.sql
├── 0002_create_vehicle_categories_table.up.sql
├── 0002_create_vehicle_categories_table.down.sql
├── 0003_create_vehicles_table.up.sql
├── 0003_create_vehicles_table.down.sql
... and so on
```

**Migration Principles**:
- Each .up.sql file is idempotent (safe to run multiple times)
- Each .down.sql fully reverses the .up.sql
- Migrations are append-only (never modify existing migration files)
- If mistake found in migration, create new migration to fix (0010_fix_column_typo.up.sql)

---

## Data Validation & Constraints

### Application-Level Validation (Go)

| Field | Rule |
|-------|------|
| User.email | Must be valid RFC 5322 email format |
| User.password | Minimum 8 characters, hashed via bcrypt |
| User.name | 1-100 characters, non-empty |
| Vehicle.year | 1980-2050 |
| Vehicle.fuel_type | One of: Gasoline, Diesel, Ethanol, Hybrid, Electric |
| Vehicle.transmission | One of: Manual, Automatic, CVT |
| Vehicle.seats | 2-8 (typical car range) |
| Question.type | One of: SINGLE_CHOICE, MULTIPLE_CHOICE, SCALE |
| Question.weight | Decimal 0.0-1.0 |
| AnswerOption.score_profile | Valid JSON with numeric values 0.0-1.0 |
| RecommendationItem.score | Decimal 0-100 |

---

## Query Patterns

### Common Queries

**Get user with active session**:
```sql
SELECT u.* FROM users u
INNER JOIN sessions s ON u.id = s.user_id
WHERE s.id = $1 AND s.expires_at > NOW();
```

**Get active vehicles in category**:
```sql
SELECT * FROM vehicles
WHERE category_id = $1 AND active = true
ORDER BY brand, model;
```

**Get user's latest answers (for recommendation)**:
```sql
SELECT q.id, q.text, q.weight, ao.text as answer_text, ao.score_profile
FROM questions q
LEFT JOIN user_answers ua ON ua.question_id = q.id AND ua.user_id = $1
LEFT JOIN answer_options ao ON ua.answer_option_id = ao.id
WHERE q.active = true
ORDER BY q.display_order;
```

**Get recommendation details with items**:
```sql
SELECT r.id, r.user_id, r.summary, r.created_at,
       ri.id as item_id, ri.vehicle_id, ri.score, ri.reason,
       v.brand, v.model, v.version, v.year, v.fuel_type, v.transmission,
       v.price_range, v.consumption_city, v.consumption_highway
FROM recommendations r
LEFT JOIN recommendation_items ri ON ri.recommendation_id = r.id
LEFT JOIN vehicles v ON ri.vehicle_id = v.id
WHERE r.id = $1
ORDER BY ri.score DESC;
```

---

## Performance Considerations

### Indexes Strategy

| Table | Indexes |
|-------|---------|
| users | (email) UNIQUE, (active) |
| vehicles | (category_id), (active), (brand, model) |
| questions | (active), (display_order) |
| user_answers | (user_id, question_id), (user_id, created_at DESC) |
| recommendations | (user_id, created_at DESC) |
| recommendation_items | (recommendation_id, score DESC) |
| sessions | (expires_at) for cleanup |

### Query Optimization Notes

- **Recommendation calculation**: Fetch user answers + vehicle data in two queries (N+1 acceptable for MVP, optimize if needed)
- **User dashboard**: Aggregate counts for statistics (total recommendations, most recent recommendation)
- **Pagination**: Use LIMIT + OFFSET for list endpoints (avoid large offset; consider keyset pagination if performance issue)

---

## Soft-Delete Pattern

Three core tables support soft-delete via `active` boolean:
- `users` - deactivate user account
- `vehicles` - hide vehicle from recommendations
- `vehicle_categories` - hide category but keep relationships

**Query Pattern**:
```sql
-- Always filter: WHERE active = true
SELECT * FROM vehicles WHERE active = true;

-- Delete: Set active = false, keep historical data
UPDATE vehicles SET active = false WHERE id = $1;

-- Restore: Set active = true
UPDATE vehicles SET active = true WHERE id = $1;
```

---

## Scaling Considerations (Phase 2+)

If system grows beyond MVP needs:

1. **Session Storage**: Migrate from in-memory or sessions table to Redis (faster lookups)
2. **Read Replicas**: PostgreSQL read replicas for recommendation queries (read-heavy)
3. **Caching**: Cache frequently-accessed data (categories, questions) in memory (Go)
4. **Materialized Views**: Pre-calculate recommendation scores for popular preference profiles
5. **Partitioning**: Partition `user_answers` and `recommendation_items` by date if billions of rows

For MVP (estimated <10k users during TCC period), single PostgreSQL instance sufficient.

---

## Testing & Data Setup

### Test Database Schema

Same as production. Tests use:
- Docker PostgreSQL container for integration tests
- SQLite in-memory for unit test isolation (if desired)
- Test fixtures in SQL files to seed initial data

### Sample Data for Development

Sample vehicles, categories, and questions provided in seed migration file (0_seed_data.sql, optional):
```sql
INSERT INTO vehicle_categories (name, description) VALUES
  ('Sedan', 'Passenger cars with trunk'),
  ('SUV', 'Sport Utility Vehicles'),
  ...
```

---

**Schema Status**: ✅ COMPLETE - Ready for Migration Implementation