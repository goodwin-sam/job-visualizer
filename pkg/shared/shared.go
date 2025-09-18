package shared

import (
	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var Program ProgramData
var StartWindow fyne.Window
var MainWindow fyne.Window

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
type LatLong struct {
	Latitude  float64
	Longitude float64
}

type ProgramData struct {
	InputFiles      []string
	OutputDirectory string
	CacheDirectory  string
}

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

type FilterEntries struct {
	KeywordEntry      string
	LocationEntry     string
	MinSalaryEntry    string
	WorkFromHomeEntry bool
}

type JsonLocation struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckErrorWarn(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
