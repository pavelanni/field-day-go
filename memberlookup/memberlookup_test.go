package memberlookup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMembers(t *testing.T) {
	tmp := t.TempDir()
	csvPath := filepath.Join(tmp, "members.csv")
	content := "callsign,first_name,last_name,email\nW1AW,Hiram,Maxim,hiram@arrl.org\nK1ABC,John,Doe,john@example.com\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ml, err := LoadCSV(csvPath)
	if err != nil {
		t.Fatalf("LoadCSV failed: %v", err)
	}

	_, ok1 := ml.Lookup("W1AW")
	_, ok2 := ml.Lookup("K1ABC")
	if !ok1 || !ok2 {
		t.Errorf("expected both W1AW and K1ABC to be found")
	}
}

func TestLookupFound(t *testing.T) {
	tmp := t.TempDir()
	csvPath := filepath.Join(tmp, "members.csv")
	content := "callsign,first_name,last_name,email\nW1AW,Hiram,Maxim,hiram@arrl.org\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ml, _ := LoadCSV(csvPath)

	m, ok := ml.Lookup("W1AW")
	if !ok {
		t.Fatal("expected to find W1AW")
	}
	if m.FirstName != "Hiram" {
		t.Errorf("expected Hiram, got %q", m.FirstName)
	}
	if m.LastName != "Maxim" {
		t.Errorf("expected Maxim, got %q", m.LastName)
	}
	if m.Email != "hiram@arrl.org" {
		t.Errorf("expected hiram@arrl.org, got %q", m.Email)
	}
}

func TestLookupNotFound(t *testing.T) {
	tmp := t.TempDir()
	csvPath := filepath.Join(tmp, "members.csv")
	content := "callsign,first_name,last_name,email\nW1AW,Hiram,Maxim,hiram@arrl.org\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ml, _ := LoadCSV(csvPath)

	_, ok := ml.Lookup("XYZ999")
	if ok {
		t.Error("expected not found for XYZ999")
	}
}

func TestLookupCaseInsensitive(t *testing.T) {
	tmp := t.TempDir()
	csvPath := filepath.Join(tmp, "members.csv")
	content := "callsign,first_name,last_name,email\nW1AW,Hiram,Maxim,hiram@arrl.org\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ml, _ := LoadCSV(csvPath)

	_, ok := ml.Lookup("w1aw")
	if !ok {
		t.Error("expected case-insensitive match for w1aw")
	}
}

func TestLoadCSVEmptyFile(t *testing.T) {
	tmp := t.TempDir()
	csvPath := filepath.Join(tmp, "members.csv")
	if err := os.WriteFile(csvPath, []byte("callsign,first_name,last_name,email\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ml, err := LoadCSV(csvPath)
	if err != nil {
		t.Fatalf("LoadCSV failed: %v", err)
	}
	_, ok := ml.Lookup("W1AW")
	if ok {
		t.Errorf("expected 0 members, but found W1AW")
	}
}

func TestLookupEmptyCallsign(t *testing.T) {
	tmp := t.TempDir()
	csvPath := filepath.Join(tmp, "members.csv")
	content := "callsign,first_name,last_name,email\nW1AW,Hiram,Maxim,hiram@arrl.org\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ml, _ := LoadCSV(csvPath)

	_, ok := ml.Lookup("")
	if ok {
		t.Error("expected not found for empty callsign")
	}
}
