package buildcontainers

import (
	"job-visualizer/pkg/shared"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
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

func TestCreateListItem(t *testing.T) {
	item := createListItem()
	label, ok := item.(*widget.Label)
	if !ok {
		t.Fatalf("createListItem did not return a *widget.Label, got %T", item)
	}
	if label.Text != "list items here" {
		t.Errorf("Expected label text 'list items here', got '%s'", label.Text)
	}
}

func TestBuildKeywordContainer(t *testing.T) {
	testFilterContainer(
		t,
		buildKeywordContainer,
		func(wd *shared.GuiWindowData) { wd.KeywordEntryWidget = nil },
		func(wd *shared.GuiWindowData, val string) { wd.Filters.KeywordEntry = val },
		func(wd *shared.GuiWindowData) string { return wd.Filters.KeywordEntry },
		"test-keyword",
	)
}

func TestBuildLocationContainer(t *testing.T) {
	testFilterContainer(
		t,
		buildLocationContainer,
		func(wd *shared.GuiWindowData) { wd.LocationEntryWidget = nil },
		func(wd *shared.GuiWindowData, val string) { wd.Filters.LocationEntry = val },
		func(wd *shared.GuiWindowData) string { return wd.Filters.LocationEntry },
		"test-location",
	)
}

func TestBuildMinSalaryContainer(t *testing.T) {
	testFilterContainer(
		t,
		buildMinSalaryContainer,
		func(wd *shared.GuiWindowData) { wd.MinSalaryEntryWidget = nil },
		func(wd *shared.GuiWindowData, val string) { wd.Filters.MinSalaryEntry = val },
		func(wd *shared.GuiWindowData) string { return wd.Filters.MinSalaryEntry },
		"12345",
	)
}

func TestBuildRightSplit(t *testing.T) {
	windowData := &shared.GuiWindowData{}
	container := BuildRightSplit(windowData)
	if container == nil {
		t.Fatal("Expected non-nil container")
	}

	objs := container.Objects
	if len(objs) != 2 {
		t.Fatalf("Expected 2 objects in container, got %d", len(objs))
	}

	label, ok := objs[0].(*widget.Label)
	if !ok {
		t.Fatalf("First object is not *widget.Label, got %T", objs[0])
	}
	if label.Text != "Select a job to display details" {
		t.Errorf("Expected label text 'Select a job to display details', got '%s'", label.Text)
	}

	button, ok := objs[1].(*widget.Button)
	if !ok {
		t.Fatalf("Second object is not *widget.Button, got %T", objs[1])
	}
	if button.Text != "Quit" {
		t.Errorf("Expected button text 'Quit', got '%s'", button.Text)
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

func testFilterContainer(
	t *testing.T,
	buildFunc func(*shared.GuiWindowData) *fyne.Container,
	resetWidget func(*shared.GuiWindowData),
	setFilterValue func(*shared.GuiWindowData, string),
	getFilterValue func(*shared.GuiWindowData) string,
	testValue string,
) {
	windowData := &shared.GuiWindowData{}
	resetWidget(windowData)
	setFilterValue(windowData, "")
	container := buildFunc(windowData)
	if container == nil {
		t.Fatal("Expected non-nil container")
	}
	objs := container.Objects
	if len(objs) != 2 {
		t.Fatalf("Expected 2 objects in container, got %d", len(objs))
	}
	entry, ok := objs[0].(*widget.Entry)
	if !ok {
		t.Fatalf("First object is not *widget.Entry, got %T", objs[0])
	}
	button, ok := objs[1].(*widget.Button)
	if !ok {
		t.Fatalf("Second object is not *widget.Button, got %T", objs[1])
	}
	entry.SetText(testValue)
	test.Tap(button)
	if getFilterValue(windowData) != testValue {
		t.Errorf("Expected filter to be set to '%s', got '%s'", testValue, getFilterValue(windowData))
	}
}
