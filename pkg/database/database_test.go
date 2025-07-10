package database

import (
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateDatabase(t *testing.T) {
	tempDirectory := t.TempDir()
	shared.Program.OutputDirectory = tempDirectory
	db := CreateDatabase()
	if db == nil {
		t.Errorf("Expected database, got nil")
	}
	defer db.Close()

	dbPath := filepath.Join(tempDirectory, "job_data.sqlite")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Expected database file to exist at %s, but it does not", dbPath)
	}

	err := db.Ping()
	if err != nil {
		t.Errorf("Expected database to be reachable, but got error: %v", err)
	}
}

func TestSetupDatabase(t *testing.T) {
	tempDirectory := t.TempDir()
	shared.Program.OutputDirectory = tempDirectory
	db := CreateDatabase()
	defer db.Close()

	SetupDatabase(db)

	tables := []string{"job_data", "qualifications", "links"}
	for _, tableName := range tables {
		var count int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		err := db.QueryRow(query, tableName).Scan(&count)
		if err != nil {
			t.Errorf("Error checking for table %s: %v", tableName, err)
		}
		if count == 0 {
			t.Errorf("Expected table %s to exist, but it does not", tableName)
		}
	}

	// Get the actual table structure from the database
	rows, err := db.Query("PRAGMA table_info(job_data)")
	if err != nil {
		t.Errorf("Error getting table info for job_data: %v", err)
	}
	defer rows.Close()

	// Define what columns we expect to find in the job_data table
	expectedColumns := map[string]string{
		"id":             "INTEGER",
		"location":       "TEXT",
		"job_title":      "TEXT",
		"company_name":   "TEXT",
		"description":    "TEXT",
		"date_posted":    "TEXT",
		"salary":         "INT",
		"work_from_home": "TEXT",
		"qualifications": "TEXT",
		"links":          "TEXT",
		"country":        "TEXT",
	}

	foundColumns := make(map[string]string)
	for rows.Next() {
		var columnIndex int
		var name, typ string
		var notnull int
		var defaultValue interface{}
		var primaryKey int

		err := rows.Scan(&columnIndex, &name, &typ, &notnull, &defaultValue, &primaryKey)
		if err != nil {
			t.Errorf("Error scanning table info: %v", err)
		}
		foundColumns[name] = typ
	}

	for expectedName, expectedType := range expectedColumns {
		if foundType, exists := foundColumns[expectedName]; !exists {
			t.Errorf("Expected column %s to exist in job_data table, but it does not", expectedName)
		} else if foundType != expectedType {
			t.Errorf("Expected column %s to be of type %s, but got %s", expectedName, expectedType, foundType)
		}
	}
}
