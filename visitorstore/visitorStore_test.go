package visitorstore

import (
	"os"
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

	// Test case: Save visitor with missing required field
	v = Visitor{
		LastName: "Doe",
		Email:    "test@example.com",
	}
	err = vs.SaveVisitor(v)
	if err == nil {
		t.Error("Expected error when saving visitor with missing FirstName")
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