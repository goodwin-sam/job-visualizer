package filter

import (
	"job-visualizer/pkg/shared"
	"reflect"
	"testing"
)

func TestFilterWorkFromHome(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		expect bool
	}{
		{
			name:   "Work from home is Yes",
			job:    shared.JobData{WorkFromHome: "Yes"},
			expect: true,
		},
		{
			name:   "Work from home is No",
			job:    shared.JobData{WorkFromHome: "No"},
			expect: false,
		},
		{
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

func TestFilterKeyword(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		filter string
		expect bool
	}{
		{
			name:   "Match in JobTitle",
			job:    shared.JobData{JobTitle: "Software Engineer"},
			filter: "engineer",
			expect: true,
		},
		{
			name:   "Match in CompanyName",
			job:    shared.JobData{CompanyName: "Tech Corp"},
			filter: "tech",
			expect: true,
		},
		{
			name:   "Match in Description",
			job:    shared.JobData{Description: "Great team environment"},
			filter: "team",
			expect: true,
		},
		{
			name:   "Match in Qualifications",
			job:    shared.JobData{Qualifications: "Bachelor's degree required"},
			filter: "bachelor",
			expect: true,
		},
		{
			name:   "No match",
			job:    shared.JobData{JobTitle: "Manager", CompanyName: "Retail Inc", Description: "Sales", Qualifications: "Experience"},
			filter: "developer",
			expect: false,
		},
		{
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

func TestFilterLocation(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		filter string
		expect bool
	}{
		{
			name:   "Exact match",
			job:    shared.JobData{Location: "New York"},
			filter: "New York",
			expect: true,
		},
		{
			name:   "Partial match",
			job:    shared.JobData{Location: "San Francisco"},
			filter: "Francisco",
			expect: true,
		},
		{
			name:   "No match",
			job:    shared.JobData{Location: "London"},
			filter: "Paris",
			expect: false,
		},
		{
			name:   "Case insensitivity",
			job:    shared.JobData{Location: "Berlin"},
			filter: "berLIN",
			expect: true,
		},
		{
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

func TestFilterMinSalary(t *testing.T) {
	cases := []struct {
		name   string
		job    shared.JobData
		filter string
		expect bool
	}{
		{
			name:   "Salary above minimum",
			job:    shared.JobData{Salary: 60000},
			filter: "50000",
			expect: true,
		},
		{
			name:   "Salary equal to minimum",
			job:    shared.JobData{Salary: 50000},
			filter: "50000",
			expect: false,
		},
		{
			name:   "Salary below minimum",
			job:    shared.JobData{Salary: 40000},
			filter: "50000",
			expect: false,
		},
		{
			name:   "Invalid filter input",
			job:    shared.JobData{Salary: 60000},
			filter: "notanumber",
			expect: false,
		},
		{
			name:   "Negative minimum salary",
			job:    shared.JobData{Salary: 1000},
			filter: "-1000",
			expect: true,
		},
		{
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

func TestFilterIndividualJob(t *testing.T) {
	job := makeJob("Engineer", "Tech", "Remote work", "BS", "New York", 100000, "Yes")

	cases := []struct {
		name    string
		filters shared.FilterEntries
		expect  bool
	}{
		{
			name:    "Keyword match",
			filters: shared.FilterEntries{KeywordEntry: "engineer"},
			expect:  true,
		},
		{
			name:    "Location match",
			filters: shared.FilterEntries{LocationEntry: "new york"},
			expect:  true,
		},
		{
			name:    "Min salary match",
			filters: shared.FilterEntries{MinSalaryEntry: "90000"},
			expect:  true,
		},
		{
			name:    "Work from home match",
			filters: shared.FilterEntries{WorkFromHomeEntry: true},
			expect:  true,
		},
		{
			name:    "No match (wrong keyword)",
			filters: shared.FilterEntries{KeywordEntry: "manager"},
			expect:  false,
		},
		{
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
			name:      "No filters",
			filters:   shared.FilterEntries{},
			expectIdx: []int{0, 1, 2},
		},
		{
			name:      "Keyword filter (Engineer)",
			filters:   shared.FilterEntries{KeywordEntry: "engineer"},
			expectIdx: []int{0},
		},
		{
			name:      "Work from home filter",
			filters:   shared.FilterEntries{WorkFromHomeEntry: true},
			expectIdx: []int{0, 2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			shared.WindowData.Filters = c.filters
			result := FilterJobs(jobs)
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
