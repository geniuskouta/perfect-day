# REST API Endpoint Contracts

## API Overview

All endpoints follow the pattern: `/api/v1/<resource>`

**Base URL**: `http://localhost:8080` (local development)
**Content-Type**: `application/json`
**Authentication**: Session-based (for local development)

## Global Response Format

### Success Response
```json
{
  "data": { ... },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "0.1.0"
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": { ... }
  },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "0.1.0"
  }
}
```

## Authentication Endpoints

### POST /api/v1/auth/login
**Purpose**: Authenticate user (username-based, matching CLI)

**Request**:
```json
{
  "username": "kouta"
}
```

**Response 200**:
```json
{
  "data": {
    "user": {
      "username": "kouta",
      "timezone": "Asia/Tokyo",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "session": {
      "id": "session-123",
      "expires_at": "2024-01-02T00:00:00Z"
    }
  }
}
```

**Response 401**:
```json
{
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User not found"
  }
}
```

### GET /api/v1/auth/me
**Purpose**: Get current authenticated user info

**Headers**: `Cookie: session_id=session-123`

**Response 200**:
```json
{
  "data": {
    "username": "kouta",
    "timezone": "Asia/Tokyo",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response 401**:
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Not authenticated"
  }
}
```

## Perfect Days Endpoints

### GET /api/v1/perfect-days
**Purpose**: List and search perfect days

**Query Parameters**:
- `q` (string): Search query
- `areas` (string): Comma-separated areas
- `user` (string): Filter by username
- `from` (string): Date from (YYYY-MM-DD)
- `to` (string): Date to (YYYY-MM-DD)
- `sort` (string): Sort by (date, created_at, title)
- `order` (string): Sort order (asc, desc)
- `limit` (int): Results per page (default: 10)
- `offset` (int): Results offset (default: 0)

**Response 200**:
```json
{
  "data": {
    "perfect_days": [
      {
        "id": "uuid-123",
        "title": "Perfect Tokyo Day",
        "description": "Amazing day exploring Tokyo",
        "username": "kouta",
        "date": "2024-01-01",
        "areas": ["Shibuya", "Harajuku"],
        "activities": [
          {
            "name": "Visit Senso-ji Temple",
            "location": {
              "type": "google_place",
              "place_id": "ChIJ...",
              "name": "Senso-ji Temple",
              "area": "Asakusa"
            },
            "start_time": "09:00",
            "duration": 120,
            "commentary": "Beautiful morning visit"
          }
        ],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "total": 42,
      "offset": 0,
      "limit": 10,
      "has_more": true
    }
  }
}
```

### POST /api/v1/perfect-days
**Purpose**: Create new perfect day

**Headers**: `Cookie: session_id=session-123`

**Request**:
```json
{
  "title": "Perfect Tokyo Day",
  "description": "Amazing day exploring Tokyo",
  "date": "2024-01-01",
  "activities": [
    {
      "name": "Visit Senso-ji Temple",
      "description": "Historic temple visit",
      "location": {
        "type": "google_place",
        "place_id": "ChIJ...",
        "name": "Senso-ji Temple",
        "area": "Asakusa"
      },
      "start_time": "09:00",
      "duration": 120,
      "commentary": "Beautiful morning visit"
    }
  ]
}
```

**Response 201**:
```json
{
  "data": {
    "id": "uuid-123",
    "title": "Perfect Tokyo Day",
    "description": "Amazing day exploring Tokyo",
    "username": "kouta",
    "date": "2024-01-01",
    "areas": ["Asakusa"],
    "activities": [...],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response 400**:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "title": "Title is required"
    }
  }
}
```

### GET /api/v1/perfect-days/{id}
**Purpose**: Get specific perfect day

**Response 200**: Same as POST 201 response

**Response 404**:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Perfect day not found"
  }
}
```

### PUT /api/v1/perfect-days/{id}
**Purpose**: Update perfect day

**Headers**: `Cookie: session_id=session-123`

**Request**: Same as POST request

**Response 200**: Same as POST 201 response

**Response 403**:
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Can only edit your own perfect days"
  }
}
```

### DELETE /api/v1/perfect-days/{id}
**Purpose**: Soft delete perfect day

**Headers**: `Cookie: session_id=session-123`

**Response 204**: No content

**Response 403**: Same as PUT 403 response

## Users Endpoints

### GET /api/v1/users/{username}
**Purpose**: Get public user profile

**Response 200**:
```json
{
  "data": {
    "username": "kouta",
    "timezone": "Asia/Tokyo",
    "created_at": "2024-01-01T00:00:00Z",
    "stats": {
      "total_perfect_days": 15,
      "total_activities": 87,
      "favorite_areas": ["Shibuya", "Harajuku", "Shinjuku"]
    }
  }
}
```

### GET /api/v1/users/{username}/perfect-days
**Purpose**: Get user's public perfect days

**Query Parameters**: Same as GET /api/v1/perfect-days

**Response 200**: Same as GET /api/v1/perfect-days

## Places Endpoints

### GET /api/v1/places/search
**Purpose**: Proxy to Google Places API

**Query Parameters**:
- `q` (string, required): Search query

**Response 200**:
```json
{
  "data": {
    "places": [
      {
        "place_id": "ChIJ...",
        "name": "Senso-ji Temple",
        "address": "2 Chome-3-1 Asakusa, Taito City, Tokyo",
        "latitude": 35.7147651,
        "longitude": 139.7966831
      }
    ]
  }
}
```

**Response 400**:
```json
{
  "error": {
    "code": "MISSING_QUERY",
    "message": "Search query is required"
  }
}
```

### GET /api/v1/areas
**Purpose**: Get popular areas from existing data

**Response 200**:
```json
{
  "data": {
    "areas": [
      {
        "name": "Shibuya",
        "count": 25,
        "recent_activity": "2024-01-01T00:00:00Z"
      },
      {
        "name": "Harajuku",
        "count": 18,
        "recent_activity": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

## System Endpoints

### GET /api/v1/health
**Purpose**: API health check

**Response 200**:
```json
{
  "data": {
    "status": "healthy",
    "uptime_seconds": 3600,
    "version": "0.1.0",
    "checks": {
      "storage": "ok",
      "google_places": "ok"
    }
  }
}
```

### GET /api/v1/version
**Purpose**: API version information

**Response 200**:
```json
{
  "data": {
    "version": "0.1.0",
    "build": "abc123",
    "go_version": "1.21.0",
    "built_at": "2024-01-01T00:00:00Z"
  }
}
```

## HTTP Status Codes

- **200 OK**: Successful GET/PUT
- **201 Created**: Successful POST
- **204 No Content**: Successful DELETE
- **400 Bad Request**: Invalid request data
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Permission denied
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

## Notes

- All timestamps are in ISO 8601 format (UTC)
- Session cookies are httpOnly and secure in production
- CORS headers allow localhost origins for development
- Rate limiting: 100 requests/minute per IP (development)
- Request/response logging for debugging
- Graceful error handling with consistent format