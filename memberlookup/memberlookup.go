package memberlookup

import (
	"encoding/csv"
	"errors"
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
			return nil, errors.New("missing required column: " + col)
		}
	}

	maxIdx := 0
	for _, idx := range colIdx {
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	members := make(map[string]Member)
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		if len(record) <= maxIdx {
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

// Len returns the number of members loaded.
func (l *Lookup) Len() int {
	return len(l.members)
}
