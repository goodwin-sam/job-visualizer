package mapping

import (
	"job-visualizer/pkg/shared"
	"strconv"
	"testing"

	"github.com/morikuni/go-geoplot"
)

func TestDisplayHoverword(t *testing.T) {
	tests := []struct {
		name     string
		jobs     []shared.JobData
		expected string
	}{
		{
			name:     "No jobs",
			jobs:     []shared.JobData{},
			expected: "No jobs available at this location.",
		},
		{
			name:     "One job",
			jobs:     []shared.JobData{newJob("CompanyA", "")},
			expected: "CompanyA",
		},
		{
			name:     "Two jobs",
			jobs:     []shared.JobData{newJob("CompanyA", ""), newJob("CompanyB", "")},
			expected: "CompanyA and CompanyB",
		},
		{
			name:     "Three jobs",
			jobs:     []shared.JobData{newJob("CompanyA", ""), newJob("CompanyB", ""), newJob("CompanyC", "")},
			expected: "CompanyA, CompanyB and CompanyC",
		},
		{
			name:     "Four jobs",
			jobs:     []shared.JobData{newJob("CompanyA", ""), newJob("CompanyB", ""), newJob("CompanyC", ""), newJob("CompanyD", "")},
			expected: "CompanyA, CompanyB, CompanyC and 1 more",
		},
		{
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

func TestDisplayDescription(t *testing.T) {
	tests := []struct {
		name     string
		jobs     []shared.JobData
		expected string
	}{
		{
			name:     "No jobs",
			jobs:     []shared.JobData{},
			expected: expectedDescriptionForJobs([]shared.JobData{}),
		},
		{
			name:     "One job",
			jobs:     []shared.JobData{newJob("CompanyA", "Engineer")},
			expected: expectedDescriptionForJobs([]shared.JobData{newJob("CompanyA", "Engineer")}),
		},
		{
			name:     "Two jobs",
			jobs:     []shared.JobData{newJob("A", "X"), newJob("B", "Y")},
			expected: expectedDescriptionForJobs([]shared.JobData{newJob("A", "X"), newJob("B", "Y")}),
		},
		{
			name:     "Eleven jobs (truncation)",
			jobs:     makeJobs(11),
			expected: expectedDescriptionForJobs(makeJobs(11)),
		},
		{
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

func assertMapNotNil(t *testing.T, m *geoplot.Map) {
	if m == nil {
		t.Fatal("Map is nil")
	}
}

func TestCreateBaseMap(t *testing.T) {
	m := createBaseMap()
	assertMapNotNil(t, m)
	if m.Center == nil || m.Center.Latitude != 42.361145 || m.Center.Longitude != -71.057083 {
		t.Errorf("Center = %+v, want Latitude=42.361145, Longitude=-71.057083", m.Center)
	}
	if m.Zoom != 7 {
		t.Errorf("Zoom = %d, want 7", m.Zoom)
	}
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

func TestCreateGeoplotMap(t *testing.T) {
	jobs := []shared.JobData{}
	m := createGeoplotMap(jobs)
	assertMapNotNil(t, m)
}

func newJob(company, title string) shared.JobData {
	return shared.JobData{CompanyName: company, JobTitle: title}
}

func makeJobs(n int) []shared.JobData {
	jobs := make([]shared.JobData, n)
	for i := range n {
		jobs[i] = shared.JobData{CompanyName: "C" + string(rune('A'+i)), JobTitle: "T" + string(rune('A'+i))}
	}
	return jobs
}
