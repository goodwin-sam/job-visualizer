package filter

import (
	"job-visualizer/pkg/shared"
	"strconv"
	"strings"
)

func FilterJobs(jobs []shared.JobData) []shared.JobData {
	filters := shared.WindowData.Filters
	if filters.KeywordEntry != "" || filters.LocationEntry != "" || filters.MinSalaryEntry != "" ||
		filters.WorkFromHomeEntry {
		var filteredJobs []shared.JobData
		for _, job := range jobs {
			filteredJobs = filterIndividualJob(job, filteredJobs)
		}
		return filteredJobs
	}
	return jobs
}

func filterIndividualJob(job shared.JobData, filteredJobs []shared.JobData) []shared.JobData {
	filters := shared.WindowData.Filters
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

func filterLocation(job shared.JobData, filterInput string) bool {
	filterMatch := false
	jobLocation := strings.ToLower(job.Location)
	filter := strings.ToLower(filterInput)
	if strings.Contains(jobLocation, filter) {
		filterMatch = true
	}
	return filterMatch
}

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

func filterWorkFromHome(job shared.JobData) bool {
	return job.WorkFromHome == "Yes"
}
