package visitorstore

import (
	"fmt"

	"github.com/asdine/storm"
)

// Visitor contains information about Field Day visitors
type Visitor struct {
	ID        int    `storm:"id,increment"`
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
	Callsign  string `schema:"callsign"`
	Email     string `schema:"email"`
	Nfarl     bool   `schema:"nfarl"`
	Contactme bool   `schema:"contactme"`
	Youth     bool   `schema:"youth"`
	Firsttime bool   `schema:"firsttime"`
}

type VisitorStore struct {
	db *storm.DB
}

func NewVisitorStore(dbFile string) (*VisitorStore, error) {
	db, err := storm.Open(dbFile)
	if err != nil {
		return nil, err
	}
	err = db.Init(&Visitor{})
	if err != nil {
		return nil, err
	}
	return &VisitorStore{db}, nil
}

func (vs *VisitorStore) SaveVisitor(v Visitor) error {
	if v.FirstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}
	return vs.db.Save(&v)
}
func (vs *VisitorStore) ListVisitors() ([]Visitor, error) {
	var visitors []Visitor
	err := vs.db.AllByIndex("ID", &visitors)
	if err != nil {
		return nil, err
	}
	return visitors, nil
}
