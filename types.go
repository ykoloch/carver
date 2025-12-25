package main

import (
	"bytes"
	"sync"
)

type fileFormat struct {
	headers [][]byte
	header  []byte
	tail    []byte
	ext     string
}

var (
	// JPEG signatures - different markers after FF D8
	JPEG_JFIF  = []byte{0xFF, 0xD8, 0xFF, 0xE0} // most common - JFIF
	JPEG_EXIF  = []byte{0xFF, 0xD8, 0xFF, 0xE1} // EXIF metadata
	JPEG_HUFF  = []byte{0xFF, 0xD8, 0xFF, 0xC4} // Huffman table first
	JPEG_QUANT = []byte{0xFF, 0xD8, 0xFF, 0xDB} // Quantization table first
	JPEG_SOF   = []byte{0xFF, 0xD8, 0xFF, 0xC0} // Start of Frame (baseline)

	JPEG_SIGS = [][]byte{
		JPEG_JFIF,
		JPEG_EXIF,
		JPEG_HUFF,
		JPEG_QUANT,
		JPEG_SOF,
	}

	JPEG_TAIL = []byte{0xFF, 0xD9}

	//PNG_SIGNATURE = []byte{0x89, 0x50, 0x4E, 0x47}
	PNG_SIGNATURE = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	PNG_TAIL      = []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

	PDF_SIGNATURE = []byte{0x25, 0x50, 0x44, 0x46, 0x2d}
	PDF_TAIL      = []byte{0x25, 0x25, 0x45, 0x4F, 0x46}

	fileFormats []fileFormat
)

func (ff *fileFormat) hasMultipleHeaders() bool {
	return len(ff.headers) > 0
}

func init() {
	jpeg := fileFormat{
		headers: JPEG_SIGS,
		tail:    JPEG_TAIL,
		ext:     JPEG_EXT,
	}
	fileFormats = append(fileFormats, jpeg)

	png := fileFormat{
		header: PNG_SIGNATURE,
		tail:   PNG_TAIL,
		ext:    PNG_EXT,
	}
	fileFormats = append(fileFormats, png)

	pdf := fileFormat{
		header: PDF_SIGNATURE,
		tail:   PDF_TAIL,
		ext:    PNG_EXT,
	}
	fileFormats = append(fileFormats, pdf)
}

// processMultiHeaded deals with file formats which have multiple headers/signatures
func (ff *fileFormat) processLongHeaded(chunk []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	usedTails := make(map[int]bool)
	for _, sig := range ff.headers {
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
			tailPos := bytes.Index(chunk[absHeaderPos:], ff.tail)
			if tailPos < 0 {
				break // tail not found, the file end is in the the chunk, rare
			}

			absTailPos := absHeaderPos + tailPos
			if usedTails[absTailPos] {
				searchOffset = absTailPos + len(ff.tail)
				continue
			}
			usedTails[absTailPos] = true
			_ = saveFile(chunk[absHeaderPos:absTailPos+len(ff.tail)], ff.ext)

			// shift search after tail
			searchOffset = absTailPos + len(ff.tail)
		}
	}
}

func (ff *fileFormat) process(chunk []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	searchOffset := 0
	for searchOffset < len(chunk) {
		headerPos := bytes.Index(chunk[searchOffset:], ff.header)
		if headerPos < 0 {
			break
		}
		absHeaderPos := searchOffset + headerPos

		tailPos := bytes.Index(chunk[absHeaderPos:], ff.tail)
		if tailPos < 0 {
			break
		}

		absTailPos := absHeaderPos + tailPos
		_ = saveFile(chunk[absHeaderPos:absTailPos+len(ff.tail)], ff.ext)

		searchOffset = absTailPos + len(ff.tail)
	}
}
