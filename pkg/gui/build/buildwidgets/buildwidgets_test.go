// package buildwidgets provides tests for GUI widget building functionality
package buildwidgets

import (
	"testing"

	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2"
)

// TestBuildLabel tests the BuildLabel function with various text styling options
func TestBuildLabel(t *testing.T) {
	// defines test cases for different label styling combinations
	testCases := []struct {
		name   string
		text   string
		bold   bool
		italic bool
	}{
		{"plain", "Test Label", false, false},
		{"bold", "Bold Label", true, false},
		{"italic", "Italic Label", false, true},
		{"bolditalic", "Bold Italic Label", true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// tests label creation with specific styling
			label := BuildLabel(tc.text, tc.bold, tc.italic)

			// verifies text content
			if label.Text != tc.text {
				t.Errorf("expected text %q, got %q", tc.text, label.Text)
			}

			// verifies center alignment
			if label.Alignment != fyne.TextAlignCenter {
				t.Errorf("expected alignment %v, got %v", fyne.TextAlignCenter, label.Alignment)
			}

			// verifies bold styling
			if label.TextStyle.Bold != tc.bold {
				t.Errorf("expected bold %v, got %v", tc.bold, label.TextStyle.Bold)
			}

			// verifies italic styling
			if label.TextStyle.Italic != tc.italic {
				t.Errorf("expected italic %v, got %v", tc.italic, label.TextStyle.Italic)
			}
		})
	}
}

// TestRemoveActiveFilters tests the removeActiveFilters function to ensure it clears all filter values
func TestRemoveActiveFilters(t *testing.T) {
	// sets up window data with active filter values
	windowData := &shared.GuiWindowData{}
	windowData.Filters.KeywordEntry = "developer"
	windowData.Filters.LocationEntry = "remote"
	windowData.Filters.MinSalaryEntry = "100000"
	windowData.Filters.WorkFromHomeEntry = true

	// tests filter removal
	removeActiveFilters(windowData)

	// verifies all filter entries are cleared
	if windowData.Filters.KeywordEntry != "" {
		t.Errorf("expected KeywordEntry to be reset, got %q", windowData.Filters.KeywordEntry)
	}
	if windowData.Filters.LocationEntry != "" {
		t.Errorf("expected LocationEntry to be reset, got %q", windowData.Filters.LocationEntry)
	}
	if windowData.Filters.MinSalaryEntry != "" {
		t.Errorf("expected MinSalaryEntry to be reset, got %q", windowData.Filters.MinSalaryEntry)
	}
	if windowData.Filters.WorkFromHomeEntry != false {
		t.Errorf("expected WorkFromHomeEntry to be reset to false, got %v", windowData.Filters.WorkFromHomeEntry)
	}
}

// TestBuildRemoteCheckbox tests the BuildRemoteCheckbox function to ensure it properly handles checkbox state changes
func TestBuildRemoteCheckbox(t *testing.T) {
	// sets up window data with unchecked work-from-home filter
	windowData := &shared.GuiWindowData{}
	windowData.Filters.WorkFromHomeEntry = false
	cb := BuildRemoteCheckbox(windowData)

	// tests checkbox checked state
	cb.OnChanged(true)
	if windowData.Filters.WorkFromHomeEntry != true {
		t.Errorf("expected WorkFromHomeEntry to be true after checking, got %v", windowData.Filters.WorkFromHomeEntry)
	}

	// tests checkbox unchecked state
	cb.OnChanged(false)
	if windowData.Filters.WorkFromHomeEntry != false {
		t.Errorf("expected WorkFromHomeEntry to be false after unchecking, got %v", windowData.Filters.WorkFromHomeEntry)
	}
}
