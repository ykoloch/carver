package main

import (
	"fmt"
	"io"
	"os"
)

func scanPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not open %v: %w", path, err)
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("can not get file stats %v: %w", f.Name(), err)
	}
	println("the device size is", fInfo.Size())

	startData, err := checkDevice(f)
	if err != nil {
		return fmt.Errorf("can not read the device %v: %w", f.Name(), err)
	}

	if verbose {
		fmt.Println("initial data read:")
		for i, b := range startData {
			fmt.Printf("%02X ", b)
			if (i+1)%16 == 0 {
				fmt.Println()
			}
		}
	}

	return nil
}

// checkDevice reads first 512 bytes (normal sector size) of a given device
// and thus checks availability of the device as well as presence of any meaningful
// data at the device
func checkDevice(f io.Reader) ([]byte, error) {
	buf := make([]byte, SECTOR_SIZE)
	// todo: log number of bytes read
	_, err := f.Read(buf)
	return buf, err
}
