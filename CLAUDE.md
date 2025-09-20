# Perfect Day Terminal App - Claude Code Context

## Project Overview
Terminal-based application for sharing perfect day experiences. Users can create, document, and share detailed day plans including activities, locations, time spent, and personal commentary. Built as a Golang CLI application with local JSON storage and Google Places API integration.

## Technical Stack
- **Language**: Go 1.21+
- **CLI Framework**: cobra
- **Storage**: Local JSON files
- **External APIs**: Google Places API
- **Testing**: Go standard testing package
- **Platform**: Cross-platform terminal (Linux, macOS, Windows)

## Project Structure
```
src/
├── models/          # Data structures (User, PerfectDay, Activity, Location)
├── storage/         # JSON file persistence layer
├── places/          # Google Places API integration
├── search/          # Search and filtering functionality
├── cli/             # Cobra CLI commands
└── lib/             # Shared utilities

tests/
├── contract/        # CLI contract tests
├── integration/     # Full workflow tests
└── unit/           # Package-specific tests
```

## Core Entities
- **User**: Username-based identification with timezone support
- **PerfectDay**: Complete day plan with title, description, areas, activities
- **Activity**: Individual activity with location, timing, and commentary
- **Location**: Google Places or custom text location with area

## Key Requirements
- Username-based authentication (simple, no passwords)
- Google Places API integration with custom text fallback
- Public sharing of all perfect day entries
- Soft delete functionality (hidden but preserved)
- Search by area and activity content
- Timeline format for displaying perfect days
- Local JSON file storage with structured directories

## Constitutional Principles
- Library-first architecture (models, storage, places, search as packages)
- CLI commands for each library function
- TDD approach with tests before implementation
- Real dependencies in tests (actual files, API calls)
- Structured JSON logging for debugging
- Version 0.1.0 with BUILD increments

## Current Implementation Status
- **Phase 0**: Research completed (Go tooling, Google Places API, storage strategy)
- **Phase 1**: Design completed (data model, CLI contracts, quickstart guide)
- **Phase 2**: Ready for task generation (/tasks command)

## Recent Changes
- Created feature specification with 13 functional requirements
- Completed research phase with technology decisions
- Designed data model with validation rules and JSON schema
- Defined CLI command contracts with interactive flows
- Created comprehensive quickstart guide

## Development Notes
- Use table-driven tests following Go idioms
- Implement graceful fallback when Google Places API unavailable
- Ensure timezone awareness for Tokyo persona (Kouta)
- Build single binary with cross-platform support
- Follow UNIX philosophy for CLI design (stdin/stdout, composable)

This file is auto-updated during the planning phase and should not be manually edited outside the designated manual sections.