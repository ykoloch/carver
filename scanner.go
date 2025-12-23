package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

var (
	fileCount atomic.Int32
	seen      sync.Map
)

// fileSeq is a sequence of bytes bounded by the head's and tail's positions;
// this sequence represents content of the file to be recoverd
type fileSeq struct {
	headPos, tailPos int64
	data             []byte
}

// scan sequentially gets chunks of data from the device and sends
// them one by one further on for processing
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
	// process the data chunk by chunk
	for {
		_, err := sr.ReadAt(buf, offset)
		if err != nil {
			if err == io.EOF {
				println("##### EOF")
				break
			}
			return fmt.Errorf("can not read chunk %d: %w", offset, err)
		}

		// todo: multiple headers in the same chunk
		// will work when we implement processing when the file spreads
		// beyond current chunk
		//headPos := offset + int64(headPosition)
		// todo: goroutine, log error
		processJPEG(buf)
		// todo: process png, pdf etc.

		offset += CHUNK_SIZE
	}

	return nil
}

// processJPEG
func processJPEG(chunk []byte) {
	// no jpegs in given chunk
	// but what if jpeg spawns 2 chunks?
	if bytes.Index(chunk, JPEG_TAIL) < 0 {
		return
	}

	for _, sig := range JPEG_SIGS {
		headerPos := bytes.Index(chunk, sig)
		if headerPos < 0 {
			return
		} else {
			tailPosition := bytes.Index(chunk[headerPos:], JPEG_TAIL)
			if tailPosition < 0 {
				// todo: the tail in the next chunk?
				return
			}
			// todo: process error
			_ = storeJPEG(chunk[headerPos : tailPosition+2])
			// base case for the recursion
			processJPEG(chunk[tailPosition+2:])
		}
	}
}

// todo: use generics? make this function universal for extracting all file types?
func storeJPEG(data []byte) error {
	// todo: what if the target directory doesn't exist
	num := fileCount.Add(1)
	fCount := fmt.Sprintf("%d.%s", num, JPEG_EXT)
	fName := filepath.Join(output, fCount)

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)

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
