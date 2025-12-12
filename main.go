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
	flag.StringVar(&output, "output", DEFAULT_OUTPUT_DIR, fmt.Sprintf("path to the output directory; default is %q", DEFAULT_OUTPUT_DIR))
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.UintVar(&workers, "workers", DEFAULT_WORKERS, fmt.Sprintf("number of workers; default is %d", DEFAULT_WORKERS))
	flag.StringVar(&device, "device", "", "path to the device to be scanned")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Carver - File recovery tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -device /dev/sdb1 -output ./recovered\n", os.Args[0])
	}

	flag.Parse()
	if device == "" {
		flag.Usage()
		os.Exit(-1)
	}
}
