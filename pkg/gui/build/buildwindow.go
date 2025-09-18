package build

import (
	"job-visualizer/pkg/gui/build/buildcontainers"
	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildStartWindow(window fyne.Window, startButton *widget.Button, progressBar *widget.ProgressBar, programData *shared.ProgramData) fyne.Window {
	startContainer := buildcontainers.BuildStartContainer(window, startButton, progressBar, programData)
	window.SetContent(startContainer)
	return window
}

func BuildMainWindow(window fyne.Window, jobs []shared.JobData, windowData *shared.GuiWindowData, mappingService interface{}) fyne.Window {
	contentPane := buildMainContent(jobs, windowData, mappingService)
	window.SetContent(contentPane)
	return window
}

func buildMainContent(jobs []shared.JobData, windowData *shared.GuiWindowData, mappingService interface{}) *container.Split {
	leftSplit := buildcontainers.BuildLeftSplit(jobs, windowData, mappingService)
	rightSplit := buildcontainers.BuildRightSplit(windowData)
	contentPane := container.NewHSplit(leftSplit, rightSplit)

	return contentPane
}
