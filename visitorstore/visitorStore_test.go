package visitorstore

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewVisitorStore(t *testing.T) {
	// Test case: successful database initialization
	dbFile := "test_new.db"
	defer os.Remove(dbFile) // Clean up

	vs, err := NewVisitorStore(dbFile)
	if err != nil {
		t.Errorf("NewVisitorStore() error = %v, want nil", err)
		return
	}
	if vs == nil {
		t.Errorf("NewVisitorStore() = nil, want non-nil")
		return
	}
	if vs.db == nil {
		t.Errorf("NewVisitorStore().db = nil, want non-nil")
		return
	}
	vs.db.Close()

	// Test case: invalid database path
	dbFile = "/invalid/path/test.db"
	vs, err = NewVisitorStore(dbFile)
	if err == nil {
		t.Errorf("NewVisitorStore() error = nil, want non-nil for invalid path")
	}
	if vs != nil {
		t.Errorf("NewVisitorStore() = %v, want nil for invalid path", vs)
	}
}

func TestSaveVisitor(t *testing.T) {
	dbFile := "test_save.db"
	defer os.Remove(dbFile)

	vs, err := NewVisitorStore(dbFile)
	if err != nil {
		t.Fatalf("Failed to create visitor store: %v", err)
	}
	defer vs.db.Close()

	// Test case: Save a valid visitor
	v := Visitor{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Callsign:  "W1AW",
		CreatedAt: time.Now(),
	}
	err = vs.SaveVisitor(v)
	if err != nil {
		t.Errorf("Failed to save visitor: %v", err)
	}

	// Test case: Save visitor with missing required field (validation is now in Validate())
	v = Visitor{
		FirstName: "",
		LastName:  "Doe",
		Email:     "test@example.com",
	}
	err = v.Validate()
	if err == nil {
		t.Error("Expected validation error when FirstName is missing")
	}

	// Test case: Save visitor with CreatedAt auto-populated
	v = Visitor{
		FirstName: "Jane",
		LastName:  "Smith",
	}
	err = vs.SaveVisitor(v)
	if err != nil {
		t.Errorf("Failed to save visitor with auto-populated CreatedAt: %v", err)
	}
}

func TestListVisitors(t *testing.T) {
	dbFile := "test_list.db"
	defer os.Remove(dbFile)

	vs, err := NewVisitorStore(dbFile)
	if err != nil {
		t.Fatalf("Failed to create visitor store: %v", err)
	}
	defer vs.db.Close()

	// Test case: empty visitor list
	visitors, err := vs.ListVisitors()
	if err != nil {
		t.Errorf("ListVisitors() error = %v, want nil", err)
	}
	if len(visitors) != 0 {
		t.Errorf("ListVisitors() = %v, want empty slice", visitors)
	}

	// Test case: non-empty visitor list
	visitor1 := Visitor{FirstName: "Alice", LastName: "Smith", CreatedAt: time.Now()}
	visitor2 := Visitor{FirstName: "Bob", LastName: "Johnson", CreatedAt: time.Now()}

	err = vs.SaveVisitor(visitor1)
	if err != nil {
		t.Fatal(err)
	}
	err = vs.SaveVisitor(visitor2)
	if err != nil {
		t.Fatal(err)
	}

	visitors, err = vs.ListVisitors()
	if err != nil {
		t.Errorf("ListVisitors() error = %v, want nil", err)
	}
	if len(visitors) != 2 {
		t.Errorf("ListVisitors() = %d visitors, want 2", len(visitors))
	}

	// Verify visitors are ordered by ID
	if visitors[0].ID != 1 || visitors[0].FirstName != "Alice" {
		t.Errorf("ListVisitors()[0] = %+v, want Alice with ID 1", visitors[0])
	}
	if visitors[1].ID != 2 || visitors[1].FirstName != "Bob" {
		t.Errorf("ListVisitors()[1] = %+v, want Bob with ID 2", visitors[1])
	}
}

func TestTotalVisitors(t *testing.T) {
	dbFile := "test_total.db"
	defer os.Remove(dbFile)

	vs, err := NewVisitorStore(dbFile)
	if err != nil {
		t.Fatalf("Failed to create visitor store: %v", err)
	}
	defer vs.db.Close()

	// Test case: empty database
	total, err := vs.TotalVisitors()
	if err != nil {
		t.Errorf("TotalVisitors() error = %v, want nil", err)
	}
	if total != 0 {
		t.Errorf("TotalVisitors() = %d, want 0", total)
	}

	// Test case: after adding visitors
	visitors := []Visitor{
		{FirstName: "Alice", LastName: "Smith", CreatedAt: time.Now()},
		{FirstName: "Bob", LastName: "Johnson", CreatedAt: time.Now()},
		{FirstName: "Charlie", LastName: "Brown", CreatedAt: time.Now()},
	}

	for _, v := range visitors {
		err = vs.SaveVisitor(v)
		if err != nil {
			t.Fatal(err)
		}
	}

	total, err = vs.TotalVisitors()
	if err != nil {
		t.Errorf("TotalVisitors() error = %v, want nil", err)
	}
	if total != 3 {
		t.Errorf("TotalVisitors() = %d, want 3", total)
	}
}

func TestImportFromCSV(t *testing.T) {
	dbFile := "test_import.db"
	defer os.Remove(dbFile)

	vs, err := NewVisitorStore(dbFile)
	if err != nil {
		t.Fatalf("Failed to create visitor store: %v", err)
	}
	defer vs.db.Close()

	// Create a test CSV file
	csvFile := "test_data.csv"
	csvContent := `Callsign,Contactme,CreatedAt,Email,FirstName,Firsttime,ID,LastName,Nfarl,Youth,id
W1AW,true,2025-06-28T10:30:35.463982438-04:00,w1aw@arrl.org,Test,false,1,User,true,false,
N0CALL,false,2025-06-28T10:35:09.316541124-04:00,,Example,true,2,Person,false,true,
`
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	defer os.Remove(csvFile)

	// Test CSV import
	err = vs.ImportFromCSV(csvFile)
	if err != nil {
		t.Errorf("ImportFromCSV() error = %v, want nil", err)
	}

	// Verify imported data
	total, err := vs.TotalVisitors()
	if err != nil {
		t.Errorf("TotalVisitors() after import error = %v", err)
	}
	if total != 2 {
		t.Errorf("TotalVisitors() after import = %d, want 2", total)
	}

	visitors, err := vs.ListVisitors()
	if err != nil {
		t.Errorf("ListVisitors() after import error = %v", err)
	}

	if len(visitors) != 2 {
		t.Fatalf("Expected 2 visitors, got %d", len(visitors))
	}

	// Check first visitor
	if visitors[0].FirstName != "Test" || visitors[0].Callsign != "W1AW" {
		t.Errorf("First visitor = %+v, want Test/W1AW", visitors[0])
	}

	// Check second visitor
	if visitors[1].FirstName != "Example" || visitors[1].Youth != true {
		t.Errorf("Second visitor = %+v, want Example with Youth=true", visitors[1])
	}

	// Test case: non-existent CSV file
	err = vs.ImportFromCSV("nonexistent.csv")
	if err == nil {
		t.Error("ImportFromCSV() with non-existent file should return error")
	}
}

func TestVisitorValidateCallsign(t *testing.T) {
	tests := []struct {
		name      string
		callsign  string
		wantError bool
	}{
		// Valid callsigns
		{name: "valid US callsign", callsign: "W1AW", wantError: false},
		{name: "valid US callsign lowercase", callsign: "w1aw", wantError: false},
		{name: "valid with portable indicator", callsign: "KB6NU/M", wantError: false},
		{name: "valid UK callsign", callsign: "2E0ABC", wantError: false},
		{name: "valid Australian callsign", callsign: "VK2ABC", wantError: false},
		{name: "valid short callsign", callsign: "W1A", wantError: false},
		{name: "empty callsign (optional)", callsign: "", wantError: false},
		{name: "whitespace only (optional)", callsign: "   ", wantError: false},

		// Invalid callsigns
		{name: "no letters", callsign: "123", wantError: true},
		{name: "no digit", callsign: "INVALID", wantError: true},
		{name: "suffix too long", callsign: "W1TOOLONG", wantError: true},
		{name: "special characters", callsign: "W1AW!", wantError: true},
		{name: "missing suffix", callsign: "W1", wantError: true},
		{name: "only prefix", callsign: "ABC", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Visitor{Callsign: tt.callsign}
			err := v.ValidateCallsign()

			if tt.wantError && err == nil {
				t.Errorf("ValidateCallsign() error = nil, want error for callsign %q", tt.callsign)
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateCallsign() error = %v, want nil for callsign %q", err, tt.callsign)
			}
		})
	}
}

func TestVisitorValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		// Valid emails
		{name: "valid simple email", email: "user@example.com", wantError: false},
		{name: "valid with subdomain", email: "user@mail.example.com", wantError: false},
		{name: "valid with plus sign", email: "user+tag@example.com", wantError: false},
		{name: "valid with dots", email: "first.last@example.com", wantError: false},
		{name: "valid UK domain", email: "test@example.co.uk", wantError: false},
		{name: "valid mixed case", email: "User@Example.COM", wantError: false},
		{name: "empty email (optional)", email: "", wantError: false},
		{name: "whitespace only (optional)", email: "   ", wantError: false},

		// Invalid emails
		{name: "no at sign", email: "notanemail", wantError: true},
		{name: "missing domain", email: "user@", wantError: true},
		{name: "missing local part", email: "@example.com", wantError: true},
		{name: "missing TLD", email: "user@domain", wantError: true},
		{name: "spaces in email", email: "user @example.com", wantError: true},
		{name: "multiple at signs", email: "user@@example.com", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Visitor{Email: tt.email}
			err := v.ValidateEmail()

			if tt.wantError && err == nil {
				t.Errorf("ValidateEmail() error = nil, want error for email %q", tt.email)
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateEmail() error = %v, want nil for email %q", err, tt.email)
			}
		})
	}
}

func TestVisitorValidate(t *testing.T) {
	tests := []struct {
		name       string
		visitor    Visitor
		wantError  bool
		errorField string
	}{
		{
			name: "valid with all fields",
			visitor: Visitor{
				FirstName: "John",
				Callsign:  "W1AW",
				Email:     "john@example.com",
			},
			wantError: false,
		},
		{
			name: "valid with empty optionals",
			visitor: Visitor{
				FirstName: "Jane",
				Callsign:  "",
				Email:     "",
			},
			wantError: false,
		},
		{
			name: "invalid callsign",
			visitor: Visitor{
				FirstName: "Bob",
				Callsign:  "INVALID",
				Email:     "bob@example.com",
			},
			wantError:  true,
			errorField: "Callsign",
		},
		{
			name: "invalid email",
			visitor: Visitor{
				FirstName: "Alice",
				Callsign:  "W1AW",
				Email:     "notanemail",
			},
			wantError:  true,
			errorField: "Email",
		},
		{
			name: "both invalid - should return first error (callsign)",
			visitor: Visitor{
				FirstName: "Charlie",
				Callsign:  "123",
				Email:     "notanemail",
			},
			wantError:  true,
			errorField: "Callsign",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.visitor.Validate()

			if tt.wantError && err == nil {
				t.Errorf("Validate() error = nil, want error")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
			if tt.wantError && err != nil {
				if verr, ok := err.(*ValidationError); ok {
					if verr.Field != tt.errorField {
						t.Errorf("Validate() error field = %q, want %q", verr.Field, tt.errorField)
					}
				} else {
					t.Errorf("Validate() error type = %T, want *ValidationError", err)
				}
			}
		})
	}
}

func TestNormalization(t *testing.T) {
	dbFile := "test_normalization.db"
	defer os.Remove(dbFile)

	vs, err := NewVisitorStore(dbFile)
	if err != nil {
		t.Fatalf("Failed to create visitor store: %v", err)
	}
	defer vs.db.Close()

	tests := []struct {
		name              string
		input             Visitor
		expectedCallsign  string
		expectedEmail     string
		expectedFirstName string
	}{
		{
			name: "uppercase callsign",
			input: Visitor{
				FirstName: "Test",
				Callsign:  "w1aw",
				Email:     "test@example.com",
			},
			expectedCallsign:  "W1AW",
			expectedEmail:     "test@example.com",
			expectedFirstName: "Test",
		},
		{
			name: "lowercase email",
			input: Visitor{
				FirstName: "Test",
				Callsign:  "W1AW",
				Email:     "Test@EXAMPLE.COM",
			},
			expectedCallsign:  "W1AW",
			expectedEmail:     "test@example.com",
			expectedFirstName: "Test",
		},
		{
			name: "trim whitespace",
			input: Visitor{
				FirstName: "  Test  ",
				LastName:  "  User  ",
				Callsign:  "  W1AW  ",
				Email:     "  test@example.com  ",
			},
			expectedCallsign:  "W1AW",
			expectedEmail:     "test@example.com",
			expectedFirstName: "Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normalize FirstName/LastName (done in handler in production)
			tt.input.FirstName = strings.TrimSpace(tt.input.FirstName)
			tt.input.LastName = strings.TrimSpace(tt.input.LastName)

			// Validate normalizes callsign and email
			if err := tt.input.Validate(); err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}

			err := vs.SaveVisitor(tt.input)
			if err != nil {
				t.Fatalf("SaveVisitor() error = %v, want nil", err)
			}

			visitors, err := vs.ListVisitors()
			if err != nil {
				t.Fatalf("ListVisitors() error = %v", err)
			}

			if len(visitors) == 0 {
				t.Fatal("No visitors found after save")
			}

			last := visitors[len(visitors)-1]
			if last.Callsign != tt.expectedCallsign {
				t.Errorf("Callsign = %q, want %q", last.Callsign, tt.expectedCallsign)
			}
			if last.Email != tt.expectedEmail {
				t.Errorf("Email = %q, want %q", last.Email, tt.expectedEmail)
			}
			if last.FirstName != tt.expectedFirstName {
				t.Errorf("FirstName = %q, want %q", last.FirstName, tt.expectedFirstName)
			}
		})
	}
}
