// package mapping provides tests for geographic mapping and visualization functionality
package mapping

import (
	"job-visualizer/pkg/shared"
	"strconv"
	"testing"

	"github.com/morikuni/go-geoplot"
)

// TestDisplayHoverword tests the displayHoverword function with various job counts
func TestDisplayHoverword(t *testing.T) {
	tests := []struct {
		name     string
		jobs     []shared.JobData
		expected string
	}{
		{
			// tests empty job list handling
			name:     "No jobs",
			jobs:     []shared.JobData{},
			expected: "No jobs available at this location.",
		},
		{
			// tests single job display
			name:     "One job",
			jobs:     []shared.JobData{newJob("CompanyA", "")},
			expected: "CompanyA",
		},
		{
			// tests two jobs with "and" conjunction
			name:     "Two jobs",
			jobs:     []shared.JobData{newJob("CompanyA", ""), newJob("CompanyB", "")},
			expected: "CompanyA and CompanyB",
		},
		{
			// tests three jobs with comma and "and" conjunction
			name:     "Three jobs",
			jobs:     []shared.JobData{newJob("CompanyA", ""), newJob("CompanyB", ""), newJob("CompanyC", "")},
			expected: "CompanyA, CompanyB and CompanyC",
		},
		{
			// tests four jobs with truncation after third job
			name:     "Four jobs",
			jobs:     []shared.JobData{newJob("CompanyA", ""), newJob("CompanyB", ""), newJob("CompanyC", ""), newJob("CompanyD", "")},
			expected: "CompanyA, CompanyB, CompanyC and 1 more",
		},
		{
			// tests five jobs with truncation after third job
			name:     "Five jobs",
			jobs:     []shared.JobData{newJob("A", ""), newJob("B", ""), newJob("C", ""), newJob("D", ""), newJob("E", "")},
			expected: "A, B, C and 2 more",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := displayHoverword(tt.jobs)
			if result != tt.expected {
				t.Errorf("displayHoverword() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// expectedDescriptionForJobs generates expected description text for job lists
func expectedDescriptionForJobs(jobs []shared.JobData) string {
	desc := ""
	for i, job := range jobs {
		if i > 10 {
			desc += " ...and " + strconv.Itoa(len(jobs)-10) + " more jobs at this location.\n"
			break
		}
		desc += "Company Name: " + job.CompanyName + "\nJob Title: " + job.JobTitle + "\n\n"
	}
	return desc
}

// TestDisplayDescription tests the displayDescription function with various job counts
func TestDisplayDescription(t *testing.T) {
	tests := []struct {
		name     string
		jobs     []shared.JobData
		expected string
	}{
		{
			// tests empty job list handling
			name:     "No jobs",
			jobs:     []shared.JobData{},
			expected: expectedDescriptionForJobs([]shared.JobData{}),
		},
		{
			// tests single job description formatting
			name:     "One job",
			jobs:     []shared.JobData{newJob("CompanyA", "Engineer")},
			expected: expectedDescriptionForJobs([]shared.JobData{newJob("CompanyA", "Engineer")}),
		},
		{
			// tests multiple job descriptions
			name:     "Two jobs",
			jobs:     []shared.JobData{newJob("A", "X"), newJob("B", "Y")},
			expected: expectedDescriptionForJobs([]shared.JobData{newJob("A", "X"), newJob("B", "Y")}),
		},
		{
			// tests truncation at 11 jobs (shows first 10 + count)
			name:     "Eleven jobs (truncation)",
			jobs:     makeJobs(11),
			expected: expectedDescriptionForJobs(makeJobs(11)),
		},
		{
			// tests truncation at 12 jobs (shows first 10 + count)
			name:     "Twelve jobs (truncation)",
			jobs:     makeJobs(12),
			expected: expectedDescriptionForJobs(makeJobs(12)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := displayDescription(tt.jobs)
			if result != tt.expected {
				t.Errorf("displayDescription() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// assertMapNotNil verifies that a map object is not nil
func assertMapNotNil(t *testing.T, m *geoplot.Map) {
	if m == nil {
		t.Fatal("Map is nil")
	}
}

// TestCreateBaseMap tests the creation of a base map with correct center and zoom settings
func TestCreateBaseMap(t *testing.T) {
	ms := NewMappingService()
	m := ms.createBaseMap()

	// verifies map was created successfully
	assertMapNotNil(t, m)

	// verifies map center is set to Boston coordinates
	if m.Center == nil || m.Center.Latitude != 42.361145 || m.Center.Longitude != -71.057083 {
		t.Errorf("Center = %+v, want Latitude=42.361145, Longitude=-71.057083", m.Center)
	}

	// verifies zoom level is set correctly
	if m.Zoom != 7 {
		t.Errorf("Zoom = %d, want 7", m.Zoom)
	}

	// verifies map area boundaries are defined
	if m.Area == nil {
		t.Error("Area is nil")
	} else {
		from := m.Area.From
		to := m.Area.To
		if from == nil || to == nil {
			t.Error("Area.From or Area.To is nil")
		}
	}
}

// TestCreateGeoplotMap tests the creation of a complete geoplot map with job data
func TestCreateGeoplotMap(t *testing.T) {
	ms := NewMappingService()
	jobs := []shared.JobData{}

	// tests map creation with empty job list
	m := ms.createGeoplotMap(jobs)

	// verifies map was created successfully
	assertMapNotNil(t, m)
}

// newJob creates a test job with specified company and title
func newJob(company, title string) shared.JobData {
	return shared.JobData{CompanyName: company, JobTitle: title}
}

// makeJobs creates a slice of n test jobs with sequential company and title names
func makeJobs(n int) []shared.JobData {
	jobs := make([]shared.JobData, n)
	for i := range n {
		jobs[i] = shared.JobData{CompanyName: "C" + string(rune('A'+i)), JobTitle: "T" + string(rune('A'+i))}
	}
	return jobs
}
