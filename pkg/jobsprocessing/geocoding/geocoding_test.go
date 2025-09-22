// package geocoding provides tests for geocoding and location processing functionality
package geocoding

import (
	"path/filepath"
	"testing"

	"job-visualizer/pkg/shared"
)

// TestStandardizeLocations tests the standardizeLocations function with various location formats
func TestStandardizeLocations(t *testing.T) {
	// creates test jobs with various location formats
	jobs := []shared.JobData{
		{Location: "123 Main St., New York!"},
		{Location: "456 Elm St, San Francisco."},
		{Location: "789 Broadway, New-York"},
	}

	// defines expected standardized location strings
	expected := []string{
		"main st new york",
		"elm st san francisco",
		"broadway new york",
	}

	// tests location standardization
	standardized := standardizeLocations(jobs)
	for i, job := range standardized {
		if job.StandardizedLocation != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], job.StandardizedLocation)
		}
	}
}

// TestAssignLatLongs tests the assignLatLongs function with cached location data
func TestAssignLatLongs(t *testing.T) {
	// creates test jobs with standardized locations
	jobs := []shared.JobData{
		{StandardizedLocation: "new york"},
		{StandardizedLocation: "san francisco"},
	}

	// creates test cache with known coordinates
	cache := map[string]shared.LatLong{
		"new york":      {Latitude: 40.7128, Longitude: -74.0060},
		"san francisco": {Latitude: 37.7749, Longitude: -122.4194},
	}

	// defines expected latitude/longitude values
	expected := []shared.LatLong{
		{Latitude: 40.7128, Longitude: -74.0060},
		{Latitude: 37.7749, Longitude: -122.4194},
	}

	// tests latitude/longitude assignment from cache
	result := assignLatLongs(jobs, cache)
	for i, job := range result {
		if job.LatLong != expected[i] {
			t.Errorf("Job %d: expected LatLong %+v, got %+v", i, expected[i], job.LatLong)
		}
	}
}

// TestAddLocationToCache tests the addLocationToCache function with JSON location data
func TestAddLocationToCache(t *testing.T) {
	// creates empty cache and test location data
	cache := make(map[string]shared.LatLong)
	locations := []shared.JsonLocation{
		{Lat: "51.5074", Lon: "-0.1278"}, // London coordinates
	}

	// tests adding location to cache
	addLocationToCache("london", locations, cache)

	// verifies location was added to cache correctly
	expected := shared.LatLong{Latitude: 51.5074, Longitude: -0.1278}
	if val, ok := cache["london"]; !ok {
		t.Errorf("Expected cache to have key 'london'")
	} else if val != expected {
		t.Errorf("Expected %+v, got %+v", expected, val)
	}
}

// TestSaveAndLoadCacheToFile tests cache persistence by saving and loading cache data
func TestSaveAndLoadCacheToFile(t *testing.T) {
	// sets up temporary directory for test file
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test_cache.json")

	// creates test cache with sample location data
	originalCache := map[string]shared.LatLong{
		"paris": {Latitude: 48.8566, Longitude: 2.3522},
		"tokyo": {Latitude: 35.6895, Longitude: 139.6917},
	}

	// tests saving cache to file
	saveCacheToFile(filename, originalCache)

	// tests loading cache from file
	loadedCache := make(map[string]shared.LatLong)
	loadCacheFromFile(filename, loadedCache)

	// verifies cache data integrity after save/load cycle
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
