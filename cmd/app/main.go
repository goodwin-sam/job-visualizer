package main

import (
	"fmt"
	"job-visualizer/pkg/arguments"
	"job-visualizer/pkg/gui"
	"job-visualizer/pkg/shared"
	"os"
	"path/filepath"
)

func main() {
	shared.Program.CacheDirectory = getAppCacheDirectory()
	isHeadless := arguments.CheckCLIArguments()
	gui.RunGUIorHeadless(isHeadless)
}

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
