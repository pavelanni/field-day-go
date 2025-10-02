# Data Model: Input Validation

**Feature**: Input Validation for Callsign and Email Fields
**Date**: 2025-10-01

## Visitor Entity

### Existing Structure

```go
type Visitor struct {
    ID        int       `storm:"id,increment"`
    CreatedAt time.Time `storm:"index"`
    FirstName string    `schema:"firstname"`
    LastName  string    `schema:"lastname"`
    Callsign  string    `schema:"callsign"`
    Email     string    `schema:"email"`
    Nfarl     bool      `schema:"nfarl"`
    Contactme bool      `schema:"contactme"`
    Youth     bool      `schema:"youth"`
    Firsttime bool      `schema:"firsttime"`
}
```

### New Validation Methods

#### Validate()

**Signature**: `func (v *Visitor) Validate() error`

**Purpose**: Validates all Visitor fields, returns first validation error encountered

**Logic**:

1. Validate FirstName is not empty (existing requirement)
1. If Callsign is non-empty after trimming, validate format
1. If Email is non-empty after trimming, validate format
1. Return nil if all validations pass

**Returns**:

- `nil` if all fields valid
- `*ValidationError` with field name and message if validation fails

**Example**:

```go
v := Visitor{FirstName: "John", Callsign: "123"}
err := v.Validate()
// Returns: &ValidationError{Field: "Callsign", Message: "Callsign must be in amateur radio format (e.g., W1AW, KB6NU/M)"}
```

#### ValidateCallsign()

**Signature**: `func (v *Visitor) ValidateCallsign() error`

**Purpose**: Validates callsign field format

**Logic**:

1. Trim whitespace and convert to uppercase
1. If empty string, return nil (optional field)
1. Match against regex: `^[A-Z0-9]{1,3}[0-9][A-Z]{1,4}(/[A-Z0-9]+)?$`
1. Return error with user-friendly message if invalid

**Returns**:

- `nil` if valid or empty
- `*ValidationError` if invalid format

#### ValidateEmail()

**Signature**: `func (v *Visitor) ValidateEmail() error`

**Purpose**: Validates email field format

**Logic**:

1. Trim whitespace and convert to lowercase
1. If empty string, return nil (optional field)
1. Match against regex: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
1. Return error with user-friendly message if invalid

**Returns**:

- `nil` if valid or empty
- `*ValidationError` if invalid format

### Normalization

Normalization occurs in `SaveVisitor()` before validation:

```go
func (vs *VisitorStore) SaveVisitor(v Visitor) error {
    // Normalize inputs
    v.FirstName = strings.TrimSpace(v.FirstName)
    v.LastName = strings.TrimSpace(v.LastName)
    v.Callsign = strings.ToUpper(strings.TrimSpace(v.Callsign))
    v.Email = strings.ToLower(strings.TrimSpace(v.Email))

    // Validate
    if err := v.Validate(); err != nil {
        return err
    }

    // Existing save logic...
}
```

## ValidationError Type

### Structure

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### Purpose

Carries field-specific validation error information for user display

### Fields

- `Field`: Name of the field that failed validation (e.g., "Callsign", "Email")
- `Message`: User-friendly error message explaining expected format

### Usage

```go
err := visitor.Validate()
if err != nil {
    if valErr, ok := err.(*ValidationError); ok {
        // Display field-specific error to user
        fmt.Printf("Error in %s: %s\n", valErr.Field, valErr.Message)
    }
}
```

## Package-Level Variables

### Regex Patterns

Pre-compiled at package initialization for performance:

```go
var (
    callsignRegex = regexp.MustCompile(`^[A-Z0-9]{1,3}[0-9][A-Z]{1,4}(/[A-Z0-9]+)?$`)
    emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)
```

## Validation Rules Summary

| Field | Required | Validation Rule | Normalization |
|-------|----------|----------------|---------------|
| FirstName | Yes | Non-empty string | Trim whitespace |
| LastName | No | None | Trim whitespace |
| Callsign | No | Amateur radio format (if provided) | Trim, uppercase |
| Email | No | Email format (if provided) | Trim, lowercase |
| Nfarl, Contactme, Youth, Firsttime | No | Boolean (no validation needed) | None |

## Database Schema

No changes to existing SQLite schema. Visitor table structure unchanged:

```sql
CREATE TABLE visitors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT,
    callsign TEXT,
    email TEXT,
    nfarl BOOLEAN,
    contactme BOOLEAN,
    youth BOOLEAN,
    firsttime BOOLEAN
);
```

## Backward Compatibility

- Existing records remain valid (no migration required)
- Validation only applies to new form submissions
- Reading existing records requires no changes
- Normalized data is compatible with existing queries
