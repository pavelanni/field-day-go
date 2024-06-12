package visitorstore

import (
	"testing"

	"github.com/asdine/storm"
)

func TestSaveVisitor(t *testing.T) {
	// Test case: Save a new visitor
	db, err := storm.Open("test.db")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	vs := &VisitorStore{db}
	v := Visitor{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	err = vs.SaveVisitor(v)
	if err != nil {
		t.Errorf("Failed to save visitor: %v", err)
	}

	// Test case: Save an existing visitor
	v.ID = 1
	err = vs.SaveVisitor(v)
	if err != nil {
		t.Errorf("Failed to save visitor: %v", err)
	}

	// Test case: Save a visitor with missing fields
	v = Visitor{}
	err = vs.SaveVisitor(v)
	if err == nil {
		t.Error("Expected error when saving visitor with missing fields")
	}
	err = vs.db.Drop(&v)
	if err != nil {
		t.Errorf("Failed to drop bucket Visitor: %v", err)
	}
}

func TestListVisitors(t *testing.T) {
	// Test case: empty visitor list
	db, err := storm.Open("test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	vs := &VisitorStore{db: db}
	visitors, err := vs.ListVisitors()
	if err != nil {
		t.Errorf("ListVisitors() error = %v, want nil", err)
	}
	if len(visitors) != 0 {
		t.Errorf("ListVisitors() = %v, want empty slice", visitors)
	}

	// Test case: non-empty visitor list
	visitor1 := Visitor{ID: 1, FirstName: "Alice", LastName: "Smith"}
	visitor2 := Visitor{ID: 2, FirstName: "Bob", LastName: "Johnson"}
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
		t.Errorf("ListVisitors() = %v, want slice of length 2", visitors)
	}
	if visitors[0].ID != visitor1.ID || visitors[0].FirstName != visitor1.FirstName || visitors[0].LastName != visitor1.LastName {
		t.Errorf("ListVisitors()[0] = %v, want %v", visitors[0], visitor1)
	}
	if visitors[1].ID != visitor2.ID || visitors[1].FirstName != visitor2.FirstName || visitors[1].LastName != visitor2.LastName {
		t.Errorf("ListVisitors()[1] = %v, want %v", visitors[1], visitor2)
	}

	// Test case: error from db.AllByIndex
	err = vs.db.Drop(&visitor1)
	if err != nil {
		t.Fatal(err)
	}
	err = vs.db.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = vs.ListVisitors()
	if err == nil {
		t.Errorf("ListVisitors() error = nil, want non-nil")
	}
}
