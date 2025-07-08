package gui

import (
	"testing"

	"fyne.io/fyne/v2/app"
)

func TestCreateGuiWindow(t *testing.T) {
	application = app.New()
	title := "test window"
	window := createGuiWindow(title)

	if window == nil {
		t.Fatal("Expected window to be created, but got nil")
	}
	if window.Title() != title {
		t.Errorf("Expected window title '%s', but got '%s'", title, window.Title())
	}
}
