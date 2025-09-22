// package filter provides tests for job filtering functionality
package filter

import (
	"job-visualizer/pkg/shared"
	"reflect"
	"testing"
)

// TestFilterWorkFromHome tests the filterWorkFromHome function with various work-from-home values
func TestFilterWorkFromHome(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		expect bool
	}{
		{
			// tests positive work-from-home case
			name:   "Work from home is Yes",
			job:    shared.JobData{WorkFromHome: "Yes"},
			expect: true,
		},
		{
			// tests negative work-from-home case
			name:   "Work from home is No",
			job:    shared.JobData{WorkFromHome: "No"},
			expect: false,
		},
		{
			// tests empty work-from-home case
			name:   "Work from home is empty",
			job:    shared.JobData{WorkFromHome: ""},
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := filterWorkFromHome(c.job)
			if result != c.expect {
				t.Errorf("%s: expected %v, got %v", c.name, c.expect, result)
			}
		})
	}
}

// TestFilterKeyword tests the filterKeyword function with various keyword matching scenarios
func TestFilterKeyword(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		filter string
		expect bool
	}{
		{
			// tests keyword matching in job title
			name:   "Match in JobTitle",
			job:    shared.JobData{JobTitle: "Software Engineer"},
			filter: "engineer",
			expect: true,
		},
		{
			// tests keyword matching in company name
			name:   "Match in CompanyName",
			job:    shared.JobData{CompanyName: "Tech Corp"},
			filter: "tech",
			expect: true,
		},
		{
			// tests keyword matching in job description
			name:   "Match in Description",
			job:    shared.JobData{Description: "Great team environment"},
			filter: "team",
			expect: true,
		},
		{
			// tests keyword matching in qualifications
			name:   "Match in Qualifications",
			job:    shared.JobData{Qualifications: "Bachelor's degree required"},
			filter: "bachelor",
			expect: true,
		},
		{
			// tests no match scenario across all fields
			name:   "No match",
			job:    shared.JobData{JobTitle: "Manager", CompanyName: "Retail Inc", Description: "Sales", Qualifications: "Experience"},
			filter: "developer",
			expect: false,
		},
		{
			// tests case-insensitive matching
			name:   "Case insensitivity",
			job:    shared.JobData{JobTitle: "Data Scientist"},
			filter: "data scientist",
			expect: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := filterKeyword(c.job, c.filter)
			if result != c.expect {
				t.Errorf("%s: expected %v, got %v", c.name, c.expect, result)
			}
		})
	}
}

// TestFilterLocation tests the filterLocation function with various location matching scenarios
func TestFilterLocation(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		filter string
		expect bool
	}{
		{
			// tests exact location matching
			name:   "Exact match",
			job:    shared.JobData{Location: "New York"},
			filter: "New York",
			expect: true,
		},
		{
			// tests partial location matching
			name:   "Partial match",
			job:    shared.JobData{Location: "San Francisco"},
			filter: "Francisco",
			expect: true,
		},
		{
			// tests no match scenario
			name:   "No match",
			job:    shared.JobData{Location: "London"},
			filter: "Paris",
			expect: false,
		},
		{
			// tests case-insensitive location matching
			name:   "Case insensitivity",
			job:    shared.JobData{Location: "Berlin"},
			filter: "berLIN",
			expect: true,
		},
		{
			// tests empty filter (should match all)
			name:   "Empty filter",
			job:    shared.JobData{Location: "Tokyo"},
			filter: "",
			expect: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := filterLocation(c.job, c.filter)
			if result != c.expect {
				t.Errorf("%s: expected %v, got %v", c.name, c.expect, result)
			}
		})
	}
}

// TestFilterMinSalary tests the filterMinSalary function with various salary comparison scenarios
func TestFilterMinSalary(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		filter string
		expect bool
	}{
		{
			// tests salary above minimum threshold
			name:   "Salary above minimum",
			job:    shared.JobData{Salary: 60000},
			filter: "50000",
			expect: true,
		},
		{
			// tests salary equal to minimum (should be excluded)
			name:   "Salary equal to minimum",
			job:    shared.JobData{Salary: 50000},
			filter: "50000",
			expect: false,
		},
		{
			// tests salary below minimum threshold
			name:   "Salary below minimum",
			job:    shared.JobData{Salary: 40000},
			filter: "50000",
			expect: false,
		},
		{
			// tests invalid filter input handling
			name:   "Invalid filter input",
			job:    shared.JobData{Salary: 60000},
			filter: "notanumber",
			expect: false,
		},
		{
			// tests negative minimum salary handling
			name:   "Negative minimum salary",
			job:    shared.JobData{Salary: 1000},
			filter: "-1000",
			expect: true,
		},
		{
			// tests zero minimum salary handling
			name:   "Zero minimum salary",
			job:    shared.JobData{Salary: 0},
			filter: "0",
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := filterMinSalary(c.job, c.filter)
			if result != c.expect {
				t.Errorf("%s: expected %v, got %v", c.name, c.expect, result)
			}
		})
	}
}

// makeJob creates a test job with specified field values
func makeJob(title, company, desc, qual, loc string, salary int, wfh string) shared.JobData {
	return shared.JobData{
		JobTitle:       title,
		CompanyName:    company,
		Description:    desc,
		Qualifications: qual,
		Location:       loc,
		Salary:         salary,
		WorkFromHome:   wfh,
	}
}

// TestFilterIndividualJob tests the filterIndividualJob function with various filter combinations
func TestFilterIndividualJob(t *testing.T) {
	job := makeJob("Engineer", "Tech", "Remote work", "BS", "New York", 100000, "Yes")

	cases := []struct {
		name    string
		filters shared.FilterEntries
		expect  bool
	}{
		{
			// tests keyword filter matching
			name:    "Keyword match",
			filters: shared.FilterEntries{KeywordEntry: "engineer"},
			expect:  true,
		},
		{
			// tests location filter matching
			name:    "Location match",
			filters: shared.FilterEntries{LocationEntry: "new york"},
			expect:  true,
		},
		{
			// tests minimum salary filter matching
			name:    "Min salary match",
			filters: shared.FilterEntries{MinSalaryEntry: "90000"},
			expect:  true,
		},
		{
			// tests work-from-home filter matching
			name:    "Work from home match",
			filters: shared.FilterEntries{WorkFromHomeEntry: true},
			expect:  true,
		},
		{
			// tests keyword filter no match
			name:    "No match (wrong keyword)",
			filters: shared.FilterEntries{KeywordEntry: "manager"},
			expect:  false,
		},
		{
			// tests salary filter no match
			name:    "No match (salary too high)",
			filters: shared.FilterEntries{MinSalaryEntry: "200000"},
			expect:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := filterIndividualJob(job, nil, c.filters)
			if c.expect && len(result) != 1 {
				t.Errorf("%s: expected job to be included, but it was not", c.name)
			}
			if !c.expect && len(result) != 0 {
				t.Errorf("%s: expected job to be excluded, but it was included", c.name)
			}
		})
	}
}

// TestFilterJobs tests the FilterJobs function with multiple jobs and various filter combinations
func TestFilterJobs(t *testing.T) {
	jobs := []shared.JobData{
		makeJob("Engineer", "Tech", "Remote work", "BS", "New York", 100000, "Yes"),
		makeJob("Manager", "Biz", "Office", "MBA", "San Francisco", 90000, "No"),
		makeJob("Analyst", "DataCorp", "Flexible", "BA", "London", 80000, "Yes"),
	}

	cases := []struct {
		name      string
		filters   shared.FilterEntries
		expectIdx []int
	}{
		{
			// tests filtering with no active filters (should return all jobs)
			name:      "No filters",
			filters:   shared.FilterEntries{},
			expectIdx: []int{0, 1, 2},
		},
		{
			// tests keyword filtering for specific job title
			name:      "Keyword filter (Engineer)",
			filters:   shared.FilterEntries{KeywordEntry: "engineer"},
			expectIdx: []int{0},
		},
		{
			// tests work-from-home filtering (should return remote jobs)
			name:      "Work from home filter",
			filters:   shared.FilterEntries{WorkFromHomeEntry: true},
			expectIdx: []int{0, 2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := FilterJobs(jobs, c.filters)
			var expected []shared.JobData
			for _, idx := range c.expectIdx {
				expected = append(expected, jobs[idx])
			}
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("%s: expected %v, got %v", c.name, expected, result)
			}
		})
	}
}
