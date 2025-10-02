# Tasks: Input Validation for Callsign and Email Fields

**Input**: Design documents from `/specs/001-add-input-validation/`
**Prerequisites**: plan.md (required), research.md, data-model.md, quickstart.md

## Execution Flow (main)

1. Load plan.md from feature directory
   - If not found: ERROR "No implementation plan found"
   - Extract: tech stack, libraries, structure
1. Load optional design documents:
   - data-model.md: Extract entities → model tasks
   - contracts/: Each file → contract test task
   - research.md: Extract decisions → setup tasks
1. Generate tasks by category:
   - Setup: project init, dependencies, linting
   - Tests: contract tests, integration tests
   - Core: models, services, CLI commands
   - Integration: DB, middleware, logging
   - Polish: unit tests, performance, docs
1. Apply task rules:
   - Different files = mark [P] for parallel
   - Same file = sequential (no [P])
   - Tests before implementation (TDD)
1. Number tasks sequentially (T001, T002...)
1. Generate dependency graph
1. Create parallel execution examples
1. Validate task completeness:
   - All contracts have tests?
   - All entities have models?
   - All endpoints implemented?
1. Return: SUCCESS (tasks ready for execution)

## Format: `[ID] [P?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Phase 3.1: Setup

No setup tasks needed - using existing Go project structure and dependencies.

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [ ] T001 [P] Write validation unit tests in visitorstore/visitorStore_test.go covering:
  - TestVisitorValidateCallsign with table-driven tests (W1AW, KB6NU/M, 2E0ABC valid; 123, INVALID, TOOLONG
    invalid)
  - TestVisitorValidateEmail with table-driven tests (user@example.com, test+tag@domain.co.uk valid; notanemail,
    user@, @domain invalid)
  - TestVisitorValidate for combined validation (empty optionals valid, first error returned)
  - TestNormalization for trim/uppercase/lowercase behavior
- [ ] T002 [P] Write HTTP handler integration tests in main_test.go covering:
  - TestNewVisitorHandler_InvalidCallsign (returns HTTP 400, error message displayed)
  - TestNewVisitorHandler_InvalidEmail (returns HTTP 400, error message displayed)
  - TestNewVisitorHandler_ValidInputs (returns HTTP 303, record saved with normalized values)
  - TestNewVisitorHandler_EmptyOptionalFields (returns HTTP 303, empty strings accepted)

## Phase 3.3: Core Implementation (ONLY after tests are failing)

- [ ] T003 Add ValidationError type to visitorstore/visitorStore.go with Field and Message fields, implement
  Error() method
- [ ] T004 Add package-level regex patterns to visitorstore/visitorStore.go (callsignRegex, emailRegex) using
  regexp.MustCompile
- [ ] T005 Implement Visitor.ValidateCallsign() method in visitorstore/visitorStore.go (trim, uppercase, optional
  check, regex match, error message)
- [ ] T006 Implement Visitor.ValidateEmail() method in visitorstore/visitorStore.go (trim, lowercase, optional
  check, regex match, error message)
- [ ] T007 Implement Visitor.Validate() method in visitorstore/visitorStore.go (call ValidateCallsign and
  ValidateEmail, return first error)
- [ ] T008 Update VisitorStore.SaveVisitor() in visitorstore/visitorStore.go to normalize inputs (trim
  FirstName/LastName, uppercase Callsign, lowercase Email) and call Validate() before save
- [ ] T009 Update newVisitorHandler in main.go to handle ValidationError (return HTTP 400, render error message in template)
- [ ] T010 Update templates/new.go.html to display validation error message if present

## Phase 3.4: Integration

No integration tasks needed - validation integrated into existing SaveVisitor flow.

## Phase 3.5: Polish

- [ ] T011 Run full test suite to verify all tests pass: `go test . ./visitorstore ./morse`
- [ ] T012 Execute manual testing scenarios from specs/001-add-input-validation/quickstart.md (10 scenarios + 6 edge cases)
- [ ] T013 Run performance validation: `go test ./visitorstore -bench=BenchmarkValidation` (verify <50ms target)

## Dependencies

- Tests (T001-T002) before implementation (T003-T010)
- T003 (ValidationError type) blocks T005, T006, T007
- T004 (regex patterns) blocks T005, T006
- T005, T006 block T007 (Validate calls both methods)
- T007 blocks T008 (SaveVisitor calls Validate)
- T008 blocks T009 (handler depends on SaveVisitor behavior)
- T009 blocks T010 (template depends on handler error passing)
- Implementation (T003-T010) before polish (T011-T013)

## Parallel Example

```text
# Launch T001-T002 together (different test files):
Task: "Write validation unit tests in visitorstore/visitorStore_test.go"
Task: "Write HTTP handler integration tests in main_test.go"
```

## Notes

- [P] tasks = different files, no dependencies
- Verify tests fail before implementing (TDD approach)
- Commit after each task with descriptive message
- No external dependencies added (stdlib only)
- Maintains constitutional principles (pure Go, offline-first, simplicity)

## Task Generation Rules Applied

1. **From data-model.md**:
   - ValidationError type → T003
   - Regex patterns → T004
   - ValidateCallsign method → T005
   - ValidateEmail method → T006
   - Validate method → T007
   - SaveVisitor normalization → T008

1. **From quickstart.md**:
   - Unit test scenarios → T001
   - Integration test scenarios → T002
   - Manual testing → T012
   - Performance testing → T013

1. **From plan.md flow**:
   - Handler error handling → T009
   - Template error display → T010
   - Full test suite → T011

1. **Ordering**:
   - Tests (T001-T002) before implementation (T003-T010)
   - Type definitions (T003-T004) before methods (T005-T007)
   - Core validation (T005-T007) before integration (T008-T010)
   - Implementation before polish (T011-T013)

## Validation Checklist

- [x] All entities have implementation tasks (Visitor validation methods)
- [x] All tests come before implementation (T001-T002 before T003-T010)
- [x] Parallel tasks truly independent (T001-T002 use different files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Constitutional compliance maintained (pure Go, no new dependencies)
