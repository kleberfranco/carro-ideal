# REST API Contract: Carro Ideal MVP

**Version**: 1.0  
**Status**: Phase 1 Design Complete  
**Base URL**: `/api/v1` (or `/api` for MVP)

---

## Overview

All endpoints follow RESTful conventions with standardized JSON responses. Authentication via HTTP-only session cookies.

### Response Format

**Success Response** (2xx):
```json
{
  "success": true,
  "data": {/* resource or array of resources */},
  "message": "Human-readable success message"
}
```

**Error Response** (4xx, 5xx):
```json
{
  "success": false,
  "error": "Specific error message for client",
  "code": "ERROR_CODE"
}
```

### Error Codes

| HTTP | Code | Meaning |
|------|------|---------|
| 400 | INVALID_INPUT | Validation failed (user error) |
| 400 | EMAIL_EXISTS | Email already registered |
| 400 | PASSWORD_TOO_WEAK | Password doesn't meet requirements |
| 401 | AUTH_FAILED | Email/password incorrect |
| 401 | UNAUTHENTICATED | No valid session |
| 403 | UNAUTHORIZED | Authenticated but insufficient role |
| 404 | NOT_FOUND | Resource doesn't exist |
| 500 | INTERNAL_ERROR | Server error (no details exposed) |

---

## Authentication Endpoints

### POST /api/auth/register

Register a new user account.

**Request**:
```json
{
  "name": "Maria Silva",
  "email": "maria@email.com",
  "password": "SecurePass123!"
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "user_id": 42,
    "email": "maria@email.com",
    "name": "Maria Silva"
  },
  "message": "Account created successfully"
}
```

**Error Response** (400):
```json
{
  "success": false,
  "error": "Email already registered",
  "code": "EMAIL_EXISTS"
}
```

**Validation Rules**:
- `email`: Valid RFC 5322 format, unique in database
- `name`: 1-100 characters, non-empty
- `password`: Minimum 8 characters

**HTTP-Only Cookie**: Session cookie automatically set on successful registration.

---

### POST /api/auth/login

Authenticate user and create session.

**Request**:
```json
{
  "email": "maria@email.com",
  "password": "SecurePass123!"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user_id": 42,
    "email": "maria@email.com",
    "name": "Maria Silva",
    "role": "USER"
  },
  "message": "Login successful"
}
```

**Error Response** (401):
```json
{
  "success": false,
  "error": "Email not found or password incorrect",
  "code": "AUTH_FAILED"
}
```

**Notes**:
- Session cookie set for 24-hour duration
- Cookie is HTTP-only (not accessible to JavaScript)
- Generic error message for security (doesn't distinguish "email not found" vs "wrong password")

---

### POST /api/auth/logout

Destroy user session.

**Request**: No body required (session from cookie)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {},
  "message": "Logged out successfully"
}
```

**Notes**:
- Session cookie cleared on client
- User redirected to homepage

---

### GET /api/auth/me

Get current authenticated user.

**Request**: No body (authenticated via session cookie)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user_id": 42,
    "email": "maria@email.com",
    "name": "Maria Silva",
    "role": "USER"
  },
  "message": "Current user"
}
```

**Error Response** (401):
```json
{
  "success": false,
  "error": "Unauthenticated",
  "code": "UNAUTHENTICATED"
}
```

**Usage**: Frontend calls this on page load to verify session and update user context.

---

## User Endpoints

### GET /api/user/dashboard

Get authenticated user's dashboard data.

**Request**: Authenticated session required

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "user_id": 42,
      "email": "maria@email.com",
      "name": "Maria Silva"
    },
    "stats": {
      "total_recommendations": 5,
      "last_recommendation": "2024-06-01T12:00:00Z",
      "answered_questions": true
    }
  },
  "message": "Dashboard data"
}
```

**Error Response** (401): Unauthenticated

---

### GET /api/user/recommendations

Get user's recommendation history (paginated).

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "recommendation_id": 1,
        "created_at": "2024-06-01T12:00:00Z",
        "vehicle_count": 8,
        "summary": "Based on your budget and efficiency preferences..."
      },
      {
        "recommendation_id": 2,
        "created_at": "2024-05-28T09:30:00Z",
        "vehicle_count": 7,
        "summary": "..."
      }
    ],
    "total": 5,
    "page": 1,
    "limit": 10
  },
  "message": "Recommendation history"
}
```

**Error Response** (401): Unauthenticated

---

### GET /api/user/recommendations/{id}

Get specific recommendation with detailed results.

**Path Parameters**:
- `id` (required): Recommendation ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "recommendation": {
      "recommendation_id": 1,
      "created_at": "2024-06-01T12:00:00Z",
      "summary": "Based on your preferences, these vehicles are ideal"
    },
    "items": [
      {
        "rank": 1,
        "vehicle_id": 1,
        "brand": "Toyota",
        "model": "Corolla",
        "version": "GLi 1.6",
        "year": 2023,
        "fuel_type": "Gasoline",
        "transmission": "Automatic",
        "price_range": "80000-120000",
        "seats": 5,
        "trunk_capacity": 300,
        "consumption_city": 9.5,
        "consumption_highway": 13.2,
        "strengths": "Reliable, fuel-efficient, good resale value",
        "weaknesses": "Limited cargo space",
        "score": 92.5,
        "reason": "Excellent match on budget, fuel efficiency, and reliability"
      },
      {
        "rank": 2,
        "vehicle_id": 3,
        "brand": "Honda",
        "model": "Civic",
        "score": 88.0,
        "reason": "..."
      }
    ]
  },
  "message": "Recommendation details"
}
```

**Error Response** (404): Recommendation not found  
**Error Response** (401): Unauthenticated

---

## Questionnaire Endpoints

### GET /api/questions

Get all active questions with answer options.

**Query Parameters**:
- `active` (optional, default: true): Filter by active status

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "question_id": 1,
        "text": "What is your budget (BRL)?",
        "type": "SINGLE_CHOICE",
        "weight": 0.4,
        "display_order": 1,
        "answer_options": [
          {
            "option_id": 1,
            "text": "R$ 30.000 - R$ 50.000",
            "score_profile": {
              "budget_score": 0.3,
              "efficiency_score": 0.7,
              "comfort_score": 0.4
            }
          },
          {
            "option_id": 2,
            "text": "R$ 50.000 - R$ 80.000",
            "score_profile": {
              "budget_score": 0.7,
              "efficiency_score": 0.5,
              "comfort_score": 0.6
            }
          }
        ]
      },
      {
        "question_id": 2,
        "text": "How important is fuel efficiency?",
        "type": "SINGLE_CHOICE",
        "weight": 0.3,
        "display_order": 2,
        "answer_options": [
          {
            "option_id": 5,
            "text": "Very important (top priority)",
            "score_profile": {"efficiency_score": 0.9}
          }
        ]
      }
    ]
  },
  "message": "Active questions"
}
```

---

### POST /api/questionnaire/answers

Store or update user's questionnaire answers.

**Request**:
```json
{
  "answers": [
    {
      "question_id": 1,
      "answer_option_id": 2
    },
    {
      "question_id": 2,
      "answer_option_id": 5
    }
  ]
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user_id": 42,
    "answers_saved": 2,
    "timestamp": "2024-06-01T11:30:00Z"
  },
  "message": "Answers saved successfully"
}
```

**Error Response** (400): Missing required question or invalid option  
**Error Response** (401): Unauthenticated

**Validation**:
- All answer_option_ids must reference valid AnswerOption records
- All question_ids must reference active questions
- Transaction: all answers saved or none (atomicity)

---

## Recommendation Endpoints

### POST /api/recommendations/generate

Generate recommendations based on user's latest questionnaire answers.

**Request**:
```json
{}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "recommendation_id": 42,
    "created_at": "2024-06-01T12:00:00Z",
    "vehicles": [
      {
        "rank": 1,
        "vehicle_id": 1,
        "brand": "Toyota",
        "model": "Corolla",
        "version": "GLi 1.6",
        "year": 2023,
        "fuel_type": "Gasoline",
        "transmission": "Automatic",
        "price_range": "80000-120000",
        "seats": 5,
        "trunk_capacity": 300,
        "consumption_city": 9.5,
        "consumption_highway": 13.2,
        "category": "Sedan",
        "strengths": "Reliable, fuel-efficient, good resale value",
        "weaknesses": "Limited cargo space",
        "score": 92.5,
        "reason": "Excellent match: budget fits, fuel efficiency aligns with preference, reliable brand"
      },
      {
        "rank": 2,
        "vehicle_id": 3,
        "brand": "Honda",
        "model": "Civic",
        "score": 88.0,
        "reason": "..."
      },
      {
        "rank": 3,
        "vehicle_id": 5,
        "brand": "Hyundai",
        "model": "Elantra",
        "score": 81.5,
        "reason": "..."
      }
    ]
  },
  "message": "Recommendations generated"
}
```

**Error Response** (400):
```json
{
  "success": false,
  "error": "No questions answered yet. Please complete the questionnaire first.",
  "code": "NO_ANSWERS"
}
```

**Error Response** (400):
```json
{
  "success": false,
  "error": "No vehicles match your preferences",
  "code": "NO_MATCHING_VEHICLES"
}
```

**Error Response** (401): Unauthenticated

**Behavior**:
- Fetches user's latest answers
- Loads all active vehicles
- Calculates recommendation scores for each vehicle
- Returns top 10 vehicles sorted by score descending
- Stores recommendation in database for history

---

## Vehicle Endpoints

### GET /api/vehicles

List active vehicles (paginated, filterable).

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 25): Items per page
- `category_id` (optional): Filter by category
- `fuel_type` (optional): Filter by fuel type
- `search` (optional): Search in brand/model

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "vehicle_id": 1,
        "brand": "Toyota",
        "model": "Corolla",
        "version": "GLi 1.6",
        "year": 2023,
        "fuel_type": "Gasoline",
        "transmission": "Automatic",
        "price_range": "80000-120000",
        "seats": 5,
        "trunk_capacity": 300,
        "consumption_city": 9.5,
        "consumption_highway": 13.2,
        "category": "Sedan"
      }
    ],
    "total": 150,
    "page": 1,
    "limit": 25
  },
  "message": "Active vehicles"
}
```

---

### GET /api/vehicles/{id}

Get vehicle details.

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "vehicle_id": 1,
    "brand": "Toyota",
    "model": "Corolla",
    "version": "GLi 1.6",
    "year": 2023,
    "fuel_type": "Gasoline",
    "transmission": "Automatic",
    "price_range": "80000-120000",
    "seats": 5,
    "trunk_capacity": 300,
    "consumption_city": 9.5,
    "consumption_highway": 13.2,
    "category": "Sedan",
    "description": "Reliable daily driver with excellent fuel economy",
    "strengths": "Fuel-efficient, reliable, good resale value, comfortable interior",
    "weaknesses": "Limited cargo space, not ideal for off-road driving",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-06-01T14:30:00Z"
  },
  "message": "Vehicle details"
}
```

**Error Response** (404): Vehicle not found

---

## Admin Endpoints

### GET /api/admin/vehicles

List all vehicles including inactive (admin only).

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "vehicle_id": 1,
        "brand": "Toyota",
        "model": "Corolla",
        "year": 2023,
        "price_range": "80000-120000",
        "active": true,
        "recommendations_count": 15
      },
      {
        "vehicle_id": 2,
        "brand": "Honda",
        "model": "Civic",
        "year": 2023,
        "price_range": "85000-125000",
        "active": false,
        "recommendations_count": 3
      }
    ],
    "total": 152,
    "page": 1,
    "limit": 25
  },
  "message": "All vehicles"
}
```

**Error Response** (403): Unauthorized (non-admin)

---

### POST /api/admin/vehicles

Create new vehicle (admin only).

**Request**:
```json
{
  "category_id": 1,
  "brand": "Toyota",
  "model": "Corolla",
  "version": "GLi 1.6",
  "year": 2023,
  "fuel_type": "Gasoline",
  "transmission": "Automatic",
  "price_range": "80000-120000",
  "seats": 5,
  "trunk_capacity": 300,
  "consumption_city": 9.5,
  "consumption_highway": 13.2,
  "description": "Reliable daily driver",
  "strengths": "Fuel-efficient, reliable",
  "weaknesses": "Limited cargo space"
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "vehicle_id": 153,
    "brand": "Toyota",
    "model": "Corolla",
    "created_at": "2024-06-01T15:00:00Z"
  },
  "message": "Vehicle created"
}
```

**Error Response** (400): Validation failure  
**Error Response** (403): Unauthorized

---

### PUT /api/admin/vehicles/{id}

Update vehicle (admin only).

**Request**: Same fields as POST (all optional for partial update)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "vehicle_id": 153,
    "updated_fields": ["consumption_city", "strengths"],
    "updated_at": "2024-06-01T15:05:00Z"
  },
  "message": "Vehicle updated"
}
```

---

### DELETE /api/admin/vehicles/{id}

Soft-delete vehicle (admin only).

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "vehicle_id": 153,
    "status": "inactive"
  },
  "message": "Vehicle deleted"
}
```

---

### GET /api/admin/categories

List all categories (admin only).

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "category_id": 1,
        "name": "Sedan",
        "description": "Passenger cars with trunk",
        "active": true,
        "vehicle_count": 45
      },
      {
        "category_id": 2,
        "name": "SUV",
        "description": "Sport Utility Vehicles",
        "active": true,
        "vehicle_count": 38
      }
    ]
  },
  "message": "All categories"
}
```

---

### POST /api/admin/categories

Create new category (admin only).

**Request**:
```json
{
  "name": "Pickup Truck",
  "description": "Vehicles designed for cargo transport"
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "category_id": 8,
    "name": "Pickup Truck",
    "created_at": "2024-06-01T15:10:00Z"
  },
  "message": "Category created"
}
```

---

### PUT /api/admin/categories/{id}

Update category (admin only).

**Response** (200 OK): Similar to POST

---

### DELETE /api/admin/categories/{id}

Delete category (admin only) - only if no active vehicles linked.

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "category_id": 8
  },
  "message": "Category deleted"
}
```

**Error Response** (400):
```json
{
  "success": false,
  "error": "Cannot delete category with 5 active vehicles linked",
  "code": "CATEGORY_IN_USE"
}
```

---

### GET /api/admin/questions

List all questions (admin only).

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "question_id": 1,
        "text": "What is your budget?",
        "type": "SINGLE_CHOICE",
        "weight": 0.4,
        "display_order": 1,
        "active": true,
        "answer_options_count": 4
      }
    ]
  },
  "message": "All questions"
}
```

---

### POST /api/admin/questions

Create new question (admin only).

**Request**:
```json
{
  "text": "How many passengers typically?",
  "type": "SINGLE_CHOICE",
  "weight": 0.2,
  "display_order": 3
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "question_id": 11,
    "text": "How many passengers typically?",
    "created_at": "2024-06-01T15:15:00Z"
  },
  "message": "Question created"
}
```

---

### PUT /api/admin/questions/{id}

Update question (admin only).

---

### DELETE /api/admin/questions/{id}

Delete question (admin only) - soft-delete via active flag.

---

### GET /api/admin/answer-options

List all answer options (admin only).

---

### POST /api/admin/answer-options

Create answer option (admin only).

**Request**:
```json
{
  "question_id": 1,
  "text": "2-3 passengers",
  "score_profile": {
    "comfort_score": 0.5,
    "efficiency_score": 0.7
  }
}
```

---

### PUT /api/admin/answer-options/{id}

Update answer option (admin only).

---

### DELETE /api/admin/answer-options/{id}

Delete answer option (admin only).

---

## Utility Endpoints

### GET /health

Health check (no authentication required).

**Response** (200 OK):
```json
{
  "status": "healthy",
  "timestamp": "2024-06-01T15:20:00Z"
}
```

---

## Authentication Notes

- All endpoints except `/health`, `/api/auth/register`, `/api/auth/login`, `/api/vehicles`, `/api/questions` require valid session cookie
- Admin endpoints require session with role = "ADMIN"
- Sessions expire after 24 hours
- Session cookie is HTTP-only and cannot be accessed by JavaScript (XSS protection)

---

**API Contract Status**: ✅ COMPLETE - Ready for Implementation