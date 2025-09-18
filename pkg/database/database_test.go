package database

import (
	"database/sql"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	tempDirectory := t.TempDir()
	db := CreateDatabase(tempDirectory)
	if db == nil {
		t.Fatal("Expected database, got nil")
	}
	return db
}

func checkTableExists(t *testing.T, db *sql.DB, tableName string) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
	if err != nil {
		t.Errorf("Error checking for table %s: %v", tableName, err)
	}
	if count != 1 {
		t.Errorf("Expected table %s to exist, but it does not", tableName)
	}
}

func checkRowCount(t *testing.T, db *sql.DB, tableName string, expected int) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	if err != nil {
		t.Errorf("Error counting rows in %s: %v", tableName, err)
	}
	if count != expected {
		t.Errorf("Expected %d rows in %s, got %d", expected, tableName, count)
	}
}

func createTestJobs() []shared.JobData {
	return []shared.JobData{
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
}

func TestCreateDatabase(t *testing.T) {
	tempDirectory := t.TempDir()
	db := CreateDatabase(tempDirectory)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

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
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	SetupDatabase(db)

	tables := []string{"job_data", "qualifications", "links"}
	for _, tableName := range tables {
		checkTableExists(t, db, tableName)
	}
}

func TestWriteToDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	SetupDatabase(db)

	testJobs := createTestJobs()

	WriteToDatabase(db, testJobs)

	checkRowCount(t, db, "job_data", 2)

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
	checkRowCount(t, db, "qualifications", 2)
	checkRowCount(t, db, "links", 2)
}

func TestCreateMainTable(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	createMainTable(db)

	checkTableExists(t, db, "job_data")

	// testing basic insert functionality
	_, err = db.Exec(`INSERT INTO job_data (location, job_title, company_name, date_posted) 
		VALUES (?, ?, ?, ?)`, "Test Location", "Test Job", "Test Company", "2024-01-01")
	if err != nil {
		t.Errorf("Error inserting test row: %v", err)
	}

	// verifying data was inserted
	checkRowCount(t, db, "job_data", 1)
}

func TestCreateSecondaryTables(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	createMainTable(db)

	createSecondaryTables(db)

	tables := []string{"qualifications", "links"}
	for _, tableName := range tables {
		checkTableExists(t, db, tableName)
	}

	checkRowCount(t, db, "qualifications", 0)
	checkRowCount(t, db, "links", 0)
}
