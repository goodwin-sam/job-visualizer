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
