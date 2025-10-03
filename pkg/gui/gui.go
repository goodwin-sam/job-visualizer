// package gui provides gui operations
package gui

import (
	"fmt"
	"job-visualizer/pkg/gui/build"
	"job-visualizer/pkg/mapping"
	"job-visualizer/pkg/processor"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	_ "modernc.org/sqlite"
)

// RunGUIorHeadless launches either GUI or headless mode based on user preference
func RunGUIorHeadless(programData shared.ProgramData, headless bool) {
	if headless {
		workingDirectory, err := os.Getwd()
		shared.CheckErrorWarn(err)
		files, err := os.ReadDir(workingDirectory)
		shared.CheckErrorWarn(err)
		var inputFiles []string
		for _, file := range files {
			if !file.IsDir() {
				extension := filepath.Ext(file.Name())
				if extension == ".xlsx" || extension == ".xls" {
					filePath := filepath.Join(workingDirectory, file.Name())
					inputFiles = append(inputFiles, filePath)
				}
			}

		}
		programData.InputFiles = inputFiles
		programData.OutputDirectory = workingDirectory

		allJobData := processor.ProcessJobs(programData, nil, mapping.NewMappingService())
		for i, job := range allJobData {
			if i%100 == 0 {
				fmt.Printf("%-4s | %-25s | %-55s | %-25s\n",
					"#", "Location", "Job Title", "Company Name")
				fmt.Println(strings.Repeat("-", 120))
			}
			fmt.Printf("%-4d | %-25s | %-55s | %-25s\n",
				i+1, job.Location, job.JobTitle, job.CompanyName)
		}
	} else {
		createGuiApp(programData)
	}
}

// createGuiApp initializes and runs the GUI application
func createGuiApp(programData shared.ProgramData) {
	application := app.NewWithID("job-visualizer")
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	windowData := &shared.GuiWindowData{}
	mappingService := mapping.NewMappingService()
	startWindow := createGuiWindow(application, "job-visualizer")
	startButton := widget.NewButton("Start Application", func() {
		go func() {
			allJobData := processor.ProcessJobs(programData, progressBar, mappingService)
			fyne.DoAndWait(func() {
				mainWindow := createGuiWindow(application, "job-visualizer")
				mainWindow.SetOnClosed(func() { application.Quit() })
				mainWindow = build.BuildMainWindow(mainWindow, allJobData, windowData, mappingService)
				startWindow.Hide()
				mainWindow.Show()
			})
		}()
	})
	startWindow = build.BuildStartWindow(startWindow, startButton, progressBar, &programData)
	startWindow.ShowAndRun()
}

// createGuiWindow creates a new Fyne window
func createGuiWindow(app fyne.App, title string) fyne.Window {
	Window := app.NewWindow(title)
	Window.Resize(fyne.NewSize(1000, 600))
	return Window
}
