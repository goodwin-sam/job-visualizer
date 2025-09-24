// package buildwidgets provides widget creation and event handling for the GUI
package buildwidgets

import (
	"job-visualizer/pkg/jobsprocessing/filter"
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/shared"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// BuildMainButtons creates the main window buttons for job filtering and display
func BuildMainButtons(jobs []shared.JobData, windowData *shared.GuiWindowData, mappingService interface{}) (*widget.Button, *widget.Button, *widget.Button) {
	refreshButton := widget.NewButton("Click to refresh list of jobs to original", func() {
		handleJobRefresh(jobs, windowData, mappingService)
	})
	filterButton := widget.NewButton("Click to filter the jobs", func() {
		handleJobFilter(jobs, windowData, mappingService)
	})
	selectedDetailsButton := widget.NewButton("Click to display selected job details", func() {
		if windowData.ListWidget != nil {
			windowData.ListWidget.Refresh()
		}
		if windowData.DetailsWidget != nil {
			windowData.DetailsWidget.SetText(windowData.SelectedJobDetails)
		}
	})

	return refreshButton, filterButton, selectedDetailsButton
}

// BuildStartButtons creates the start window buttons for selection and navigation
func BuildStartButtons(window fyne.Window, inputFileLabel *widget.Label, outputDirectoryLabel *widget.Label, programData *shared.ProgramData) (*widget.Button, *widget.Button, *widget.Button) {
	inputFileButton := widget.NewButton("Select Input Files", func() {
		inputFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			shared.CheckErrorWarn(err)
			if reader == nil {
				println("user cancelled file selection")
				return
			}
			defer func() {
				err = reader.Close()
				shared.CheckErrorWarn(err)
			}()
			programData.InputFiles = append(programData.InputFiles, reader.URI().Path())
			selectedFiles := strings.Join(programData.InputFiles, "\n")
			inputFileLabel.SetText(selectedFiles)

		}, window)
		inputFileDialog.Show()
	})
	outputDirectoryButton := widget.NewButton("Select output directory", func() {
		outputDirectoryDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			shared.CheckErrorWarn(err)
			if uri == nil {
				println("user cancelled directory selection")
				return
			}
			programData.OutputDirectory = uri.Path()
			outputDirectoryLabel.SetText(programData.OutputDirectory)
		}, window)
		outputDirectoryDialog.Show()
	})
	quitButton := BuildQuitButton()
	return inputFileButton, outputDirectoryButton, quitButton
}

// BuildLabel creates a styled label widget with specified formatting
func BuildLabel(text string, boldBool bool, italicBool bool) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignCenter,
		fyne.TextStyle{Bold: boldBool, Italic: italicBool})
}

// BuildRemoteCheckbox creates a checkbox for filtering remote work jobs
func BuildRemoteCheckbox(windowData *shared.GuiWindowData) *widget.Check {
	remoteCheckbox := widget.NewCheck("Remote Work: check for yes, uncheck for all", func(checked bool) {
		if checked {
			windowData.Filters.WorkFromHomeEntry = true
		} else {
			windowData.Filters.WorkFromHomeEntry = false
		}
	})
	return remoteCheckbox
}

// BuildQuitButton creates a quit button that exits the application
func BuildQuitButton() *widget.Button {
	return widget.NewButton("Quit", func() { fyne.CurrentApp().Quit() })
}

// handleJobRefresh clears all filters and refreshes the job list
func handleJobRefresh(jobs []shared.JobData, windowData *shared.GuiWindowData, mappingService interface{}) {
	removeActiveFilters(windowData)
	filteredJobs := filter.FilterJobs(jobs, windowData.Filters)
	if ms, ok := mappingService.(*mapping.MappingService); ok {
		ms.GenerateMap(filteredJobs, windowData)
	}
	windowData.FilteredJobs = &filteredJobs
	refreshEntries(windowData)
	if windowData.ListWidget != nil {
		windowData.ListWidget.Refresh()
	}
	if windowData.DetailsWidget != nil {
		windowData.DetailsWidget.SetText("Select a job to display details")
	}
}

// handleJobFilter applies current filter settings to the job list
func handleJobFilter(jobs []shared.JobData, windowData *shared.GuiWindowData, mappingService interface{}) {
	filteredJobs := filter.FilterJobs(jobs, windowData.Filters)
	if ms, ok := mappingService.(*mapping.MappingService); ok {
		ms.GenerateMap(filteredJobs, windowData)
	}
	windowData.FilteredJobs = &filteredJobs
	if windowData.ListWidget != nil {
		windowData.ListWidget.Refresh()
	}
	if windowData.DetailsWidget != nil {
		windowData.DetailsWidget.SetText("Select a job to display details")
	}
}

// removeActiveFilters clears all filter entries from the window data
func removeActiveFilters(windowData *shared.GuiWindowData) {
	windowData.Filters.KeywordEntry = ""
	windowData.Filters.LocationEntry = ""
	windowData.Filters.MinSalaryEntry = ""
	windowData.Filters.WorkFromHomeEntry = false
}

// refreshEntries resets all filter input fields to empty state
func refreshEntries(windowData *shared.GuiWindowData) {
	windowData.KeywordEntryWidget.SetText("")
	windowData.LocationEntryWidget.SetText("")
	windowData.MinSalaryEntryWidget.SetText("")
	windowData.KeywordEntryWidget.SetPlaceHolder("Enter keyword filter here")
	windowData.LocationEntryWidget.SetPlaceHolder("Enter location filter here")
	windowData.MinSalaryEntryWidget.SetPlaceHolder("Enter minimum salary filter here")
}
