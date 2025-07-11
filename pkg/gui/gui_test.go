package gui

import (
	"job-visualizer/pkg/shared"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func TestCreateGuiWindow(t *testing.T) {
	application = app.New()
	defer func() { application = nil }()

	title := "test window"
	window := createGuiWindow(title)

	if window == nil {
		t.Fatal("Expected window to be created, but got nil")
	}
	if window.Title() != title {
		t.Errorf("Expected window title '%s', but got '%s'", title, window.Title())
	}
	expectedSize := fyne.NewSize(1000, 600)
	if window.Canvas().Size() != expectedSize {
		t.Errorf("Expected window size %v, but got %v", expectedSize, window.Canvas().Size())
	}
}

func TestProcessJobsReturnsSlice(t *testing.T) {
	tempDir := t.TempDir()
	shared.Program.OutputDirectory = tempDir
	shared.Program.ResourcesDirectory = tempDir
	result := processJobs(nil)
	if result == nil {
		t.Fatal("Expected processJobs to return a non-nil slice, but got nil")
	}
}
