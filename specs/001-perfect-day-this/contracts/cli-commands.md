# CLI Command Contracts

## Command Overview

All commands follow the pattern: `perfectday <command> [options] [arguments]`

## Global Options
- `--format` string: Output format (table, json, yaml) [default: table]
- `--config` string: Config file path [default: ~/.perfectday/config.json]
- `--help`, `-h`: Show help
- `--version`, `-v`: Show version

## Command: `perfectday init`

**Purpose**: Initialize Perfect Day configuration for a new user
**Requirements**: FR-010 (user identification)

**Usage**: `perfectday init [username]`

**Arguments**:
- `username` (optional): Username to register. If not provided, prompt interactively.

**Flags**:
- `--timezone` string: Set default timezone [default: system timezone]

**Success Output** (format=table):
```
✓ Perfect Day initialized for user: kouta
✓ Configuration saved to: ~/.perfectday/config.json
✓ Data directory created: ~/.perfectday/days/kouta/

Get started: perfectday create
```

**Success Output** (format=json):
```json
{
  "status": "success",
  "username": "kouta",
  "config_path": "~/.perfectday/config.json",
  "data_directory": "~/.perfectday/days/kouta"
}
```

**Error Cases**:
- Username already exists: Exit code 1
- Invalid username format: Exit code 1
- Config directory not writable: Exit code 1

## Command: `perfectday create`

**Purpose**: Create a new perfect day entry
**Requirements**: FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007

**Usage**: `perfectday create`

**Interactive Flow**:
1. Prompt for day title (5-100 chars)
2. Prompt for day description (10-1000 chars)
3. Prompt for date (defaults to today)
4. Prompt for main areas (comma-separated)
5. Activity loop:
   - Prompt for location (with Google Places search)
   - Prompt for start time
   - Prompt for end time
   - Prompt for personal comment
   - Ask "Add another activity? (y/n)"
6. Save and display summary

**Success Output**:
```
✓ Perfect day created: "Morning in Shibuya"
✓ Saved: ~/.perfectday/days/kouta/2025-01-15_morning-in-shibuya.json

View: perfectday view morning-in-shibuya
Edit: perfectday edit morning-in-shibuya
Share: perfectday list --format=json
```

**Error Cases**:
- Validation errors: Exit code 1, show specific field errors
- Google Places API error: Continue with custom text entry
- File write error: Exit code 1

## Command: `perfectday list`

**Purpose**: List perfect days with filtering
**Requirements**: FR-008, FR-011 (public viewing)

**Usage**: `perfectday list [options]`

**Flags**:
- `--user` string: Filter by username (defaults to all users)
- `--area` string: Filter by area
- `--limit` int: Limit results [default: 20]
- `--sort` string: Sort order (date-desc, date-asc, title-asc) [default: date-desc]

**Success Output** (format=table):
```
TITLE                   DATE        AREAS           ACTIVITIES  AUTHOR
Morning in Shibuya      2025-01-15  Shibuya         3          kouta
Tokyo Foodie Adventure  2025-01-20  Ginza, Tsukiji  5          kouta
Art Gallery Hopping     2025-02-01  Roppongi        4          alice

3 perfect days found. Use 'perfectday view <title>' for details.
```

**Success Output** (format=json):
```json
{
  "perfect_days": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Morning in Shibuya",
      "date": "2025-01-15",
      "areas": ["Shibuya"],
      "activity_count": 3,
      "username": "kouta",
      "created_at": "2025-01-15T20:00:00+09:00"
    }
  ],
  "total": 1,
  "limit": 20
}
```

## Command: `perfectday search`

**Purpose**: Search perfect days by area and activity
**Requirements**: FR-013

**Usage**: `perfectday search <query> [options]`

**Arguments**:
- `query` (required): Search term for areas, locations, or activity comments

**Flags**:
- `--area` string: Restrict search to specific area
- `--activity-only`: Search only in activity locations and comments
- `--limit` int: Limit results [default: 10]

**Success Output**: Same format as `list` command

## Command: `perfectday view`

**Purpose**: Display detailed view of a perfect day
**Requirements**: FR-009 (timeline format)

**Usage**: `perfectday view <day-identifier>`

**Arguments**:
- `day-identifier` (required): Day title (fuzzy matched) or UUID

**Success Output** (format=table):
```
Morning in Shibuya (2025-01-15) by kouta
════════════════════════════════════════

A perfect morning exploring the bustling heart of Tokyo, from quiet coffee
to vibrant streets.

Areas: Shibuya, Harajuku

TIMELINE
────────────────────────────────────────
08:00-09:30  Blue Bottle Coffee Shibuya (Shibuya)
             Perfect quiet spot to start the day. The flat white was
             exceptional and the morning light through the windows was beautiful.

10:00-11:00  Shibuya Sky observation deck (Shibuya)
             Amazing 360-degree view of Tokyo. Worth the wait in line for
             the panoramic photos.

Total time: 2.5 hours across 2 activities
```

## Command: `perfectday edit`

**Purpose**: Edit an existing perfect day
**Requirements**: FR-012 (edit capability)

**Usage**: `perfectday edit <day-identifier>`

**Interactive Flow**: Same as `create` but pre-filled with existing data

## Command: `perfectday delete`

**Purpose**: Soft delete a perfect day
**Requirements**: FR-012 (soft delete)

**Usage**: `perfectday delete <day-identifier>`

**Flags**:
- `--confirm`: Skip confirmation prompt

**Success Output**:
```
✓ Perfect day deleted: "Morning in Shibuya"
Note: This is a soft delete. Data is preserved but hidden from public view.
```

## Error Handling Standards

**Exit Codes**:
- 0: Success
- 1: User error (validation, not found, etc.)
- 2: System error (file I/O, permissions, etc.)
- 3: External service error (Google Places API)

**Error Message Format**:
```
Error: <brief description>

<detailed explanation with suggested action>

For help: perfectday <command> --help
```

**Common Error Patterns**:
- Field validation: Show field name and validation rule
- File operations: Show path and permissions suggestion
- API errors: Show service name and fallback options
- Network errors: Suggest offline mode where applicable