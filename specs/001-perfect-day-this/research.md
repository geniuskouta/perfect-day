# Research Phase: Perfect Day Terminal App

## Go CLI Framework Selection

**Decision**: Use cobra CLI framework
**Rationale**: Industry standard for Go CLI apps, excellent command structure, built-in help system, used by kubectl, helm, etc.
**Alternatives considered**:
- flag package (too basic for complex commands)
- urfave/cli (less features than cobra)
- kingpin (less maintained)

## Google Places API Integration

**Decision**: Use Google Places API Text Search
**Rationale**: Meets FR-004 requirement, returns business profiles, has Go client library
**Alternatives considered**:
- OpenStreetMap Nominatim (free but less business data)
- Foursquare Places API (good for venues but requires different setup)
- Manual location entry only (doesn't meet Google Business Profile requirement)

**Implementation approach**:
- Use googlemaps/maps package for Go
- Implement text search for business lookup
- Cache results locally to minimize API calls
- Fallback to custom text entry per FR-005

## Data Storage Strategy

**Decision**: Local JSON files with structured directories
**Rationale**: Simple, version-controllable, no server needed, meets single-user terminal requirement
**Alternatives considered**:
- SQLite (overkill for simple data)
- CSV files (harder to handle nested data)
- YAML/TOML (JSON has better Go support)

**Structure**:
```
~/.perfectday/
├── users.json           # Username registry
├── days/
│   ├── username/
│   │   ├── 2025-01-15.json
│   │   └── 2025-01-20.json
└── config.json         # App configuration
```

## Terminal UI/UX Approach

**Decision**: Interactive prompts with structured output
**Rationale**: Follows UNIX philosophy, scriptable, accessible
**Alternatives considered**:
- TUI with tcell/bubbletea (too complex for MVP)
- Simple flags only (poor UX for complex input)

**User flow**:
- `perfectday create` - interactive prompt for day creation
- `perfectday list` - table format output
- `perfectday search --area "Shibuya"` - filtered results
- `perfectday view <day-id>` - detailed timeline view

## Time Handling Strategy

**Decision**: Use time.Time with location awareness
**Rationale**: Native Go time support, handles timezones (important for Tokyo persona)
**Alternatives considered**:
- String timestamps (error-prone)
- Unix timestamps (poor readability)

**Format**: RFC3339 with timezone for storage, human-readable for display

## Search Implementation

**Decision**: In-memory filtering with indexing
**Rationale**: Fast for hundreds of entries, simple implementation
**Alternatives considered**:
- Full-text search with Bleve (overkill for MVP)
- SQL FTS (adds database dependency)

**Search criteria**: Area name, activity location, activity comments

## Error Handling Strategy

**Decision**: Structured errors with user-friendly messages
**Rationale**: Terminal users need clear guidance
**Implementation**:
- Wrap errors with context
- Distinguish user errors from system errors
- Provide actionable error messages

## Configuration Management

**Decision**: JSON config file with sensible defaults
**Rationale**: Simple, familiar format
**Location**: ~/.perfectday/config.json
**Settings**: API keys, default timezone, output preferences

## Testing Strategy

**Decision**: Table-driven tests with real file operations
**Rationale**: Go idiom, tests real behavior, not mocks
**Approach**:
- Use temp directories for file operations
- Test with real (limited) Google Places API calls
- Golden file testing for CLI output
- Integration tests for full command workflows

## Deployment Strategy

**Decision**: Single binary with go build, distributed via GitHub releases
**Rationale**: Go strength, easy user installation
**Process**:
- Cross-compile for major platforms
- Include installation instructions
- Version embedded in binary

## All Technical Context Requirements Resolved

✅ Language/Version: Go 1.21+
✅ Primary Dependencies: cobra, googlemaps/maps, standard library
✅ Storage: Local JSON files
✅ Testing: Go testing package with table-driven tests
✅ Target Platform: Cross-platform terminal
✅ Performance Goals: Achievable with local storage
✅ Constraints: Addressed with offline viewing capability
✅ Scale/Scope: Well-suited for local JSON storage

No NEEDS CLARIFICATION items remain.