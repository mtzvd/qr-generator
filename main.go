package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

// Constants
const (
	maxURLLength = 2048
	minQRSize    = 100
	maxQRSize    = 4096
	unitSize     = 6
)

// Store regular expression for reuse
var filenameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Helper function to check and handle errors
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// customUsage prints usage message
func customUsage() {
	programName := filepath.Base(os.Args[0]) // Get the base name of the binary
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", programName)
	fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  %s -u 'https://www.example.com' -s 256 -l M -f png -d /path/to/save\n", programName)
	fmt.Fprintf(flag.CommandLine.Output(), "  %s -u 'https://www.example.com' -s 512 -l Q -f svg\n", programName)
}

// generateSVG generates svg vector image as string
func generateSVG(qr *qrcode.QRCode) string {
	var builder strings.Builder

	bitmap := qr.Bitmap()
	dim := len(bitmap)

	// Use fmt.Fprintf for direct writing to builder
	fmt.Fprintf(&builder, "<svg width=\"%d\" height=\"%d\" xmlns=\"http://www.w3.org/2000/svg\">\n", dim*unitSize, dim*unitSize)
	for y := 0; y < dim; y++ {
		for x := 0; x < dim; x++ {
			if bitmap[y][x] {
				fmt.Fprintf(&builder, "<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"#000\"/>\n", x*unitSize, y*unitSize, unitSize, unitSize)
			}
		}
	}
	builder.WriteString("</svg>")

	return builder.String()
}

// sanitizeFilename clears string from characters unsafe for filenames
func sanitizeFilename(input string) string {
	return filenameSanitizer.ReplaceAllString(input, "_")
}

func main() {
	// Parse command string flags
	urlFlag := flag.String("u", "", "URL to generate QR code for (max URL length 2048)")
	levelFlag := flag.String("l", "M", "Correction level (L, M, Q, H)")
	formatFlag := flag.String("f", "png", "Output format (png, svg)")
	sizeFlag := flag.Int("s", 256, "Size of the QR code (default 256, min 100, max 4096)")
	dirFlag := flag.String("d", ".", "Directory to save the file (default is current directory)")
	flag.Parse()

	// Display defaults if no flags provided
	flag.Usage = customUsage
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	// Check URL length
	if len(*urlFlag) == 0 || len(*urlFlag) > maxURLLength {
		fmt.Printf("Error: URL is required and must be less than %d characters.\n", maxURLLength)
		return
	}

	// Check QR size
	if *sizeFlag < minQRSize || *sizeFlag > maxQRSize {
		fmt.Printf("Error: Size of the QR code must be between %d and %d.\n", minQRSize, maxQRSize)
		return
	}

	// Connect stadard correction levels to constants
	var level qrcode.RecoveryLevel
	switch *levelFlag {
	case "L":
		level = qrcode.Low
	case "M":
		level = qrcode.Medium
	case "Q":
		level = qrcode.High
	case "H":
		level = qrcode.Highest
	default:
		fmt.Println("Invalid correction level. Choose from L, M, Q, H.")
		return
	}

	//Generate QRcode
	qr, err := qrcode.New(*urlFlag, level)
	checkError(err)

	// Print QRcode to console
	fmt.Println(qr.ToSmallString(false))

	// Prepare filename
	dir, err := filepath.Abs(*dirFlag)
	checkError(err)
	currentTime := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("qrcode%s%s.%s", currentTime, sanitizeFilename(*urlFlag), *formatFlag)
	filePath := filepath.Join(dir, filename)

	// Save file in selected format
	switch *formatFlag {
	case "png":
		err = qr.WriteFile(*sizeFlag, filePath)
	case "svg":
		svgStr := generateSVG(qr)
		err = os.WriteFile(filePath, []byte(svgStr), 0644)
	default:
		fmt.Println("Invalid format. Choose from png or svg.")
		return
	}
	checkError(err)

	fmt.Println("QR code saved as:", filePath)
}
