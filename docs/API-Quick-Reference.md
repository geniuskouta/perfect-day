# Perfect Day API - Quick Reference

## Base URL
`http://localhost:8080/api/v1`

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/version` | API version |
| GET | `/perfect-days` | List perfect days |
| POST | `/perfect-days` | Create perfect day |
| GET | `/perfect-days/{id}` | Get perfect day |
| PUT | `/perfect-days/{id}` | Update perfect day |
| DELETE | `/perfect-days/{id}` | Delete perfect day |

## Quick Examples

### Create Perfect Day
```bash
curl -X POST http://localhost:8080/api/v1/perfect-days \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Amazing Tokyo Day",
    "date": "2025-01-15",
    "activities": [
      {
        "name": "Visit Temple",
        "location": {"type": "custom_text", "name": "Senso-ji", "area": "Asakusa"},
        "start_time": "09:00",
        "duration": 120,
        "commentary": "Beautiful experience"
      }
    ]
  }'
```

### List Perfect Days
```bash
# Basic list
curl http://localhost:8080/api/v1/perfect-days

# With filtering
curl "http://localhost:8080/api/v1/perfect-days?user=kouta&limit=5&sort=date"

# Search
curl "http://localhost:8080/api/v1/perfect-days?q=tokyo&areas=Shibuya"
```

### Get Perfect Day
```bash
curl http://localhost:8080/api/v1/perfect-days/{id}
```

### Update Perfect Day
```bash
curl -X PUT http://localhost:8080/api/v1/perfect-days/{id} \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Title", "date": "2025-01-15", "activities": []}'
```

### Delete Perfect Day
```bash
curl -X DELETE http://localhost:8080/api/v1/perfect-days/{id}
```

## Response Format
All responses return JSON with `data` and `meta` fields:
```json
{
  "data": { /* actual response data */ },
  "meta": {
    "timestamp": "2025-01-01T00:00:00Z",
    "version": "0.1.0"
  }
}
```

## Common Query Parameters
- `q` - Search query
- `user` - Filter by username
- `areas` - Filter by area
- `from` / `to` - Date range (YYYY-MM-DD)
- `sort` - Sort by (`date`, `created_at`, `title`)
- `order` - Sort order (`asc`, `desc`)
- `limit` - Results per page
- `offset` - Results offset

## Location Types
```json
// Custom text location
{"type": "custom_text", "name": "My Cafe", "area": "Shibuya"}

// Google Place location
{"type": "google_place", "place_id": "ChIJ...", "name": "Blue Bottle", "area": "Shibuya"}
```

## Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

## HTTP Status Codes
- `200` - Success (GET/PUT)
- `201` - Created (POST)
- `204` - No Content (DELETE)
- `400` - Bad Request
- `404` - Not Found
- `500` - Server Error