package excel

import (
	"job-visualizer/pkg/shared"
	"path/filepath"
	"slices"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestOpenExcelFile(t *testing.T) {
	tempDirectory := t.TempDir()
	testFile1 := excelize.NewFile()
	sheet := "Jobs"
	_, err := testFile1.NewSheet(sheet)
	shared.CheckErrorWarn(err)
	_ = testFile1.SetSheetRow(sheet, "A1", &[]string{"Location", "Job Title", "Company Name"})
	_ = testFile1.SetSheetRow(sheet, "A2", &[]string{"New York", "Software Engineer", "Tech Corp"})

	testFile2 := excelize.NewFile()
	_, err = testFile2.NewSheet(sheet)
	shared.CheckErrorWarn(err)
	_ = testFile2.SetSheetRow(sheet, "A1", &[]string{"Location", "Job Title", "Company Name"})
	_ = testFile2.SetSheetRow(sheet, "A2", &[]string{"San Francisco", "Data Scientist", "Data Inc."})
	_ = testFile2.SetSheetRow(sheet, "A3", &[]string{"Los Angeles", "Product Manager", "Creative Solutions"})

	testFile1Path := filepath.Join(tempDirectory, "test1.xlsx")
	testFile2Path := filepath.Join(tempDirectory, "test2.xlsx")
	err = testFile1.SaveAs(testFile1Path)
	if err != nil {
		t.Fatalf("Failed to save test file 1: %v", err)
	}
	err = testFile2.SaveAs(testFile2Path)
	if err != nil {
		t.Fatalf("Failed to save test file 2: %v", err)
	}
	shared.Program.InputFiles = []string{testFile1Path, testFile2Path}

	files := OpenExcelFile()

	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}
	for i, file := range files {
		sheets := file.GetSheetList()
		if !slices.Contains(sheets, sheet) {
			t.Errorf("File %d does not contain expected sheet '%s'", i+1, sheet)
		}
	}
}

func TestGetAllRows(t *testing.T) {
	file := excelize.NewFile()
	sheet := "Jobs"
	_, err := file.NewSheet(sheet)
	shared.CheckErrorWarn(err)
	_ = file.SetSheetRow(sheet, "A1", &[]string{"Location", "Job Title", "Company Name"})
	_ = file.SetSheetRow(sheet, "A2", &[]string{"New York", "Software Engineer", "Tech Corp"})
	_ = file.SetSheetRow(sheet, "A3", &[]string{"San Francisco", "Data Scientist", "Data Inc."})
	_ = file.SetSheetRow(sheet, "A4", &[]string{"Los Angeles", "Product Manager", "Creative Solutions"})
	files := []*excelize.File{file}
	rows := GetAllRows(files)

	expected := [][]string{
		{"New York", "Software Engineer", "Tech Corp"},
		{"San Francisco", "Data Scientist", "Data Inc."},
		{"Los Angeles", "Product Manager", "Creative Solutions"},
	}
	if len(rows) != len(expected) {
		t.Fatalf("Expected %d rows, got %d", len(expected), len(rows))
	}
	for i, row := range rows {
		if len(row) != len(expected[i]) {
			t.Fatalf("Row %d: Expected %d columns, got %d", i+1, len(expected[i]), len(row))
		}
		for j, cell := range row {
			if cell != expected[i][j] {
				t.Errorf("Row %d, Column %d: Expected '%s', got '%s'", i+1, j+1, expected[i][j], cell)
			}
		}
	}
}
