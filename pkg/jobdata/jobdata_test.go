package jobdata

import (
	"job-visualizer/pkg/shared"
	"testing"
)

func TestProcessRows(t *testing.T) {
	tests := []struct {
		name         string
		rows         [][]string
		allJobData   []shared.JobData
		expectedJobs int
	}{
		{
			name: "valid rows",
			rows: [][]string{
				{"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", "MaxSalary", "MinSalary", "HourlyOrYearly", "JobTitle"},
				{"Tech Corp", "2024-01-01", "123", "USA", "Boston, MA", "field5", "120000", "80000", "yearly", "Software Engineer"},
				{"Web Co", "2024-01-02", "124", "USA", "New York, NY", "field5", "100000", "60000", "yearly", "Frontend Developer"},
			},
			allJobData:   []shared.JobData{},
			expectedJobs: 2,
		},
		{
			name: "skip rows with insufficient columns",
			rows: [][]string{
				{"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", "MaxSalary", "MinSalary", "HourlyOrYearly", "JobTitle"},
				{"Tech Corp", "2024-01-01", "123", "USA", "Boston, MA", "field5", "120000", "80000", "yearly", "Software Engineer"},
				{"Incomplete"}, // This row should be skipped
				{"Web Co", "2024-01-02", "124", "USA", "New York, NY", "field5", "100000", "60000", "yearly", "Frontend Developer"},
			},
			allJobData:   []shared.JobData{},
			expectedJobs: 2,
		},
		{
			name: "empty rows",
			rows: [][]string{
				{"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", "MaxSalary", "MinSalary", "HourlyOrYearly", "JobTitle"},
			},
			allJobData:   []shared.JobData{},
			expectedJobs: 0,
		},
		{
			name: "append to existing data",
			rows: [][]string{
				{"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", "MaxSalary", "MinSalary", "HourlyOrYearly", "JobTitle"},
				{"Tech Corp", "2024-01-01", "123", "USA", "Boston, MA", "field5", "120000", "80000", "yearly", "Software Engineer"},
			},
			allJobData: []shared.JobData{
				{CompanyName: "Existing Corp", JobTitle: "Existing Job"},
			},
			expectedJobs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessRows(tt.rows, tt.allJobData)
			if len(result) != tt.expectedJobs {
				t.Errorf("ProcessRows() returned %d jobs, expected %d", len(result), tt.expectedJobs)
			}

			// Check that the original data is preserved when appending
			if len(tt.allJobData) > 0 {
				if result[0].CompanyName != tt.allJobData[0].CompanyName {
					t.Errorf("Original job data was not preserved")
				}
			}

			// Check that new jobs have correct data
			if len(result) > len(tt.allJobData) {
				newJobIndex := len(tt.allJobData)
				if result[newJobIndex].CompanyName != tt.rows[1][0] {
					t.Errorf("Expected company name %s, got %s", tt.rows[1][0], result[newJobIndex].CompanyName)
				}
				if result[newJobIndex].JobTitle != tt.rows[1][9] {
					t.Errorf("Expected job title %s, got %s", tt.rows[1][9], result[newJobIndex].JobTitle)
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
			name: "yearly salary",
			row: []string{
				"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", 
				"120000", "80000", "yearly", "JobTitle",
			},
			expected: 100000, // (120000 + 80000) / 2
		},
		{
			name: "hourly salary",
			row: []string{
				"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", 
				"60", "40", "hourly", "JobTitle",
			},
			expected: 100000, // (60 + 40) / 2 * 40 * 50 = 50 * 40 * 50 = 100000
		},
		{
			name: "yearly salary with decimals",
			row: []string{
				"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", 
				"125000.5", "85000.5", "yearly", "JobTitle",
			},
			expected: 105000, // (125000.5 + 85000.5) / 2 = 105000.5, converted to int = 105000
		},
		{
			name: "hourly salary with decimals",
			row: []string{
				"CompanyName", "DatePosted", "JobId", "Country", "Location", "Field5", 
				"62.5", "37.5", "hourly", "JobTitle",
			},
			expected: 100000, // (62.5 + 37.5) / 2 * 40 * 50 = 50 * 40 * 50 = 100000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calcSalary(tt.row)
			if result != tt.expected {
				t.Errorf("calcSalary() = %d, expected %d", result, tt.expected)
			}
		})
	}
}