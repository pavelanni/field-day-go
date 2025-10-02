<!--
Sync Impact Report:
Version: 0.0.0 → 1.0.0
Initial constitution creation for Field Day Registration Kiosk project.

Modified principles: N/A (initial creation)
Added sections: All sections (initial creation)
Removed sections: None

Templates requiring updates:
✅ .specify/templates/plan-template.md - reviewed, compatible
✅ .specify/templates/spec-template.md - reviewed, compatible
✅ .specify/templates/tasks-template.md - reviewed, compatible
✅ .specify/templates/agent-file-template.md - reviewed, compatible

Follow-up TODOs: None
-->

# Field Day Registration Kiosk Constitution

## Core Principles

### I. Single Binary Deployment

The application MUST compile to a single, self-contained binary with all templates, static files, and assets embedded using Go's embed package. External dependencies at runtime are prohibited except for the SQLite database file provided as a command-line argument.

**Rationale**: Field Day events operate in remote locations without reliable internet connectivity. Single-binary deployment eliminates dependency management issues, simplifies installation on single-board computers (Raspberry Pi, Orange Pi), and enables offline operation.

### II. Pure Go Implementation

All dependencies MUST be pure Go implementations without CGO requirements. This principle applies to database drivers (modernc.org/sqlite), audio generation, and all other libraries.

**Rationale**: CGO-free compilation enables fast cross-compilation from development machines to ARM targets, eliminates GCC toolchain dependencies on build machines, and simplifies the build process. This directly supports the rapid development and deployment cycle needed for annual Field Day events.

### III. Simplicity Over Framework Complexity

The application MUST use Go's standard library (`net/http`, `html/template`) as the foundation, adding minimal external dependencies only when clear value is demonstrated. Full-featured web frameworks are prohibited.

**Rationale**: This is an extremely simple application (visitor registration form with database storage). Framework overhead adds complexity without benefit. Standard library code is more maintainable, has better long-term stability, and is easier for contributors to understand. Current exceptions: `gorilla/schema` for form parsing, Tailwind CSS for responsive mobile styling.

### IV. Offline-First Operation

All features MUST function without network connectivity. Time synchronization MUST use hardware RTC (real-time clock) modules. Data export and import MUST work with local files.

**Rationale**: Field Day events occur in parks and remote locations without reliable internet. The application must maintain accurate timestamps (via RTC) and operate entirely offline for 24+ hours. Network-dependent features compromise the core use case.

### V. Mobile-Responsive Interface

The web interface MUST provide full functionality on mobile devices with touch-friendly controls and responsive layouts using Tailwind CSS.

**Rationale**: Field Day 2025 deployment includes cloud hosting for remote visitor registration via mobile devices. The interface must work equally well on kiosk displays and smartphones without separate codebases.

## Testing Requirements

### Test Coverage Standards

1. **Unit Tests**: Required for all packages (`visitorstore`, `morse`). Tests MUST cover data validation, database operations, and business logic.
2. **Integration Tests**: Required for HTTP handlers. Tests MUST verify form submission, data persistence, and response generation.
3. **Test Isolation**: Each test MUST create and clean up its own temporary database. No shared state between tests.
4. **Pre-Commit Testing**: Run `go test . ./visitorstore ./morse` before all commits. All tests MUST pass.

**Rationale**: The application is deployed in unattended kiosk mode at events. Bugs disrupt the user experience for radio operators and visitors. Comprehensive testing prevents regression and ensures reliability.

## Development Workflow

### Build and Development Commands

All common operations MUST be documented in CLAUDE.md and accessible via Makefile or direct `go` commands:

- `make build` - Build for current platform
- `make build-raspi` - Cross-compile for ARM7
- `go test . ./visitorstore ./morse` - Run test suite
- `npx tailwindcss -i input.css -o static/css/tailwind.css --watch` - CSS development

### Deployment Process

Production deployment follows this sequence:
1. Run test suite - all tests MUST pass
2. Build binary via `make build` or `make build-raspi`
3. Test binary with `./fieldday test.db`
4. Install to production via `make install` (copies to `/opt/fieldday`, configures systemd)

## Code Organization

### Package Structure

- **main.go**: HTTP server, route handlers, embedded templates/static files
- **visitorstore/**: SQLite database operations, visitor struct, schema management
- **morse/**: Morse code generation, WAV audio file creation
- **cmd/**: Utility programs (importcsv, exportbolt legacy tool)
- **templates/**: Go HTML templates
- **static/**: CSS, JavaScript, images

**Separation of Concerns**: HTTP handling (main.go), data operations (visitorstore), domain logic (morse) MUST remain in separate packages with clear responsibilities.

### Data Flow Pattern

Standard request flow:
1. HTTP request → handler in main.go
2. Form data decoded via gorilla/schema → Visitor struct
3. Visitor struct validated and saved via visitorstore package
4. Response template rendered with embedded file system
5. Morse code audio generated on-demand via /morse-audio endpoint

## Constraints

### Performance Targets

- Form submission response time: < 500ms (including database write)
- Morse audio generation: < 200ms for typical callsigns
- List page rendering: < 1s for 500 visitors

**Rationale**: These targets ensure responsive kiosk operation on single-board computers (Raspberry Pi 3+, Orange Pi Zero 3) under typical Field Day conditions.

### Platform Support

- **Primary**: Linux ARM7 (Raspberry Pi), Linux ARM64 (Orange Pi Zero 3)
- **Development**: macOS ARM64 (Apple Silicon)
- **Secondary**: Linux x86_64

Cross-compilation and cross-platform compatibility are mandatory due to development on macOS and deployment on ARM Linux boards.

## Governance

### Amendment Process

1. Proposed changes MUST be documented with rationale in a pull request
2. Changes affecting core principles (I-V) require explicit justification of how they serve Field Day deployment needs
3. `CONSTITUTION_VERSION` MUST be incremented per semantic versioning:
   - MAJOR: Removing or redefining core principles
   - MINOR: Adding new principles or expanding guidance
   - PATCH: Clarifications, wording improvements, non-semantic changes

### Compliance Verification

- All feature specifications MUST reference constitution principles in their design rationale
- Code reviews MUST verify adherence to Single Binary Deployment, Pure Go Implementation, and Simplicity principles
- Deviations from principles MUST be documented in the feature plan's Complexity Tracking section with justification

### Runtime Development Guidance

Agent-specific development instructions are maintained in project root:
- **CLAUDE.md**: Guidance for Claude Code (architecture, commands, patterns)
- **AGENTS.md**: Cross-agent guidance (if applicable)

These files provide implementation-level guidance while this constitution establishes non-negotiable architectural principles.

**Version**: 1.0.0 | **Ratified**: 2025-10-01 | **Last Amended**: 2025-10-01
