package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/pflag"
	bolt "go.etcd.io/bbolt"
)

func main() {
	var dbPath, csvPath, bucketName string
	pflag.StringVar(&dbPath, "db", "", "Path to BoltDB file")
	pflag.StringVar(&csvPath, "out", "", "Path to output CSV file (default: <db>.csv)")
	pflag.StringVar(&bucketName, "bucket", "Visitor", "Bucket name (default: 'Visitor')")
	pflag.Parse()

	if dbPath == "" {
		log.Fatal("Please provide the path to the BoltDB file using --db flag.")
	}

	if csvPath == "" {
		base := filepath.Base(dbPath)
		ext := filepath.Ext(base)
		csvPath = base[:len(base)-len(ext)] + ".csv"
	}

	db, err := bolt.Open(dbPath, 0444, nil)
	if err != nil {
		log.Fatalf("Failed to open BoltDB: %v", err)
	}
	defer db.Close()

	var records []map[string]interface{}
	var keys []string

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Bucket '%s' not found", bucketName)
		}
		return b.ForEach(func(k, v []byte) error {
			var m map[string]interface{}
			if err := json.Unmarshal(v, &m); err != nil {
				log.Printf("Skipping key %s: %v", string(k), err)
				return nil
			}
			m["id"] = string(k)
			records = append(records, m)
			for field := range m {
				if !contains(keys, field) {
					keys = append(keys, field)
				}
			}
			return nil
		})
	})
	if err != nil {
		log.Fatalf("Error reading bucket: %v", err)
	}

	// Sort keys for consistent CSV column order
	sort.Strings(keys)

	f, err := os.Create(csvPath)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(keys); err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	for _, rec := range records {
		row := make([]string, len(keys))
		for i, k := range keys {
			if v, ok := rec[k]; ok {
				row[i] = fmt.Sprintf("%v", v)
			} else {
				row[i] = ""
			}
		}
		if err := w.Write(row); err != nil {
			log.Printf("Failed to write row: %v", err)
		}
	}

	fmt.Printf("Exported %d records to %s\n", len(records), csvPath)
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
