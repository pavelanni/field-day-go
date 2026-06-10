package visitorstore

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var (
	// callsignRegex matches international amateur radio callsigns:
	// 1-3 alphanumeric prefix + 1 digit + 1-4 letter suffix + optional /portable indicator
	callsignRegex = regexp.MustCompile(`^[A-Z0-9]{1,3}[0-9][A-Z]{1,4}(/[A-Z0-9]+)?$`)
	// emailRegex matches simplified RFC 5322 email format
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Visitor contains information about Field Day visitors
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

// ValidateCallsign validates the callsign field format
func (v *Visitor) ValidateCallsign() error {
	// Normalize: trim and uppercase
	v.Callsign = strings.ToUpper(strings.TrimSpace(v.Callsign))

	// Empty callsign is valid (optional field)
	if v.Callsign == "" {
		return nil
	}

	// Validate format
	if !callsignRegex.MatchString(v.Callsign) {
		return &ValidationError{
			Field:   "Callsign",
			Message: "Callsign must be in amateur radio format (e.g., W1AW, KB6NU/M)",
		}
	}

	return nil
}

// ValidateEmail validates the email field format
func (v *Visitor) ValidateEmail() error {
	// Normalize: trim and lowercase
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))

	// Empty email is valid (optional field)
	if v.Email == "" {
		return nil
	}

	// Validate format
	if !emailRegex.MatchString(v.Email) {
		return &ValidationError{
			Field:   "Email",
			Message: "Email address must be in valid format (e.g., user@example.com)",
		}
	}

	return nil
}

// Validate validates all Visitor fields, returns first error encountered
func (v *Visitor) Validate() error {
	// Validate required field
	v.FirstName = strings.TrimSpace(v.FirstName)
	if v.FirstName == "" {
		return &ValidationError{
			Field:   "FirstName",
			Message: "First name is required",
		}
	}

	// Validate callsign
	if err := v.ValidateCallsign(); err != nil {
		return err
	}

	// Validate email
	if err := v.ValidateEmail(); err != nil {
		return err
	}

	return nil
}

type VisitorStore struct {
	db *sql.DB
}

func NewVisitorStore(dbFile string) (*VisitorStore, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	vs := &VisitorStore{db}

	// Initialize the database schema
	if err := vs.initSchema(); err != nil {
		return nil, err
	}

	return vs, nil
}

func (vs *VisitorStore) SaveVisitor(v Visitor) error {
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now()
	}

	query := `INSERT INTO visitors (created_at, first_name, last_name, callsign, email, nfarl, contactme, youth, firsttime)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := vs.db.Exec(query, v.CreatedAt.Format(time.RFC3339), v.FirstName, v.LastName, v.Callsign, v.Email, v.Nfarl, v.Contactme, v.Youth, v.Firsttime)
	return err
}
func (vs *VisitorStore) ListVisitors() ([]Visitor, error) {
	query := `SELECT id, created_at, first_name, last_name, callsign, email, nfarl, contactme, youth, firsttime
			  FROM visitors ORDER BY id`

	rows, err := vs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visitors []Visitor
	for rows.Next() {
		var v Visitor
		var createdAtStr string
		err := rows.Scan(&v.ID, &createdAtStr, &v.FirstName, &v.LastName, &v.Callsign, &v.Email, &v.Nfarl, &v.Contactme, &v.Youth, &v.Firsttime)
		if err != nil {
			return nil, err
		}

		// Parse the timestamp string back to time.Time
		v.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}

		visitors = append(visitors, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visitors, nil
}

func (vs *VisitorStore) TotalVisitors() (int, error) {
	var count int
	err := vs.db.QueryRow("SELECT COUNT(*) FROM visitors").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// HealthCheck verifies DB connectivity.
func (vs *VisitorStore) HealthCheck() error {
	return vs.db.Ping()
}

// initSchema creates the visitors table if it doesn't exist
func (vs *VisitorStore) initSchema() error {
	query := `CREATE TABLE IF NOT EXISTS visitors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT,
		callsign TEXT,
		email TEXT,
		nfarl BOOLEAN DEFAULT FALSE,
		contactme BOOLEAN DEFAULT FALSE,
		youth BOOLEAN DEFAULT FALSE,
		firsttime BOOLEAN DEFAULT FALSE
	)`

	_, err := vs.db.Exec(query)
	return err
}

// ImportFromCSV imports visitor data from a CSV file
func (vs *VisitorStore) ImportFromCSV(csvFile string) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return err
	}

	// Create a map for column indices
	colMap := make(map[string]int)
	for i, col := range header {
		colMap[col] = i
	}

	// Read and import data
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		v := Visitor{}

		// Parse CreatedAt
		if createdAtStr := record[colMap["CreatedAt"]]; createdAtStr != "" {
			if parsedTime, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
				v.CreatedAt = parsedTime
			} else {
				v.CreatedAt = time.Now()
			}
		} else {
			v.CreatedAt = time.Now()
		}

		v.FirstName = record[colMap["FirstName"]]
		v.LastName = record[colMap["LastName"]]
		v.Callsign = record[colMap["Callsign"]]
		v.Email = record[colMap["Email"]]

		// Parse boolean fields
		v.Nfarl, _ = strconv.ParseBool(record[colMap["Nfarl"]])
		v.Contactme, _ = strconv.ParseBool(record[colMap["Contactme"]])
		v.Youth, _ = strconv.ParseBool(record[colMap["Youth"]])
		v.Firsttime, _ = strconv.ParseBool(record[colMap["Firsttime"]])

		// Skip empty first names
		if strings.TrimSpace(v.FirstName) == "" {
			continue
		}

		if err := vs.SaveVisitor(v); err != nil {
			return fmt.Errorf("failed to save visitor %s: %w", v.FirstName, err)
		}
	}

	return nil
}

// Close closes the database connection
func (vs *VisitorStore) Close() error {
	return vs.db.Close()
}
