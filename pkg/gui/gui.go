package gui

import (
	"fmt"
	"job-visualizer/pkg/database"
	"job-visualizer/pkg/excel"
	"job-visualizer/pkg/gui/build"
	"job-visualizer/pkg/jobdata"
	"job-visualizer/pkg/jobdata/processing"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	_ "modernc.org/sqlite"
)

var application fyne.App

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

		allJobData := processJobs(programData, nil)
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

func createGuiApp(programData shared.ProgramData) {
	application = app.NewWithID("job-visualizer")
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	windowData := &shared.GuiWindowData{}
	startWindow := createGuiWindow("job-visualizer")
	startButton := widget.NewButton("Start Application", func() {
		go func() {
			allJobData := processJobs(programData, progressBar)
			fyne.DoAndWait(func() {
				mainWindow := createGuiWindow("job-visualizer")
				mainWindow.SetOnClosed(func() { application.Quit() })
				mainWindow = build.BuildMainWindow(mainWindow, allJobData, windowData)
				startWindow.Hide()
				mainWindow.Show()
			})
		}()
	})
	startWindow = build.BuildStartWindow(startWindow, startButton, progressBar, &programData)
	startWindow.ShowAndRun()
}

func createGuiWindow(title string) fyne.Window {
	Window := application.NewWindow(title)
	Window.Resize(fyne.NewSize(1000, 600))
	return Window
}

func processJobs(programData shared.ProgramData, progressBar *widget.ProgressBar) []shared.JobData {
	files := excel.OpenExcelFile(programData.InputFiles)
	rows := excel.GetAllRows(files)
	allJobData := jobdata.ProcessRows(rows, []shared.JobData{})
	if progressBar != nil {
		allJobData = processing.ProcessLatLongs(allJobData, programData.CacheDirectory, progressBar)
	} else {
		allJobData = processing.ProcessLatLongs(allJobData, programData.CacheDirectory, nil)
	}
	// Note: mapping.GenerateMap will be called from the GUI with windowData
	jobsDatabase := database.CreateDatabase(programData.OutputDirectory)
	database.SetupDatabase(jobsDatabase)
	database.WriteToDatabase(jobsDatabase, allJobData)
	return allJobData
}
