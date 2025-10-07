// package build provides tests for GUI window building functionality
package build

import (
	"testing"

	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// TestBuildStartWindow tests the BuildStartWindow function to ensure it sets up the initial window correctly
func TestBuildStartWindow(t *testing.T) {
	// creates test window and required components
	window := test.NewWindow(nil)
	startButton := widget.NewButton("Start", nil)
	progressBar := widget.NewProgressBar()
	programData := &shared.ProgramData{}

	// tests start window building
	resultWindow := BuildStartWindow(window, startButton, progressBar, programData)

	// verifies same window is returned
	if resultWindow != window {
		t.Errorf("Expected returned window to be the same as input window")
	}

	// verifies window content is set
	if window.Content() == nil {
		t.Errorf("Expected window content to be set, got nil")
	}
}

// TestBuildMainWindow tests the BuildMainWindow function to ensure it sets up the main job display window correctly
func TestBuildMainWindow(t *testing.T) {
	// creates test window and data structures
	window := test.NewWindow(nil)
	jobs := []shared.JobData{} // empty slice for simplicity
	windowData := &shared.GuiWindowData{}

	// tests main window building
	resultWindow := BuildMainWindow(window, jobs, windowData, nil)

	// verifies same window is returned
	if resultWindow != window {
		t.Errorf("Expected returned window to be the same as input window")
	}

	// verifies window content is set
	if window.Content() == nil {
		t.Errorf("Expected window content to be set, got nil")
	}
}

// TestBuildMainContent tests the buildMainContent function to ensure it creates the correct split layout
func TestBuildMainContent(t *testing.T) {
	// creates test data for content building
	jobs := []shared.JobData{}
	windowData := &shared.GuiWindowData{}

	// tests main content building
	contentPane := buildMainContent(jobs, windowData, nil)

	// verifies content pane was created
	if contentPane == nil {
		t.Errorf("Expected contentPane to be non-nil")
		return
	}

	// verifies left split panel exists
	if contentPane.Leading == nil {
		t.Errorf("Expected left split (Leading) to be non-nil")
	}

	// verifies right split panel exists
	if contentPane.Trailing == nil {
		t.Errorf("Expected right split (Trailing) to be non-nil")
	}
}
