# Carver 🔍

File recovery tool that scans raw storage devices (disks, flash drives, SD cards) and recovers deleted files by searching for file signatures (magic bytes).

## Why?

Sometimes you accidentally delete files, format a drive, or deal with corrupted filesystems. Traditional recovery tools are either expensive, bloated, or require GUIs. Carver is:

- **Simple** - just 2 flags, does one thing well
- **Fast** - concurrent scanning with Go goroutines
- **Portable** - single binary, no dependencies
- **Transparent** - shows progress and tells you what it finds

Built during military leave to keep coding skills sharp. Turned out pretty useful.

## Supported File Formats

- **JPEG** - supports all common signatures (JFIF, EXIF, Huffman-first, etc.)
- **PNG** - full 8-byte signature validation
- **PDF** - finds document boundaries
- **GIF** - both GIF87a and GIF89a
- **ZIP** - includes variable-length EOCD parsing

### Why are DOCX and XLSX saved as ZIP?

Modern Microsoft Office formats (DOCX, XLSX, PPTX) are actually **ZIP archives** with XML inside. They share the same signature: `50 4B 03 04`.

Currently, Carver saves all of them with `.zip` extension. You can identify the real type using:
```bash
file *.zip
```

Output:
```
1.zip: Microsoft Word 2007+
2.zip: Zip archive data
3.zip: Microsoft Excel 2007+
```

Adding separate DOCX/XLSX detection is planned, but for now - this works.

## Installation
```bash
git clone <your-gitlab-url>
cd carver
go build
```

Binary will be in current directory: `./carver`

## Usage
```bash
# Recover files from USB drive
sudo ./carver -device /dev/sdb1 -output ./recovered

# Use default output directory (./recovered)
sudo ./carver -device /dev/sdc1

# Scan disk image file
./carver -device ./disk.img -output ./found
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-device` | *(required)* | Path to device or image file to scan |
| `-output` | `./recovered` | Directory for recovered files |

### Output

Files are saved as: `1.jpeg`, `2.png`, `3.zip`, etc.

Progress indicator shows:
```
Progress: 42.3% (3200 / 7644 MB)
```

Final summary:
```
Done! Recovered 847 files
```

## How It Works

1. **Reads device in 64MB chunks** - efficient for large drives
2. **Searches for file signatures** - e.g., JPEG starts with `FF D8 FF`
3. **Finds file boundaries** - looks for end-of-file markers
4. **Extracts complete files** - saves header → tail
5. **Handles edge cases** - variable-length structures (ZIP EOCD), multiple signatures (JPEG), deduplication

### Technical Details

- **Concurrent scanning** - each file format processed in parallel
- **Deduplication** - won't save the same file twice (e.g., JPEG with 5 different signatures)
- **Smart tail parsing** - ZIP End of Central Directory includes comment length field

## Limitations

- **Files spanning chunks** - if file is split across 64MB boundary, won't recover (rare)
- **Fragmented files** - only works for contiguous files
- **No filesystem metadata** - recovered files lose original names/timestamps
- **Linux-specific** - uses `unix.BLKGETSIZE` ioctl (PRs for Windows/macOS welcome)

## Requirements

- Go 1.21+
- Linux (for block device access)
- `sudo` for reading raw devices

## Roadmap

- [ ] MP4 video recovery (no fixed tail signature)
- [ ] DOCX/XLSX type detection (ZIP subformat analysis)
- [ ] Cross-platform support (Windows, macOS)
- [ ] Verbose mode with per-file logging
- [ ] Statistics (files by type, success rate)

## License

MIT

## Author

Built by [Yurii](https://linkedin.com/in/your-profile) - Ukrainian military medical officer in ZSU, 8+ years Go developer.

Created during leave to maintain coding skills and build something useful for data recovery scenarios.

---

**⚠️ Warning:** Always work on copies or unmounted devices to avoid further data loss.