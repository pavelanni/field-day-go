# Implementation Plan: Input Validation for Callsign and Email Fields

**Branch**: `001-add-input-validation` | **Date**: 2025-10-01 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-add-input-validation/spec.md`

## Execution Flow (/plan command scope)

1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
1. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
1. Fill the Constitution Check section based on the content of the constitution document.
1. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
1. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
1. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
1. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
1. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
1. STOP - Ready for /tasks command

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:

- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary

Add input validation for amateur radio callsign and email address fields in the visitor registration form. Validation
will normalize inputs (trim whitespace, uppercase callsigns, lowercase emails), validate formats using regex patterns,
and return user-friendly error messages when validation fails. Both fields remain optional but must be valid if
provided.

## Technical Context

**Language/Version**: Go 1.23.0
**Primary Dependencies**: Standard library (`regexp`, `strings`), existing gorilla/schema for form parsing
**Storage**: SQLite via modernc.org/sqlite (existing)
**Testing**: go test with table-driven tests for validation logic
**Target Platform**: Linux ARM7/ARM64 (Raspberry Pi, Orange Pi), macOS ARM64 (development)
**Project Type**: single (monolithic Go application with embedded templates)
**Performance Goals**: Validation execution <50ms, form submission <500ms (existing target maintained)
**Constraints**: Pure Go implementation (no CGO), offline operation, backward compatibility with existing DB records
**Scale/Scope**: Single-file change to visitorstore package, new validation functions, ~100-150 LOC added

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify feature design adheres to project constitution (`.specify/memory/constitution.md`):

### Core Principles Compliance

- [x] **I. Single Binary Deployment**: Feature uses standard library only, no new external runtime deps
- [x] **II. Pure Go Implementation**: Regex validation via stdlib `regexp` package, CGO-free
- [x] **III. Simplicity Over Framework Complexity**: Uses stdlib, adds minimal code to existing visitorstore package
- [x] **IV. Offline-First Operation**: Validation is local, no network required
- [x] **V. Mobile-Responsive Interface**: No UI changes, form submission flow unchanged

### Testing Compliance

- [x] Unit tests planned for validation functions (table-driven tests)
- [x] Integration tests planned for HTTP handler error responses
- [x] Tests use isolated temporary databases (existing pattern)
- [x] Test suite passes before commit (`go test . ./visitorstore ./morse`)

### Architecture Compliance

- [x] Package separation maintained (validation logic in visitorstore package)
- [x] Standard data flow pattern followed (HTTP → decode → **validate** → store → render)
- [x] Performance targets considered (validation <50ms easily achievable with regex)

*All checks pass - no constitutional violations.*

## Project Structure

### Documentation (this feature)

```text
specs/001-add-input-validation/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)

```text
visitorstore/
├── visitorStore.go      # Modified: Add Validate() method to Visitor, validation helpers
└── visitorStore_test.go # Modified: Add validation test cases

main.go                  # Modified: Call validation before SaveVisitor, handle errors
main_test.go             # Modified: Add test cases for validation error responses

templates/
└── new.go.html          # Modified: Display validation error messages
```

**Structure Decision**: Single project structure - all changes confined to existing visitorstore package and main.go
handler. No new packages needed for this focused enhancement.

## Phase 0: Outline & Research

**Status**: ✅ Complete

### Research Tasks Completed

1. **Amateur radio callsign patterns**: Analyzed international callsign formats (US, UK, Australia, Japan)
1. **Email validation**: Evaluated RFC 5322 patterns, chose simplified approach
1. **Go regex best practices**: Selected pre-compiled patterns for performance
1. **Error message UX**: Defined user-friendly messaging for amateur radio audience

### Output

- ✅ `research.md` created with decisions, rationale, and alternatives
- ✅ Callsign regex pattern: `^[A-Z0-9]{1,3}[0-9][A-Z]{1,4}(/[A-Z0-9]+)?$`
- ✅ Email regex pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
- ✅ All NEEDS CLARIFICATION resolved (none existed in spec)

## Phase 1: Design & Contracts

**Status**: ✅ Complete

### Design Artifacts Created

1. **Data Model** (`data-model.md`):
   - ✅ Visitor entity with new Validate(), ValidateCallsign(), ValidateEmail() methods
   - ✅ ValidationError type for field-specific errors
   - ✅ Pre-compiled regex patterns at package level
   - ✅ Normalization strategy (trim, uppercase callsign, lowercase email)
   - ✅ Backward compatibility analysis
1. **API Contracts**:
   - ✅ No new HTTP endpoints (existing POST `/new` modified)
   - ✅ New response behavior: HTTP 400 with error message on validation failure
   - ✅ Unchanged: HTTP 303 redirect on success, HTTP 500 on server error
1. **Manual Test Plan** (`quickstart.md`):
   - ✅ 10 primary test scenarios covering valid/invalid inputs
   - ✅ 6 edge cases (long callsigns, special characters, TLD validation)
   - ✅ Normalization tests (lowercase, whitespace)
   - ✅ Performance testing approach
   - ✅ Database verification queries for each scenario
1. **Agent File Update**:
   - ⏭️ Skipped - no new technologies introduced (stdlib only)
   - Existing CLAUDE.md already documents Go testing patterns

### Key Design Decisions

- Validation in visitorstore package (maintains package separation)
- First error returned (acceptable for MVP, can enhance to collect all errors later)
- Optional fields validated only if non-empty
- Pre-compiled regex for performance (<1ms validation time)

## Phase 2: Task Planning Approach

**Status**: ✅ Complete (planning only, tasks.md created by /tasks command)

### Task Generation Strategy

The `/tasks` command will generate tasks from Phase 1 design documents:

**From data-model.md**:

- Validation method implementation tasks (Validate, ValidateCallsign, ValidateEmail)
- ValidationError type creation
- Package-level regex pattern definitions

**From quickstart.md**:

- Unit test tasks for each validation method (table-driven tests)
- Integration test tasks for HTTP handler error responses
- Test cases covering 10 scenarios + 6 edge cases

**From modified flow in main.go**:

- Update SaveVisitor to normalize inputs
- Update newVisitorHandler to call validation before save
- Update error handling to return HTTP 400 with error message

### Ordering Strategy

**TDD Approach**:

1. **Phase 3.1**: Setup (none needed - using existing project)
1. **Phase 3.2**: Write failing tests FIRST (MUST complete before 3.3)
   - Validation unit tests in visitorstore_test.go [P]
   - HTTP handler integration tests in main_test.go [P]
1. **Phase 3.3**: Implementation (only after tests fail)
   - ValidationError type
   - Regex patterns
   - Validation methods
   - Update SaveVisitor normalization
   - Update newVisitorHandler error handling
1. **Phase 3.5**: Polish and verification
   - Run full test suite
   - Execute quickstart.md manual tests
   - Performance validation

**Parallelization**: Tests for different packages can run in parallel [P] since they use independent test databases

### Estimated Output

**8-10 tasks** in tasks.md:

- 2 test tasks (unit + integration) [P]
- 5-6 implementation tasks (sequential within visitorstore)
- 2 polish tasks (test suite + quickstart)

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation

*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md following constitutional principles)
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |

## Progress Tracking

**Phase Status**:

- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [x] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:

- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved (spec had none)
- [x] Complexity deviations documented (none - no violations)

**Artifacts Generated**:

- [x] plan.md (this file)
- [x] research.md (amateur radio callsign and email validation patterns)
- [x] data-model.md (Visitor validation methods, ValidationError type)
- [x] quickstart.md (10 test scenarios + 6 edge cases)
- [x] tasks.md (generated by /tasks command)

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
