// package mapping handles displaying the map of the jobs
package mapping

import (
	"fmt"
	"job-visualizer/pkg/shared"
	"net/http"

	"github.com/morikuni/go-geoplot"
	"github.com/skratchdot/open-golang/open"
)

// MappingService handles the creation and serving of interactive job location maps
type MappingService struct {
	geoplotMap *geoplot.Map
}

// NewMappingService creates a new instance of the MappingService
func NewMappingService() *MappingService {
	return &MappingService{}
}

// GenerateMap creates an interactive map from job data and starts the HTTP server, entry point
func (ms *MappingService) GenerateMap(jobs []shared.JobData, windowData *shared.GuiWindowData) []shared.JobData {
	ms.geoplotMap = ms.createGeoplotMap(jobs)
	windowData.Server = ms.createHttpServer(windowData)
	openWebpage()
	return jobs
}

// mapPage serves the main map page HTML
func (ms *MappingService) mapPage(writer http.ResponseWriter, request *http.Request) {
	if ms.geoplotMap != nil {
		writer.Header().Set("Content-Type", "text/html")
		_, err := fmt.Fprint(writer, `
            <html>
                <head>
                    <title>job-visualizer Map</title>
                </head>
                <body style="margin:0;padding:0;">
                    <iframe src="/innerMap" style="width:100vw;height:100vh;border:none;"></iframe>
                </body>
            </html>
        `)
		shared.CheckErrorWarn(err)
	} else {
		http.Error(writer, "Map not ready", http.StatusServiceUnavailable)
	}
}

// innerMap serves the actual interactive map
func (ms *MappingService) innerMap(writer http.ResponseWriter, request *http.Request) {
	if ms.geoplotMap != nil {
		err := geoplot.ServeMap(writer, request, ms.geoplotMap)
		shared.CheckErrorWarn(err)
	} else {
		http.Error(writer, "Map not ready", http.StatusServiceUnavailable)
	}
}

// createGeoplotMap creates a complete interactive map with job location markers
func (ms *MappingService) createGeoplotMap(jobs []shared.JobData) *geoplot.Map {
	geoplotMap := ms.createBaseMap()
	ms.createMarkers(jobs, geoplotMap)

	return geoplotMap
}

// createBaseMap creates the base map configuration centered on Boston
func (ms *MappingService) createBaseMap() *geoplot.Map {
	boston := &geoplot.LatLng{
		Latitude:  42.361145,
		Longitude: -71.057083,
	}
	geoplotMap := &geoplot.Map{
		Center: boston,
		Zoom:   7,
		Area: &geoplot.Area{
			From: boston.Offset(-0.1, -0.1),
			To:   boston.Offset(0.2, 0.2),
		},
	}
	return geoplotMap
}

// createMarkers adds job location markers to the map, grouping multiple jobs at same location
func (ms *MappingService) createMarkers(jobs []shared.JobData, geoplotMap *geoplot.Map) {
	commonLocations := make(map[shared.LatLong][]shared.JobData)
	for _, job := range jobs {
		if _, ok := commonLocations[job.LatLong]; ok {
			commonLocations[job.LatLong] = append(commonLocations[job.LatLong], job)
		} else {
			commonLocations[job.LatLong] = append(make([]shared.JobData, 0), job)
		}
	}
	for key, value := range commonLocations {
		latitude := key.Latitude
		longitude := key.Longitude
		coordinates := &geoplot.LatLng{
			Latitude:  latitude,
			Longitude: longitude,
		}
		icon := geoplot.ColorIcon(255, 255, 0)
		geoplotMap.AddMarker(&geoplot.Marker{
			LatLng:  coordinates,
			Popup:   displayDescription(value),
			Tooltip: displayHoverword(value),
			Icon:    icon,
		})
	}
}

// displayHoverword creates hover text for map markers
func displayHoverword(markerJobs []shared.JobData) string {
	hoverword := ""
	jobLength := len(markerJobs)
	switch jobLength {
	case 0:
		hoverword = "No jobs available at this location."
	case 1:
		hoverword = markerJobs[0].CompanyName
	case 2:
		hoverword = fmt.Sprintf("%s and %s", markerJobs[0].CompanyName, markerJobs[1].CompanyName)
	case 3:
		hoverword = fmt.Sprintf("%s, %s and %s", markerJobs[0].CompanyName, markerJobs[1].CompanyName,
			markerJobs[2].CompanyName)
	default:
		hoverword = fmt.Sprintf("%s, %s, %s and %d more", markerJobs[0].CompanyName, markerJobs[1].CompanyName,
			markerJobs[2].CompanyName, jobLength-3)
	}
	return hoverword
}

// displayDescription creates popup content for map markers showing job details
func displayDescription(markerJobs []shared.JobData) string {
	description := ""
	for i, job := range markerJobs {
		if i > 10 {
			description += fmt.Sprintf(" ...and %d more jobs at this location.\n", len(markerJobs)-10)
			break
		} else {
			description += fmt.Sprintf("Company Name: %s\nJob Title: %s\n\n", job.CompanyName, job.JobTitle)
		}
	}
	return description
}

// createHttpServer sets up HTTP server for map serving
func (ms *MappingService) createHttpServer(windowData *shared.GuiWindowData) *http.Server {
	if windowData.Server != nil {
		err := windowData.Server.Close()
		shared.CheckErrorWarn(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/map", ms.mapPage)
	mux.HandleFunc("/innerMap", ms.innerMap)

	server := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		err := server.ListenAndServe()
		shared.CheckErrorWarn(err)
	}()
	return server
}

// openWebpage opens the map in the user's default browser
func openWebpage() {
	url := "http://localhost:8080/map"
	err := open.Run(url)
	shared.CheckErrorWarn(err)
}
