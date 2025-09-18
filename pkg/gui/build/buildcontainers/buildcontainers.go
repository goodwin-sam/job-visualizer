package buildcontainers

import (
	"fmt"
	"job-visualizer/pkg/gui/build/buildwidgets"
	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildStartContainer(window fyne.Window, startButton *widget.Button, progressBar *widget.ProgressBar, programData *shared.ProgramData) *fyne.Container {
	startLabel := buildwidgets.BuildLabel("Welcome to job-visualizer, choose your input files and output file location",
		true, false)
	inputFileLabel := buildwidgets.BuildLabel("No input files selected", false, false)
	outputDirectoryLabel := buildwidgets.BuildLabel("No output file selected", false, false)
	inputFileButton, outputDirectoryButton, quitButton := buildwidgets.BuildStartButtons(window, inputFileLabel, outputDirectoryLabel, programData)

	inputFilesBox := container.NewVBox(inputFileLabel, inputFileButton)
	outputDirectoryBox := container.NewVBox(outputDirectoryLabel, outputDirectoryButton)
	inputOutputContainers := container.NewHSplit(inputFilesBox, outputDirectoryBox)

	return container.NewVBox(startLabel, inputOutputContainers, startButton, quitButton)
}

func BuildLeftSplit(jobs []shared.JobData, windowData *shared.GuiWindowData) *container.Split {
	createJobList(windowData)
	refreshButton, filterButton, selectedDetailsButton := buildwidgets.BuildMainButtons(jobs, windowData)
	selectedDetailsLabel := buildwidgets.BuildLabel("Select a job and click button to display details",
		true, false)

	filterContainer := buildFilterContainer(windowData)

	jobScroll := container.NewScroll(windowData.ListWidget)
	filterVBox := container.NewVBox(refreshButton, filterContainer, filterButton)
	selectedDetailsContainer := container.NewBorder(
		selectedDetailsLabel,
		selectedDetailsButton,
		nil,
		nil,
		jobScroll)
	leftSplit := container.NewVSplit(filterVBox, selectedDetailsContainer)
	return leftSplit
}

func BuildRightSplit(windowData *shared.GuiWindowData) *fyne.Container {
	detailsLabel := buildwidgets.BuildLabel("Select a job to display details", true, false)
	detailsLabel.Wrapping = fyne.TextWrapWord
	windowData.DetailsWidget = detailsLabel
	rightPane := container.NewVBox(windowData.DetailsWidget, buildwidgets.BuildQuitButton())
	return rightPane
}

func createJobList(windowData *shared.GuiWindowData) {
	getDataLen := func() int {
		if windowData.FilteredJobs == nil {
			return 0
		}
		return len(*windowData.FilteredJobs)
	}

	updateListItem := func(itemNum widget.ListItemID, listItem fyne.CanvasObject) {
		itemName := (*windowData.FilteredJobs)[itemNum].CompanyName
		listItem.(*widget.Label).SetText(itemName)
	}
	windowData.ListWidget = widget.NewList(getDataLen, createListItem, updateListItem)
	windowData.ListWidget.OnSelected = func(i int) {
		windowData.SelectedJobDetails = formatJobDetails(i, *windowData)
	}
}

func createListItem() fyne.CanvasObject {
	return widget.NewLabel("list items here")
}

func formatJobDetails(i int, window shared.GuiWindowData) string {
	jobData := *window.FilteredJobs
	job := jobData[i]
	formattedDetails := fmt.Sprintf("Company Name:\n%s\n\nJob Title:\n%s\n\nLocation:\n%s\n\nDate Posted:"+
		"\n%s\n\nSalary:\n%d\n\nWork From Home:\n%s\n\nQualifications:\n%s\n\nLinks:\n%s\n\n",
		job.CompanyName, job.JobTitle, job.Location, job.DatePosted, job.Salary, job.WorkFromHome, job.Qualifications,
		job.Links)
	return formattedDetails
}

func buildFilterContainer(windowData *shared.GuiWindowData) *fyne.Container {
	keywordContainer := buildKeywordContainer(windowData)
	locationContainer := buildLocationContainer(windowData)
	minSalaryContainer := buildMinSalaryContainer(windowData)
	remoteCheckbox := buildwidgets.BuildRemoteCheckbox(windowData)
	filterContainer := container.NewVBox(keywordContainer, locationContainer, minSalaryContainer, remoteCheckbox)
	return filterContainer
}
