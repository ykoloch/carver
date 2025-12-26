package main

import (
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
	defer func() { _ = f.Close() }()

	// get number of device's blocks
	blocks, err := unix.IoctlGetInt(int(f.Fd()), unix.BLKGETSIZE)
	if err != nil {
		return fmt.Errorf("can not get size of %v: %w", path, err)
	}
	devSize := blocks * SECTOR_SIZE

	sr := io.NewSectionReader(f, 0, int64(devSize))
	buf := make([]byte, CHUNK_SIZE)
	var offset int64
	// process the data chunk by chunk
	for {
		_, err := sr.ReadAt(buf, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("can not read chunk %d: %w", offset, err)
		}
		percent := float64(offset) / float64(devSize) * 100
		fmt.Printf("\r\033[34mProgress:\033[0m %.1f%% (%d / %d MB)",
			percent,
			offset/(1<<20),
			devSize/(1<<20))

		wg.Add(len(fileFormats))
		for _, ff := range fileFormats {
			go ff.process(buf, wg)
		}
		wg.Wait()

		offset += CHUNK_SIZE
	}
	fmt.Printf("\n\033[32mDone!\033[0m Recovered %d files\n", fileCount.Load())
	return nil
}

func saveFile(data []byte, ext string) error {
	num := fileCount.Add(1)
	fCount := fmt.Sprintf("%d.%s", num, ext)
	fName := filepath.Join(output, fCount)

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", closeErr)
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
