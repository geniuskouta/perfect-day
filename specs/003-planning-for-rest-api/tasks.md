# Tasks: Perfect Day Dual Interface (CLI + REST API)

**Input**: Design documents from `/specs/003-planning-for-rest-api/`
**Prerequisites**: plan.md (✓), data-model.md (✓), contracts/ (✓), quickstart.md (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21+, cobra CLI, gin/echo API, Google Places API, JSON storage
   → Structure: Dual-interface (CLI + REST API) with shared core packages
2. Load design documents:
   → data-model.md: 4 entities (User, PerfectDay, Activity, Location)
   → contracts/cli-commands.md: 7 CLI commands (init, create, list, search, view, edit, delete)
   → quickstart.md: 3 main user scenarios for integration tests
3. Generate tasks by category per TDD requirements:
   → Refactoring: Extract shared packages from existing code
   → Contract tests: CLI commands + REST API endpoints
   → Core implementation: Models, storage, places, search
   → Interface implementation: CLI refactor + API handlers
   → Integration tests: End-to-end workflows
4. Apply task rules: Different files = [P], tests before implementation
5. Number tasks sequentially (T001-T053)
6. Generate dependency graph and parallel execution examples
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Dual-interface architecture at repository root:
- **Shared Core**: `pkg/` with auth, models, storage, places, search, config
- **CLI Interface**: `cmd/perfectday-cli/` (refactored from src/cli/)
- **API Interface**: `cmd/perfectday-api/` and `internal/api/`
- **Tests**: `tests/` with contract, api, integration, unit subdirectories

## Phase 3.1: Project Restructure & Setup

- [ ] T001 Create new package structure (pkg/, cmd/, internal/api/) per plan.md
- [ ] T002 Initialize Go modules and add gin/echo dependency to go.mod
- [ ] T003 [P] Configure golangci-lint with .golangci.yml for both CLI and API
- [ ] T004 [P] Extract models package from src/models/ to pkg/models/
- [ ] T005 [P] Extract storage package from src/storage/ to pkg/storage/
- [ ] T006 [P] Extract places package from src/places/ to pkg/places/
- [ ] T007 [P] Extract search package from src/search/ to pkg/search/
- [ ] T008 [P] Extract config package from src/lib/ to pkg/config/

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### CLI Contract Tests (Existing Interface)
- [ ] T009 [P] Contract test `perfectday init` command in tests/contract/test_init_command.go
- [ ] T010 [P] Contract test `perfectday create` command in tests/contract/test_create_command.go
- [ ] T011 [P] Contract test `perfectday list` command in tests/contract/test_list_command.go
- [ ] T012 [P] Contract test `perfectday search` command in tests/contract/test_search_command.go
- [ ] T013 [P] Contract test `perfectday view` command in tests/contract/test_view_command.go
- [ ] T014 [P] Contract test `perfectday edit` command in tests/contract/test_edit_command.go
- [ ] T015 [P] Contract test `perfectday delete` command in tests/contract/test_delete_command.go

### REST API Contract Tests (New Interface)
- [ ] T016 [P] API contract test POST /api/v1/auth/login in tests/api/test_auth_login.go
- [ ] T017 [P] API contract test GET /api/v1/auth/me in tests/api/test_auth_me.go
- [ ] T018 [P] API contract test GET /api/v1/perfect-days in tests/api/test_perfect_days_list.go
- [ ] T019 [P] API contract test POST /api/v1/perfect-days in tests/api/test_perfect_days_create.go
- [ ] T020 [P] API contract test GET /api/v1/perfect-days/{id} in tests/api/test_perfect_days_get.go
- [ ] T021 [P] API contract test PUT /api/v1/perfect-days/{id} in tests/api/test_perfect_days_update.go
- [ ] T022 [P] API contract test DELETE /api/v1/perfect-days/{id} in tests/api/test_perfect_days_delete.go
- [ ] T023 [P] API contract test GET /api/v1/users/{username} in tests/api/test_users_profile.go
- [ ] T024 [P] API contract test GET /api/v1/places/search in tests/api/test_places_search.go
- [ ] T025 [P] API contract test GET /api/v1/health in tests/api/test_health.go

### Integration Tests (User Workflows)
- [ ] T026 [P] Integration test create perfect day workflow (CLI) in tests/integration/test_create_workflow_cli.go
- [ ] T027 [P] Integration test create perfect day workflow (API) in tests/integration/test_create_workflow_api.go
- [ ] T028 [P] Integration test browse and search perfect days (CLI) in tests/integration/test_browse_workflow_cli.go
- [ ] T029 [P] Integration test browse and search perfect days (API) in tests/integration/test_browse_workflow_api.go
- [ ] T030 [P] Integration test edit and delete perfect days (CLI) in tests/integration/test_manage_workflow_cli.go
- [ ] T031 [P] Integration test edit and delete perfect days (API) in tests/integration/test_manage_workflow_api.go

## Phase 3.3: Shared Core Implementation (ONLY after tests are failing)

### Data Models (Shared)
- [ ] T032 [P] User model with validation in pkg/models/user.go
- [ ] T033 [P] PerfectDay model with validation in pkg/models/perfect_day.go
- [ ] T034 [P] Activity model with validation in pkg/models/activity.go
- [ ] T035 [P] Location model with validation in pkg/models/location.go

### Storage Layer (Shared)
- [ ] T036 [P] JSON file storage interface in pkg/storage/storage.go
- [ ] T037 User storage operations (CRUD) in pkg/storage/user_storage.go
- [ ] T038 PerfectDay storage operations (CRUD with soft delete) in pkg/storage/perfect_day_storage.go
- [ ] T039 [P] Storage indexing for search functionality in pkg/storage/indexer.go

### External Services (Shared)
- [ ] T040 [P] Google Places API client in pkg/places/places_client.go
- [ ] T041 [P] Places search with fallback to custom text in pkg/places/search.go

### Search and Filtering (Shared)
- [ ] T042 [P] Search service for area and activity filtering in pkg/search/search_service.go
- [ ] T043 [P] Search result ranking and pagination in pkg/search/ranking.go

### Authentication (Shared)
- [ ] T044 [P] User authentication service in pkg/auth/auth_service.go
- [ ] T045 [P] Session management for API in pkg/auth/session.go

## Phase 3.4: Interface Implementation

### CLI Interface (Refactored)
- [ ] T046 CLI main entry point with shared packages in cmd/perfectday-cli/main.go
- [ ] T047 Refactor CLI commands to use shared packages in cmd/perfectday-cli/commands/

### REST API Interface (New)
- [ ] T048 API server setup with gin framework in cmd/perfectday-api/main.go
- [ ] T049 API routes configuration in internal/api/routes/routes.go
- [ ] T050 Authentication middleware in internal/api/middleware/auth.go
- [ ] T051 CORS and logging middleware in internal/api/middleware/cors.go
- [ ] T052 [P] Auth handlers (login, me) in internal/api/handlers/auth.go
- [ ] T053 [P] Perfect days handlers (CRUD) in internal/api/handlers/perfect_days.go
- [ ] T054 [P] Users handlers (profiles) in internal/api/handlers/users.go
- [ ] T055 [P] Places proxy handlers in internal/api/handlers/places.go
- [ ] T056 [P] Health check handler in internal/api/handlers/health.go

## Phase 3.5: Integration & Polish

- [ ] T057 Error handling and JSON response formatting in internal/api/utils/response.go
- [ ] T058 [P] API documentation generation (basic for local dev)
- [ ] T059 [P] Unit tests for shared models in tests/unit/test_models.go
- [ ] T060 [P] Unit tests for storage operations in tests/unit/test_storage.go
- [ ] T061 [P] Unit tests for search functionality in tests/unit/test_search.go
- [ ] T062 [P] Performance tests for large datasets in tests/unit/test_performance.go
- [ ] T063 Cross-platform build scripts for both CLI and API in scripts/build.sh
- [ ] T064 Local development setup documentation in docs/local-dev.md
- [ ] T065 API usage examples and Postman collection in docs/api-examples.md

## Dependencies

**Setup before everything**: T001-T008 → all other tasks

**Tests before implementation**:
- T009-T031 (all tests) → T032-T056 (all implementation)
- Contract tests must fail first, then implementation makes them pass

**Package extraction dependencies**:
- T004-T008 (extract packages) → T009-T031 (tests using packages) → T032-T056 (implementation)

**Model dependencies**:
- T032-T035 (models) → T036-T039 (storage) → T042-T043 (search)
- T040-T041 (places) can run parallel with storage

**Interface dependencies**:
- T044-T045 (auth) → T050 (auth middleware)
- T036-T043 (shared services) → T046-T056 (interface implementations)
- T046-T047 (CLI refactor) can run parallel with T048-T056 (API implementation)

**Polish after everything**: T057-T058 → T059-T065

## Parallel Example

### Phase 3.1 - Launch package extraction together:
```bash
# All package extractions can run in parallel (different target directories)
Task: "Extract models package from src/models/ to pkg/models/"
Task: "Extract storage package from src/storage/ to pkg/storage/"
Task: "Extract places package from src/places/ to pkg/places/"
Task: "Extract search package from src/search/ to pkg/search/"
```

### Phase 3.2 - Launch contract tests together:
```bash
# All contract tests can run in parallel (different files)
Task: "Contract test perfectday init command in tests/contract/test_init_command.go"
Task: "API contract test POST /api/v1/auth/login in tests/api/test_auth_login.go"
Task: "API contract test GET /api/v1/perfect-days in tests/api/test_perfect_days_list.go"
Task: "Integration test create perfect day workflow (CLI) in tests/integration/test_create_workflow_cli.go"
```

### Phase 3.3 - Launch model creation in parallel:
```bash
# All models can be created in parallel (different files)
Task: "User model with validation in pkg/models/user.go"
Task: "PerfectDay model with validation in pkg/models/perfect_day.go"
Task: "Activity model with validation in pkg/models/activity.go"
Task: "Location model with validation in pkg/models/location.go"
```

### Phase 3.4 - Launch API handlers in parallel:
```bash
# All API handlers can be created in parallel (different files)
Task: "Auth handlers (login, me) in internal/api/handlers/auth.go"
Task: "Perfect days handlers (CRUD) in internal/api/handlers/perfect_days.go"
Task: "Users handlers (profiles) in internal/api/handlers/users.go"
Task: "Places proxy handlers in internal/api/handlers/places.go"
```

## Notes
- [P] tasks = different files, no dependencies, can run simultaneously
- Verify all contract tests fail before starting implementation (TDD requirement)
- Each task modifies only one file to avoid conflicts
- Commit after completing each task
- Use table-driven tests following Go conventions
- Include error cases and edge conditions in all tests
- Both CLI and API share the same data storage (JSON files)
- API uses session-based auth for local development simplicity

## Task Generation Rules Applied

1. **From CLI Contracts**: 7 CLI commands → 7 contract test tasks [P] + CLI refactor tasks
2. **From API Design**: 11 API endpoints → 11 contract test tasks [P] + handler implementation tasks
3. **From Data Model**: 4 entities → 4 model creation tasks [P] + storage tasks
4. **From User Stories**: 3 main workflows → 6 integration tests [P] (CLI + API versions)
5. **From Architecture**: Package extraction → shared core → dual interfaces
6. **Ordering**: Setup → Package extraction → Tests → Models → Storage → Services → Interfaces → Integration → Polish

## Validation Checklist

- [x] All CLI contracts have corresponding tests (T009-T015)
- [x] All API endpoints have corresponding tests (T016-T025)
- [x] All entities have model tasks (T032-T035)
- [x] All tests come before implementation (T009-T031 before T032-T056)
- [x] Parallel tasks truly independent (different files, no shared dependencies)
- [x] Each task specifies exact file path
- [x] No [P] task modifies same file as another [P] task
- [x] TDD cycle enforced (tests must fail first)
- [x] Dual-interface architecture properly separated (CLI + API)
- [x] Shared core packages properly extracted and reused