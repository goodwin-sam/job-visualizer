// package excel handles reading and parsing Excel files containing job data
package excel

import (
	"job-visualizer/pkg/shared"

	"github.com/xuri/excelize/v2"
)

// OpenExcelFile opens multiple Excel files and returns file handles
func OpenExcelFile(inputFiles []string) []*excelize.File {
	var allFiles []*excelize.File
	for _, filePath := range inputFiles {
		file, err := excelize.OpenFile(filePath)
		shared.CheckError(err)
		allFiles = append(allFiles, file)
	}
	return allFiles
}

// GetAllRows extracts all data rows from the "Jobs" worksheet, skipping the header
func GetAllRows(files []*excelize.File) [][]string {
	var allRows [][]string
	for _, file := range files {
		rows, err := file.GetRows("Jobs")
		shared.CheckError(err)
		allRows = append(allRows, rows[1:]...)
	}
	return allRows
}
