# Feature Specification: Input Validation for Callsign and Email Fields

**Feature Branch**: `001-add-input-validation`

**Created**: 2025-10-01

**Status**: Draft

**Input**: User description: "Add input validation for the call sign and email inputs. Analyze typical patterns
for amateur radio callsigns for that. Usually they consist of several letters and numbers, sometimes include a slash (/).
For email validation use standard rules for email addresses."

## User Scenarios & Testing

### Primary User Story

When a visitor submits the Field Day registration form with invalid data (malformed callsign or email), the system
should reject the submission with clear error messages explaining what format is expected, preventing invalid data from
being stored in the database.

### Acceptance Scenarios

1. **Given** a visitor enters a valid amateur radio callsign (e.g., "W1AW", "KB6NU/M", "2E0ABC"), **When** they submit
   the form, **Then** the system accepts and saves the visitor record
1. **Given** a visitor enters an invalid callsign (e.g., "123", "!!!INVALID", "TOOLONGCALLSIGN"), **When** they submit
   the form, **Then** the system rejects the submission and displays an error message explaining valid callsign format
1. **Given** a visitor enters a valid email address (e.g., "user@example.com"), **When** they submit the form, **Then**
   the system accepts and saves the visitor record
1. **Given** a visitor enters an invalid email (e.g., "notanemail", "user@", "@example.com"), **When** they submit the
   form, **Then** the system rejects the submission and displays an error message
1. **Given** a visitor leaves callsign and email fields empty (optional fields), **When** they submit the form with a
   valid first name, **Then** the system accepts the registration (these fields are optional)

### Edge Cases

- What happens when a callsign contains lowercase letters? System should accept and normalize to uppercase
- What happens when a callsign has leading/trailing whitespace? System should trim whitespace before validation
- What happens when multiple validation errors occur simultaneously (invalid callsign AND invalid email)? System
  should display all relevant error messages
- What happens with international callsign formats (e.g., special characters, varying lengths)? System should
  accommodate common international patterns
- What happens when email contains uppercase letters? System should accept and normalize to lowercase

## Requirements

### Functional Requirements

- **FR-001**: System MUST validate amateur radio callsign format when a non-empty callsign is provided
- **FR-002**: System MUST accept callsigns matching standard amateur radio patterns: 1-2 character prefix
  (letters/numbers), followed by a digit, followed by 1-4 letters, optionally followed by "/" and a suffix
- **FR-003**: System MUST accept common callsign variants including portable indicators (e.g., /P, /M, /MM), vanity
  callsigns, and international prefixes
- **FR-004**: System MUST validate email address format when a non-empty email is provided
- **FR-005**: System MUST accept email addresses conforming to standard RFC 5322 simplified pattern: local-part@domain
  with valid characters
- **FR-006**: System MUST return clear, actionable error messages when validation fails, specifying which field is
  invalid and what format is expected
- **FR-007**: System MUST allow empty/blank values for callsign and email fields (these are optional fields)
- **FR-008**: System MUST trim leading and trailing whitespace from callsign and email inputs before validation
- **FR-009**: System MUST normalize callsigns to uppercase for consistency
- **FR-010**: System MUST normalize email addresses to lowercase for consistency

### Non-Functional Requirements

- **NFR-001**: Validation errors MUST be returned within 50ms to maintain form responsiveness
- **NFR-002**: Error messages MUST be user-friendly and avoid technical jargon
- **NFR-003**: Validation logic MUST be testable in isolation (unit tests)
- **NFR-004**: Validation MUST NOT break backward compatibility with existing valid records in the database

### Key Entities

- **Visitor**: The existing entity that will gain validation methods for its Callsign and Email fields
- **ValidationError**: Error type that will carry field-specific validation failure information

## Review & Acceptance Checklist

### Content Quality

- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness

- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

## Execution Status

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [x] Review checklist passed
