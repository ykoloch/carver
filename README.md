# carver

File recovery tool. Scans raw storage devices and extracts deleted files by signature matching.

## implementation

- concurrent chunk processing with goroutines
- 64MB buffer size
- signature-based file boundary detection
- variable-length structure parsing (ZIP EOCD)

## supported formats

| format | header | tail | notes |
|--------|--------|------|-------|
| jpeg | `FF D8 FF` | `FF D9` | 5 signature variants (JFIF, EXIF, Huffman, Quantization, SOF) |
| png | `89 50 4E 47 0D 0A 1A 0A` | `49 45 4E 44 AE 42 60 82` | full 8-byte validation |
| pdf | `25 50 44 46 2D` | `25 25 45 4F 46` | |
| gif | `47 49 46 38` | `00 3B` | GIF87a and GIF89a |
| zip | `50 4B 03 04` | `50 4B 05 06` | EOCD comment field parsing |

note: DOCX/XLSX/PPTX share ZIP signature. use `file` command for differentiation.

## usage

```bash
# recover all formats
sudo ./carver -device /dev/sdb1 -output ./recovered

# specific formats only
sudo ./carver -device /dev/sdb1 -jpeg -png

# disk image
./carver -device ./disk.img
```

## flags

```
-device string
    path to device or image file (required)
-output string
    output directory (default "./recovered")
-jpeg
    recover JPEG files
-png
    recover PNG files
-pdf
    recover PDF files
-gif
    recover GIF files
-zip
    recover ZIP/DOCX/XLSX files
```

if no format flags specified, all formats are enabled.

## build

```bash
git clone <repo>
cd carver
go build
```

requirements: Go 1.21+, Linux (uses `unix.BLKGETSIZE` ioctl)

## algorithm

1. read device in 64MB chunks via `io.SectionReader`
2. spawn goroutine per file format
3. search for header signatures using `bytes.Index`
4. locate corresponding tail marker
5. extract data slice: `chunk[header:tail+tailSize]`
6. write to disk with sequential naming

## limitations

- fragmented files not supported
- files spanning chunk boundaries not recovered (edge case)
- no metadata preservation (filenames, timestamps)
- Linux only (block device ioctl)

## output

```
Progress: 42.3% (3200 / 7644 MB)
Done! Recovered 847 files
```

files: `1.jpeg`, `2.png`, `3.zip`, etc.

## license

MIT