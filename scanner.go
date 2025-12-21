package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func scan(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not open %v: %w", path, err)
	}
	defer f.Close()

	// get number of device's blocks
	blocks, err := unix.IoctlGetInt(int(f.Fd()), unix.BLKGETSIZE)
	if err != nil {
		return fmt.Errorf("can not get size of %v: %w", path, err)
	}
	devSize := blocks * SECTOR_SIZE
	println("the device size is", devSize, "bytes")

	// it's just a babystep just to check if we can work with
	// the data accessed
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

	sr := io.NewSectionReader(f, 0, int64(devSize))
	buf := make([]byte, CHUNK_SIZE)
	var offset int64
	for {
		_, err := sr.ReadAt(buf, offset)
		if err != nil {
			if err == io.EOF {
				println("##### EOF")
				break
			}
			return fmt.Errorf("can not read chunk %d: %w", offset, err)
		}
		if i := bytes.Index(buf, JPEG_SIGNATURE); i >= 0 {
			headPos := offset + int64(i)
			println("found jpeg sig at", headPos)
			// todo: goroutine
			extract(buf, headPos)
		}
		offset += CHUNK_SIZE
	}

	return nil
}

// jpeg, one chunk for now; it's just a POC
func extract(data []byte, headPos int64) error {
	fName := filepath.Join(output, "1.jpeg")
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		return err
	}
	defer f.Close()

	// search for jpeg tail
	if i := bytes.Index(data[headPos:], JPEG_TAIL); i < 0 {
		return nil
	} else {
		_, err = f.Write(data[headPos:i])
		if err != nil {
			return err
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
