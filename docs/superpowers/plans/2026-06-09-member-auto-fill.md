# Member Auto-Fill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** When a visitor enters their callsign on the registration form, auto-fill their name, email, and NFARL membership status from a pre-loaded member list.

**Architecture:** A new `memberlookup` package loads a CSV of club members at server startup into an in-memory map. A new HTTP endpoint `/member-lookup?callsign=XYZ` returns JSON with the member's details. A small inline JS snippet on the registration form calls this endpoint on callsign blur and populates the form fields.

**Tech Stack:** Go (stdlib `encoding/csv`, `encoding/json`, `net/http`), plain JavaScript (fetch API), Go templates.

---

## File Structure

| File | Responsibility |
|------|---------------|
| `memberlookup/memberlookup.go` | New package: CSV loading, in-memory map, lookup by callsign |
| `memberlookup/memberlookup_test.go` | Tests for CSV parsing and lookup |
| `main.go` | Add `--members` flag, wire up `memberlookup`, add `/member-lookup` endpoint |
| `templates/new.go.html` | Move callsign field to top, add `id` attributes, add inline JS for auto-fill |
| `static/js/member-lookup.js` | Small JS module for the fetch + form-fill logic |

---

### Task 1: Create `memberlookup` package with CSV loading and lookup

**Files:**
- Create: `memberlookup/memberlookup.go`
- Test: `memberlookup/memberlookup_test.go`

- [ ] **Step 1: Write tests for the memberlookup package**

```go
package memberlookup

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoadMembers(t *testing.T) {
    // Create temp CSV
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

    if len(ml.members) != 2 {
        t.Errorf("expected 2 members, got %d", len(ml.members))
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
    if len(ml.members) != 0 {
        t.Errorf("expected 0 members, got %d", len(ml.members))
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./memberlookup/ -v`
Expected: FAIL (package does not exist yet)

- [ ] **Step 3: Implement the memberlookup package**

```go
package memberlookup

import (
    "encoding/csv"
    "os"
    "strings"
)

// Member holds information about a club member.
type Member struct {
    Callsign  string
    FirstName string
    LastName  string
    Email     string
}

// Lookup holds an in-memory map of callsign -> Member.
type Lookup struct {
    members map[string]Member
}

// LoadCSV reads a CSV file and returns a Lookup populated with its members.
// Expected CSV columns: callsign,first_name,last_name,email
func LoadCSV(path string) (*Lookup, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    r := csv.NewReader(f)
    r.TrimLeadingSpace = true

    // Read and validate header
    header, err := r.Read()
    if err != nil {
        return nil, err
    }
    colIdx := make(map[string]int)
    for i, h := range header {
        colIdx[strings.TrimSpace(strings.ToLower(h))] = i
    }

    required := []string{"callsign", "first_name", "last_name", "email"}
    for _, col := range required {
        if _, ok := colIdx[col]; !ok {
            return nil, &csv.ParseError{Line: 1, Column: 0, Err: err}
        }
    }

    members := make(map[string]Member)
    for {
        record, err := r.Read()
        if err != nil {
            break // EOF or parse error
        }
        if len(record) < len(required) {
            continue
        }

        callsign := strings.TrimSpace(record[colIdx["callsign"]])
        if callsign == "" {
            continue
        }

        members[strings.ToUpper(callsign)] = Member{
            Callsign:  strings.ToUpper(callsign),
            FirstName: strings.TrimSpace(record[colIdx["first_name"]]),
            LastName:  strings.TrimSpace(record[colIdx["last_name"]]),
            Email:     strings.TrimSpace(record[colIdx["email"]]),
        }
    }

    return &Lookup{members: members}, nil
}

// Lookup returns the Member for the given callsign, or false if not found.
// Callsign matching is case-insensitive.
func (l *Lookup) Lookup(callsign string) (Member, bool) {
    if callsign == "" {
        return Member{}, false
    }
    m, ok := l.members[strings.ToUpper(callsign)]
    return m, ok
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./memberlookup/ -v`
Expected: All 6 tests PASS

- [ ] **Step 5: Commit**

```bash
git add memberlookup/
git commit -s -m "feat: add memberlookup package for CSV-based member lookup"
```

---

### Task 2: Wire up member lookup in main.go

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Add `--members` flag and wire up lookup**

Changes to `main.go`:

1. Add `memberlookup` import:
```go
import (
    // ... existing imports ...
    "github.com/pavelanni/field-day-go/memberlookup"
)
```

2. Add `members` field to `Server` struct:
```go
type Server struct {
    store     *visitorstore.VisitorStore
    members   *memberlookup.Lookup
    templates map[string]*template.Template
    year      string
}
```

3. Update `NewServer` to accept members path:
```go
func NewServer(dbFile, membersFile, year string) (*Server, error) {
    store, err := visitorstore.NewVisitorStore(dbFile)
    if err != nil {
        return nil, err
    }

    var members *memberlookup.Lookup
    if membersFile != "" {
        members, err = memberlookup.LoadCSV(membersFile)
        if err != nil {
            return nil, fmt.Errorf("loading members CSV: %w", err)
        }
        log.Printf("Loaded %d club members from %s", len(members.members), membersFile)
    } else {
        log.Println("No members CSV specified; member auto-fill disabled")
    }

    s := &Server{store: store, members: members, year: year}
    if err := s.initTemplates(); err != nil {
        return nil, err
    }
    return s, nil
}
```

4. Add `fmt` to imports.

5. Add `/member-lookup` endpoint in `run()`:
```go
mux.HandleFunc("/member-lookup", s.memberLookupHandler)
```

6. Add the handler:
```go
func (s *Server) memberLookupHandler(w http.ResponseWriter, r *http.Request) {
    if s.members == nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{}`))
        return
    }

    callsign := r.URL.Query().Get("callsign")
    if callsign == "" {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{}`))
        return
    }

    member, ok := s.members.Lookup(callsign)
    if !ok {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{}`))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "first_name": member.FirstName,
        "last_name":  member.LastName,
        "email":      member.Email,
    })
}
```

7. Add `"encoding/json"` to imports.

8. Update `main()` to add the `--members` flag:
```go
var flagMembers string
flagMembers = os.Getenv("FD_MEMBERS")
// ...
pflag.StringVar(&flagMembers, "members", flagMembers, "Path to club members CSV file")
// ...
server, err := NewServer(flagDB, flagMembers, flagYear)
```

- [ ] **Step 2: Verify build compiles**

Run: `go build -o fieldday .`
Expected: Success, no errors

- [ ] **Step 3: Commit**

```bash
git add main.go
git commit -s -m "feat: add --members flag and /member-lookup endpoint"
```

---

### Task 3: Update registration form template

**Files:**
- Modify: `templates/new.go.html`

- [ ] **Step 1: Move callsign field to top and add IDs for JS targeting**

Reorder the form fields so callsign is first. Add `id` attributes to all form fields (they already have them, but verify). The form should look like:

```html
{{define "signupForm"}}
<form action="new" method="POST" class="space-y-2 bg-gray-50 p-4 rounded shadow text-base" id="signup-form">
    <div>
        <label for="callsign" class="block font-semibold mb-1">Call sign (if licensed)</label>
        <input type="text" name="callsign" id="callsign" placeholder="Call sign" class="w-full border rounded px-3 py-1.5 text-base focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500" autocomplete="off">
    </div>
    <div>
        <label for="firstname" class="block font-semibold mb-1">First name (required)</label>
        <input type="text" name="firstname" id="firstname" placeholder="First name (required)" required class="w-full border rounded px-3 py-1.5 text-base focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
    </div>
    <div>
        <label for="lastname" class="block font-semibold mb-1">Last name (optional)</label>
        <input type="text" name="lastname" id="lastname" placeholder="Last name" class="w-full border rounded px-3 py-1.5 text-base focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
    </div>
    <div>
        <label for="email" class="block font-semibold mb-1">Email address (if you want to be contacted)</label>
        <input type="email" name="email" id="email" placeholder="Email" class="w-full border rounded px-3 py-1.5 text-base focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
    </div>
    <div class="space-y-1 mt-1">
        <div class="flex items-center gap-2">
            <input type="checkbox" name="nfarl" id="nfarl" class="h-5 w-5 mr-2 align-middle">
            <label for="nfarl" class="text-gray-800 cursor-pointer select-none align-middle">NFARL member?</label>
        </div>
        <div class="flex items-center gap-2">
            <input type="checkbox" name="contactme" id="contactme" class="h-5 w-5 mr-2 align-middle">
            <label for="contactme" class="text-gray-800 cursor-pointer select-none align-middle">May we contact you?</label>
        </div>
        <div class="flex items-center gap-2">
            <input type="checkbox" name="youth" id="youth" class="h-5 w-5 mr-2 align-middle">
            <label for="youth" class="text-gray-800 cursor-pointer select-none align-middle">Younger than 18?</label>
        </div>
        <div class="flex items-center gap-2">
            <input type="checkbox" name="firsttime" id="firsttime" class="h-5 w-5 mr-2 align-middle">
            <label for="firsttime" class="text-gray-800 cursor-pointer select-none align-middle">First time at Field Day?</label>
        </div>
    </div>
    <div>
        <p class="text-sm text-gray-600 mt-1">By clicking the Submit button, you agree to the <a href="/privacy" class="underline text-blue-600 hover:text-blue-800">privacy policy</a></p>
    </div>
    <button type="submit" class="w-full bg-blue-600 text-white font-bold py-1.5 px-3 rounded text-base hover:bg-blue-700 transition focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2">
        Submit
    </button>
</form>
{{end}}
```

- [ ] **Step 2: Verify build and templates compile**

Run: `go build -o fieldday .`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add templates/new.go.html
git commit -s -m "feat: move callsign field to top of registration form"
```

---

### Task 4: Add JavaScript for auto-fill

**Files:**
- Create: `static/js/member-lookup.js`

- [ ] **Step 1: Create the JS module**

```javascript
(function () {
    const callsignInput = document.getElementById("callsign");
    if (!callsignInput) return;

    let debounceTimer;

    callsignInput.addEventListener("blur", function () {
        const callsign = callsignInput.value.trim();
        if (!callsign) return;

        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(function () {
            fetch("/member-lookup?callsign=" + encodeURIComponent(callsign))
                .then(function (resp) { return resp.json(); })
                .then(function (data) {
                    if (!data || !data.first_name) return;

                    document.getElementById("firstname").value = data.first_name;
                    document.getElementById("lastname").value = data.last_name || "";
                    document.getElementById("email").value = data.email || "";
                    document.getElementById("nfarl").checked = true;
                })
                .catch(function () {});
        }, 300);
    });
})();
```

- [ ] **Step 2: Include the script in the `new.go.html` template**

Add this line at the bottom of the `signupForm` template definition, after the closing `</form>` tag but before `{{end}}`:

```html
<script src="/static/js/member-lookup.js" defer></script>
```

- [ ] **Step 3: Verify build**

Run: `go build -o fieldday .`
Expected: Success

- [ ] **Step 4: Commit**

```bash
git add static/js/member-lookup.js templates/new.go.html
git commit -s -m "feat: add JS auto-fill for member lookup on callsign blur"
```

---

### Task 5: Integration test and manual verification

- [ ] **Step 1: Create a test members CSV**

```csv
callsign,first_name,last_name,email
W1AW,Hiram,Maxim,hiram@arrl.org
K1ABC,John,Doe,john@example.com
```

- [ ] **Step 2: Run the server with the members file**

```bash
./fieldday --db test.db --members test_members.csv
```

Expected output: `Loaded 2 club members from test_members.csv`

- [ ] **Step 3: Test the endpoint directly**

```bash
curl "http://localhost:3000/member-lookup?callsign=W1AW"
```

Expected: `{"email":"hiram@arrl.org","first_name":"Hiram","last_name":"Maxim"}`

```bash
curl "http://localhost:3000/member-lookup?callsign=XYZ999"
```

Expected: `{}`

- [ ] **Step 4: Run all tests**

```bash
go test ./...
```

Expected: All tests pass

- [ ] **Step 5: Clean up test files**

```bash
rm -f test.db test_members.csv
```

- [ ] **Step 6: Commit (if any changes needed)**

---

### Task 6: Update issue #16 checklist

- [ ] **Step 1: Update the GitHub issue checklist**

Comment on issue #16 with the completed checklist:

- [x] Callsign field is the first field on the registration form
- [x] A member list can be loaded (CSV format)
- [x] Entering a known callsign auto-fills the remaining fields
- [x] Unknown callsigns leave all fields blank for manual entry
- [x] Works offline (no external API calls)
