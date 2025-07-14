package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pavelanni/field-day-go/visitorstore"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run cmd/importcsv/main.go <csv_file> <db_file>")
	}
	
	csvFile := os.Args[1]
	dbFile := os.Args[2]
	
	// Create new visitor store
	store, err := visitorstore.NewVisitorStore(dbFile)
	if err != nil {
		log.Fatalf("Failed to create visitor store: %v", err)
	}
	
	// Import CSV data
	fmt.Printf("Importing data from %s to %s...\n", csvFile, dbFile)
	err = store.ImportFromCSV(csvFile)
	if err != nil {
		log.Fatalf("Failed to import CSV: %v", err)
	}
	
	// Show total visitors
	total, err := store.TotalVisitors()
	if err != nil {
		log.Fatalf("Failed to get total visitors: %v", err)
	}
	
	fmt.Printf("Successfully imported %d visitors!\n", total)
}