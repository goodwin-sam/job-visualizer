package buildcontainers

import (
	"job-visualizer/pkg/shared"
	"strings"
	"testing"
)

func TestFormatJobDetails(t *testing.T) {
	jobs := []shared.JobData{
		{
			CompanyName:    "TestCorp",
			JobTitle:       "Software Engineer",
			Location:       "Remote",
			DatePosted:     "2024-06-01",
			Salary:         120000,
			WorkFromHome:   "Yes",
			Qualifications: "Go, Docker",
			Links:          "https://testcorp.com/jobs/1",
		},
	}
	windowData := shared.GuiWindowData{
		FilteredJobs: &jobs,
	}
	result := formatJobDetails(0, windowData)

	if !containsAllSubstrings(result, []string{
		"TestCorp",
		"Software Engineer",
		"Remote",
		"2024-06-01",
		"120000",
		"Yes",
		"Go, Docker",
		"https://testcorp.com/jobs/1",
	}) {
		t.Errorf("formatJobDetails output missing expected content: %s", result)
	}
}

func containsAllSubstrings(output string, requiredSubstrings []string) bool {
	for _, substring := range requiredSubstrings {
		if !strings.Contains(output, substring) {
			return false
		}
	}
	return true
}
