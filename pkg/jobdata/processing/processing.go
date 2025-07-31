package processing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"io"
	"job-visualizer/pkg/shared"
	"net/http"
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func ProcessLatLongs(jobs []shared.JobData, progressBar *widget.ProgressBar) []shared.JobData {
	cacheFilename := filepath.Join(shared.Program.CacheDirectory, "cached_locations.json")
	cachedLocations := make(map[string]shared.LatLong)
	loadCacheFromFile(cacheFilename, cachedLocations)
	jobs = standardizeLocations(jobs)
	cacheLatLongs(jobs, cachedLocations, progressBar)
	jobs = assignLatLongs(jobs, cachedLocations)
	saveCacheToFile(cacheFilename, cachedLocations)
	return jobs
}

func loadCacheFromFile(filename string, cachedLocations map[string]shared.LatLong) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(filename)
	shared.CheckErrorWarn(err)
	defer func() {
		err = file.Close()
		shared.CheckErrorWarn(err)
	}()
	dec := json.NewDecoder(file)
	err = dec.Decode(&cachedLocations)
	shared.CheckErrorWarn(err)
}

func standardizeLocations(jobs []shared.JobData) []shared.JobData {
	reNumbers := regexp.MustCompile(`[0-9]+`)
	rePunctuation := regexp.MustCompile(`[^-\w\s]`)
	for i, job := range jobs {
		location := job.Location
		location = strings.ToLower(location)
		location = strings.ReplaceAll(location, "-", " ")
		location = reNumbers.ReplaceAllString(location, "")
		location = rePunctuation.ReplaceAllString(location, "")
		location = strings.Join(strings.Fields(location), " ")
		jobs[i].StandardizedLocation = location
	}
	return jobs
}

func cacheLatLongs(jobs []shared.JobData, cachedLocations map[string]shared.LatLong, progressBar *widget.ProgressBar) {
	for i, job := range jobs {
		fmt.Printf("\rcaching job locations (%d/%d)", i+1, len(jobs))
		if progressBar != nil {
			fyne.Do(func() {
				progressBar.SetValue(float64(i+1) / float64(len(jobs)))
			})
		}
		if _, ok := cachedLocations[job.StandardizedLocation]; ok {
		} else {
			responseBody := getNominatimResponse(job.StandardizedLocation)

			var locations []shared.JsonLocation
			err := json.Unmarshal(responseBody, &locations)
			shared.CheckErrorWarn(err)
			if len(locations) > 0 {
				addLocationToCache(job.StandardizedLocation, locations, cachedLocations)
			}
		}
	}
	fmt.Println()
}

func assignLatLongs(jobs []shared.JobData, cachedLocations map[string]shared.LatLong) []shared.JobData {
	for i, job := range jobs {
		if coordinates, ok := cachedLocations[job.StandardizedLocation]; ok {
			jobs[i].LatLong = coordinates
		}
	}
	return jobs
}

func saveCacheToFile(filename string, cachedLocations map[string]shared.LatLong) {
	file, err := os.Create(filename)
	shared.CheckErrorWarn(err)
	defer func() {
		err = file.Close()
		shared.CheckErrorWarn(err)
	}()
	enc := json.NewEncoder(file)
	err = enc.Encode(cachedLocations)
	shared.CheckErrorWarn(err)
}

func getNominatimResponse(location string) []byte {
	encodedLocation := url.QueryEscape(location)
	apiUrl := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", encodedLocation)
	response, err := http.Get(apiUrl)
	shared.CheckErrorWarn(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		shared.CheckErrorWarn(err)
	}(response.Body)
	responseBody, err := io.ReadAll(response.Body)
	shared.CheckErrorWarn(err)
	return responseBody
}

func addLocationToCache(jobLocation string, locations []shared.JsonLocation, cachedLocations map[string]shared.LatLong) {
	latitude, err := strconv.ParseFloat(locations[0].Lat, 64)
	shared.CheckErrorWarn(err)
	longitude, err := strconv.ParseFloat(locations[0].Lon, 64)
	shared.CheckErrorWarn(err)
	cachedLocations[jobLocation] = shared.LatLong{
		Latitude:  latitude,
		Longitude: longitude,
	}
}
