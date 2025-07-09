package processing

import (
	"job-visualizer/pkg/shared"
	"testing"
)

func TestStandardizeLocations(t *testing.T) {
	tests := []struct {
		name     string
		jobs     []shared.JobData
		expected []shared.JobData
	}{
		{
			name: "basic location standardization",
			jobs: []shared.JobData{
				{Location: "Boston, MA"},
				{Location: "New York, NY"},
			},
			expected: []shared.JobData{
				{Location: "Boston, MA", StandardizedLocation: "boston ma"},
				{Location: "New York, NY", StandardizedLocation: "new york ny"},
			},
		},
		{
			name: "remove numbers and punctuation",
			jobs: []shared.JobData{
				{Location: "Boston123, MA!@#"},
				{Location: "New-York_456, NY$%^"},
			},
			expected: []shared.JobData{
				{Location: "Boston123, MA!@#", StandardizedLocation: "boston ma"},
				{Location: "New-York_456, NY$%^", StandardizedLocation: "newyork_ ny"},
			},
		},
		{
			name: "handle extra spaces",
			jobs: []shared.JobData{
				{Location: "  Boston   ,    MA  "},
				{Location: "New    York,   NY"},
			},
			expected: []shared.JobData{
				{Location: "  Boston   ,    MA  ", StandardizedLocation: "boston ma"},
				{Location: "New    York,   NY", StandardizedLocation: "new york ny"},
			},
		},
		{
			name: "empty location",
			jobs: []shared.JobData{
				{Location: ""},
			},
			expected: []shared.JobData{
				{Location: "", StandardizedLocation: ""},
			},
		},
		{
			name: "mixed case with numbers and punctuation",
			jobs: []shared.JobData{
				{Location: "SaN FrAnCiScO, CA 94102!"},
			},
			expected: []shared.JobData{
				{Location: "SaN FrAnCiScO, CA 94102!", StandardizedLocation: "san francisco ca"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := standardizeLocations(tt.jobs)
			if len(result) != len(tt.expected) {
				t.Errorf("standardizeLocations() returned %d jobs, expected %d", len(result), len(tt.expected))
				return
			}
			for i, job := range result {
				if job.StandardizedLocation != tt.expected[i].StandardizedLocation {
					t.Errorf("standardizeLocations() job %d StandardizedLocation = %q, expected %q", 
						i, job.StandardizedLocation, tt.expected[i].StandardizedLocation)
				}
			}
		})
	}
}

func TestAssignLatLongs(t *testing.T) {
	tests := []struct {
		name            string
		jobs            []shared.JobData
		cachedLocations map[string]shared.LatLong
		expected        []shared.JobData
	}{
		{
			name: "assign existing coordinates",
			jobs: []shared.JobData{
				{StandardizedLocation: "boston ma"},
				{StandardizedLocation: "new york ny"},
			},
			cachedLocations: map[string]shared.LatLong{
				"boston ma":    {Latitude: 42.3601, Longitude: -71.0589},
				"new york ny":  {Latitude: 40.7128, Longitude: -74.0060},
			},
			expected: []shared.JobData{
				{StandardizedLocation: "boston ma", LatLong: shared.LatLong{Latitude: 42.3601, Longitude: -71.0589}},
				{StandardizedLocation: "new york ny", LatLong: shared.LatLong{Latitude: 40.7128, Longitude: -74.0060}},
			},
		},
		{
			name: "missing coordinates",
			jobs: []shared.JobData{
				{StandardizedLocation: "boston ma"},
				{StandardizedLocation: "unknown location"},
			},
			cachedLocations: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3601, Longitude: -71.0589},
			},
			expected: []shared.JobData{
				{StandardizedLocation: "boston ma", LatLong: shared.LatLong{Latitude: 42.3601, Longitude: -71.0589}},
				{StandardizedLocation: "unknown location", LatLong: shared.LatLong{Latitude: 0, Longitude: 0}},
			},
		},
		{
			name: "empty cache",
			jobs: []shared.JobData{
				{StandardizedLocation: "boston ma"},
			},
			cachedLocations: map[string]shared.LatLong{},
			expected: []shared.JobData{
				{StandardizedLocation: "boston ma", LatLong: shared.LatLong{Latitude: 0, Longitude: 0}},
			},
		},
		{
			name: "empty jobs",
			jobs: []shared.JobData{},
			cachedLocations: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3601, Longitude: -71.0589},
			},
			expected: []shared.JobData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := assignLatLongs(tt.jobs, tt.cachedLocations)
			if len(result) != len(tt.expected) {
				t.Errorf("assignLatLongs() returned %d jobs, expected %d", len(result), len(tt.expected))
				return
			}
			for i, job := range result {
				if job.LatLong.Latitude != tt.expected[i].LatLong.Latitude || 
				   job.LatLong.Longitude != tt.expected[i].LatLong.Longitude {
					t.Errorf("assignLatLongs() job %d LatLong = %+v, expected %+v", 
						i, job.LatLong, tt.expected[i].LatLong)
				}
			}
		})
	}
}

func TestAddLocationToCache(t *testing.T) {
	tests := []struct {
		name            string
		jobLocation     string
		locations       []shared.JsonLocation
		initialCache    map[string]shared.LatLong
		expectedCache   map[string]shared.LatLong
	}{
		{
			name:        "add new location",
			jobLocation: "boston ma",
			locations: []shared.JsonLocation{
				{Lat: "42.3601", Lon: "-71.0589"},
			},
			initialCache: map[string]shared.LatLong{},
			expectedCache: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3601, Longitude: -71.0589},
			},
		},
		{
			name:        "add to existing cache",
			jobLocation: "new york ny",
			locations: []shared.JsonLocation{
				{Lat: "40.7128", Lon: "-74.0060"},
			},
			initialCache: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3601, Longitude: -71.0589},
			},
			expectedCache: map[string]shared.LatLong{
				"boston ma":   {Latitude: 42.3601, Longitude: -71.0589},
				"new york ny": {Latitude: 40.7128, Longitude: -74.0060},
			},
		},
		{
			name:        "overwrite existing location",
			jobLocation: "boston ma",
			locations: []shared.JsonLocation{
				{Lat: "42.3602", Lon: "-71.0590"},
			},
			initialCache: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3601, Longitude: -71.0589},
			},
			expectedCache: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3602, Longitude: -71.0590},
			},
		},
		{
			name:        "multiple locations use first",
			jobLocation: "boston ma",
			locations: []shared.JsonLocation{
				{Lat: "42.3601", Lon: "-71.0589"},
				{Lat: "42.3602", Lon: "-71.0590"},
			},
			initialCache: map[string]shared.LatLong{},
			expectedCache: map[string]shared.LatLong{
				"boston ma": {Latitude: 42.3601, Longitude: -71.0589},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cachedLocations := make(map[string]shared.LatLong)
			for k, v := range tt.initialCache {
				cachedLocations[k] = v
			}

			addLocationToCache(tt.jobLocation, tt.locations, cachedLocations)

			if len(cachedLocations) != len(tt.expectedCache) {
				t.Errorf("addLocationToCache() cache size = %d, expected %d", len(cachedLocations), len(tt.expectedCache))
				return
			}

			for key, expectedVal := range tt.expectedCache {
				if actualVal, exists := cachedLocations[key]; !exists {
					t.Errorf("addLocationToCache() missing key %q in cache", key)
				} else if actualVal.Latitude != expectedVal.Latitude || actualVal.Longitude != expectedVal.Longitude {
					t.Errorf("addLocationToCache() cache[%q] = %+v, expected %+v", key, actualVal, expectedVal)
				}
			}
		})
	}
}