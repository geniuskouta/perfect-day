# Tasks: Perfect Day PWA Timeline Builder

**Input**: Design documents from `/Users/knakano/workspace/private/perfect-day/specs/004-pwa-app/`
**Prerequisites**: plan.md (✓), contracts/ (empty), data-model.md (N/A - frontend only)

## Execution Flow
Based on plan.md analysis:
- **Tech Stack**: React 18 + Vite + Tailwind CSS + shadcn/ui + Material Icons
- **Project Type**: Frontend PWA with API integration
- **Structure**: Single `ui/` directory at repository root
- **API Integration**: Existing REST API endpoints
- **Focus**: Mobile-first, touch-optimized timeline interface

## Phase 3.1: Setup & Project Structure
- [ ] T001 Create React + Vite project structure in ui/ directory
- [ ] T002 Install and configure Tailwind CSS + PostCSS
- [ ] T003 [P] Install and configure shadcn/ui component library
- [ ] T004 [P] Install Material Icons and configure icon system
- [ ] T005 [P] Configure ESLint and Prettier for code formatting
- [ ] T006 Set up Vite PWA plugin and basic configuration

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T007 [P] Timeline component test in ui/src/components/timeline/__tests__/Timeline.test.tsx
- [ ] T008 [P] Activity card component test in ui/src/components/activity/__tests__/ActivityCard.test.tsx
- [ ] T009 [P] Activity modal test in ui/src/components/activity/__tests__/ActivityModal.test.tsx
- [ ] T010 [P] Perfect day page test in ui/src/pages/__tests__/PerfectDayPage.test.tsx
- [ ] T011 [P] API service test in ui/src/services/__tests__/apiService.test.ts

## Phase 3.3: Core Components (ONLY after tests are failing)
- [ ] T012 [P] Timeline base component in ui/src/components/timeline/Timeline.tsx
- [ ] T013 [P] TimeSlot component in ui/src/components/timeline/TimeSlot.tsx
- [ ] T014 [P] ActivityCard component in ui/src/components/activity/ActivityCard.tsx
- [ ] T015 ActivityModal component in ui/src/components/activity/ActivityModal.tsx
- [ ] T016 [P] PerfectDayCard component in ui/src/components/perfect-day/PerfectDayCard.tsx
- [ ] T017 [P] Icon components using Material Icons in ui/src/components/ui/Icons.tsx

## Phase 3.4: Pages & Routing
- [ ] T018 React Router setup in ui/src/App.tsx
- [ ] T019 Timeline builder page in ui/src/pages/TimelineBuilderPage.tsx
- [ ] T020 [P] Perfect day discovery page in ui/src/pages/DiscoveryPage.tsx
- [ ] T021 [P] Perfect day view page in ui/src/pages/PerfectDayViewPage.tsx
- [ ] T022 Navigation component in ui/src/components/navigation/Navigation.tsx

## Phase 3.5: API Integration & State Management
- [ ] T023 [P] API service for perfect days in ui/src/services/perfectDayService.ts
- [ ] T024 [P] API service for authentication in ui/src/services/authService.ts
- [ ] T025 [P] API service for places search in ui/src/services/placesService.ts
- [ ] T026 React Context for perfect day state in ui/src/context/PerfectDayContext.tsx
- [ ] T027 React Context for auth state in ui/src/context/AuthContext.tsx
- [ ] T028 Custom hooks for API calls in ui/src/hooks/useApi.ts

## Phase 3.6: PWA Features
- [ ] T029 Service worker implementation in ui/public/sw.js
- [ ] T030 PWA manifest configuration in ui/public/manifest.json
- [ ] T031 [P] IndexedDB offline storage in ui/src/services/offlineService.ts
- [ ] T032 [P] Install prompt component in ui/src/components/pwa/InstallPrompt.tsx
- [ ] T033 Background sync for offline actions in ui/src/services/syncService.ts

## Phase 3.7: Mobile Optimization & Polish
- [ ] T034 [P] Touch gesture handling in ui/src/hooks/useTouch.ts
- [ ] T035 [P] Mobile-specific CSS optimizations in ui/src/styles/mobile.css
- [ ] T036 [P] Loading states and skeleton components in ui/src/components/ui/Loading.tsx
- [ ] T037 [P] Error boundary component in ui/src/components/ui/ErrorBoundary.tsx
- [ ] T038 Sharing functionality in ui/src/components/sharing/ShareDialog.tsx

## Phase 3.8: Integration & E2E Tests
- [ ] T039 [P] E2E test for timeline creation in ui/e2e/timeline-creation.spec.ts
- [ ] T040 [P] E2E test for perfect day sharing in ui/e2e/sharing.spec.ts
- [ ] T041 [P] PWA installation test in ui/e2e/pwa-install.spec.ts
- [ ] T042 [P] Offline functionality test in ui/e2e/offline.spec.ts
- [ ] T043 Performance testing and optimization

## Dependencies
- Setup (T001-T006) before everything else
- Tests (T007-T011) before implementation (T012-T038)
- T012 (Timeline) blocks T019 (Timeline page)
- T015 (ActivityModal) blocks T019 (Timeline page)
- T023-T025 (API services) block T026-T028 (state management)
- T029-T033 (PWA features) after core implementation
- E2E tests (T039-T043) after all implementation

## Parallel Execution Examples
```
# Phase 3.1 Setup (after T001-T002):
Task: "Install and configure shadcn/ui component library"
Task: "Install Material Icons and configure icon system"
Task: "Configure ESLint and Prettier for code formatting"

# Phase 3.2 Tests:
Task: "Timeline component test in ui/src/components/timeline/__tests__/Timeline.test.tsx"
Task: "Activity card component test in ui/src/components/activity/__tests__/ActivityCard.test.tsx"
Task: "Activity modal test in ui/src/components/activity/__tests__/ActivityModal.test.tsx"
Task: "Perfect day page test in ui/src/pages/__tests__/PerfectDayPage.test.tsx"
Task: "API service test in ui/src/services/__tests__/apiService.test.ts"

# Phase 3.3 Core Components:
Task: "Timeline base component in ui/src/components/timeline/Timeline.tsx"
Task: "TimeSlot component in ui/src/components/timeline/TimeSlot.tsx"
Task: "ActivityCard component in ui/src/components/activity/ActivityCard.tsx"
Task: "PerfectDayCard component in ui/src/components/perfect-day/PerfectDayCard.tsx"
Task: "Icon components using Material Icons in ui/src/components/ui/Icons.tsx"
```

## Project Structure
```
ui/
├── public/
│   ├── manifest.json
│   └── sw.js
├── src/
│   ├── components/
│   │   ├── timeline/
│   │   ├── activity/
│   │   ├── perfect-day/
│   │   ├── navigation/
│   │   ├── sharing/
│   │   ├── pwa/
│   │   └── ui/
│   ├── pages/
│   ├── services/
│   ├── context/
│   ├── hooks/
│   ├── styles/
│   └── App.tsx
├── e2e/
└── __tests__/
```

## Success Criteria
- ✅ Mobile-optimized timeline interface with tap interactions
- ✅ PWA features working (install, offline, service worker)
- ✅ Integration with existing Perfect Day REST API
- ✅ shadcn/ui components with Material Icons
- ✅ Touch-friendly interface for mobile devices
- ✅ Sharing functionality for perfect days
- ✅ Offline-first data persistence with sync

## Notes
- All file paths are absolute from repository root
- [P] tasks target different files and can run in parallel
- Tests must fail before implementation begins (TDD)
- Mobile-first responsive design throughout
- PWA features require HTTPS for production deployment