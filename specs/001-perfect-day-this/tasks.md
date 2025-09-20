# Tasks: Perfect Day Sharing App

**Input**: Design documents from `/specs/001-perfect-day-this/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21+, cobra CLI, Google Places API, JSON storage
   → Structure: Single project (src/, tests/)
2. Load design documents:
   → data-model.md: 4 entities (User, PerfectDay, Activity, Location)
   → contracts/: CLI commands (init, create, list, search, view, edit, delete)
   → quickstart.md: 3 main user scenarios for integration tests
3. Generate tasks by category per TDD requirements
4. Apply task rules: Different files = [P], tests before implementation
5. Number tasks sequentially (T001-T035)
6. Generate dependency graph and parallel execution examples
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Single project structure at repository root:
- **Source**: `src/` with packages: models, storage, places, search, cli
- **Tests**: `tests/` with contract, integration, unit subdirectories

## Phase 3.1: Setup

- [ ] T001 Create Go project structure with src/{models,storage,places,search,cli} and tests/{contract,integration,unit} directories
- [ ] T002 Initialize Go module with cobra CLI and googlemaps dependencies in go.mod
- [ ] T003 [P] Configure golangci-lint with .golangci.yml for code quality

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (CLI Commands)
- [ ] T004 [P] Contract test `perfectday init` command in tests/contract/test_init_command.go
- [ ] T005 [P] Contract test `perfectday create` command in tests/contract/test_create_command.go
- [ ] T006 [P] Contract test `perfectday list` command in tests/contract/test_list_command.go
- [ ] T007 [P] Contract test `perfectday search` command in tests/contract/test_search_command.go
- [ ] T008 [P] Contract test `perfectday view` command in tests/contract/test_view_command.go
- [ ] T009 [P] Contract test `perfectday edit` command in tests/contract/test_edit_command.go
- [ ] T010 [P] Contract test `perfectday delete` command in tests/contract/test_delete_command.go

### Integration Tests (User Stories)
- [ ] T011 [P] Integration test create perfect day workflow in tests/integration/test_create_workflow.go
- [ ] T012 [P] Integration test browse and search perfect days in tests/integration/test_browse_workflow.go
- [ ] T013 [P] Integration test edit and delete perfect days in tests/integration/test_manage_workflow.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Data Models
- [ ] T014 [P] User model with validation in src/models/user.go
- [ ] T015 [P] PerfectDay model with validation in src/models/perfect_day.go
- [ ] T016 [P] Activity model with validation in src/models/activity.go
- [ ] T017 [P] Location model with validation in src/models/location.go

### Storage Layer
- [ ] T018 [P] JSON file storage interface in src/storage/storage.go
- [ ] T019 User storage operations (CRUD) in src/storage/user_storage.go
- [ ] T020 PerfectDay storage operations (CRUD with soft delete) in src/storage/perfect_day_storage.go
- [ ] T021 [P] Storage indexing for search functionality in src/storage/indexer.go

### External Services
- [ ] T022 [P] Google Places API client in src/places/places_client.go
- [ ] T023 [P] Places search with fallback to custom text in src/places/search.go

### Search and Filtering
- [ ] T024 [P] Search service for area and activity filtering in src/search/search_service.go
- [ ] T025 [P] Search result ranking and pagination in src/search/ranking.go

### CLI Commands Implementation
- [ ] T026 Root CLI command with global flags in src/cli/root.go
- [ ] T027 `perfectday init` command implementation in src/cli/init.go
- [ ] T028 `perfectday create` command with interactive prompts in src/cli/create.go
- [ ] T029 `perfectday list` command with filtering in src/cli/list.go
- [ ] T030 `perfectday search` command implementation in src/cli/search.go
- [ ] T031 `perfectday view` command with timeline display in src/cli/view.go
- [ ] T032 `perfectday edit` command with pre-filled prompts in src/cli/edit.go
- [ ] T033 `perfectday delete` command with confirmation in src/cli/delete.go

## Phase 3.4: Integration

- [ ] T034 Main CLI entry point connecting all commands in main.go
- [ ] T035 Error handling and structured logging across all packages
- [ ] T036 Configuration loading and validation in src/cli/config.go

## Phase 3.5: Polish

- [ ] T037 [P] Unit tests for User model validation in tests/unit/test_user_model.go
- [ ] T038 [P] Unit tests for storage operations in tests/unit/test_storage.go
- [ ] T039 [P] Unit tests for search functionality in tests/unit/test_search.go
- [ ] T040 [P] Performance tests for large datasets in tests/unit/test_performance.go
- [ ] T041 Cross-platform build script with goreleaser in .goreleaser.yml
- [ ] T042 [P] Update project README with installation and usage
- [ ] T043 Manual testing using quickstart.md scenarios
- [ ] T044 Code cleanup and documentation improvements

## Dependencies

**Setup before everything**: T001-T003 → all other tasks

**Tests before implementation**:
- T004-T013 (all tests) → T014-T036 (all implementation)
- Contract tests must fail first, then implementation makes them pass

**Model dependencies**:
- T014-T017 (models) → T018-T021 (storage) → T024-T025 (search)
- T022-T023 (places) can run parallel with storage

**CLI dependencies**:
- T018-T025 (services) → T026-T033 (CLI commands)
- T026 (root) → T027-T033 (individual commands)
- T027-T033 → T034 (main entry)

**Polish after everything**: T034-T036 → T037-T044

## Parallel Example

### Phase 3.2 - Launch all contract tests together:
```bash
# All contract tests can run in parallel (different files)
Task: "Contract test perfectday init command in tests/contract/test_init_command.go"
Task: "Contract test perfectday create command in tests/contract/test_create_command.go"
Task: "Contract test perfectday list command in tests/contract/test_list_command.go"
Task: "Contract test perfectday search command in tests/contract/test_search_command.go"
```

### Phase 3.3 - Launch model creation in parallel:
```bash
# All models can be created in parallel (different files)
Task: "User model with validation in src/models/user.go"
Task: "PerfectDay model with validation in src/models/perfect_day.go"
Task: "Activity model with validation in src/models/activity.go"
Task: "Location model with validation in src/models/location.go"
```

## Notes
- [P] tasks = different files, no dependencies, can run simultaneously
- Verify all contract tests fail before starting implementation (TDD requirement)
- Each task modifies only one file to avoid conflicts
- Commit after completing each task
- Use table-driven tests following Go conventions
- Include error cases and edge conditions in all tests

## Task Generation Rules Applied

1. **From Contracts**: 7 CLI commands → 7 contract test tasks [P] + 7 implementation tasks
2. **From Data Model**: 4 entities → 4 model creation tasks [P] + storage tasks
3. **From User Stories**: 3 main workflows → 3 integration tests [P]
4. **From Research**: Go tooling → setup tasks, Google Places → API integration
5. **Ordering**: Setup → Tests → Models → Storage → Services → CLI → Integration → Polish

## Validation Checklist

- [x] All CLI contracts have corresponding tests (T004-T010)
- [x] All entities have model tasks (T014-T017)
- [x] All tests come before implementation (T004-T013 before T014-T036)
- [x] Parallel tasks truly independent (different files, no shared dependencies)
- [x] Each task specifies exact file path
- [x] No [P] task modifies same file as another [P] task
- [x] TDD cycle enforced (tests must fail first)
- [x] Constitutional compliance (library architecture, real dependencies in tests)