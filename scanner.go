package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

var fileCount atomic.Int32

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
		// ok, jpeg signature found, get its POSITION inside given chunk
		// todo: multiple headers in the same chunk
		if headPosition := bytes.Index(buf, JPEG_SIGNATURE); headPosition >= 0 {
			// will work when we implement processing when the file spreads
			// beyond current chunk
			//headPos := offset + int64(headPosition)
			// todo: goroutine
			extract(buf[headPosition:], int64(headPosition))
		}
		offset += CHUNK_SIZE
	}

	return nil
}

// jpeg, one chunk for now; it's just a POC
func extract(data []byte, headPos int64) error {
	// todo: what if the target directory doesn't exist
	num := fileCount.Add(1)
	fCount := fmt.Sprintf("%d.jpeg", num)
	fName := filepath.Join(output, fCount)

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// search for jpeg tail
	//if tailPosition := bytes.LastIndex(data, JPEG_TAIL); tailPosition < 0 {
	if tailPosition := bytes.Index(data, JPEG_TAIL); tailPosition < 0 {
		return nil
	} else {
		// tailPosition+2 - include the tail bytes themselves
		_, err = f.Write(data[:tailPosition+2])
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
