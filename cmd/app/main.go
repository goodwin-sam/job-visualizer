package main

import (
	"job-visualizer/pkg/arguments"
	"job-visualizer/pkg/gui"
	"job-visualizer/pkg/shared"
	"os"
)

func main() {
	resourceDirectory, err := os.Getwd()
	shared.CheckError(err)
	shared.Program.ResourcesDirectory = resourceDirectory
	isHeadless := arguments.CheckCLIArguments()
	gui.RunGUIorHeadless(isHeadless)
}
