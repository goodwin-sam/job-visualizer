// package jobsprocessing provides tests for job data processing functionality
package jobsprocessing

import (
	"job-visualizer/pkg/shared"
	"testing"
)

// createTestJob creates a test job with specified field values
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

// createTestRows creates test Excel rows with header and data rows
func createTestRows(header bool, dataRows ...[]string) [][]string {
	rows := [][]string{}
	if header {
		rows = append(rows, []string{"Company Name", "Date Posted", "Job ID", "Country", "Location", "Col5", "Max Salary", "Min Salary", "Hourly/Yearly", "Job Title"})
	}
	rows = append(rows, dataRows...)
	return rows
}

// createSalaryTestRow creates a test row with salary-related fields for testing
func createSalaryTestRow(maxSalary, minSalary, hourlyYearly string) []string {
	return []string{"", "", "", "", "", "", maxSalary, minSalary, hourlyYearly, ""}
}

// assertJobData verifies that job data matches expected row values
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

// TestProcessRows tests the ProcessRows function with various input scenarios
func TestProcessRows(t *testing.T) {
	tests := []struct {
		name       string
		rows       [][]string
		allJobData []shared.JobData
		expected   int
	}{
		{
			// tests handling of empty input rows
			name:       "Empty rows should return original data",
			rows:       [][]string{},
			allJobData: []shared.JobData{createTestJob("Existing Company", "Existing Job", "", "", "", 0)},
			expected:   1,
		},
		{
			// tests handling of header-only rows
			name:       "Single header row should return original data",
			rows:       createTestRows(true),
			allJobData: []shared.JobData{createTestJob("Existing Company", "Existing Job", "", "", "", 0)},
			expected:   1,
		},
		{
			// tests processing of valid job data with multiple rows
			name: "Valid job data should be processed correctly",
			rows: createTestRows(true,
				[]string{"Tech Corp", "2024-01-15", "123", "USA", "New York", "Col5", "120000", "80000", "yearly", "Software Engineer"},
				[]string{"Data Inc", "2024-01-16", "124", "Canada", "Toronto", "Col5", "100000", "70000", "yearly", "Data Scientist"},
			),
			allJobData: []shared.JobData{},
			expected:   2,
		},
		{
			// tests skipping of rows with insufficient columns
			name: "Row with insufficient columns should be skipped",
			rows: createTestRows(true,
				[]string{"Tech Corp", "2024-01-15", "123", "USA", "New York"},
				[]string{"Data Inc", "2024-01-16", "124", "Canada", "Toronto", "Col5", "100000", "70000", "yearly", "Data Scientist"},
			),
			allJobData: []shared.JobData{},
			expected:   1,
		},
		{
			// tests appending new jobs to existing job data
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

// TestCalcSalary tests the calcSalary function with various salary calculation scenarios
func TestCalcSalary(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		expected int
	}{
		{
			// tests average calculation for yearly salaries
			name:     "Yearly salary calculation",
			row:      createSalaryTestRow("120000", "80000", "yearly"),
			expected: 100000, // (120000 + 80000) / 2
		},
		{
			// tests hourly to yearly conversion (40 hours/week * 50 weeks/year)
			name:     "Hourly salary calculation",
			row:      createSalaryTestRow("50", "30", "hourly"),
			expected: 80000, // (50 + 30) / 2 * 40 * 50 = 40 * 40 * 50 = 80000
		},
		{
			// tests handling of zero salary values
			name:     "Zero salary values",
			row:      createSalaryTestRow("0", "0", "yearly"),
			expected: 0,
		},
		{
			// tests handling of identical min and max salaries
			name:     "Same min and max salary",
			row:      createSalaryTestRow("75000", "75000", "yearly"),
			expected: 75000,
		},
		{
			// tests decimal salary handling with truncation
			name:     "Decimal salary values",
			row:      createSalaryTestRow("125000.50", "75000.25", "yearly"),
			expected: 100000, // (125000.50 + 75000.25) / 2 = 100000.375, truncated to int
		},
		{
			// tests hourly conversion with decimal values
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
