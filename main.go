package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	device, output string
	verbose        bool
	workers        uint
)

func main() {
	flag.StringVar(&output, "output", DEFAULT_OUTPUT_DIR, "path to the output directory")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.UintVar(&workers, "workers", DEFAULT_WORKERS, "number of workers")
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

	if err := scan(device); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}
}
