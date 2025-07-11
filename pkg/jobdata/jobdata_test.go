package jobdata

import (
	"job-visualizer/pkg/shared"
	"testing"
)

func createTestJob(company, title, country, location, datePosted string, salary int) shared.JobData {
	return shared.JobData{
		CompanyName: company,
		JobTitle:    title,
		Country:     country,
		Location:    location,
		DatePosted:  datePosted,
		Salary:      salary,
	}
}

func createTestRows(header bool, dataRows ...[]string) [][]string {
	rows := [][]string{}
	if header {
		rows = append(rows, []string{"Company Name", "Date Posted", "Job ID", "Country", "Location", "Col5", "Max Salary", "Min Salary", "Hourly/Yearly", "Job Title"})
	}
	rows = append(rows, dataRows...)
	return rows
}

func createSalaryTestRow(maxSalary, minSalary, hourlyYearly string) []string {
	return []string{"", "", "", "", "", "", maxSalary, minSalary, hourlyYearly, ""}
}

func assertJobData(t *testing.T, job shared.JobData, expectedRow []string) {
	if job.CompanyName != expectedRow[0] {
		t.Errorf("CompanyName = %v, want %v", job.CompanyName, expectedRow[0])
	}
	if job.DatePosted != expectedRow[1] {
		t.Errorf("DatePosted = %v, want %v", job.DatePosted, expectedRow[1])
	}
	if job.Country != expectedRow[3] {
		t.Errorf("Country = %v, want %v", job.Country, expectedRow[3])
	}
	if job.Location != expectedRow[4] {
		t.Errorf("Location = %v, want %v", job.Location, expectedRow[4])
	}
	if job.JobTitle != expectedRow[9] {
		t.Errorf("JobTitle = %v, want %v", job.JobTitle, expectedRow[9])
	}
}

func TestProcessRows(t *testing.T) {
	tests := []struct {
		name       string
		rows       [][]string
		allJobData []shared.JobData
		expected   int
	}{
		{
			name:       "Empty rows should return original data",
			rows:       [][]string{},
			allJobData: []shared.JobData{createTestJob("Existing Company", "Existing Job", "", "", "", 0)},
			expected:   1,
		},
		{
			name:       "Single header row should return original data",
			rows:       createTestRows(true),
			allJobData: []shared.JobData{createTestJob("Existing Company", "Existing Job", "", "", "", 0)},
			expected:   1,
		},
		{
			name: "Valid job data should be processed correctly",
			rows: createTestRows(true,
				[]string{"Tech Corp", "2024-01-15", "123", "USA", "New York", "Col5", "120000", "80000", "yearly", "Software Engineer"},
				[]string{"Data Inc", "2024-01-16", "124", "Canada", "Toronto", "Col5", "100000", "70000", "yearly", "Data Scientist"},
			),
			allJobData: []shared.JobData{},
			expected:   2,
		},
		{
			name: "Row with insufficient columns should be skipped",
			rows: createTestRows(true,
				[]string{"Tech Corp", "2024-01-15", "123", "USA", "New York"},
				[]string{"Data Inc", "2024-01-16", "124", "Canada", "Toronto", "Col5", "100000", "70000", "yearly", "Data Scientist"},
			),
			allJobData: []shared.JobData{},
			expected:   1,
		},
		{
			name: "Should append to existing job data",
			rows: createTestRows(true,
				[]string{"Tech Corp", "2024-01-15", "123", "USA", "New York", "Col5", "120000", "80000", "yearly", "Software Engineer"},
			),
			allJobData: []shared.JobData{createTestJob("Existing Company", "Existing Job", "", "", "", 0)},
			expected:   2,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := ProcessRows(testCase.rows, testCase.allJobData)

			if len(result) != testCase.expected {
				t.Errorf("ProcessRows() returned %d jobs, expected %d", len(result), testCase.expected)
			}

			if len(testCase.rows) > 1 && len(testCase.rows[1]) >= 10 {
				newJobIndex := len(testCase.allJobData)
				if newJobIndex < len(result) {
					newJob := result[newJobIndex]
					expectedRow := testCase.rows[1] // First data row
					assertJobData(t, newJob, expectedRow)
				}
			}
		})
	}
}

func TestCalcSalary(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		expected int
	}{
		{
			name:     "Yearly salary calculation",
			row:      createSalaryTestRow("120000", "80000", "yearly"),
			expected: 100000, // (120000 + 80000) / 2
		},
		{
			name:     "Hourly salary calculation",
			row:      createSalaryTestRow("50", "30", "hourly"),
			expected: 80000, // (50 + 30) / 2 * 40 * 50 = 40 * 40 * 50 = 80000
		},
		{
			name:     "Zero salary values",
			row:      createSalaryTestRow("0", "0", "yearly"),
			expected: 0,
		},
		{
			name:     "Same min and max salary",
			row:      createSalaryTestRow("75000", "75000", "yearly"),
			expected: 75000,
		},
		{
			name:     "Decimal salary values",
			row:      createSalaryTestRow("125000.50", "75000.25", "yearly"),
			expected: 100000, // (125000.50 + 75000.25) / 2 = 100000.375, truncated to int
		},
		{
			name:     "Hourly with decimal values",
			row:      createSalaryTestRow("45.50", "35.25", "hourly"),
			expected: 80750, // (45.50 + 35.25) / 2 * 40 * 50 = 40.375 * 40 * 50 = 80750
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := calcSalary(testCase.row)
			if result != testCase.expected {
				t.Errorf("calcSalary() = %v, want %v", result, testCase.expected)
			}
		})
	}
}
