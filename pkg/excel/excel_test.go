package excel

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestGetAllRows(t *testing.T) {
	file := excelize.NewFile()
	sheet := "Jobs"
	file.NewSheet(sheet)
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
