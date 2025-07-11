package main

import (
	"job-visualizer/pkg/gui"
	"job-visualizer/pkg/headless"
	"job-visualizer/pkg/shared"
	"os"
)

func main() {
	resourceDirectory, err := os.Getwd()
	shared.CheckError(err)
	shared.Program.ResourcesDirectory = resourceDirectory
	isHeadless := headless.CheckCLIArguments()
	gui.RunGUIorHeadless(isHeadless)
}
