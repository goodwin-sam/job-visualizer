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

func TestWriteToDatabase(t *testing.T) {
	tempDirectory := t.TempDir()
	shared.Program.OutputDirectory = tempDirectory
	db := CreateDatabase()
	defer db.Close()

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
	err := db.QueryRow("SELECT COUNT(*) FROM job_data").Scan(&jobCount)
	if err != nil {
		t.Errorf("Error counting jobs in database: %v", err)
	}
	if jobCount != 2 {
		t.Errorf("Expected 2 jobs in database, got %d", jobCount)
	}

	var location, jobTitle, companyName, description, datePosted, workFromHome, country string
	var salary int
	err = db.QueryRow(
		`SELECT location, job_title, company_name, description, date_posted, salary, work_from_home, country 
		 FROM job_data WHERE company_name = ?`, "Tech Corp",
	).Scan(&location, &jobTitle, &companyName, &description, &datePosted, &salary, &workFromHome, &country)
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
	if companyName != "Tech Corp" {
		t.Errorf("Expected company name 'Tech Corp', got '%s'", companyName)
	}
	if description != "Build amazing software" {
		t.Errorf("Expected description 'Build amazing software', got '%s'", description)
	}
	if datePosted != "2024-01-15" {
		t.Errorf("Expected date posted '2024-01-15', got '%s'", datePosted)
	}
	if salary != 80000 {
		t.Errorf("Expected salary 80000, got %d", salary)
	}
	if workFromHome != "Yes" {
		t.Errorf("Expected work from home 'Yes', got '%s'", workFromHome)
	}
	if country != "USA" {
		t.Errorf("Expected country 'USA', got '%s'", country)
	}

	// verifying the qualifications table
	var qualCount int
	err = db.QueryRow("SELECT COUNT(*) FROM qualifications").Scan(&qualCount)
	if err != nil {
		t.Errorf("Error counting qualifications: %v", err)
	}
	if qualCount != 2 {
		t.Errorf("Expected 2 qualifications entries, got %d", qualCount)
	}

	// verifying the links table
	var linkCount int
	err = db.QueryRow("SELECT COUNT(*) FROM links").Scan(&linkCount)
	if err != nil {
		t.Errorf("Error counting links: %v", err)
	}
	if linkCount != 2 {
		t.Errorf("Expected 2 links entries, got %d", linkCount)
	}

	// verifying the foreign key relationships
	var qualifications string
	err = db.QueryRow("SELECT q.qualifications FROM qualifications q JOIN job_data j ON q.id = j.id WHERE j.company_name = ?", "Tech Corp").Scan(&qualifications)
	if err != nil {
		t.Errorf("Error querying qualifications with join: %v", err)
	}
	if qualifications != "Go, SQL, Git" {
		t.Errorf("Expected qualifications 'Go, SQL, Git', got '%s'", qualifications)
	}

	var links string
	err = db.QueryRow("SELECT l.links FROM links l JOIN job_data j ON l.id = j.id WHERE j.company_name = ?", "Data Inc").Scan(&links)
	if err != nil {
		t.Errorf("Error querying links with join: %v", err)
	}
	if links != "https://datainc.com/careers" {
		t.Errorf("Expected links 'https://datainc.com/careers', got '%s'", links)
	}
}
