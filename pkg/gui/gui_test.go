// package gui provides tests for GUI functionality and job processing
package gui

import (
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/shared"
	"testing"
)

// TestProcessJobs tests that the processJobs function returns a valid slice
func TestProcessJobs(t *testing.T) {
	// sets up test environment with temporary directory
	tempDir := t.TempDir()
	programData := shared.ProgramData{
		OutputDirectory: tempDir,
	}
	mappingService := mapping.NewMappingService()

	// tests processJobs function with minimal data
	result := processJobs(programData, nil, mappingService)

	// verifies function returns a non-nil slice
	if result == nil {
		t.Fatal("Expected processJobs to return a non-nil slice, but got nil")
	}
}
