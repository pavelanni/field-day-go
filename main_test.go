package main

import (
	"encoding/json"
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
	server, err := NewServer(dbFile, "", thisYear)
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
	server, err := NewServer(dbFile, "", thisYear)
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

	server, err := NewServer(dbFile, "", thisYear)
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
	if !strings.Contains(body, "visitors so far in") {
		t.Errorf("newVisitorHandler GET should contain visitor count display")
	}
}

func TestNewVisitorHandler_POST(t *testing.T) {
	// Create test database
	dbFile := "test_handler_post.db"
	defer os.Remove(dbFile)

	server, err := NewServer(dbFile, "", thisYear)
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

	server, err := NewServer(dbFile, "", thisYear)
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

	server, err := NewServer(dbFile, "", thisYear)
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
	server, err := NewServer("test_morse.db", "", thisYear)
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
	server, err := NewServer("test_morse2.db", "", thisYear)
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
	server, err := NewServer(dbFile, "", thisYear)
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

func TestMemberLookupHandler(t *testing.T) {
	tests := []struct {
		name        string
		callsign    string
		membersFile string
		wantBody    string
		wantJSON    map[string]string
	}{
		{
			name:        "members nil returns empty object",
			callsign:    "W1AW",
			membersFile: "",
			wantBody:    "{}",
		},
		{
			name:        "empty callsign returns empty object",
			callsign:    "",
			membersFile: "memberlookup/testdata/members.csv",
			wantBody:    "{}",
		},
		{
			name:        "unknown callsign returns empty object",
			callsign:    "ZZZZZZ",
			membersFile: "memberlookup/testdata/members.csv",
			wantBody:    "{}",
		},
		{
			name:        "known callsign returns member data",
			callsign:    "W1AW",
			membersFile: "memberlookup/testdata/members.csv",
			wantJSON: map[string]string{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "john@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbFile := "test_member_lookup.db"
			defer os.Remove(dbFile)

			server, err := NewServer(dbFile, tt.membersFile, thisYear)
			if err != nil {
				t.Fatal(err)
			}
			defer server.store.Close()

			reqURL := "/member-lookup"
			if tt.callsign != "" {
				reqURL += "?callsign=" + tt.callsign
			}
			req, err := http.NewRequest("GET", reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.memberLookupHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("memberLookupHandler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("memberLookupHandler should return application/json, got: %s", contentType)
			}

			body := strings.TrimSpace(rr.Body.String())

			if tt.wantBody != "" {
				if body != tt.wantBody {
					t.Errorf("expected body %q, got %q", tt.wantBody, body)
				}
				return
			}

			if tt.wantJSON != nil {
				var got map[string]string
				if err := json.Unmarshal([]byte(body), &got); err != nil {
					t.Fatalf("failed to parse JSON response: %v", err)
				}
				for k, v := range tt.wantJSON {
					if got[k] != v {
						t.Errorf("expected %s=%q, got %q", k, v, got[k])
					}
				}
			}
		})
	}
}
