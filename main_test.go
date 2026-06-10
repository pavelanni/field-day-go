package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Create test server
	dbFile := "test_home.db"
	defer os.Remove(dbFile)
	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.homeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("homeHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Field Day 2026") {
		t.Errorf("homeHandler should contain 'Field Day 2026', got: %s", body)
	}
}

func TestConfirmHandler(t *testing.T) {
	dbFile := "test_confirm.db"
	defer os.Remove(dbFile)
	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	req, err := http.NewRequest("GET", "/confirmation?callsign=W1AW&name=Test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.confirmHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("confirmHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "W1AW") {
		t.Errorf("confirmHandler should contain callsign 'W1AW', got: %s", body)
	}
	if !strings.Contains(body, "Test") {
		t.Errorf("confirmHandler should contain name 'Test', got: %s", body)
	}
}

func TestNewVisitorHandler_GET(t *testing.T) {
	// Create test database
	dbFile := "test_handler.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	req, err := http.NewRequest("GET", "/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("newVisitorHandler GET returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Field Day 2026") {
		t.Errorf("newVisitorHandler GET should contain 'Field Day 2026', got: %s", body)
	}
}

func TestNewVisitorHandler_POST(t *testing.T) {
	// Create test database
	dbFile := "test_handler_post.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	// Prepare form data
	form := url.Values{}
	form.Add("firstname", "Test")
	form.Add("lastname", "User")
	form.Add("callsign", "W1AW")
	form.Add("email", "test@example.com")

	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("newVisitorHandler POST returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	// Check redirect location
	location := rr.Header().Get("Location")
	if !strings.Contains(location, "/confirmation") {
		t.Errorf("newVisitorHandler POST should redirect to confirmation page, got: %s", location)
	}

	// Verify visitor was saved
	total, err := server.store.TotalVisitors()
	if err != nil {
		t.Errorf("Failed to get total visitors: %v", err)
	}
	if total != 1 {
		t.Errorf("Expected 1 visitor after POST, got %d", total)
	}
}

func TestNewVisitorHandler_POST_MissingFirstName(t *testing.T) {
	// Create test database
	dbFile := "test_handler_missing.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	// Prepare form data without firstname
	form := url.Values{}
	form.Add("lastname", "User")
	form.Add("callsign", "W1AW")

	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("newVisitorHandler POST with missing firstname should return 400, got %v", status)
	}

	// Should display error message
	body := rr.Body.String()
	if !strings.Contains(body, "First name") {
		t.Errorf("Response should contain 'First name' error, got: %s", body)
	}
}

func TestListHandler(t *testing.T) {
	// Create test database
	dbFile := "test_list_handler.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	req, err := http.NewRequest("GET", "/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.listHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("listHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "List of Field Day visitors") {
		t.Errorf("listHandler should contain visitor list title, got: %s", body)
	}
}

func TestMorseAudioHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/morse-audio?callsign=W1AW", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// morse handler doesn't need DB but uses server method now
	server, err := NewServer("test_morse.db", thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()
	handler := http.HandlerFunc(server.morseAudioHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("morseAudioHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "audio/wav" {
		t.Errorf("morseAudioHandler should return audio/wav, got: %s", contentType)
	}

	// Check that we got some audio data
	if rr.Body.Len() == 0 {
		t.Error("morseAudioHandler should return audio data")
	}
}

func TestMorseAudioHandler_MissingCallsign(t *testing.T) {
	req, err := http.NewRequest("GET", "/morse-audio", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server, err := NewServer("test_morse2.db", thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()
	handler := http.HandlerFunc(server.morseAudioHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("morseAudioHandler without callsign should return 400, got %v", status)
	}
}

func TestPrivacyHandler(t *testing.T) {
	dbFile := "test_privacy.db"
	defer os.Remove(dbFile)
	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	req, err := http.NewRequest("GET", "/privacy", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.privacyHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("privacyHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Field Day 2026") {
		t.Errorf("privacyHandler should contain 'Field Day 2026', got: %s", body)
	}
}

func TestNewVisitorHandler_InvalidCallsign(t *testing.T) {
	dbFile := "test_invalid_callsign.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	// Test with invalid callsign
	form := url.Values{}
	form.Add("firstname", "Test")
	form.Add("lastname", "User")
	form.Add("callsign", "123") // Invalid: no letters
	form.Add("email", "test@example.com")

	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	// Should return HTTP 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("newVisitorHandler with invalid callsign returned status %v, want %v", status, http.StatusBadRequest)
	}

	// Should display error message
	body := rr.Body.String()
	if !strings.Contains(body, "Callsign") {
		t.Errorf("Response should contain 'Callsign' error, got: %s", body)
	}

	// Verify visitor was NOT saved
	total, err := server.store.TotalVisitors()
	if err != nil {
		t.Errorf("Failed to get total visitors: %v", err)
	}
	if total != 0 {
		t.Errorf("Expected 0 visitors after invalid submission, got %d", total)
	}
}

func TestNewVisitorHandler_InvalidEmail(t *testing.T) {
	dbFile := "test_invalid_email.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	// Test with invalid email
	form := url.Values{}
	form.Add("firstname", "Test")
	form.Add("lastname", "User")
	form.Add("callsign", "W1AW")
	form.Add("email", "notanemail") // Invalid: missing @ and domain

	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	// Should return HTTP 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("newVisitorHandler with invalid email returned status %v, want %v", status, http.StatusBadRequest)
	}

	// Should display error message
	body := rr.Body.String()
	if !strings.Contains(body, "Email") {
		t.Errorf("Response should contain 'Email' error, got: %s", body)
	}

	// Verify visitor was NOT saved
	total, err := server.store.TotalVisitors()
	if err != nil {
		t.Errorf("Failed to get total visitors: %v", err)
	}
	if total != 0 {
		t.Errorf("Expected 0 visitors after invalid submission, got %d", total)
	}
}

func TestNewVisitorHandler_ValidInputs(t *testing.T) {
	dbFile := "test_valid_inputs.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	// Test with valid inputs that need normalization
	form := url.Values{}
	form.Add("firstname", "  Test  ") // Whitespace to trim
	form.Add("lastname", "User")
	form.Add("callsign", "w1aw")          // Lowercase to uppercase
	form.Add("email", "Test@EXAMPLE.COM") // Mixed case to lowercase

	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	// Should redirect successfully
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("newVisitorHandler with valid inputs returned status %v, want %v", status, http.StatusSeeOther)
	}

	// Verify visitor was saved with normalized values
	visitors, err := server.store.ListVisitors()
	if err != nil {
		t.Fatalf("Failed to list visitors: %v", err)
	}
	if len(visitors) != 1 {
		t.Fatalf("Expected 1 visitor, got %d", len(visitors))
	}

	v := visitors[0]
	if v.FirstName != "Test" {
		t.Errorf("FirstName = %q, want %q (trimmed)", v.FirstName, "Test")
	}
	if v.Callsign != "W1AW" {
		t.Errorf("Callsign = %q, want %q (uppercase)", v.Callsign, "W1AW")
	}
	if v.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q (lowercase)", v.Email, "test@example.com")
	}
}

func TestNewVisitorHandler_EmptyOptionalFields(t *testing.T) {
	dbFile := "test_empty_optional.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, thisYear)
	if err != nil {
		t.Fatal(err)
	}
	defer server.store.Close()

	// Test with empty optional fields
	form := url.Values{}
	form.Add("firstname", "Test")
	form.Add("lastname", "User")
	form.Add("callsign", "") // Empty optional field
	form.Add("email", "")    // Empty optional field

	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.newVisitorHandler)
	handler.ServeHTTP(rr, req)

	// Should redirect successfully (empty optional fields are valid)
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("newVisitorHandler with empty optionals returned status %v, want %v", status, http.StatusSeeOther)
	}

	// Verify visitor was saved
	total, err := server.store.TotalVisitors()
	if err != nil {
		t.Errorf("Failed to get total visitors: %v", err)
	}
	if total != 1 {
		t.Errorf("Expected 1 visitor, got %d", total)
	}
}
