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
		t.Fatal("Expected database, got nil")
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

	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	SetupDatabase(db)

	tables := []string{"job_data", "qualifications", "links"}
	for _, tableName := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
		if err != nil {
			t.Errorf("Error checking for table %s: %v", tableName, err)
		}
		if count != 1 {
			t.Errorf("Expected table %s to exist, but it does not", tableName)
		}
	}
}

func TestWriteToDatabase(t *testing.T) {
	tempDirectory := t.TempDir()
	shared.Program.OutputDirectory = tempDirectory
	db := CreateDatabase()
	defer db.Close()

	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	SetupDatabase(db)

	testJobs := []shared.JobData{
		{
			Location:       "Boston, MA",
			JobTitle:       "Software Engineer",
			CompanyName:    "Tech Corp",
			Description:    "Build amazing software",
			DatePosted:     "2024-01-15",
			Salary:         80000,
			WorkFromHome:   "Yes",
			Qualifications: "Go, SQL, Git",
			Links:          "https://techcorp.com/jobs",
			Country:        "USA",
		},
		{
			Location:       "San Francisco, CA",
			JobTitle:       "Data Scientist",
			CompanyName:    "Data Inc",
			Description:    "Analyze big data",
			DatePosted:     "2024-01-20",
			Salary:         95000,
			WorkFromHome:   "No",
			Qualifications: "Python, R, Machine Learning",
			Links:          "https://datainc.com/careers",
			Country:        "USA",
		},
	}

	WriteToDatabase(db, testJobs)

	var jobCount int
	err = db.QueryRow("SELECT COUNT(*) FROM job_data").Scan(&jobCount)
	if err != nil {
		t.Errorf("Error counting jobs: %v", err)
	}
	if jobCount != 2 {
		t.Errorf("Expected 2 jobs, got %d", jobCount)
	}

	// Check one specific job was inserted correctly
	var location, jobTitle string
	err = db.QueryRow("SELECT location, job_title FROM job_data WHERE company_name = ?", "Tech Corp").Scan(&location, &jobTitle)
	if err != nil {
		t.Errorf("Error querying job data: %v", err)
	}

	// verifying the main table
	if location != "Boston, MA" {
		t.Errorf("Expected location 'Boston, MA', got '%s'", location)
	}
	if jobTitle != "Software Engineer" {
		t.Errorf("Expected job title 'Software Engineer', got '%s'", jobTitle)
	}

	// verifying the related tables
	var qualCount, linkCount int
	err = db.QueryRow("SELECT COUNT(*) FROM qualifications").Scan(&qualCount)
	if err != nil {
		t.Errorf("Error counting qualifications: %v", err)
	}
	if qualCount != 2 {
		t.Errorf("Expected 2 qualifications entries, got %d", qualCount)
	}
	err = db.QueryRow("SELECT COUNT(*) FROM links").Scan(&linkCount)
	if err != nil {
		t.Errorf("Error counting links: %v", err)
	}
	if linkCount != 2 {
		t.Errorf("Expected 2 links entries, got %d", linkCount)
	}
}

func TestCreateMainTable(t *testing.T) {
	tempDirectory := t.TempDir()
	shared.Program.OutputDirectory = tempDirectory
	db := CreateDatabase()
	defer db.Close()

	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	createMainTable(db)

	// checking table exists
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='job_data'").Scan(&tableCount)
	if err != nil {
		t.Errorf("Error checking for job_data table: %v", err)
	}
	if tableCount != 1 {
		t.Errorf("Expected job_data table to exist, but it does not")
	}

	// testing basic insert functionality
	_, err = db.Exec(`INSERT INTO job_data (location, job_title, company_name, date_posted) 
		VALUES (?, ?, ?, ?)`, "Test Location", "Test Job", "Test Company", "2024-01-01")
	if err != nil {
		t.Errorf("Error inserting test row: %v", err)
	}

	// verifying data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM job_data").Scan(&count)
	if err != nil {
		t.Errorf("Error counting rows: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

func TestCreateSecondaryTables(t *testing.T) {
	tempDirectory := t.TempDir()
	shared.Program.OutputDirectory = tempDirectory
	db := CreateDatabase()
	defer db.Close()

	createMainTable(db)

	createSecondaryTables(db)

	tables := []string{"qualifications", "links"}
	for _, tableName := range tables {
		var tableCount int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&tableCount)
		if err != nil {
			t.Errorf("Error checking for %s table: %v", tableName, err)
		}
		if tableCount != 1 {
			t.Errorf("Expected %s table to exist, but it does not", tableName)
		}
	}

	var qualCount int
	err := db.QueryRow("SELECT COUNT(*) FROM qualifications").Scan(&qualCount)
	if err != nil {
		t.Errorf("Error checking qualifications table: %v", err)
	}

	var linkCount int
	err = db.QueryRow("SELECT COUNT(*) FROM links").Scan(&linkCount)
	if err != nil {
		t.Errorf("Error checking links table: %v", err)
	}
}
