package processing

import (
	"path/filepath"
	"testing"

	"job-visualizer/pkg/shared"
)

func TestStandardizeLocations(t *testing.T) {
	jobs := []shared.JobData{
		{Location: "123 Main St., New York!"},
		{Location: "456 Elm St, San Francisco."},
		{Location: "789 Broadway, New-York"},
	}

	expected := []string{
		"main st new york",
		"elm st san francisco",
		"broadway new york",
	}

	standardized := standardizeLocations(jobs)
	for i, job := range standardized {
		if job.StandardizedLocation != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], job.StandardizedLocation)
		}
	}
}

func TestAssignLatLongs(t *testing.T) {
	jobs := []shared.JobData{
		{StandardizedLocation: "new york"},
		{StandardizedLocation: "san francisco"},
	}

	cache := map[string]shared.LatLong{
		"new york":      {Latitude: 40.7128, Longitude: -74.0060},
		"san francisco": {Latitude: 37.7749, Longitude: -122.4194},
	}

	expected := []shared.LatLong{
		{Latitude: 40.7128, Longitude: -74.0060},
		{Latitude: 37.7749, Longitude: -122.4194},
	}

	result := assignLatLongs(jobs, cache)
	for i, job := range result {
		if job.LatLong != expected[i] {
			t.Errorf("Job %d: expected LatLong %+v, got %+v", i, expected[i], job.LatLong)
		}
	}
}

func TestAddLocationToCache(t *testing.T) {
	cache := make(map[string]shared.LatLong)
	locations := []shared.JsonLocation{
		{Lat: "51.5074", Lon: "-0.1278"}, // London
	}
	addLocationToCache("london", locations, cache)

	expected := shared.LatLong{Latitude: 51.5074, Longitude: -0.1278}
	if val, ok := cache["london"]; !ok {
		t.Errorf("Expected cache to have key 'london'")
	} else if val != expected {
		t.Errorf("Expected %+v, got %+v", expected, val)
	}
}

func TestSaveAndLoadCacheToFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test_cache.json")

	originalCache := map[string]shared.LatLong{
		"paris": {Latitude: 48.8566, Longitude: 2.3522},
		"tokyo": {Latitude: 35.6895, Longitude: 139.6917},
	}

	saveCacheToFile(filename, originalCache)

	loadedCache := make(map[string]shared.LatLong)
	loadCacheFromFile(filename, loadedCache)

	if len(originalCache) != len(loadedCache) {
		t.Fatalf("Expected cache length %d, got %d", len(originalCache), len(loadedCache))
	}
	for k, v := range originalCache {
		if loaded, ok := loadedCache[k]; !ok {
			t.Errorf("Key %s missing in loaded cache", k)
		} else if loaded != v {
			t.Errorf("For key %s, expected %+v, got %+v", k, v, loaded)
		}
	}
}
