package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	device, output string
)

func main() {
	flag.StringVar(&output, "output", DEFAULT_OUTPUT_DIR, "path to the output directory")
	flag.StringVar(&device, "device", "", "path to the device to be scanned")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Carver - File recovery tool\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\nExample:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -device /dev/sdb1 -output ./recovered\n", os.Args[0])
	}

	flag.Parse()
	if device == "" {
		flag.Usage()
		os.Exit(-1)
	}

	err := os.MkdirAll(output, 0755)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "can not create output directory %s: %v\n", output, err)
		os.Exit(-1)
	}

	if err := scan(device); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
