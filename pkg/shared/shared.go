// package shared contains core data structures and utilities used throughout the application
package shared

import (
	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne/v2/widget"
)

var Program ProgramData

// JobData represents a single job posting with all relevant information
type JobData struct {
	Location             string
	StandardizedLocation string
	JobTitle             string
	CompanyName          string
	Description          string
	DatePosted           string
	Salary               int
	WorkFromHome         string
	Qualifications       string
	Links                string
	Country              string
	LatLong              LatLong
}

// LatLong represents geographic coordinates, used for job location
type LatLong struct {
	Latitude  float64
	Longitude float64
}

// ProgramData holds input and output file paths and cache directory
type ProgramData struct {
	InputFiles      []string
	OutputDirectory string
	CacheDirectory  string
}

// GuiWindowData manages GUI state and widget references
type GuiWindowData struct {
	ListWidget           *widget.List
	KeywordEntryWidget   *widget.Entry
	LocationEntryWidget  *widget.Entry
	MinSalaryEntryWidget *widget.Entry
	DetailsWidget        *widget.Label
	FilteredJobs         *[]JobData
	SelectedJobDetails   string
	Filters              FilterEntries
	Server               *http.Server
}

// FilterEntries stores user filter criteria from the GUI
type FilterEntries struct {
	KeywordEntry      string
	LocationEntry     string
	MinSalaryEntry    string
	WorkFromHomeEntry bool
}

// JsonLocation is bridge between the API response and the LatLong struct
type JsonLocation struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// CheckError terminates the program if an error occurs
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CheckErrorWarn prints an error but continues program execution
func CheckErrorWarn(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
