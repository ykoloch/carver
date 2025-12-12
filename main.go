package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	DEFAULT_OUTPUT_DIR = "./recovered"
	DEFAULT_WORKERS    = 4
)

var (
	device, output string
	verbose        bool
	workers        uint
)

func main() {
	output = *flag.String("output", DEFAULT_OUTPUT_DIR, fmt.Sprintf("path to the output directory; default is %q", DEFAULT_OUTPUT_DIR))
	device = *flag.String("device", "", "path to the device to be scanned")
	if len(device) < 1 {
		os.Stderr.WriteString("error: the device should be specified")
		os.Exit(-1)
	}
	verbose = *flag.Bool("verbose", false, "verbose mode")
	workers = *flag.Uint("workers", DEFAULT_WORKERS, "number of workers")

	flag.Parse()

	println(output)
}
