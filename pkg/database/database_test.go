package database

import (
	"database/sql"
	"job-visualizer/pkg/shared"
	"testing"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func TestCreateMainTable(t *testing.T) {
	// Create temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test creating main table
	createMainTable(db)

	// Verify table exists and has correct structure
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='job_data'").Scan(&tableName)
	if err != nil {
		t.Fatalf("job_data table was not created: %v", err)
	}

	if tableName != "job_data" {
		t.Fatalf("Expected table name 'job_data', got %s", tableName)
	}

	// Test table structure by attempting to insert a row
	_, err = db.Exec(`INSERT INTO job_data (location, job_title, company_name, description, date_posted, salary, work_from_home, qualifications, links, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Boston, MA", "Software Engineer", "Tech Corp", "Build apps", "2024-01-01", 100000, "No", "Go, Python", "http://example.com", "USA")
	if err != nil {
		t.Errorf("Failed to insert test data into job_data table: %v", err)
	}

	// Verify the insert worked
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM job_data").Scan(&count)
	if err != nil {
		t.Errorf("Failed to count rows: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row in job_data, got %d", count)
	}
}

func TestCreateSecondaryTables(t *testing.T) {
	// Create temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Need main table first for foreign key references
	createMainTable(db)
	createSecondaryTables(db)

	// Verify qualifications table exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='qualifications'").Scan(&tableName)
	if err != nil {
		t.Errorf("qualifications table was not created: %v", err)
	}

	// Verify links table exists
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='links'").Scan(&tableName)
	if err != nil {
		t.Errorf("links table was not created: %v", err)
	}

	// Test table structure by inserting test data
	// First insert into main table
	result, err := db.Exec(`INSERT INTO job_data (location, job_title, company_name, description, date_posted, salary, work_from_home, qualifications, links, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Boston, MA", "Software Engineer", "Tech Corp", "Build apps", "2024-01-01", 100000, "No", "Go, Python", "http://example.com", "USA")
	if err != nil {
		t.Fatalf("Failed to insert test data into job_data: %v", err)
	}

	jobID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}

	// Test qualifications table
	_, err = db.Exec("INSERT INTO qualifications (id, qualifications) VALUES (?, ?)", jobID, "Go, Python")
	if err != nil {
		t.Errorf("Failed to insert into qualifications table: %v", err)
	}

	// Test links table
	_, err = db.Exec("INSERT INTO links (id, links) VALUES (?, ?)", jobID, "http://example.com")
	if err != nil {
		t.Errorf("Failed to insert into links table: %v", err)
	}

	// Verify foreign key constraints work
	var qualCount, linkCount int
	err = db.QueryRow("SELECT COUNT(*) FROM qualifications WHERE id = ?", jobID).Scan(&qualCount)
	if err != nil {
		t.Errorf("Failed to count qualifications: %v", err)
	}
	if qualCount != 1 {
		t.Errorf("Expected 1 qualification row, got %d", qualCount)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM links WHERE id = ?", jobID).Scan(&linkCount)
	if err != nil {
		t.Errorf("Failed to count links: %v", err)
	}
	if linkCount != 1 {
		t.Errorf("Expected 1 link row, got %d", linkCount)
	}
}

func TestWriteToDatabase(t *testing.T) {
	// Create temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Set up database tables
	createMainTable(db)
	createSecondaryTables(db)

	// Test data
	testJobs := []shared.JobData{
		{
			Location:       "Boston, MA",
			JobTitle:       "Software Engineer",
			CompanyName:    "Tech Corp",
			Description:    "Build applications",
			DatePosted:     "2024-01-01",
			Salary:         100000,
			WorkFromHome:   "No",
			Qualifications: "Go, Python",
			Links:          "http://example.com",
			Country:        "USA",
		},
		{
			Location:       "New York, NY",
			JobTitle:       "Frontend Developer",
			CompanyName:    "Web Co",
			Description:    "Create UIs",
			DatePosted:     "2024-01-02",
			Salary:         90000,
			WorkFromHome:   "Yes",
			Qualifications: "JavaScript, React",
			Links:          "http://example2.com",
			Country:        "USA",
		},
	}

	// Write test data
	WriteToDatabase(db, testJobs)

	// Verify data was written to main table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM job_data").Scan(&count)
	if err != nil {
		t.Errorf("Failed to count job_data rows: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 rows in job_data, got %d", count)
	}

	// Verify data was written to secondary tables
	err = db.QueryRow("SELECT COUNT(*) FROM qualifications").Scan(&count)
	if err != nil {
		t.Errorf("Failed to count qualifications rows: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 rows in qualifications, got %d", count)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM links").Scan(&count)
	if err != nil {
		t.Errorf("Failed to count links rows: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 rows in links, got %d", count)
	}

	// Verify specific data integrity
	var jobTitle, companyName, qualifications, links string
	err = db.QueryRow("SELECT job_title, company_name, qualifications, links FROM job_data WHERE location = ?", "Boston, MA").Scan(&jobTitle, &companyName, &qualifications, &links)
	if err != nil {
		t.Errorf("Failed to query specific job data: %v", err)
	}

	if jobTitle != "Software Engineer" {
		t.Errorf("Expected job title 'Software Engineer', got %s", jobTitle)
	}
	if companyName != "Tech Corp" {
		t.Errorf("Expected company name 'Tech Corp', got %s", companyName)
	}
	if qualifications != "Go, Python" {
		t.Errorf("Expected qualifications 'Go, Python', got %s", qualifications)
	}
	if links != "http://example.com" {
		t.Errorf("Expected links 'http://example.com', got %s", links)
	}
}