# Research: Input Validation Patterns

**Feature**: Input Validation for Callsign and Email Fields
**Date**: 2025-10-01

## Amateur Radio Callsign Patterns

### Decision

Use regex pattern `^[A-Z0-9]{1,3}[0-9][A-Z]{1,4}(/[A-Z0-9]+)?$` for callsign validation after normalization to uppercase.

### Rationale

Amateur radio callsigns follow an international standard structure defined by the ITU (International
Telecommunication Union):

- **Prefix**: 1-3 alphanumeric characters indicating country/region (e.g., W, K, N for USA, G for UK, VK for Australia)
- **Digit**: Single required digit (0-9) indicating zone or region within the country
- **Suffix**: 1-4 letters forming the unique identifier
- **Optional portable indicator**: Slash followed by suffix (e.g., /P for portable, /M for mobile, /MM for maritime
  mobile, /QRP for low power)

### Examples of Valid Callsigns

- **US callsigns**: W1AW, K2ABC, N9XYZ, AA9A, W1A
- **UK callsigns**: G4ABC, M0XYZ, 2E0ABC
- **Australian**: VK2ABC, VK3XYZ
- **Japanese**: JA1ABC, JR1XYZ
- **With indicators**: KB6NU/M, W1AW/P, G4ABC/MM

### Invalid Patterns (Should Reject)

- `123` - No letters in prefix
- `INVALID` - No required digit
- `TOOLONGCALLSIGN123` - Suffix too long (>4 characters)
- `!!!` - Special characters not allowed (except slash in portable indicator)
- `W1` - Missing suffix

### Alternatives Considered

1. **ITU callsign database validation**: Would require external data source, violates offline-first principle
1. **Country-specific regex patterns**: Too complex, maintenance burden, international events need flexibility
1. **Lenient "any alphanumeric + slash" pattern**: Too permissive, would accept clearly invalid inputs

## Email Address Validation

### Decision

Use simplified RFC 5322 pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

### Rationale

Full RFC 5322 email validation is extremely complex and includes edge cases rarely encountered in practice. A
simplified pattern covers 99.9% of real-world email addresses while being maintainable and performant.

### Pattern Components

- **Local part** (before @): Letters, digits, and common special characters (. _ % + -)
- **@ symbol**: Required separator
- **Domain part**: Letters, digits, dots, hyphens
- **TLD**: Minimum 2 characters (covers all current TLDs)

### Examples of Valid Emails

- `user@example.com`
- `test.user+tag@domain.co.uk`
- `firstname_lastname@company-name.org`
- `user123@sub.domain.com`

### Invalid Patterns (Should Reject)

- `notanemail` - No @ symbol
- `user@` - Missing domain
- `@example.com` - Missing local part
- `user@domain` - Missing TLD
- `user @example.com` - Whitespace (will be trimmed before validation)

### Alternatives Considered

1. **Full RFC 5322 validation**: Regex would be 100+ characters, supports edge cases like quoted strings and IP
   addresses in domain - unnecessary complexity
1. **Third-party email validation library**: Adds dependency, violates simplicity principle
1. **No validation**: Would allow clearly invalid inputs like "notanemail"

## Go Regex Performance

### Decision

Pre-compile regex patterns at package initialization using `regexp.MustCompile()`

### Rationale

- Regex compilation is expensive (~microseconds per compile)
- Patterns are static and never change at runtime
- Pre-compilation amortizes cost to startup time
- Validation performance target (<50ms) easily met with pre-compiled patterns

### Implementation Pattern

```go
var (
    callsignRegex = regexp.MustCompile(`^[A-Z0-9]{1,3}[0-9][A-Z]{1,4}(/[A-Z0-9]+)?$`)
    emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)
```

### Alternatives Considered

1. **Compile on each validation call**: 10-100x slower, violates performance requirement
1. **String matching without regex**: More complex code, harder to maintain, error-prone
1. **Lazy compilation with sync.Once**: Unnecessary complexity for two simple patterns

## Error Message UX

### Decision

Use clear, actionable error messages that explain expected format without technical jargon

### Example Messages

- Callsign: "Callsign must be in amateur radio format (e.g., W1AW, KB6NU/M)"
- Email: "Email address must be in valid format (e.g., user@example.com)"

### Rationale

- Field Day participants are amateur radio operators familiar with callsign formats
- Error messages should guide users toward correct input
- Avoid technical terms like "regex", "pattern", "RFC 5322"

### Alternatives Considered

1. **Technical error messages**: "Callsign does not match pattern `[A-Z0-9]{1,3}[0-9][A-Z]{1,4}`" - confusing to users
1. **Minimal messages**: "Invalid callsign" - doesn't explain what's expected
1. **No error messages**: Return generic "validation failed" - doesn't help user correct input

## Input Normalization

### Decision

Apply normalization before validation:

1. Trim leading/trailing whitespace with `strings.TrimSpace()`
1. Convert callsign to uppercase with `strings.ToUpper()`
1. Convert email to lowercase with `strings.ToLower()`

### Rationale

- Users may type callsigns in lowercase (easier on mobile keyboards)
- Email addresses are case-insensitive by RFC specification
- Whitespace is never meaningful in these fields
- Normalization improves usability without compromising validation

### Storage

- Store normalized values in database for consistency
- Existing database records already uppercase (manual entry via kiosk)
- No migration needed - normalization is forward-compatible

## Backward Compatibility

### Analysis

Existing database records (Field Day 2023-2025):

- Callsigns entered via kiosk are already uppercase (keyboard config)
- Email addresses may have mixed case
- Both fields may contain whitespace from copy-paste

### Compatibility Strategy

- Validation only applies to new submissions (form POST)
- Existing records remain unchanged
- Reading/displaying existing records requires no changes
- Normalization on save ensures future consistency

### Risk Assessment

**Risk**: Low - validation is additive, doesn't affect existing functionality
**Mitigation**: Comprehensive test coverage, including edge cases from production data
