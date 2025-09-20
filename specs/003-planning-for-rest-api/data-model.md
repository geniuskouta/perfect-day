# Data Model: Perfect Day Terminal App

## Core Entities

### User
Represents a user of the Perfect Day system.

**Fields**:
- `username` (string, required): Unique identifier for the user (FR-010)
- `created_at` (time.Time, required): Account creation timestamp
- `timezone` (string, optional): Default timezone for the user (defaults to system)

**Validation Rules**:
- Username must be 3-50 characters, alphanumeric + underscore/dash
- Username must be unique across the system
- Timezone must be valid IANA timezone string

**Relationships**:
- One-to-many with PerfectDay

**State Transitions**: None (static after creation)

### PerfectDay
Represents a complete day plan shared by a user.

**Fields**:
- `id` (string, required): UUID v4 identifier
- `username` (string, required): Foreign key to User
- `title` (string, required): Day title, 5-100 characters (FR-001)
- `description` (string, required): Summary description, 10-1000 characters (FR-002)
- `date` (time.Time, required): The date this perfect day represents
- `main_areas` ([]string, required): List of main areas visited (FR-003)
- `activities` ([]Activity, required): List of activities in chronological order
- `is_deleted` (bool, optional): Soft delete flag (FR-012), defaults to false
- `created_at` (time.Time, required): Creation timestamp
- `updated_at` (time.Time, required): Last modification timestamp
- `version` (int, required): Version number for optimistic locking

**Validation Rules**:
- Title must be 5-100 characters, non-empty after trimming
- Description must be 10-1000 characters, non-empty after trimming
- Date must not be in the future (beyond today)
- Must have at least 1 activity
- Activities must not have overlapping time ranges
- Main areas must have at least 1 entry, max 10 entries
- Each main area must be 2-50 characters

**Relationships**:
- Many-to-one with User
- One-to-many with Activity (embedded)

**State Transitions**:
- Draft → Published (when first saved)
- Published → Updated (when modified, FR-012)
- Any → Deleted (soft delete, FR-012)

### Activity
Represents an individual activity within a perfect day.

**Fields**:
- `id` (string, required): UUID v4 identifier
- `location` (Location, required): Where the activity took place
- `start_time` (time.Time, required): When the activity started (FR-007)
- `end_time` (time.Time, required): When the activity ended (FR-007)
- `comment` (string, required): Personal commentary about the activity (FR-006)
- `order` (int, required): Position in the day's timeline (1-based)

**Validation Rules**:
- Start time must be before end time
- Duration must be at least 5 minutes, at most 8 hours
- Comment must be 5-500 characters, non-empty after trimming
- Order must be positive integer
- Activities within a day must have unique, sequential order values

**Relationships**:
- Many-to-one with PerfectDay (embedded)
- One-to-one with Location (embedded)

**State Transitions**: None (managed as part of PerfectDay updates)

### Location
Represents a location where an activity took place.

**Fields**:
- `type` (string, required): "google_business" or "custom_text"
- `google_place_id` (string, optional): Google Places API place ID (FR-004)
- `google_name` (string, optional): Name from Google Business Profile
- `google_address` (string, optional): Address from Google Business Profile
- `custom_name` (string, optional): User-entered location name (FR-005)
- `area` (string, required): Geographic area/neighborhood

**Validation Rules**:
- Type must be either "google_business" or "custom_text"
- If type is "google_business": google_place_id, google_name, and google_address are required
- If type is "custom_text": custom_name is required (5-100 characters)
- Area is always required (2-50 characters)
- Exactly one of google_name or custom_name must be populated

**Relationships**:
- One-to-one with Activity (embedded)

**State Transitions**: None (managed as part of Activity updates)

## Data Storage Structure

### File System Layout
```
~/.perfectday/
├── config.json                    # App configuration
├── users.json                     # User registry
└── days/                           # Perfect days by user
    ├── kouta/
    │   ├── 2025-01-15_morning-in-shibuya.json
    │   └── 2025-01-20_tokyo-foodie-adventure.json
    └── alice/
        └── 2025-02-01_art-gallery-hopping.json
```

### JSON Schema Examples

**User Document (users.json)**:
```json
{
  "users": [
    {
      "username": "kouta",
      "created_at": "2025-09-15T10:00:00+09:00",
      "timezone": "Asia/Tokyo"
    }
  ]
}
```

**PerfectDay Document**:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "username": "kouta",
  "title": "Morning in Shibuya",
  "description": "A perfect morning exploring the bustling heart of Tokyo, from quiet coffee to vibrant streets.",
  "date": "2025-01-15T00:00:00+09:00",
  "main_areas": ["Shibuya", "Harajuku"],
  "activities": [
    {
      "id": "activity-1",
      "location": {
        "type": "google_business",
        "google_place_id": "ChIJ_____example",
        "google_name": "Blue Bottle Coffee Shibuya",
        "google_address": "Tokyo, Shibuya City, Shibuya, 2 Chome−24−12",
        "area": "Shibuya"
      },
      "start_time": "2025-01-15T08:00:00+09:00",
      "end_time": "2025-01-15T09:30:00+09:00",
      "comment": "Perfect quiet spot to start the day. The flat white was exceptional and the morning light through the windows was beautiful.",
      "order": 1
    },
    {
      "id": "activity-2",
      "location": {
        "type": "custom_text",
        "custom_name": "Shibuya Sky observation deck",
        "area": "Shibuya"
      },
      "start_time": "2025-01-15T10:00:00+09:00",
      "end_time": "2025-01-15T11:00:00+09:00",
      "comment": "Amazing 360-degree view of Tokyo. Worth the wait in line for the panoramic photos.",
      "order": 2
    }
  ],
  "is_deleted": false,
  "created_at": "2025-01-15T20:00:00+09:00",
  "updated_at": "2025-01-15T20:00:00+09:00",
  "version": 1
}
```

**Config Document (config.json)**:
```json
{
  "google_places_api_key": "",
  "default_timezone": "Asia/Tokyo",
  "data_directory": "~/.perfectday",
  "output_format": "table",
  "version": "0.1.0"
}
```

## Indexing Strategy

For efficient search (FR-013), maintain in-memory indices:

**Area Index**: Map of area → list of PerfectDay IDs
**Activity Index**: Map of location/comment keywords → list of PerfectDay IDs
**User Index**: Map of username → list of PerfectDay IDs

These indices are rebuilt on startup and updated on each data modification.

## Data Migration Strategy

Version field in PerfectDay enables safe schema evolution:
- Backward compatibility maintained for at least 2 major versions
- Migration scripts handle version upgrades
- Clear deprecation warnings before breaking changes