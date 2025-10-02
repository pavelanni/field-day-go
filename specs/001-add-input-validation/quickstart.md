# Quickstart: Manual Testing for Input Validation

**Feature**: Input Validation for Callsign and Email Fields
**Date**: 2025-10-01

## Prerequisites

1. Build and run the application:

   ```bash
   make build
   ./fieldday test-validation.db
   ```

1. Open browser to `http://localhost:3000/new`

## Test Scenarios

### Scenario 1: Valid Callsign

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "John"
1. Enter Callsign: "W1AW"
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Redirects to confirmation page
- ✅ Callsign stored as "W1AW" (uppercase)
- ✅ Morse code plays for W1AW

**Verification**:

```bash
sqlite3 test-validation.db "SELECT first_name, callsign FROM visitors WHERE callsign='W1AW';"
```

### Scenario 2: Invalid Callsign

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Jane"
1. Enter Callsign: "123"
1. Click Submit

**Expected Result**:

- ❌ Form submission rejected
- ✅ HTTP 400 Bad Request returned
- ✅ Error message displayed: "Callsign: Callsign must be in amateur radio format (e.g., W1AW, KB6NU/M)"
- ✅ Form fields retain entered values
- ✅ No database record created

**Verification**:

```bash
sqlite3 test-validation.db "SELECT COUNT(*) FROM visitors WHERE first_name='Jane';"
# Should return 0
```

### Scenario 3: Valid Email

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Alice"
1. Enter Email: "alice@example.com"
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Email stored as "alice@example.com" (lowercase)

**Verification**:

```bash
sqlite3 test-validation.db "SELECT first_name, email FROM visitors WHERE email='alice@example.com';"
```

### Scenario 4: Invalid Email

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Bob"
1. Enter Email: "notanemail"
1. Click Submit

**Expected Result**:

- ❌ Form submission rejected
- ✅ HTTP 400 Bad Request returned
- ✅ Error message displayed: "Email: Email address must be in valid format (e.g., user@example.com)"

### Scenario 5: Empty Optional Fields

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Charlie"
1. Leave Callsign empty
1. Leave Email empty
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Both fields stored as empty strings
- ✅ Morse code plays "73" (default when no callsign)

### Scenario 6: Lowercase Callsign Normalization

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "David"
1. Enter Callsign: "w1aw" (lowercase)
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Callsign stored as "W1AW" (normalized to uppercase)

**Verification**:

```bash
sqlite3 test-validation.db "SELECT callsign FROM visitors WHERE first_name='David';"
# Should return: W1AW
```

### Scenario 7: Whitespace Trimming

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Eve"
1. Enter Callsign: "  W1AW  " (with spaces)
1. Enter Email: "  eve@test.com  " (with spaces)
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Callsign stored as "W1AW" (trimmed and uppercase)
- ✅ Email stored as "eve@test.com" (trimmed and lowercase)

**Verification**:

```bash
sqlite3 test-validation.db "SELECT callsign, email FROM visitors WHERE first_name='Eve';"
```

### Scenario 8: Multiple Validation Errors

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Frank"
1. Enter Callsign: "INVALID" (no digit)
1. Enter Email: "notanemail"
1. Click Submit

**Expected Result**:

- ❌ Form submission rejected
- ✅ HTTP 400 Bad Request returned
- ✅ Error message displayed for first validation failure (Callsign)
- ⚠️ Note: Current implementation returns first error only (acceptable MVP behavior)

### Scenario 9: Callsign with Portable Indicator

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Grace"
1. Enter Callsign: "KB6NU/M"
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Callsign stored as "KB6NU/M"

**Verification**:

```bash
sqlite3 test-validation.db "SELECT callsign FROM visitors WHERE first_name='Grace';"
```

### Scenario 10: International Callsign

**Steps**:

1. Navigate to `/new`
1. Enter FirstName: "Henry"
1. Enter Callsign: "2E0ABC" (UK callsign)
1. Click Submit

**Expected Result**:

- ✅ Form submits successfully
- ✅ Callsign stored as "2E0ABC"

## Edge Cases to Test

### Edge Case 1: Very Long Invalid Callsign

- Input: "TOOLONGCALLSIGNINVALID" (>8 characters)
- Expected: Validation error

### Edge Case 2: Special Characters in Callsign

- Input: "W1AW!" (exclamation mark)
- Expected: Validation error

### Edge Case 3: Email Without TLD

- Input: "user@domain" (missing .com/.org/etc)
- Expected: Validation error

### Edge Case 4: Email With Multiple Dots

- Input: "user.name@sub.domain.com"
- Expected: ✅ Valid

### Edge Case 5: Email With Plus Sign

- Input: "user+tag@example.com"
- Expected: ✅ Valid

### Edge Case 6: Mixed Case Email

- Input: "User@Example.COM"
- Expected: ✅ Valid, stored as "user@example.com"

## Performance Testing

### Validation Speed Test

Run with test runner to measure validation performance:

```bash
go test ./visitorstore -bench=BenchmarkValidation -benchtime=10s
```

**Expected**: <50ms per validation (target: <1ms typical)

## Cleanup

After testing:

```bash
rm test-validation.db
```

## Success Criteria

All scenarios pass with expected results:

- [ ] Scenario 1: Valid callsign accepted
- [ ] Scenario 2: Invalid callsign rejected with clear error
- [ ] Scenario 3: Valid email accepted
- [ ] Scenario 4: Invalid email rejected with clear error
- [ ] Scenario 5: Empty optional fields accepted
- [ ] Scenario 6: Lowercase callsign normalized to uppercase
- [ ] Scenario 7: Whitespace trimmed from inputs
- [ ] Scenario 8: Multiple errors handled (first error returned)
- [ ] Scenario 9: Portable indicators accepted
- [ ] Scenario 10: International callsigns accepted
- [ ] All edge cases pass
- [ ] Performance target met (<50ms validation)
