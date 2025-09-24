// package buildcontainers provides tests for filter container building functionality
package buildcontainers

import (
	"job-visualizer/pkg/shared"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// TestBuildKeywordContainer tests the buildKeywordContainer function using the generic filter test
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

// TestBuildLocationContainer tests the buildLocationContainer function using the generic filter test
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

// TestBuildMinSalaryContainer tests the buildMinSalaryContainer function using the generic filter test
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

// testFilterContainer is a generic test function for filter container building functions
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

	// tests container building
	container := buildFunc(windowData)
	if container == nil {
		t.Fatal("Expected non-nil container")
	}

	// verifies container structure (entry + button)
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

	// tests filter functionality
	entry.SetText(testValue)
	test.Tap(button)
	if getFilterValue(windowData) != testValue {
		t.Errorf("Expected filter to be set to '%s', got '%s'", testValue, getFilterValue(windowData))
	}
}
