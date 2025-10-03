// package processor provides tests for the job data processing pipeline
package processor

import (
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/shared"
	"testing"

	"fyne.io/fyne/v2/widget"
)

// TestProcessJobs tests the ProcessJobs function with mock data
func TestProcessJobs(t *testing.T) {
	// creates test program data
	programData := shared.ProgramData{
		InputFiles:      []string{},
		OutputDirectory: t.TempDir(),
		CacheDirectory:  t.TempDir(),
	}

	// creates mock progress bar
	progressBar := widget.NewProgressBar()

	// creates mapping service
	mappingService := mapping.NewMappingService()

	// tests ProcessJobs with empty input files
	result := ProcessJobs(programData, progressBar, mappingService)

	// verifies function returns a valid slice (even if empty)
	if result == nil {
		t.Error("Expected ProcessJobs to return a valid slice, got nil")
	}

	// verifies function handles empty input gracefully
	if len(result) != 0 {
		t.Errorf("Expected empty result with no input files, got %d jobs", len(result))
	}
}

// TestProcessJobsWithoutProgressBar tests ProcessJobs without a progress bar
func TestProcessJobsWithoutProgressBar(t *testing.T) {
	// creates test program data
	programData := shared.ProgramData{
		InputFiles:      []string{},
		OutputDirectory: t.TempDir(),
		CacheDirectory:  t.TempDir(),
	}

	// creates mapping service
	mappingService := mapping.NewMappingService()

	// tests ProcessJobs without progress bar
	result := ProcessJobs(programData, nil, mappingService)

	// verifies function returns a valid slice
	if result == nil {
		t.Error("Expected ProcessJobs to return a valid slice, got nil")
	}
}
