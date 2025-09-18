// package jobsprocessing processes raw job data from Excel rows into structured JobData objects
package jobsprocessing

import (
	"job-visualizer/pkg/shared"
	"strconv"
)

// ProcessRows converts Excel rows into JobData structs, skipping invalid rows
func ProcessRows(rows [][]string, allJobData []shared.JobData) []shared.JobData {
	if len(rows) < 2 {
		return allJobData
	}
	for _, row := range rows[1:] {
		if len(row) < 10 {
			continue
		}
		job := shared.JobData{}
		job.CompanyName = row[0]
		job.DatePosted = row[1]
		// job.JobId = row[2]
		job.Country = row[3]
		job.Location = row[4]
		job.Salary = calcSalary(row)
		job.JobTitle = row[9]
		allJobData = append(allJobData, job)
	}
	return allJobData
}

// calcSalary calculates average salary, converting hourly to yearly if needed
func calcSalary(row []string) int {
	maxSalaryString := row[6]
	minSalaryString := row[7]
	hourlyOrYearly := row[8]

	maxSalary, err := strconv.ParseFloat(maxSalaryString, 64)
	shared.CheckError(err)
	minSalary, err := strconv.ParseFloat(minSalaryString, 64)
	shared.CheckError(err)
	salaryFloat := ((maxSalary + minSalary) / 2)
	if hourlyOrYearly == "hourly" {
		salaryFloat = salaryFloat * 40 * 50
	}
	salary := int(salaryFloat)
	return salary

}
