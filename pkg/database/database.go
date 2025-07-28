package database

import (
	"database/sql"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() *sql.DB {
	databasePath := filepath.Join(shared.Program.OutputDirectory, "job_data.sqlite")
	err := os.Remove(databasePath)
	shared.CheckErrorWarn(err)
	db, err := sql.Open("sqlite3", databasePath)
	shared.CheckError(err)
	_, err = db.Exec("PRAGMA journal_mode=DELETE;")
	shared.CheckError(err)
	return db
}

func SetupDatabase(db *sql.DB) {
	createMainTable(db)
	createSecondaryTables(db)
}

func WriteToDatabase(db *sql.DB, allJobData []shared.JobData) {
	insertQueryJobData := `INSERT INTO job_data (location, job_title, company_name, description, date_posted, salary,
		work_from_home, qualifications, links, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	insertQueryQualifications := `INSERT OR IGNORE INTO qualifications (id, qualifications) VALUES (?, ?);`
	insertQueryLinks := `INSERT OR IGNORE INTO links (id, links) VALUES (?, ?);`
	for _, job := range allJobData {
		result, err := db.Exec(insertQueryJobData, job.Location, job.JobTitle, job.CompanyName, job.Description, job.DatePosted,
			job.Salary, job.WorkFromHome, job.Qualifications, job.Links, job.Country)
		shared.CheckError(err)
		id, err := result.LastInsertId()
		shared.CheckError(err)
		_, err = db.Exec(insertQueryQualifications, id, job.Qualifications)
		shared.CheckError(err)
		_, err = db.Exec(insertQueryLinks, id, job.Links)
		shared.CheckError(err)
	}
}

func createMainTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS job_data(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		location TEXT NOT NULL,
		job_title TEXT NOT NULL,
		company_name TEXT NOT NULL,
		description TEXT,
		date_posted TEXT NOT NULL,
		salary INT,
		work_from_home TEXT,
		qualifications TEXT,
		links TEXT,
		country TEXT
	);`)
	shared.CheckError(err)
}

func createSecondaryTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS qualifications(
		id INTEGER PRIMARY KEY,
		qualifications TEXT NOT NULL,
		FOREIGN KEY (id) REFERENCES job_data(id)
	);`)
	shared.CheckError(err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links(
		id INTEGER PRIMARY KEY,
		links TEXT NOT NULL,
		FOREIGN KEY (id) REFERENCES job_data(id)
	);`)
	shared.CheckError(err)
}
