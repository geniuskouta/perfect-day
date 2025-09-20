# Implementation Plan: Perfect Day Sharing App


**Branch**: `001-perfect-day-this` | **Date**: 2025-09-15 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/Users/knakano/workspace/private/perfect-day/specs/001-perfect-day-this/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
4. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
5. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, or `GEMINI.md` for Gemini CLI).
6. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
7. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
8. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Perfect Day sharing application with dual interface: CLI for local-first personal use and REST API for web/mobile integration. Users (persona: Kouta) can create, document, and share detailed day plans including activities, locations, time spent, and personal commentary. Built as a Golang application with shared core logic, local JSON storage, and Google Places API integration.

**Interface Options:**
- **CLI**: Local terminal interface for personal workflow, offline capability
- **REST API**: HTTP API for web frontends, mobile apps, and team collaboration

## Technical Context
**Language/Version**: Go 1.21+
**Primary Dependencies**: Google Places API, JSON file handling, CLI framework (cobra), HTTP server (gin/echo)
**Storage**: Local JSON files for data persistence (shared between CLI and API)
**Testing**: Go standard testing package
**Target Platform**: Cross-platform terminal + HTTP API (Linux, macOS, Windows)
**Project Type**: dual-interface - terminal CLI + REST API with shared core
**Performance Goals**: <100ms response time for local operations, <2s for Google Places API calls, <200ms for API endpoints
**Constraints**: Local-first design, offline CLI capability, simple file-based storage, no Docker complexity
**Scale/Scope**: Single/multi-user depending on interface, hundreds of perfect days, dozens of activities per day

**Architecture**: Shared business logic with dual interfaces
- **CLI**: Direct file access, personal workflow, offline-first
- **API**: HTTP server, multi-user capable, real-time access
- **Core**: Shared packages for models, storage, places, search

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Simplicity**:
- Projects: 1 (cli only)
- Using framework directly? Yes (cobra CLI, native JSON)
- Single data model? Yes (unified structs for persistence and display)
- Avoiding patterns? Yes (direct file operations, no unnecessary abstractions)

**Architecture**:
- EVERY feature as library? Yes (models, storage, google-places, search as packages)
- Libraries listed: models (data structures), storage (JSON persistence), places (Google API), search (filtering)
- CLI per library: perfectday create, perfectday list, perfectday search, perfectday export --format=json
- Library docs: llms.txt format planned? Yes

**Testing (NON-NEGOTIABLE)**:
- RED-GREEN-Refactor cycle enforced? Yes (test files before implementation)
- Git commits show tests before implementation? Yes
- Order: Contract→Integration→E2E→Unit strictly followed? Yes
- Real dependencies used? Yes (real JSON files, real Google Places API calls in tests)
- Integration tests for: JSON storage, Google Places API, CLI commands
- FORBIDDEN: Implementation before test, skipping RED phase - Acknowledged

**Observability**:
- Structured logging included? Yes (JSON logging for debugging)
- Frontend logs → backend? N/A (single CLI app)
- Error context sufficient? Yes (detailed error messages with context)

**Versioning**:
- Version number assigned? 0.1.0 (initial implementation)
- BUILD increments on every change? Yes
- Breaking changes handled? Migration scripts for JSON schema changes

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Dual-Interface Architecture (CLI + REST API)

# Shared Core Logic
pkg/
├── auth/          # User authentication
├── models/        # Data structures (User, PerfectDay, Activity, Location)
├── storage/       # JSON file persistence
├── places/        # Google Places API integration
├── search/        # Search and filtering functionality
└── config/        # Configuration management

# Interface Implementations
cmd/
├── perfectday-cli/    # CLI application (existing)
└── perfectday-api/    # REST API server (new)

# API-specific code
internal/api/
├── handlers/      # HTTP request handlers
├── middleware/    # Auth, CORS, logging middleware
├── routes/        # Route definitions
└── server/        # HTTP server setup

# Legacy structure (to be refactored into pkg/)
src/
├── models/        # → move to pkg/models/
├── storage/       # → move to pkg/storage/
├── places/        # → move to pkg/places/
├── search/        # → move to pkg/search/
├── cli/           # → move to cmd/perfectday-cli/
└── lib/           # → move to pkg/config/

tests/
├── contract/      # CLI interface tests
├── api/           # REST API tests
├── integration/   # End-to-end workflow tests
└── unit/          # Package-specific tests
```

**Structure Decision**: Dual-interface with shared core - Refactor existing CLI into shared packages, add REST API interface

## REST API Design

### Core Endpoints
```
# Authentication
POST   /api/v1/auth/login          # Username-based login (matching CLI)
GET    /api/v1/auth/me             # Current user info

# Perfect Days
GET    /api/v1/perfect-days        # List/search perfect days
POST   /api/v1/perfect-days        # Create new perfect day
GET    /api/v1/perfect-days/{id}   # Get specific perfect day
PUT    /api/v1/perfect-days/{id}   # Update perfect day
DELETE /api/v1/perfect-days/{id}   # Soft delete perfect day

# Users & Public Profiles
GET    /api/v1/users/{username}                    # Public user profile
GET    /api/v1/users/{username}/perfect-days       # User's public perfect days

# Places Integration
GET    /api/v1/places/search?q=tokyo               # Google Places proxy
GET    /api/v1/areas                               # Popular areas from data

# System
GET    /api/v1/health             # API health check
GET    /api/v1/version            # API version info
```

### API Response Formats
```json
// GET /api/v1/perfect-days
{
  "perfect_days": [...],
  "total": 42,
  "offset": 0,
  "limit": 10
}

// Error responses
{
  "error": "Perfect day not found",
  "code": "NOT_FOUND",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Implementation Strategy

**Phase 1: Shared Library Refactor**
- Extract CLI business logic into `pkg/` packages
- Ensure data models work for both JSON files and HTTP responses
- Create unified configuration system (env vars + config files)

**Phase 2: Basic API Server**
- Minimal HTTP server with Gin framework
- Core CRUD endpoints using existing storage layer
- Simple authentication (username-based, matching CLI)
- Local development focus (no Docker, no complex deployment)

**Phase 3: API Enhancement**
- Search endpoint with same logic as CLI search
- Google Places proxy for web clients
- Error handling and proper HTTP status codes
- API documentation (basic, for local development)

### Local Development Approach
- **Single binary**: `go run cmd/perfectday-api/main.go`
- **Shared data**: API and CLI use same JSON files
- **Simple auth**: No JWT complexity, session-based for local dev
- **No containers**: Direct Go compilation and execution
- **Hot reload**: Use `air` or similar for rapid iteration

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `/scripts/bash/update-agent-context.sh claude` for your AI assistant
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- **Refactoring First**: Move existing code into shared packages
- **Dual Interface**: Create both CLI and REST API implementations
- Each API endpoint → contract test task [P]
- Each CLI command → contract test task [P] (existing)
- Each entity (User, PerfectDay, Activity, Location) → model creation task [P]
- Each user story from spec → integration test scenario (CLI + API)
- Storage operations → JSON persistence tests and implementation
- Google Places integration → API client tests and implementation
- Search functionality → search/filter tests and implementation
- HTTP handlers → API endpoint implementation to pass contract tests

**Ordering Strategy**:
- TDD order: Contract tests → Integration tests → Unit tests → Implementation
- Refactor order: Extract shared packages → CLI refactor → API implementation
- Dependency order: Models → Storage → External APIs → Search → CLI + API
- Mark [P] for parallel execution where dependencies allow
- Group related functionality (e.g., all model tests, all storage operations)

**Estimated Output**: 45-50 numbered, ordered tasks in tasks.md covering:
- **Refactoring**: 6 tasks to extract shared packages from existing CLI
- **Contract Tests**: 8 CLI + 8 API endpoint tests [P]
- **Models**: 4 model creation tasks (1 per entity) [P]
- **Storage**: 6 storage operation tasks (CRUD + indexing)
- **External APIs**: 3 Google Places integration tasks
- **Search**: 4 search/filter functionality tasks
- **CLI**: 8 CLI command refactor tasks
- **API**: 8 REST API handler implementation tasks
- **Integration**: 3 test scenarios (CLI + API workflows)
- **Deployment**: 3 build tasks (CLI binary, API server, local dev setup)

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*