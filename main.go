package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	device, output string
	enableJPEG     bool
	enablePNG      bool
	enablePDF      bool
	enableGIF      bool
	enableZIP      bool
)

func main() {
	flag.StringVar(&output, "output", DEFAULT_OUTPUT_DIR, "path to the output directory")
	flag.StringVar(&device, "device", "", "path to the device to be scanned")
	flag.BoolVar(&enableJPEG, "jpeg", false, "recover JPEG files")
	flag.BoolVar(&enablePNG, "png", false, "recover PNG files")
	flag.BoolVar(&enablePDF, "pdf", false, "recover PDF files")
	flag.BoolVar(&enableGIF, "gif", false, "recover GIF files")
	flag.BoolVar(&enableZIP, "zip", false, "recover ZIP/DOCX/XLSX files")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Carver - File recovery tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nFile Formats:\n")
		fmt.Fprintf(os.Stderr, "  If no format flags are specified, all formats will be recovered.\n")
		fmt.Fprintf(os.Stderr, "  Use format flags to recover only specific file types.\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -device /dev/sdb1 -output ./recovered\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -device /dev/sdb1 -jpeg -png\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -device ./disk.img -zip\n", os.Args[0])
	}

	flag.Parse()

	if device == "" {
		flag.Usage()
		os.Exit(-1)
	}

	// if no formats specified, enable all
	if !enableJPEG && !enablePNG && !enablePDF && !enableGIF && !enableZIP {
		enableJPEG = true
		enablePNG = true
		enablePDF = true
		enableGIF = true
		enableZIP = true
	}

	// initialize file formats after parsing flags
	initFileFormats()

	err := os.MkdirAll(output, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not create output directory %s: %v\n", output, err)
		os.Exit(-1)
	}

	if err := scan(device); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
