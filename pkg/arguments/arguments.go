// package arguments handles command-line argument parsing for the application.
package arguments

import (
	"flag"
)

// CheckCLIArguments parses command-line flags and returns whether headless mode is requested.
func CheckCLIArguments() bool {
	headless := flag.Bool("headless", false, "Run in headless mode (no GUI)")
	flag.Parse()
	return *headless
}
