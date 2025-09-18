package excel

import (
	"job-visualizer/pkg/shared"

	"github.com/xuri/excelize/v2"
)

func OpenExcelFile(inputFiles []string) []*excelize.File {
	var allFiles []*excelize.File
	for _, filePath := range inputFiles {
		file, err := excelize.OpenFile(filePath)
		shared.CheckError(err)
		allFiles = append(allFiles, file)
	}
	return allFiles
}

func GetAllRows(files []*excelize.File) [][]string {
	var allRows [][]string
	for _, file := range files {
		rows, err := file.GetRows("Jobs")
		shared.CheckError(err)
		allRows = append(allRows, rows[1:]...)
	}
	return allRows
}
