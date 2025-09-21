// package filter provides job filtering functionality based on user criteria
package filter

import (
	"job-visualizer/pkg/shared"
	"strconv"
	"strings"
)

// FilterJobs applies user filters to a list of jobs and returns matching results
func FilterJobs(jobs []shared.JobData, filters shared.FilterEntries) []shared.JobData {
	if filters.KeywordEntry != "" || filters.LocationEntry != "" || filters.MinSalaryEntry != "" ||
		filters.WorkFromHomeEntry {
		var filteredJobs []shared.JobData
		for _, job := range jobs {
			filteredJobs = filterIndividualJob(job, filteredJobs, filters)
		}
		return filteredJobs
	}
	return jobs
}

// filterIndividualJob checks if a single job matches all applied filters
func filterIndividualJob(job shared.JobData, filteredJobs []shared.JobData, filters shared.FilterEntries) []shared.JobData {
	filterMatch := true
	if filters.KeywordEntry != "" {
		//fmt.Printf("keyword entered: %s", filters.KeywordEntry)
		filterMatch = filterKeyword(job, filters.KeywordEntry)
	}
	if filters.LocationEntry != "" && filterMatch {
		//fmt.Printf("location entered: %s", filters.LocationEntry)
		filterMatch = filterLocation(job, filters.LocationEntry)
	}
	if filters.MinSalaryEntry != "" && filterMatch {
		//fmt.Printf("min salary entered: %s", filters.MinSalaryEntry)
		filterMatch = filterMinSalary(job, filters.MinSalaryEntry)
	}
	if filters.WorkFromHomeEntry && filterMatch {
		//fmt.Println("work from home filter applied")
		filterMatch = filterWorkFromHome(job)
	}
	if filterMatch {
		filteredJobs = append(filteredJobs, job)
	}
	return filteredJobs
}

// filterKeyword searches for keywords in job title, company, description, and qualifications
func filterKeyword(job shared.JobData, filterInput string) bool {
	filterMatch := false
	filter := strings.ToLower(filterInput)
	jobTitle := strings.ToLower(job.JobTitle)
	companyName := strings.ToLower(job.CompanyName)
	description := strings.ToLower(job.Description)
	qualifications := strings.ToLower(job.Qualifications)
	if strings.Contains(jobTitle, filter) || strings.Contains(companyName, filter) ||
		strings.Contains(description, filter) || strings.Contains(qualifications, filter) {
		filterMatch = true
	}
	return filterMatch
}

// filterLocation searches for location keywords in job location field
func filterLocation(job shared.JobData, filterInput string) bool {
	filterMatch := false
	jobLocation := strings.ToLower(job.Location)
	filter := strings.ToLower(filterInput)
	if strings.Contains(jobLocation, filter) {
		filterMatch = true
	}
	return filterMatch
}

// filterMinSalary checks if job salary meets minimum salary requirement
func filterMinSalary(job shared.JobData, filter string) bool {
	filterMatch := false
	salary := job.Salary
	minSalary, err := strconv.Atoi(filter)
	if err != nil {
		shared.CheckErrorWarn(err)
		return false
	}
	if salary > minSalary {
		filterMatch = true
	}
	return filterMatch
}

// filterWorkFromHome checks if job offers remote work option
func filterWorkFromHome(job shared.JobData) bool {
	return job.WorkFromHome == "Yes"
}
