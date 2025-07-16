package gui

import (
	"job-visualizer/pkg/shared"
	"testing"
)

func TestProcessJobsReturnsSlice(t *testing.T) {
	tempDir := t.TempDir()
	shared.Program.OutputDirectory = tempDir
	shared.Program.ResourcesDirectory = tempDir
	result := processJobs(nil)
	if result == nil {
		t.Fatal("Expected processJobs to return a non-nil slice, but got nil")
	}
}
