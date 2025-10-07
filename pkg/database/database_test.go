// package database provides tests for SQLite database operations
package database

import (
	"database/sql"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a test database in a temporary directory
func setupTestDB(t *testing.T) *sql.DB {
	tempDirectory := t.TempDir()
	db := CreateDatabase(tempDirectory)
	if db == nil {
		t.Fatal("Expected database, got nil")
	}
	return db
}

// checkTableExists verifies that a table exists in the database
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

// checkRowCount verifies that a table contains the expected number of rows
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

// createTestJobs creates sample job data for testing database operations
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

// TestCreateDatabase tests database file creation and connection
func TestCreateDatabase(t *testing.T) {
	tempDirectory := t.TempDir()
	db := CreateDatabase(tempDirectory)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	// verifies database file was created
	dbPath := filepath.Join(tempDirectory, "job_data.sqlite")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Expected database file to exist at %s, but it does not", dbPath)
	}

	// verifies database connection is working
	err := db.Ping()
	if err != nil {
		t.Errorf("Expected database to be reachable, but got error: %v", err)
	}
}

// TestSetupDatabase tests table creation for all required tables
func TestSetupDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	// cleans up any existing tables
	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	// tests table creation
	SetupDatabase(db)

	// verifies all required tables were created
	tables := []string{"job_data", "qualifications", "links"}
	for _, tableName := range tables {
		checkTableExists(t, db, tableName)
	}
}

// TestWriteToDatabase tests writing job data to all database tables
func TestWriteToDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	// cleans up any existing tables
	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	// sets up database tables
	SetupDatabase(db)

	// creates test job data
	testJobs := createTestJobs()

	// tests writing job data to database
	WriteToDatabase(db, testJobs)

	// verifies correct number of jobs in main table
	checkRowCount(t, db, "job_data", 2)

	// verifies specific job data was inserted correctly
	var location, jobTitle string
	err = db.QueryRow("SELECT location, job_title FROM job_data WHERE company_name = ?", "Tech Corp").Scan(&location, &jobTitle)
	if err != nil {
		t.Errorf("Error querying job data: %v", err)
	}

	// verifies main table data integrity
	if location != "Boston, MA" {
		t.Errorf("Expected location 'Boston, MA', got '%s'", location)
	}
	if jobTitle != "Software Engineer" {
		t.Errorf("Expected job title 'Software Engineer', got '%s'", jobTitle)
	}

	// verifies related tables contain correct data
	checkRowCount(t, db, "qualifications", 2)
	checkRowCount(t, db, "links", 2)
}

// TestCreateMainTable tests creation of the main job_data table
func TestCreateMainTable(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	// cleans up any existing table
	_, err := db.Exec("DROP TABLE IF EXISTS job_data")
	if err != nil {
		t.Errorf("Error dropping existing table: %v", err)
	}

	// tests main table creation
	createMainTable(db)

	// verifies table was created
	checkTableExists(t, db, "job_data")

	// tests basic insert functionality
	_, err = db.Exec(`INSERT INTO job_data (location, job_title, company_name, date_posted) 
		VALUES (?, ?, ?, ?)`, "Test Location", "Test Job", "Test Company", "2024-01-01")
	if err != nil {
		t.Errorf("Error inserting test row: %v", err)
	}

	// verifies data was inserted successfully
	checkRowCount(t, db, "job_data", 1)
}

// TestCreateSecondaryTables tests creation of qualifications and links tables
func TestCreateSecondaryTables(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		shared.CheckErrorWarn(db.Close())
	}()

	// creates main table first (required for foreign keys)
	createMainTable(db)

	// tests secondary table creation
	createSecondaryTables(db)

	// verifies both secondary tables were created
	tables := []string{"qualifications", "links"}
	for _, tableName := range tables {
		checkTableExists(t, db, tableName)
	}

	// verifies tables start empty
	checkRowCount(t, db, "qualifications", 0)
	checkRowCount(t, db, "links", 0)
}
