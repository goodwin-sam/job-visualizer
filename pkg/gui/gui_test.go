package gui

import (
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/shared"
	"testing"
)

func TestProcessJobsReturnsSlice(t *testing.T) {
	tempDir := t.TempDir()
	programData := shared.ProgramData{
		OutputDirectory: tempDir,
	}
	mappingService := mapping.NewMappingService()
	result := processJobs(programData, nil, mappingService)
	if result == nil {
		t.Fatal("Expected processJobs to return a non-nil slice, but got nil")
	}
}
