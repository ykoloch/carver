package main

import "io"

// checkDevice reads first 512 bytes (normal sector size) of a given device
// and thus checks availability of the device as well as presence of any meaningful
// data at the device
func checkDevice(f io.Reader) ([]byte, error) {
	buf := make([]byte, SECTOR_SIZE)
	// todo: log number of bytes read
	_, err := f.Read(buf)
	return buf, err
}
