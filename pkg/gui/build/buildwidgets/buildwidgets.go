package buildwidgets

import (
	"job-visualizer/pkg/jobdata/filter"
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/shared"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func BuildMainButtons(jobs []shared.JobData) (*widget.Button, *widget.Button, *widget.Button) {
	refreshButton := widget.NewButton("Click to refresh list of jobs to original", func() {
		handleJobRefresh(jobs)
	})
	filterButton := widget.NewButton("Click to filter the jobs", func() {
		handleJobFilter(jobs)
	})
	selectedDetailsButton := widget.NewButton("Click to display selected job details", func() {
		shared.WindowData.DetailsWidget.SetText(shared.WindowData.SelectedJobDetails)
	})

	return refreshButton, filterButton, selectedDetailsButton
}

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

func BuildLabel(text string, boldBool bool, italicBool bool) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignCenter,
		fyne.TextStyle{Bold: boldBool, Italic: italicBool})
}

func BuildRemoteCheckbox() *widget.Check {
	remoteCheckbox := widget.NewCheck("Remote Work: check for yes, uncheck for all", func(checked bool) {
		if checked {
			shared.WindowData.Filters.WorkFromHomeEntry = true
		} else {
			shared.WindowData.Filters.WorkFromHomeEntry = false
		}
	})
	return remoteCheckbox
}

func BuildQuitButton() *widget.Button {
	return widget.NewButton("Quit", func() { fyne.CurrentApp().Quit() })
}

func handleJobRefresh(jobs []shared.JobData) {
	removeActiveFilters()
	filteredJobs := filter.FilterJobs(jobs)
	mapping.GenerateMap(filteredJobs)
	shared.WindowData.FilteredJobs = &filteredJobs
	refreshEntries()
}

func handleJobFilter(jobs []shared.JobData) {
	filteredJobs := filter.FilterJobs(jobs)
	mapping.GenerateMap(filteredJobs)
	shared.WindowData.FilteredJobs = &filteredJobs
}

func removeActiveFilters() {
	shared.WindowData.Filters.KeywordEntry = ""
	shared.WindowData.Filters.LocationEntry = ""
	shared.WindowData.Filters.MinSalaryEntry = ""
	shared.WindowData.Filters.WorkFromHomeEntry = false
}

func refreshEntries() {
	shared.WindowData.KeywordEntryWidget.SetText("")
	shared.WindowData.LocationEntryWidget.SetText("")
	shared.WindowData.MinSalaryEntryWidget.SetText("")
	shared.WindowData.KeywordEntryWidget.SetPlaceHolder("Enter keyword filter here")
	shared.WindowData.LocationEntryWidget.SetPlaceHolder("Enter location filter here")
	shared.WindowData.MinSalaryEntryWidget.SetPlaceHolder("Enter minimum salary filter here")
}
