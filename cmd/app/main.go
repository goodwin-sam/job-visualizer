// package main provides the entry point for the job visualizer application.
// processes job data from Excel files with both GUI and headless modes.
//
//	go run ./cmd/app                    # GUI mode
//	go run ./cmd/app --headless         # headless mode
package main

import (
	"fmt"
	"job-visualizer/pkg/arguments"
	"job-visualizer/pkg/gui"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
)

// main initializes the application and launches the appropriate interface.
func main() {
	programData := shared.ProgramData{
		CacheDirectory: getAppCacheDirectory(),
	}
	isHeadless := arguments.CheckCLIArguments()
	gui.RunGUIorHeadless(programData, isHeadless)
}

// getAppCacheDirectory creates the cache directory locally in the user's home directory as .job-visualizer
// falls back to current directory if home directory access fails.
func getAppCacheDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return "."
	}

	cacheDir := filepath.Join(homeDir, ".job-visualizer")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "."
	}

	return cacheDir
}
