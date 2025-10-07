// package buildcontainers provides tests for GUI container building functionality
package buildcontainers

import (
	"job-visualizer/pkg/shared"
	"strings"
	"testing"

	"fyne.io/fyne/v2/widget"
)

// TestFormatJobDetails tests the formatJobDetails function with sample job data
func TestFormatJobDetails(t *testing.T) {
	// creates test job data with all fields populated
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

	// tests job details formatting
	result := formatJobDetails(0, windowData)

	// verifies all job fields are included in formatted output
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

// TestCreateListItem tests the createListItem function to ensure it returns a proper list
func TestCreateListItem(t *testing.T) {
	// tests list item creation
	item := createListItem()
	label, ok := item.(*widget.Label)
	if !ok {
		t.Fatalf("createListItem did not return a *widget.Label, got %T", item)
	}

	// verifies default label text
	if label.Text != "list items here" {
		t.Errorf("Expected label text 'list items here', got '%s'", label.Text)
	}
}

// TestBuildRightSplit tests the BuildRightSplit function to ensure it creates the correct container layout
func TestBuildRightSplit(t *testing.T) {
	windowData := &shared.GuiWindowData{}
	container := BuildRightSplit(windowData)

	// verifies container was created successfully
	if container == nil {
		t.Fatal("Expected non-nil container")
	}

	// verifies container has correct number of objects
	objs := container.Objects
	if len(objs) != 2 {
		t.Fatalf("Expected 2 objects in container, got %d", len(objs))
	}

	// verifies first object is a label with correct text
	label, ok := objs[0].(*widget.Label)
	if !ok {
		t.Fatalf("First object is not *widget.Label, got %T", objs[0])
	}
	if label.Text != "Select a job to display details" {
		t.Errorf("Expected label text 'Select a job to display details', got '%s'", label.Text)
	}

	// verifies second object is a quit button
	button, ok := objs[1].(*widget.Button)
	if !ok {
		t.Fatalf("Second object is not *widget.Button, got %T", objs[1])
	}
	if button.Text != "Quit" {
		t.Errorf("Expected button text 'Quit', got '%s'", button.Text)
	}
}

// containsAllSubstrings verifies that all required substrings are present in the output
func containsAllSubstrings(output string, requiredSubstrings []string) bool {
	for _, substring := range requiredSubstrings {
		if !strings.Contains(output, substring) {
			return false
		}
	}
	return true
}
