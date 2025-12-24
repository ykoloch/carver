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
)

// scan sequentially gets chunks of data from the device and sends
// them one by one further on for processing
func scan(path string) error {
	wg := new(sync.WaitGroup)
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

		// will work when we implement processing when the file spreads
		// beyond current chunk
		// todo: goroutine, log error
		wg.Add(2)
		go processJPEG(buf, wg)
		go processPNG(buf, wg)
		wg.Wait()
		offset += CHUNK_SIZE
	}
	return nil
}

// processJPEG
func processJPEG(chunk []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	usedTails := make(map[int]bool)
	for _, sig := range JPEG_SIGS {
		searchOffset := 0
		for searchOffset < len(chunk) {
			// look for the header relatively to  searchOffset
			headerPos := bytes.Index(chunk[searchOffset:], sig)
			if headerPos < 0 {
				break // no more of this signature
			}

			// header's absolute position in the chunk
			absHeaderPos := searchOffset + headerPos

			// search tail starting from header
			tailPos := bytes.Index(chunk[absHeaderPos:], JPEG_TAIL)
			if tailPos < 0 {
				break // tail not found, the file end is in the the chunk, rare
			}

			absTailPos := absHeaderPos + tailPos
			if usedTails[absTailPos] {
				searchOffset = absTailPos + 2
				continue
			}
			usedTails[absTailPos] = true
			_ = saveFile(chunk[absHeaderPos:absTailPos+2], JPEG_EXT)

			// shift search after tail
			searchOffset = absTailPos + 2
		}
	}
}

func processPNG(chunk []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	searchOffset := 0
	for searchOffset < len(chunk) {
		headerPos := bytes.Index(chunk[searchOffset:], PNG_SIGNATURE)
		if headerPos < 0 {
			break
		}
		absHeaderPos := searchOffset + headerPos

		tailPos := bytes.Index(chunk[absHeaderPos:], PNG_TAIL)
		if tailPos < 0 {
			break
		}

		absTailPos := absHeaderPos + tailPos
		_ = saveFile(chunk[absHeaderPos:absTailPos+len(PNG_TAIL)], PNG_EXT)

		searchOffset = absTailPos + len(PNG_TAIL)
	}
}

func saveFile(data []byte, ext string) error {
	num := fileCount.Add(1)
	fCount := fmt.Sprintf("%d.%s", num, ext)
	fName := filepath.Join(output, fCount)

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
