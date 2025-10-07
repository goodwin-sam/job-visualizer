// package processor provides the main job data processing pipeline orchestration
package processor

import (
	"job-visualizer/pkg/database"
	"job-visualizer/pkg/excel"
	"job-visualizer/pkg/jobsprocessing"
	"job-visualizer/pkg/jobsprocessing/geocoding"
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2/widget"
)

// ProcessJobs handles the complete job data processing pipeline
func ProcessJobs(programData shared.ProgramData, progressBar *widget.ProgressBar, mappingService *mapping.MappingService) []shared.JobData {
	files := excel.OpenExcelFile(programData.InputFiles)
	rows := excel.GetAllRows(files)
	allJobData := jobsprocessing.ProcessRows(rows, []shared.JobData{})
	if progressBar != nil {
		allJobData = geocoding.ProcessLatLongs(allJobData, programData.CacheDirectory, progressBar)
	} else {
		allJobData = geocoding.ProcessLatLongs(allJobData, programData.CacheDirectory, nil)
	}
	allJobData = mappingService.GenerateMap(allJobData, &shared.GuiWindowData{})
	jobsDatabase := database.CreateDatabase(programData.OutputDirectory)
	database.SetupDatabase(jobsDatabase)
	database.WriteToDatabase(jobsDatabase, allJobData)
	return allJobData
}
