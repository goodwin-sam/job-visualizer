package arguments

import (
	"flag"
)

func CheckCLIArguments() bool {
	headless := flag.Bool("headless", false, "Run in headless mode (no GUI)")
	flag.Parse()
	return *headless
}
