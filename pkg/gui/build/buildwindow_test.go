package build

import (
	"testing"

	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestBuildStartWindow(t *testing.T) {
	window := test.NewWindow(nil)
	startButton := widget.NewButton("Start", nil)
	progressBar := widget.NewProgressBar()

	resultWindow := BuildStartWindow(window, startButton, progressBar)

	if resultWindow != window {
		t.Errorf("Expected returned window to be the same as input window")
	}

	if window.Content() == nil {
		t.Errorf("Expected window content to be set, got nil")
	}
}

func TestBuildMainWindow(t *testing.T) {
	window := test.NewWindow(nil)
	jobs := []shared.JobData{} // empty slice for simplicity

	resultWindow := BuildMainWindow(window, jobs)

	if resultWindow != window {
		t.Errorf("Expected returned window to be the same as input window")
	}

	if window.Content() == nil {
		t.Errorf("Expected window content to be set, got nil")
	}
}

func TestBuildMainContent(t *testing.T) {
	jobs := []shared.JobData{}

	contentPane := buildMainContent(jobs)

	if contentPane == nil {
		t.Errorf("Expected contentPane to be non-nil")
		return
	}

	if contentPane.Leading == nil {
		t.Errorf("Expected left split (Leading) to be non-nil")
	}

	if contentPane.Trailing == nil {
		t.Errorf("Expected right split (Trailing) to be non-nil")
	}
}
