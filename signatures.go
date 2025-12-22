package main

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

	PNG_SIGNATURE = []byte{0x89, 0x50, 0x4E, 0x47}
	PDF_SIGNATURE = []byte{0x25, 0x50, 0x44, 0x46}
)
