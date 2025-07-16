package buildwidgets

import (
	"testing"

	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2"
)

func TestBuildLabel(t *testing.T) {
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
			label := BuildLabel(tc.text, tc.bold, tc.italic)
			if label.Text != tc.text {
				t.Errorf("expected text %q, got %q", tc.text, label.Text)
			}
			if label.Alignment != fyne.TextAlignCenter {
				t.Errorf("expected alignment %v, got %v", fyne.TextAlignCenter, label.Alignment)
			}
			if label.TextStyle.Bold != tc.bold {
				t.Errorf("expected bold %v, got %v", tc.bold, label.TextStyle.Bold)
			}
			if label.TextStyle.Italic != tc.italic {
				t.Errorf("expected italic %v, got %v", tc.italic, label.TextStyle.Italic)
			}
		})
	}
}

func TestRemoveActiveFilters(t *testing.T) {
	shared.WindowData.Filters.KeywordEntry = "developer"
	shared.WindowData.Filters.LocationEntry = "remote"
	shared.WindowData.Filters.MinSalaryEntry = "100000"
	shared.WindowData.Filters.WorkFromHomeEntry = true

	removeActiveFilters()

	if shared.WindowData.Filters.KeywordEntry != "" {
		t.Errorf("expected KeywordEntry to be reset, got %q", shared.WindowData.Filters.KeywordEntry)
	}
	if shared.WindowData.Filters.LocationEntry != "" {
		t.Errorf("expected LocationEntry to be reset, got %q", shared.WindowData.Filters.LocationEntry)
	}
	if shared.WindowData.Filters.MinSalaryEntry != "" {
		t.Errorf("expected MinSalaryEntry to be reset, got %q", shared.WindowData.Filters.MinSalaryEntry)
	}
	if shared.WindowData.Filters.WorkFromHomeEntry != false {
		t.Errorf("expected WorkFromHomeEntry to be reset to false, got %v", shared.WindowData.Filters.WorkFromHomeEntry)
	}
}

func TestBuildRemoteCheckbox(t *testing.T) {
	shared.WindowData.Filters.WorkFromHomeEntry = false
	cb := BuildRemoteCheckbox()

	cb.OnChanged(true)
	if shared.WindowData.Filters.WorkFromHomeEntry != true {
		t.Errorf("expected WorkFromHomeEntry to be true after checking, got %v", shared.WindowData.Filters.WorkFromHomeEntry)
	}

	cb.OnChanged(false)
	if shared.WindowData.Filters.WorkFromHomeEntry != false {
		t.Errorf("expected WorkFromHomeEntry to be false after unchecking, got %v", shared.WindowData.Filters.WorkFromHomeEntry)
	}
}
