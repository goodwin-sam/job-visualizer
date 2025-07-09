package filter

import (
	"job-visualizer/pkg/shared"
	"testing"
)

func TestFilterKeyword(t *testing.T) {
	tests := []struct {
		name        string
		job         shared.JobData
		filterInput string
		expected    bool
	}{
		{
			name: "matches job title",
			job: shared.JobData{
				JobTitle:       "Software Engineer",
				CompanyName:    "Tech Corp",
				Description:    "Build applications",
				Qualifications: "Go, Python",
			},
			filterInput: "software",
			expected:    true,
		},
		{
			name: "matches company name",
			job: shared.JobData{
				JobTitle:       "Developer",
				CompanyName:    "Google",
				Description:    "Build applications",
				Qualifications: "Go, Python",
			},
			filterInput: "google",
			expected:    true,
		},
		{
			name: "matches description",
			job: shared.JobData{
				JobTitle:       "Developer",
				CompanyName:    "Tech Corp",
				Description:    "Build microservices",
				Qualifications: "Go, Python",
			},
			filterInput: "microservices",
			expected:    true,
		},
		{
			name: "matches qualifications",
			job: shared.JobData{
				JobTitle:       "Developer",
				CompanyName:    "Tech Corp",
				Description:    "Build applications",
				Qualifications: "Go, Python, JavaScript",
			},
			filterInput: "javascript",
			expected:    true,
		},
		{
			name: "no match",
			job: shared.JobData{
				JobTitle:       "Developer",
				CompanyName:    "Tech Corp",
				Description:    "Build applications",
				Qualifications: "Go, Python",
			},
			filterInput: "marketing",
			expected:    false,
		},
		{
			name: "case insensitive match",
			job: shared.JobData{
				JobTitle:       "SOFTWARE ENGINEER",
				CompanyName:    "Tech Corp",
				Description:    "Build applications",
				Qualifications: "Go, Python",
			},
			filterInput: "software",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterKeyword(tt.job, tt.filterInput)
			if result != tt.expected {
				t.Errorf("filterKeyword() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFilterLocation(t *testing.T) {
	tests := []struct {
		name        string
		job         shared.JobData
		filterInput string
		expected    bool
	}{
		{
			name: "exact match",
			job: shared.JobData{
				Location: "Boston, MA",
			},
			filterInput: "boston",
			expected:    true,
		},
		{
			name: "partial match",
			job: shared.JobData{
				Location: "Boston, MA",
			},
			filterInput: "ma",
			expected:    true,
		},
		{
			name: "no match",
			job: shared.JobData{
				Location: "Boston, MA",
			},
			filterInput: "california",
			expected:    false,
		},
		{
			name: "case insensitive",
			job: shared.JobData{
				Location: "BOSTON, MA",
			},
			filterInput: "boston",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterLocation(tt.job, tt.filterInput)
			if result != tt.expected {
				t.Errorf("filterLocation() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFilterMinSalary(t *testing.T) {
	tests := []struct {
		name     string
		job      shared.JobData
		filter   string
		expected bool
	}{
		{
			name: "salary above minimum",
			job: shared.JobData{
				Salary: 100000,
			},
			filter:   "80000",
			expected: true,
		},
		{
			name: "salary below minimum",
			job: shared.JobData{
				Salary: 70000,
			},
			filter:   "80000",
			expected: false,
		},
		{
			name: "salary equal to minimum",
			job: shared.JobData{
				Salary: 80000,
			},
			filter:   "80000",
			expected: false, // function uses > not >=
		},
		{
			name: "zero salary",
			job: shared.JobData{
				Salary: 0,
			},
			filter:   "50000",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterMinSalary(tt.job, tt.filter)
			if result != tt.expected {
				t.Errorf("filterMinSalary() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFilterWorkFromHome(t *testing.T) {
	tests := []struct {
		name     string
		job      shared.JobData
		expected bool
	}{
		{
			name: "work from home yes",
			job: shared.JobData{
				WorkFromHome: "Yes",
			},
			expected: true,
		},
		{
			name: "work from home no",
			job: shared.JobData{
				WorkFromHome: "No",
			},
			expected: false,
		},
		{
			name: "work from home empty",
			job: shared.JobData{
				WorkFromHome: "",
			},
			expected: false,
		},
		{
			name: "work from home other value",
			job: shared.JobData{
				WorkFromHome: "Sometimes",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterWorkFromHome(tt.job)
			if result != tt.expected {
				t.Errorf("filterWorkFromHome() = %v, expected %v", result, tt.expected)
			}
		})
	}
}