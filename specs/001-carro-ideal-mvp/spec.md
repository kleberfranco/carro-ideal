# Feature Specification: Carro Ideal MVP - Complete Car Recommendation System

**Feature Branch**: `001-carro-ideal-mvp`

**Created**: 2 de junho de 2026

**Status**: Phase 5 Implemented - Ready for Testing and QA

**Input**: Complete project brief with academic context, product description, technologies, architecture, requirements, entities, and user flows for TCC capstone delivery.

---

## Vision and Scope

**Vision**: Deliver a production-grade car recommendation web application that demonstrates comprehensive software engineering practices for a TCC capstone project. The system helps Brazilian users find their ideal vehicle through intelligent questionnaire-based recommendations powered by structured data, scoring algorithms, and prepared for future AI integration.

**Scope**:
- Complete MVP with user registration, authentication, questionnaire, and vehicle recommendations
- Admin panel for vehicle, category, question, and answer management
- Docker-based deployment with PostgreSQL persistence
- Bootstrap 5 responsive UI
- REST API with standardized JSON responses
- Recommendation history tracking
- Architecture prepared for future ChatGPT integration

**Out of Scope (Phase 2+)**:
- Real ChatGPT integration (architecture prepared, not implemented)
- Mobile native apps
- Real-time messaging or chat
- Advanced analytics dashboard
- Multi-language support
- Payment processing
- Vehicle inventory sync from external APIs

---

## User Personas

### Persona 1: Maria, Car Buyer (Primary User)

**Profile**: 35-year-old software engineer, first-time car buyer, urban São Paulo resident  
**Goals**: Find a reliable, affordable car suitable for city commuting and occasional highway trips  
**Pain Points**: Overwhelmed by car market choices, unsure what features matter for her use case, mistrusts car salesman recommendations  
**Tech Savvy**: High (comfortable with web apps and forms)  
**Frequency**: One-time or occasional user (buying a car every 5-10 years)

### Persona 2: João, Admin Vehicle Manager

**Profile**: 45-year-old automotive consultant managing Carro Ideal's vehicle database  
**Goals**: Keep vehicle catalog accurate, add new models as they release, adjust recommendation criteria based on market trends  
**Pain Points**: Manual data entry, needs quick database updates, wants consistent categorization  
**Tech Savvy**: Medium (comfortable with web forms, not coding)  
**Frequency**: Weekly updates, ongoing maintenance

### Persona 3: Professor Silva, TCC Evaluator

**Profile**: 50-year-old software engineering professor evaluating the capstone project  
**Goals**: Assess proper architecture, code quality, engineering practices, and documentation  
**Pain Points**: Needs to understand system design, verify layered architecture, check if README is complete  
**Tech Savvy**: Very high (expert in software engineering)  
**Frequency**: Evaluation period (once for project defense)

---

## User Scenarios & Testing (Prioritized by Business Value)

### User Story 1 - User Registration and Onboarding (Priority: P1)

**Narrative**: As a new user, I want to create an account quickly with my basic information so I can access the car recommendation system and save my preferences.

**Why this priority**: Critical path - no user can proceed without account creation. Demonstrates basic authentication and database operations, foundational for all other features.

**Independent Test**: Can be fully tested by (1) loading homepage, (2) clicking register, (3) entering email/password/name, (4) submitting form, (5) verifying account exists in database, (6) logging in successfully. Delivers immediate value: user account ready for next steps.

**Acceptance Scenarios**:

1. **Given** user is on homepage, **When** user clicks "Sign Up" button, **Then** registration form loads with fields for name, email, password, and confirm password
2. **Given** user fills valid registration form with unique email, **When** user submits form, **Then** account is created, success message shown, user is redirected to login page
3. **Given** user tries to register with already-used email, **When** user submits form, **Then** error message states "Email already registered" and form persists
4. **Given** user enters mismatched passwords, **When** user submits form, **Then** error message states "Passwords do not match"
5. **Given** user enters weak password (less than 8 characters), **When** user submits form, **Then** error states minimum length requirement

---

### User Story 2 - User Authentication and Session Management (Priority: P1)

**Narrative**: As a registered user, I want to securely log in and stay authenticated while I use the system, so my recommendations and history remain private and accessible.

**Why this priority**: Critical - prevents unauthorized access, ensures data privacy, blocks admins from entering user areas. Demonstrates security practice understanding.

**Independent Test**: Can be tested by (1) attempting login with invalid credentials, (2) logging in successfully, (3) accessing protected pages, (4) logging out, (5) verifying redirect to login. Delivers security and user session management.

**Acceptance Scenarios**:

1. **Given** user is on login page, **When** user enters correct email and password, **Then** session is created and user redirected to dashboard
2. **Given** user is on login page, **When** user enters incorrect email, **Then** error message shows "Email not found or password incorrect" (generic for security)
3. **Given** user is logged in and accesses `/api/user/dashboard`, **When** request includes valid session, **Then** dashboard data returns with HTTP 200
4. **Given** user is not logged in and tries to access `/api/user/dashboard`, **When** no valid session present, **Then** HTTP 401 returned with redirect to login
5. **Given** logged-in user clicks logout button, **When** logout endpoint called, **Then** session destroyed and user redirected to homepage
6. **Given** user has valid session, **When** session expires after 24 hours, **Then** user redirected to login on next request

---

### User Story 3 - Answer Preference Questionnaire (Priority: P1)

**Narrative**: As a user, I want to answer questions about my preferences, budget, and use case in a clear, guided questionnaire so the system can understand my needs and provide relevant recommendations.

**Why this priority**: Core feature - this is the feature differentiator. User questionnaire captures preference profile which drives all recommendations. Demonstrates questionnaire system and business logic.

**Independent Test**: Can be tested by (1) navigating to questionnaire page, (2) answering all questions, (3) submitting responses, (4) verifying responses saved in database. Delivers personalized user profile.

**Acceptance Scenarios**:

1. **Given** user is logged in and on questionnaire page, **When** page loads, **Then** all active questions display in order with corresponding answer options
2. **Given** user selects answer options for all questions, **When** user clicks "Submit Questionnaire", **Then** responses are saved and confirmation message shows
3. **Given** user answers question about budget with value "R$ 50.000 - R$ 80.000", **When** user submits, **Then** value stored correctly in database
4. **Given** user had previously answered questionnaire, **When** user visits questionnaire again, **Then** previous answers are displayed (user can review/update)
5. **Given** user answers all questions, **When** user submits, **Then** user is directed to recommendation results page automatically
6. **Given** questionnaire has questions with weights (e.g., budget 0.3, fuel 0.2), **When** user submits, **Then** weights applied in recommendation calculation

---

### User Story 4 - Generate Vehicle Recommendations (Priority: P1)

**Narrative**: As a user, I want to receive a prioritized list of recommended vehicles based on my questionnaire answers, with clear explanations of why each vehicle matches my preferences.

**Why this priority**: Core value delivery - provides immediate recommendation output. Demonstrates recommendation algorithm and scoring logic.

**Independent Test**: Can be tested by (1) completing questionnaire, (2) triggering recommendation generation, (3) receiving ranked vehicle list with scores and explanations, (4) verifying calculations. Delivers core business value.

**Acceptance Scenarios**:

1. **Given** user submits completed questionnaire, **When** recommendation endpoint (`POST /api/recommendations/generate`) called, **Then** system calculates scores for all active vehicles
2. **Given** user with "city car, budget R$ 50k, fuel efficiency important" answers, **When** recommendations generated, **Then** small city cars with good fuel economy rank highest
3. **Given** system generates recommendations, **When** response received, **Then** top 5-10 vehicles returned sorted by score descending
4. **Given** vehicle has score 85 and another has 65, **When** recommendations displayed, **Then** higher-scoring vehicle appears first
5. **Given** recommendation calculated, **When** response sent, **Then** each vehicle includes: vehicle details, score, and reason explanation
6. **Given** user's answers match vehicle attributes exactly, **When** recommendation calculated, **Then** score reflects strong match

---

### User Story 5 - View Vehicle Details (Priority: P1)

**Narrative**: As a user, I want to see detailed information about a recommended vehicle so I can evaluate its features, specifications, and strengths/weaknesses.

**Why this priority**: Follows recommendations - users need to explore vehicles to make informed decisions. Demonstrates data presentation.

**Independent Test**: Can be tested by (1) accessing recommended vehicle from list, (2) viewing detailed page, (3) verifying all fields display, (4) checking external links if applicable. Delivers exploration capability.

**Acceptance Scenarios**:

1. **Given** user views recommendation results, **When** user clicks vehicle card or name, **Then** vehicle details page loads with complete information
2. **Given** vehicle details page loads, **When** page renders, **Then** displays: brand, model, version, year, price_range, fuel_type, transmission, seats, trunk capacity, consumption data, strengths, weaknesses
3. **Given** vehicle details page open, **When** user reviews page, **Then** information clearly organized in sections (specs, consumption, features)
4. **Given** user on vehicle detail, **When** user wants to return to recommendations, **Then** back button or link returns to previous list
5. **Given** vehicle has special features or restrictions, **When** details display, **Then** all information is accurate and complete

---

### User Story 6 - View Recommendation History (Priority: P2)

**Narrative**: As a user, I want to see past recommendations I've received so I can compare results from different time periods or review previous vehicle research.

**Why this priority**: Secondary value - nice-to-have for user engagement. Demonstrates historical data retrieval.

**Independent Test**: Can be tested by (1) generating multiple recommendations over time, (2) accessing history page, (3) viewing all past recommendations. Delivers user history feature.

**Acceptance Scenarios**:

1. **Given** user is logged in, **When** user navigates to "Recommendation History", **Then** page displays list of all past recommendations with dates
2. **Given** user has multiple recommendation records, **When** history page loads, **Then** recommendations sorted by date descending (newest first)
3. **Given** user clicks on past recommendation, **When** details load, **Then** user sees which vehicles were recommended at that time
4. **Given** user selects a past recommendation, **When** page loads, **Then** can compare that recommendation with current results (optional side-by-side view)

---

### User Story 7 - Admin Vehicle Management (Priority: P2)

**Narrative**: As an admin, I want to create, read, update, and delete vehicles in the system so I can maintain accurate and current vehicle data for recommendations.

**Why this priority**: Admin feature - essential for system maintenance. Demonstrates full CRUD operations and admin authorization.

**Independent Test**: Can be tested by (1) logging in as admin, (2) accessing vehicle management, (3) creating a new vehicle, (4) editing an existing vehicle, (5) deleting a vehicle, (6) verifying database changes. Delivers admin functionality.

**Acceptance Scenarios**:

1. **Given** admin is logged in and accesses admin dashboard, **When** admin clicks "Vehicles", **Then** admin vehicle management interface loads
2. **Given** admin is in vehicle management, **When** admin clicks "Add New Vehicle", **Then** form appears with all required fields (brand, model, year, price_range, fuel_type, transmission, seats, trunk_capacity, consumption_city, consumption_highway, strengths, weaknesses)
3. **Given** admin fills new vehicle form completely, **When** admin clicks "Save", **Then** vehicle created, success message shown, vehicle appears in list
4. **Given** existing vehicle in list, **When** admin clicks "Edit", **Then** form loads with current vehicle data
5. **Given** admin edits vehicle fields and clicks "Save", **When** request submitted, **Then** database updated and list refreshed
6. **Given** vehicle exists in list, **When** admin clicks "Delete" and confirms, **Then** vehicle marked inactive and removed from recommendations (soft delete)
7. **Given** admin views vehicle list, **When** list displays, **Then** shows: brand, model, year, price range, number of recommendations using it

---

### User Story 8 - Admin Category Management (Priority: P2)

**Narrative**: As an admin, I want to manage vehicle categories (SUV, Sedan, Hatchback, etc.) so I can organize vehicles and help users think about vehicle types.

**Why this priority**: Supporting admin feature - enables vehicle organization. Demonstrates secondary CRUD entity.

**Independent Test**: Can be tested by (1) accessing category management, (2) creating/editing/deleting categories, (3) verifying categories available when creating vehicles. Delivers vehicle categorization.

**Acceptance Scenarios**:

1. **Given** admin accesses admin dashboard, **When** admin navigates to "Categories", **Then** category management interface loads
2. **Given** category management interface open, **When** admin clicks "Add Category", **Then** form appears for name and description
3. **Given** admin enters category name (e.g., "SUV") and saves, **When** submitted, **Then** category created and appears in list
4. **Given** categories exist in list, **When** admin edits a category, **Then** form populates with current data
5. **Given** category is created, **When** admin creates new vehicle, **Then** category available in vehicle form dropdown
6. **Given** category exists with vehicles linked, **When** admin deletes category, **Then** admin prompted about linked vehicles before deletion

---

### User Story 9 - Admin Question and Answer Management (Priority: P2)

**Narrative**: As an admin, I want to manage questionnaire questions and answer options so I can refine how the system captures user preferences and adjust recommendation criteria.

**Why this priority**: Supporting admin feature - core to algorithm tuning. Demonstrates complex CRUD with nested relationships.

**Independent Test**: Can be tested by (1) accessing question management, (2) creating questions with answer options, (3) editing questions, (4) testing questionnaire reflects changes. Delivers questionnaire control.

**Acceptance Scenarios**:

1. **Given** admin accesses admin dashboard, **When** admin navigates to "Questions", **Then** question management interface loads
2. **Given** question list displays, **When** admin clicks "Add Question", **Then** form appears for question text, type, weight, display order
3. **Given** admin creates new question and adds answer options, **When** saved, **Then** question appears in user questionnaire
4. **Given** question has multiple answer options, **When** admin edits a question, **Then** can add, edit, remove answer options
5. **Given** answer option has associated score_profile field, **When** admin edits option, **Then** can adjust score weights for recommendation calculation
6. **Given** admin marks question as inactive, **When** questionnaire loads for users, **Then** inactive question no longer appears
7. **Given** question has display order value, **When** user loads questionnaire, **Then** questions displayed in specified order

---

### User Story 10 - View Admin Dashboard with System Overview (Priority: P3)

**Narrative**: As an admin, I want to see a dashboard with key system statistics so I can understand platform usage and health.

**Why this priority**: Administrative convenience - helps admin understand system state. Non-critical but demonstrates data aggregation.

**Independent Test**: Can be tested by (1) logging in as admin, (2) viewing dashboard, (3) seeing statistics like user count, vehicle count, recommendation count. Delivers admin overview.

**Acceptance Scenarios**:

1. **Given** admin logs in, **When** dashboard loads, **Then** displays: total users, total vehicles, total recommendations, active users this week
2. **Given** admin on dashboard, **When** page displays, **Then** shows links to all admin management functions
3. **Given** admin accesses dashboard, **When** page renders, **Then** data is current (updated within last hour)

---

## Edge Cases and Error Scenarios

- **Empty questionnaire**: System handles when no questions configured (shows message to admin)
- **No vehicles in system**: User completes questionnaire but no vehicles match; system shows helpful message
- **Concurrent recommendations**: User generating multiple recommendations simultaneously (system queues or rejects with clear message)
- **Session expiry during questionnaire**: User loses session mid-questionnaire; system prompts to log back in and resume
- **Database connection failure**: System returns HTTP 500 with generic error message (no internal details exposed)
- **Invalid vehicle data**: Admin enters invalid price range or year; system validates and shows specific error
- **Circular reference**: Admin tries to delete category with vehicles linked; system prevents deletion and suggests options
- **Large questionnaire**: System with 50+ questions; UI remains responsive, pagination or scrolling implemented
- **Recommendation calculation timeout**: Score calculation takes too long; system shows loading indicator, implements timeout with fallback

---

## Requirements (Mandatory)

### Functional Requirements

**Authentication & Authorization**:
- **FR-001**: System MUST allow users to create accounts with email, name, and password (minimum 8 characters)
- **FR-002**: System MUST validate email format and uniqueness before account creation
- **FR-003**: System MUST hash passwords using bcrypt or similar (never store plaintext)
- **FR-004**: System MUST support session-based authentication (session ID stored server-side, cookie-based to client)
- **FR-005**: System MUST implement role-based access control (USER and ADMIN roles)
- **FR-006**: System MUST prevent unauthorized users from accessing admin endpoints (HTTP 403 Forbidden)
- **FR-007**: System MUST provide login endpoint (`POST /api/auth/login`) returning session token
- **FR-008**: System MUST provide logout endpoint (`POST /api/auth/logout`) invalidating session
- **FR-009**: System MUST provide "get current user" endpoint (`GET /api/auth/me`) for session verification

**User Profile & Questionnaire**:
- **FR-010**: System MUST allow users to view their profile and basic account information
- **FR-011**: System MUST support questionnaire with multiple question types: single-choice, multiple-choice, scale/range
- **FR-012**: System MUST store user answers to questionnaire with timestamps
- **FR-013**: System MUST allow users to update questionnaire answers (modify previous responses)
- **FR-014**: System MUST validate that all required questions are answered before submission
- **FR-015**: System MUST support question weighting (e.g., budget = 0.3, fuel = 0.2) for recommendation scoring

**Vehicle Recommendation**:
- **FR-016**: System MUST calculate recommendation scores based on user answers and vehicle attributes
- **FR-017**: System MUST rank vehicles by score descending and return top results (minimum 5, maximum 15 vehicles)
- **FR-018**: System MUST provide recommendation explanation for each vehicle (text justification)
- **FR-019**: System MUST store recommendations with timestamp and user reference for history
- **FR-020**: System MUST support score-based algorithm: matching user preference profile with vehicle attributes
- **FR-021**: System MUST handle edge case where no vehicles match user preferences gracefully

**Vehicle Management**:
- **FR-022**: System MUST support full CRUD operations for vehicles (Create, Read, Update, Delete)
- **FR-023**: System MUST store vehicle attributes: brand, model, version, year, fuel_type, transmission, price_range, seats, trunk_capacity, consumption_city, consumption_highway, strengths, weaknesses, category, active status
- **FR-024**: System MUST support soft-delete for vehicles (mark inactive rather than permanent delete)
- **FR-025**: System MUST only show active vehicles in user recommendations
- **FR-026**: System MUST prevent vehicles with critical missing data from being included in recommendations

**Category & Answer Management**:
- **FR-027**: System MUST support CRUD for vehicle categories (Sedan, SUV, Hatchback, etc.)
- **FR-028**: System MUST support CRUD for questionnaire questions
- **FR-029**: System MUST support CRUD for answer options with associated score profiles
- **FR-030**: System MUST allow admin to order questions via display sequence

**API Endpoints - Authentication**:
- **FR-031**: `POST /api/auth/register` - Create new user account
- **FR-032**: `POST /api/auth/login` - Authenticate user, return session
- **FR-033**: `POST /api/auth/logout` - Destroy session
- **FR-034**: `GET /api/auth/me` - Return current authenticated user info

**API Endpoints - User**:
- **FR-035**: `GET /api/user/dashboard` - Return user dashboard data
- **FR-036**: `GET /api/user/recommendations` - Return user's recommendation history
- **FR-037**: `GET /api/user/recommendations/{id}` - Return specific recommendation with details

**API Endpoints - Questionnaire & Recommendations**:
- **FR-038**: `GET /api/questions` - Return all active questions with answer options
- **FR-039**: `POST /api/questionnaire/answers` - Store user questionnaire answers
- **FR-040**: `POST /api/recommendations/generate` - Generate recommendations based on latest answers

**API Endpoints - Vehicles & Data**:
- **FR-041**: `GET /api/vehicles` - Return all active vehicles (paginated, filterable)
- **FR-042**: `GET /api/vehicles/{id}` - Return vehicle details

**API Endpoints - Admin**:
- **FR-043**: `GET /api/admin/vehicles` - List all vehicles (including inactive)
- **FR-044**: `POST /api/admin/vehicles` - Create new vehicle
- **FR-045**: `PUT /api/admin/vehicles/{id}` - Update vehicle
- **FR-046**: `DELETE /api/admin/vehicles/{id}` - Delete (soft-delete) vehicle
- **FR-047**: `GET /api/admin/categories` - List categories
- **FR-048**: `POST /api/admin/categories` - Create category
- **FR-049**: `PUT /api/admin/categories/{id}` - Update category
- **FR-050**: `DELETE /api/admin/categories/{id}` - Delete category
- **FR-051**: `GET /api/admin/questions` - List questions
- **FR-052**: `POST /api/admin/questions` - Create question
- **FR-053**: `PUT /api/admin/questions/{id}` - Update question
- **FR-054**: `DELETE /api/admin/questions/{id}` - Delete question
- **FR-055**: `GET /api/admin/answer-options` - List answer options
- **FR-056**: `POST /api/admin/answer-options` - Create answer option
- **FR-057**: `PUT /api/admin/answer-options/{id}` - Update answer option
- **FR-058**: `DELETE /api/admin/answer-options/{id}` - Delete answer option

**UI/Web Pages**:
- **FR-059**: System MUST provide responsive homepage with login/register links (Bootstrap 5)
- **FR-060**: System MUST provide registration form with validation feedback
- **FR-061**: System MUST provide login form with error handling
- **FR-062**: System MUST provide user dashboard with navigation options
- **FR-063**: System MUST provide questionnaire form rendering all questions in order
- **FR-064**: System MUST provide recommendation results page with vehicle list
- **FR-065**: System MUST provide vehicle details page with all information
- **FR-066**: System MUST provide recommendation history page
- **FR-067**: System MUST provide admin login page (separate from user login)
- **FR-068**: System MUST provide admin dashboard with navigation
- **FR-069**: System MUST provide admin CRUD interfaces for vehicles, categories, questions, answer options
- **FR-070**: System MUST provide responsive design working on desktop and tablet (Bootstrap 5)

**Data Persistence & Validation**:
- **FR-071**: System MUST validate all input data (email format, password length, price ranges, year validity)
- **FR-072**: System MUST return meaningful error messages for validation failures
- **FR-073**: System MUST enforce database constraints (unique emails, foreign key relationships, non-null required fields)
- **FR-074**: System MUST support database transactions for critical operations (e.g., user registration)
- **FR-075**: System MUST create database records with created_at and updated_at timestamps automatically

**API Response Format**:
- **FR-076**: All API responses MUST follow standardized JSON format: `{success: boolean, data: {}, message: string}` for success
- **FR-077**: All API error responses MUST follow format: `{success: false, error: string, code: string}` with appropriate HTTP status codes
- **FR-078**: All list responses MUST support pagination: `{success: true, data: {items: [], total: number, page: number, limit: number}}`

**Logging & Error Handling**:
- **FR-079**: System MUST log all authentication attempts (successful and failed)
- **FR-080**: System MUST log all admin actions (create, update, delete operations)
- **FR-081**: System MUST return HTTP 400 for client errors (invalid input)
- **FR-082**: System MUST return HTTP 401 for authentication failures
- **FR-083**: System MUST return HTTP 403 for authorization failures
- **FR-084**: System MUST return HTTP 404 for resource not found
- **FR-085**: System MUST return HTTP 500 for server errors (never expose internal details to client)

---

## Key Entities

### User

**Purpose**: Represents a registered user of the system who can complete questionnaires and receive recommendations.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `name` (String, 100): User full name
- `email` (String, 255, Unique): Email address
- `password_hash` (String): Bcrypt hashed password (never plaintext)
- `role` (Enum: USER, ADMIN): User role for authorization
- `active` (Boolean, default: true): Account status
- `created_at` (Timestamp): Account creation time
- `updated_at` (Timestamp): Last profile modification

**Relationships**: One-to-Many with UserAnswer, Recommendation

---

### VehicleCategory

**Purpose**: Organizes vehicles into types (Sedan, SUV, Hatchback, etc.) for user understanding and admin management.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `name` (String, 100): Category name (e.g., "SUV", "Sedan")
- `description` (Text): Category explanation
- `active` (Boolean, default: true): Category visibility
- `created_at` (Timestamp): Creation time
- `updated_at` (Timestamp): Last update time

**Relationships**: One-to-Many with Vehicle

---

### Vehicle

**Purpose**: Represents a car model with specifications and attributes used for matching with user preferences.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `category_id` (Foreign Key): Reference to VehicleCategory
- `brand` (String, 100): Manufacturer (e.g., "Toyota")
- `model` (String, 100): Model name (e.g., "Corolla")
- `version` (String, 100): Trim/version (e.g., "GLi")
- `year` (Integer, 1980-2050): Model year
- `fuel_type` (Enum: Gasoline, Diesel, Ethanol, Hybrid, Electric): Fuel source
- `transmission` (Enum: Manual, Automatic, CVT): Transmission type
- `price_range` (String, e.g., "80000-120000"): Estimated price in BRL
- `seats` (Integer, 2-8): Number of seats
- `trunk_capacity` (Integer): Trunk volume in liters
- `consumption_city` (Decimal): City fuel consumption (km/L)
- `consumption_highway` (Decimal): Highway fuel consumption (km/L)
- `description` (Text): Vehicle overview
- `strengths` (Text): Key advantages
- `weaknesses` (Text): Notable limitations
- `active` (Boolean, default: true): Visibility in recommendations
- `created_at` (Timestamp): Creation time
- `updated_at` (Timestamp): Last update time

**Relationships**: Many-to-One with VehicleCategory, One-to-Many with RecommendationItem

---

### Question

**Purpose**: Represents a questionnaire question capturing user preferences.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `text` (Text): Question text (e.g., "What is your budget?")
- `type` (Enum: SINGLE_CHOICE, MULTIPLE_CHOICE, SCALE): Answer type
- `weight` (Decimal, 0.0-1.0): Importance in recommendation calculation (default: 0.1)
- `active` (Boolean, default: true): Visibility in questionnaire
- `display_order` (Integer): Sequence in questionnaire
- `created_at` (Timestamp): Creation time
- `updated_at` (Timestamp): Last update time

**Relationships**: One-to-Many with AnswerOption, UserAnswer

---

### AnswerOption

**Purpose**: Represents a selectable answer for a question with associated scoring profile.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `question_id` (Foreign Key): Reference to Question
- `text` (String, 500): Answer option text (e.g., "R$ 50.000 - R$ 80.000")
- `score_profile` (JSON): Weights for recommendation scoring (e.g., `{"budget_score": 0.7, "efficiency_score": 0.3}`)
- `created_at` (Timestamp): Creation time
- `updated_at` (Timestamp): Last update time

**Relationships**: Many-to-One with Question, One-to-Many with UserAnswer

---

### UserAnswer

**Purpose**: Records a user's response to a questionnaire question.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `user_id` (Foreign Key): Reference to User
- `question_id` (Foreign Key): Reference to Question
- `answer_option_id` (Foreign Key): Reference to AnswerOption selected
- `created_at` (Timestamp): Answer timestamp
- `updated_at` (Timestamp): If answer modified

**Relationships**: Many-to-One with User, Question, AnswerOption

---

### Recommendation

**Purpose**: Represents a single recommendation session and results.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `user_id` (Foreign Key): Reference to User
- `summary` (Text): Overall recommendation summary/algorithm result
- `created_at` (Timestamp): Generation timestamp

**Relationships**: Many-to-One with User, One-to-Many with RecommendationItem

---

### RecommendationItem

**Purpose**: Links a specific vehicle to a recommendation with score and justification.

**Attributes**:
- `id` (UUID/Integer): Unique identifier
- `recommendation_id` (Foreign Key): Reference to Recommendation
- `vehicle_id` (Foreign Key): Reference to Vehicle
- `score` (Decimal, 0-100): Calculated match score
- `reason` (Text): Explanation of why vehicle recommended (e.g., "Good fuel efficiency and budget match")

**Relationships**: Many-to-One with Recommendation, Vehicle

---

## API Contract Detail

### Authentication

#### POST /api/auth/register
**Request**:
```json
{
  "name": "Maria Silva",
  "email": "maria@email.com",
  "password": "SecurePass123"
}
```

**Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-or-id",
    "email": "maria@email.com",
    "message": "Account created successfully"
  },
  "message": "Registration successful"
}
```

**Error Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": "Email already registered",
  "code": "EMAIL_EXISTS"
}
```

---

#### POST /api/auth/login
**Request**:
```json
{
  "email": "maria@email.com",
  "password": "SecurePass123"
}
```

**Response (200 OK)** - Session created, cookie set:
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-or-id",
    "email": "maria@email.com",
    "role": "USER"
  },
  "message": "Login successful"
}
```

**Error Response (401 Unauthorized)**:
```json
{
  "success": false,
  "error": "Email not found or password incorrect",
  "code": "AUTH_FAILED"
}
```

---

#### POST /api/auth/logout
**Request**: (no body, requires session)

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {},
  "message": "Logout successful"
}
```

---

#### GET /api/auth/me
**Request**: (requires valid session)

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-or-id",
    "name": "Maria Silva",
    "email": "maria@email.com",
    "role": "USER",
    "created_at": "2026-01-15T10:30:00Z"
  },
  "message": "User info retrieved"
}
```

**Error Response (401 Unauthorized)**: No valid session

---

### Questionnaire & Recommendations

#### GET /api/questions
**Request**: (no body, optional session)

**Response (200 OK)**:
```json
{
  "success": true,
  "data": [
    {
      "id": "q-001",
      "text": "What is your primary use?",
      "type": "SINGLE_CHOICE",
      "weight": 0.2,
      "answer_options": [
        {
          "id": "ao-001",
          "text": "City commute",
          "score_profile": {"urban_score": 0.9, "highway_score": 0.2}
        },
        {
          "id": "ao-002",
          "text": "Mixed city/highway",
          "score_profile": {"urban_score": 0.6, "highway_score": 0.6}
        }
      ]
    }
  ],
  "message": "Questions retrieved"
}
```

---

#### POST /api/questionnaire/answers
**Request** (requires session):
```json
{
  "answers": [
    {
      "question_id": "q-001",
      "answer_option_id": "ao-001"
    },
    {
      "question_id": "q-002",
      "answer_option_id": "ao-005"
    }
  ]
}
```

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-or-id",
    "answers_saved": 2,
    "timestamp": "2026-06-02T14:30:00Z"
  },
  "message": "Answers saved successfully"
}
```

---

#### POST /api/recommendations/generate
**Request** (requires session):
```json
{}
```

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "recommendation_id": "rec-uuid",
    "vehicles": [
      {
        "id": "v-001",
        "brand": "Toyota",
        "model": "Corolla",
        "version": "GLi",
        "year": 2025,
        "score": 92,
        "reason": "Excellent fuel efficiency and reliable transmission match your budget and city usage"
      },
      {
        "id": "v-002",
        "brand": "Honda",
        "model": "Civic",
        "version": "Sport",
        "year": 2025,
        "score": 87,
        "reason": "Good balance between city performance and highway capability"
      }
    ]
  },
  "message": "Recommendations generated"
}
```

---

### Vehicles

#### GET /api/vehicles
**Query Parameters**: `page=1&limit=10&category=sedan`

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "v-001",
        "brand": "Toyota",
        "model": "Corolla",
        "version": "GLi",
        "year": 2025,
        "price_range": "80000-120000",
        "fuel_type": "Gasoline",
        "transmission": "Automatic",
        "seats": 5,
        "category": "Sedan"
      }
    ],
    "total": 42,
    "page": 1,
    "limit": 10
  },
  "message": "Vehicles retrieved"
}
```

---

#### GET /api/vehicles/{id}
**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": "v-001",
    "brand": "Toyota",
    "model": "Corolla",
    "version": "GLi",
    "year": 2025,
    "category": "Sedan",
    "fuel_type": "Gasoline",
    "transmission": "Automatic",
    "price_range": "80000-120000",
    "seats": 5,
    "trunk_capacity": 450,
    "consumption_city": 12.5,
    "consumption_highway": 15.8,
    "description": "Reliable and efficient sedan perfect for daily commuting",
    "strengths": "Great fuel efficiency, reliable transmission, affordable maintenance",
    "weaknesses": "Limited trunk space, smaller interior than competitors"
  },
  "message": "Vehicle details retrieved"
}
```

---

### Admin Endpoints

#### GET /api/admin/vehicles
**Request**: (requires ADMIN role)

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "v-001",
        "brand": "Toyota",
        "model": "Corolla",
        "active": true,
        "year": 2025
      }
    ],
    "total": 127
  },
  "message": "Vehicles retrieved"
}
```

---

#### POST /api/admin/vehicles
**Request** (requires ADMIN role):
```json
{
  "brand": "Honda",
  "model": "City",
  "version": "EX",
  "year": 2025,
  "category_id": "cat-001",
  "fuel_type": "Gasoline",
  "transmission": "Automatic",
  "price_range": "60000-90000",
  "seats": 5,
  "trunk_capacity": 500,
  "consumption_city": 13.0,
  "consumption_highway": 16.5,
  "description": "Compact city car",
  "strengths": "Easy to maneuver, fuel efficient",
  "weaknesses": "Limited trunk space"
}
```

**Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "id": "v-new-uuid",
    "message": "Vehicle created successfully"
  },
  "message": "Vehicle created"
}
```

---

#### PUT /api/admin/vehicles/{id}
**Request** (requires ADMIN role): [same fields as POST, partial update supported]

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": "v-001",
    "message": "Vehicle updated successfully"
  },
  "message": "Vehicle updated"
}
```

---

#### DELETE /api/admin/vehicles/{id}
**Request** (requires ADMIN role): (no body)

**Response (200 OK)**:
```json
{
  "success": true,
  "data": {},
  "message": "Vehicle deleted successfully"
}
```

---

**Note**: Similar endpoint patterns apply for `/api/admin/categories`, `/api/admin/questions`, and `/api/admin/answer-options` (all requiring ADMIN role).

---

## Main User Journeys

### Happy Path: Complete Recommendation Flow (Primary User - Maria)

**Timeline**: Approximately 15-20 minutes for first-time user

1. **Homepage Visit** (1 minute)
   - User accesses `https://app.carro-ideal.com`
   - Sees homepage with hero message and registration/login options
   - Clicks "Sign Up"

2. **Registration** (3 minutes)
   - Form loads: name, email, password, confirm password
   - User fills: "Maria Silva", "maria@email.com", "SecurePass123"
   - Clicks "Create Account"
   - System validates: email unique ✓, password length ✓
   - Account created, redirected to login page
   - Sees success message "Account created. Please log in."

3. **Login** (1 minute)
   - User enters "maria@email.com" and "SecurePass123"
   - Clicks "Login"
   - Session created, redirected to dashboard

4. **Dashboard Navigation** (1 minute)
   - Sees dashboard with "Start Questionnaire" button
   - Clicks button

5. **Questionnaire Completion** (8-10 minutes)
   - Questions display in order:
     - "What is your primary use?" → Selects "Mixed city/highway"
     - "What is your budget?" → Selects "R$ 80.000 - R$ 120.000"
     - "Family size?" → Selects "4-5 people"
     - "Fuel type preference?" → Selects "Gasoline" (cost-effective)
     - "Trunk space importance?" → Selects "Important"
     - "Transmission preference?" → Selects "Automatic"
   - Reviews answers
   - Clicks "Submit Questionnaire"
   - Sees loading indicator "Generating recommendations..."

6. **Recommendations Received** (automatic)
   - System processes answers through recommendation algorithm
   - Displays top 8 vehicles ranked by score (92, 87, 84, 79, 76, 73, 71, 68)
   - Each shows: brand, model, score, concise reason

7. **Vehicle Exploration** (5 minutes)
   - Clicks on first-ranked vehicle (Toyota Corolla GLi, score 92)
   - Details page opens with:
     - Full specifications (year, transmission, fuel type)
     - Price range
     - Consumption data
     - Strengths and weaknesses
     - Recommendation reason with detailed explanation
   - Reviews information
   - Clicks "Back to Recommendations"
   - Clicks on second vehicle (Honda Civic, score 87)
   - Reviews its details
   - Compares mentally with Corolla

8. **Recommendation History Access**
   - User logs out or returns later
   - Logs back in
   - Clicks "Recommendation History"
   - Sees that recommendation from today is listed
   - Can access it again to compare with future recommendations

9. **Next Steps**
   - User gathers information about vehicles
   - Plans visits to dealerships
   - May return to generate new recommendations with refined preferences

---

### Alternative Path 1: User Updates Answers

**Scenario**: Maria wants to change her budget after initial questionnaire

1. User logs in to existing account
2. Navigates to questionnaire
3. Sees previous answers prefilled
4. Changes budget from "R$ 80k-120k" to "R$ 60k-100k"
5. Modifies fuel preference to "Hybrid" for efficiency
6. Submits updated questionnaire
7. System recalculates recommendations with new weights
8. Different vehicle ranking appears with explanation of changes
9. Updated recommendation saved to history with new timestamp

---

### Alternative Path 2: No Vehicles Match Preferences

**Scenario**: Admin creates very restrictive recommendation criteria, or user has extremely specific preferences

1. User completes questionnaire with specific criteria
2. System processes but finds no vehicles scoring above threshold (e.g., min 60/100)
3. System returns graceful message: "No vehicles match your exact criteria. We recommend slightly broader preferences."
4. Offers suggestions: "Adjust fuel type" or "Increase budget"
5. Provides link to modify questionnaire
6. User can adjust and regenerate or browse all vehicles manually

---

### Alternative Path 3: Session Expiry During Questionnaire

**Scenario**: User's session expires mid-questionnaire

1. User is answering questions on questionnaire form
2. After 24 hours of inactivity, session expires
3. User tries to submit questionnaire
4. System detects missing session, returns HTTP 401
5. User redirected to login page with message: "Your session expired. Please log in again."
6. User logs back in
7. Navigates to questionnaire (browser's back button or new navigation)
8. System loads previously saved answers (partial progress)
9. User can resume from question 5 of 8 (not starting from scratch)

---

## Admin Workflow

### Admin User Journey: Initial System Setup

**Timeline**: Approximately 3-4 hours for first-time admin setup (one-time activity)

1. **Admin Access**
   - Admin registers account with admin credentials during setup
   - Logs in at `/admin/login`
   - Redirected to admin dashboard

2. **Create Vehicle Categories** (30 minutes)
   - Clicks "Categories" in admin menu
   - Creates categories:
     - Sedan (description: "4-door passenger cars")
     - SUV (description: "Sport Utility Vehicles")
     - Hatchback (description: "Compact cars with flexible cargo")
     - Pickup (description: "Open cargo trucks")
     - Van (description: "Multi-passenger vehicles")
     - Hybrid/Electric (description: "Alternative fuel vehicles")
   - Each category saved to database

3. **Create Questionnaire Questions** (60 minutes)
   - Clicks "Questions" in admin menu
   - Creates questions with weight and display order:
     - Q1 (order 1, weight 0.15): "What is your primary use?" → Answers: City/Suburban/Highway/Mixed
     - Q2 (order 2, weight 0.20): "What is your budget range?" → Answers: Budget/Mid/Premium
     - Q3 (order 3, weight 0.10): "How many people usually travel?" → Answers: 1-2/3-5/5+
     - Q4 (order 4, weight 0.15): "Fuel type preference?" → Answers: Gasoline/Diesel/Ethanol/Hybrid/Electric
     - Q5 (order 5, weight 0.10): "Trunk space importance?" → Answers: Not important/Important/Critical
     - Q6 (order 6, weight 0.15): "Transmission preference?" → Answers: Manual/Automatic/CVT
     - Q7 (order 7, weight 0.15): "Urban/Highway proportion?" → Scale: 0-100 (city %)
   - Each question with answer options and score profiles

4. **Create Initial Vehicle Catalog** (120+ minutes)
   - Clicks "Vehicles" in admin menu
   - Creates 20-30 sample vehicles with complete data:
     - Toyota Corolla GLi (Sedan, 2025, Automatic, Gasoline, R$85k-110k, good all-around match)
     - Honda Civic Sport (Sedan, 2025, Automatic, Gasoline, R$90k-120k, sporty option)
     - Jeep Compass Limited (SUV, 2025, Automatic, Gasoline, R$120k-160k, hybrid use-case)
     - VW Polo (Hatchback, 2024, Automatic, Gasoline, R$60k-85k, budget-friendly)
     - (etc. - diverse range to test algorithm)
   - Each vehicle with complete specifications, strengths, weaknesses

5. **Configure Answer Scoring Profiles** (30 minutes)
   - For each answer option, sets score_profile JSON
   - Example: "Budget: R$50k-80k" answer sets `{"price_score": 1.0, "luxury_score": 0.3}`
   - This connects user preferences to vehicle attributes in algorithm

6. **Test Recommendation Algorithm**
   - Creates test user account with known preferences
   - Completes questionnaire with specific answers
   - Generates recommendations
   - Verifies algorithm produces expected results
   - Documents algorithm rules in `/docs/ALGORITHM.md`

---

### Admin User Journey: Ongoing Maintenance

**Timeline**: Weekly 1-2 hour sessions for updates

**Weekly Tasks**:
1. Review user feedback (future feature: feedback collection)
2. Add newly released vehicles to catalog
3. Update prices for existing vehicles
4. Mark old model years as inactive (soft delete)
5. Monitor recommendation accuracy (view recommendations in dashboard)
6. Adjust question weights if certain recommendations seem off

**Monthly Tasks**:
1. Generate admin report: total users, total recommendations, popular vehicles
2. Refine recommendation algorithm if needed (adjust weights, add new scoring rules)
3. Review and enhance answer scoring profiles
4. Add new vehicle categories if market adds new types

---

## Success Criteria (Measurable Outcomes)

### User-Focused Success Metrics

- **SC-001**: First-time user can create account and receive first recommendation in under 10 minutes (reduces friction)
- **SC-002**: 95% of questionnaires submitted successfully without data loss or session errors (reliability)
- **SC-003**: Recommended vehicles match user preferences with 80%+ relevance (measured by user satisfaction survey)
- **SC-004**: Vehicle detail pages load in under 500ms (performance)
- **SC-005**: Recommendation list displays in under 2 seconds from questionnaire submission (user experience)

### Platform Stability Metrics

- **SC-006**: System handles 100+ concurrent users without response degradation (scalability)
- **SC-007**: Database supports 10,000+ vehicle records and 50,000+ user recommendations without performance loss (data volume)
- **SC-008**: Authentication success rate 99.8% or higher (reliability)
- **SC-009**: API endpoints return 99.9% uptime during production week (availability)
- **SC-010**: Error logging captures 100% of server errors for debugging (observability)

### Admin Effectiveness Metrics

- **SC-011**: Admin can add new vehicle to system in under 5 minutes via web form (usability)
- **SC-012**: Admin dashboard loads with current statistics within 1 second (performance)
- **SC-013**: Category/Question/Answer CRUD operations complete without data validation errors (data integrity)

### Academic/Project Evaluation Metrics

- **SC-014**: Codebase follows layered architecture (handlers → services → repositories → database) with clear file structure and documentation (architectural excellence)
- **SC-015**: README contains all installation, execution, and deployment instructions enabling 1-command launch (documentation quality)
- **SC-016**: Git repository contains minimum 15-20 meaningfully-documented commits showing incremental development (development process)
- **SC-017**: Code passes linting standards (gofmt, golangci-lint) and has zero hardcoded values for configuration (code quality)
- **SC-018**: Database migrations are version-controlled, immutable, and fully reversible (migration best practices)
- **SC-019**: Project demonstrates understanding of recommendation algorithms through code comments and documentation (knowledge demonstration)
- **SC-020**: Public repository (GitHub/GitLab) is accessible with clear project description and active development history (academic submission)

### Feature Completion Metrics

- **SC-021**: All 15 mandatory functional requirements implemented and testable (feature completeness)
- **SC-022**: All 13 minimum pages/screens exist and render correctly (user interface)
- **SC-023**: All 40+ API endpoints functional and return correct responses (API coverage)

### Data Quality Metrics

- **SC-024**: Database contains at minimum 20 vehicles across 5+ categories for demo (sample data)
- **SC-025**: At least 7 active questionnaire questions available at launch (questionnaire depth)
- **SC-026**: All vehicles have complete, accurate data (no missing required fields) (data completeness)

---

## Assumptions

### Technical Assumptions

- **Layered Architecture**: System will be organized in Go with clear handler → service → repository → database layers per project constitution
- **PostgreSQL 15+**: Database is hosted and managed via Docker container, with persistent volumes configured
- **Session-Based Auth**: Authentication uses server-side sessions with secure cookies (JWT not required for MVP)
- **Synchronous Processing**: Recommendation generation is synchronous (processed immediately on form submit, no background jobs initially)
- **Score-Based Algorithm**: Initial recommendation algorithm uses weighted scoring (not machine learning/AI), with rules documented
- **Single Go Binary**: Application compiles to single executable and runs as single Docker container
- **No Third-Party Auth**: No OAuth2, SSO, or third-party authentication providers integrated in MVP
- **Development Environment**: Docker Compose runs entire stack (app + database) on developer machine

### Data Assumptions

- **Vehicle Data**: Approximately 20-30 sample vehicles provided for demo; admin can add more
- **User Data**: Initial demo with admin-created test users; real registration available
- **Questionnaire**: 7-8 questions defined; admin can modify questions and weights
- **Recommendation Accuracy**: Algorithm is rule-based and scores vehicles based on preference matching; 80%+ accuracy is target (user feedback required for validation)

### User Assumptions

- **Primary Audience**: Brazilian users interested in car purchasing decisions (Portuguese language assumed)
- **Technical Literacy**: Users comfortable with web forms and basic navigation; no extensive system training needed
- **Internet Connectivity**: Users have stable internet; no offline functionality required
- **Desktop/Tablet**: Primary access via desktop browsers; mobile not required for MVP
- **Time Commitment**: Users willing to spend 10-15 minutes completing questionnaire for personalized results

### Admin Assumptions

- **Single Admin Role**: One primary admin managing system during MVP phase (no admin hierarchy or team)
- **Domain Knowledge**: Admin understands Brazilian car market and can manage vehicle database
- **Maintenance Cadence**: Admin performs weekly updates to vehicle catalog, monthly algorithm reviews
- **Learning Curve**: Admin has basic web UI navigation skills; CRUD operations intuitive enough

### Scope Boundaries

- **No Real-Time Features**: Recommendations calculated on-demand, not streamed or real-time
- **No Advanced Analytics**: Dashboard shows basic stats (user count, recommendation count), not detailed analytics
- **No Internationalization**: Portuguese-only for MVP (English support deferred)
- **No Payment**: No checkout or payment processing in MVP
- **No Email Notifications**: No automated emails sent (user-triggered actions only)
- **No Social Features**: No sharing, reviews, or social comparison in MVP
- **No Mobile App**: Web app responsive but not native mobile application
- **No AI Integration**: Architecture prepared for ChatGPT, but synchronous rule-based algorithm used in MVP

### Timeline Assumptions

- **Development Duration**: 8-12 weeks estimated for full MVP delivery
- **Launch Readiness**: Project must be deployment-ready for TCC presentation and public repository publication
- **Iteration**: After MVP launch, feedback collection and Phase 2 planning for future enhancements

### Deployment Assumptions

- **Single Server Deployment**: Initial deployment on single server with single app instance and single database instance
- **Docker Compose Orchestration**: Docker Compose for local and demo deployment; Kubernetes not required
- **Public Repository**: Code published to public GitHub repository with README and documentation
- **No CI/CD Pipeline**: Initial deployment manual; automated CI/CD (GitHub Actions) optional for Phase 2

---

## Next Steps After Specification

1. **Planning Phase** (`/speckit.plan`): Break specification into work streams and define execution strategy
2. **Task Generation** (`/speckit.tasks`): Convert requirements into granular, assignable development tasks
3. **Implementation** (`/speckit.implement`): Develop features incrementally following layered architecture
4. **Quality Validation** (`/speckit.checklist`): Verify specification requirements met through testing
5. **Documentation**: Update README with execution instructions and API documentation
6. **Publication**: Push to public GitHub repository with project description and capstone context
