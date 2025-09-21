// package buildcontainers provides container layout and organization for GUI components
package buildcontainers

import (
	"job-visualizer/pkg/shared"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// buildKeywordContainer creates the keyword filter input and button
func buildKeywordContainer(windowData *shared.GuiWindowData) *fyne.Container {
	windowData.KeywordEntryWidget = widget.NewEntry()
	windowData.KeywordEntryWidget.SetPlaceHolder("Enter keyword filter here")
	keywordButton := widget.NewButton("Click to apply keyword", func() {
		windowData.Filters.KeywordEntry = windowData.KeywordEntryWidget.Text
	})
	keywordContainer := container.NewGridWithColumns(2, windowData.KeywordEntryWidget, keywordButton)
	return keywordContainer
}

// buildLocationContainer creates the location filter input and button
func buildLocationContainer(windowData *shared.GuiWindowData) *fyne.Container {
	windowData.LocationEntryWidget = widget.NewEntry()
	windowData.LocationEntryWidget.SetPlaceHolder("Enter location filter here")
	locationButton := widget.NewButton("Click to apply location", func() {
		windowData.Filters.LocationEntry = windowData.LocationEntryWidget.Text
	})
	locationContainer := container.NewGridWithColumns(2, windowData.LocationEntryWidget, locationButton)
	return locationContainer
}

// buildMinSalaryContainer creates the minimum salary filter input and button
func buildMinSalaryContainer(windowData *shared.GuiWindowData) *fyne.Container {
	windowData.MinSalaryEntryWidget = widget.NewEntry()
	windowData.MinSalaryEntryWidget.SetPlaceHolder("Enter minimum salary filter here")
	minSalaryButton := widget.NewButton("Click to apply minimum salary", func() {
		windowData.Filters.MinSalaryEntry = windowData.MinSalaryEntryWidget.Text
	})
	minSalaryContainer := container.NewGridWithColumns(2, windowData.MinSalaryEntryWidget, minSalaryButton)
	return minSalaryContainer
}
