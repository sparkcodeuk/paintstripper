package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// Handle the -force argument if file exists
func writeForceCheck(path string, force bool) {
	if _, err := os.Stat(path); !os.IsNotExist(err) && !force {
		fmt.Printf("Error: file exists [%s], use -force to override\n", path)
		os.Exit(1)
	}
}

func main() {
	// Command line argument parsing
	args := ParseArgs()

	// Init input
	var inputFd *os.File

	if args.InputPath == "" {
		inputFd = os.Stdin
	} else {
		fileInput, err := os.Stat(args.InputPath)
		if err != nil {
			log.Fatal(err)
		}

		fileMode := fileInput.Mode()
		if !fileMode.IsRegular() {
			log.Fatalf("%s isn't a file", args.InputPath)
		}

		inputFd, err = os.Open(args.InputPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Init outputs
	var outputStrippedFd *os.File
	defer outputStrippedFd.Close()

	var outputUnstrippedFd *os.File
	defer outputUnstrippedFd.Close()

	if args.WriteStripped == "" {
		var err error

		// Write stripped arg undefined; always print stripped output to console
		if args.Quiet {
			outputStrippedFd, err = os.Create(os.DevNull)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			outputStrippedFd = os.Stdout
		}

		if args.WriteUnstripped == "" {
			// Unstripped arg undefined; discard output
			outputUnstrippedFd, err = os.Create(os.DevNull)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Unstripped arg defined; write to file
			writeForceCheck(args.WriteUnstripped, args.Force)
			outputUnstrippedFd, err = os.Create(args.WriteUnstripped)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		var err error

		if args.WriteUnstripped == "" {
			// Write stripped defined, write unstripped undefined;
			writeForceCheck(args.WriteStripped, args.Force)
			outputStrippedFd, err = os.Create(args.WriteStripped)
			if err != nil {
				log.Fatal(err)
			}

			if args.Quiet {
				outputUnstrippedFd, err = os.Create(os.DevNull)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				outputUnstrippedFd = os.Stdout
			}
		} else {
			// Write stripped & unstripped defined; both written to file (no console output)
			writeForceCheck(args.WriteStripped, args.Force)
			outputStrippedFd, err = os.Create(args.WriteStripped)
			if err != nil {
				log.Fatal(err)
			}

			writeForceCheck(args.WriteUnstripped, args.Force)
			outputUnstrippedFd, err = os.Create(args.WriteUnstripped)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Init read/write, caches & initial state
	reader := bufio.NewReader(inputFd)
	const readerBytesDefaultLength = 1024
	var readerBytes []byte

	writerStripped := bufio.NewWriter(outputStrippedFd)
	defer writerStripped.Flush()

	var writerStrippedBytes bytes.Buffer
	writerStrippedBytes.Reset()

	writerUnstripped := bufio.NewWriter(outputUnstrippedFd)
	defer writerUnstripped.Flush()

	// Precompile regexs
	const rePatternDigits0To255 = "([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])"                       // 0-9, 10-99, 100-199, 200-249, 250-255
	const rePatternAnsiColor = "^\033\\[" + rePatternDigits0To255 + "(;" + rePatternDigits0To255 + "){0,4}m" // loose enough to match standard ANSI color codes
	const rePatternAnsiColorMaxLength = 22                                                                   // e.g., "^[255;255;255;255;255m"

	reAnsiColor := regexp.MustCompile(rePatternAnsiColor)

	// Main process loop
	for {
		// Reset byte array if it's grown due to additional boundary condition reads
		if len(readerBytes) != readerBytesDefaultLength {
			readerBytes = make([]byte, readerBytesDefaultLength)
		}

		// Read/EOF termination
		readerBytesLength, err := io.ReadAtLeast(reader, readerBytes, readerBytesDefaultLength)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if readerBytesLength == 0 {
					break
				}
			} else {
				log.Fatal(err)
			}
		}

		// Write unstripoed content verbatim
		writerUnstripped.Write(readerBytes[:readerBytesLength])

		// Process read bytes
		readerBytesIndex := 0
		for readerBytesIndex < readerBytesLength {
			// Color processing
			if readerBytes[readerBytesIndex] == 27 && args.Color {
				// End-of-array boundary condition check
				if (readerBytesLength - readerBytesIndex) < rePatternAnsiColorMaxLength {
					readerBytesAdditional := make([]byte, readerBytesDefaultLength)

					readerBytesAdditionalLength, err := io.ReadAtLeast(reader, readerBytesAdditional, readerBytesDefaultLength)
					if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
						log.Fatal(err)
					}

					// Write unstripoed content verbatim
					writerUnstripped.Write(readerBytesAdditional[:readerBytesAdditionalLength])

					// Append additional read to the existing readerBytes data
					readerBytes = append(readerBytes[:readerBytesLength], readerBytesAdditional[:readerBytesAdditionalLength]...)
					readerBytesLength = len(readerBytes)
				}

				match := reAnsiColor.FindIndex(readerBytes[readerBytesIndex:])

				// Found ANSI color escape code
				if match != nil {
					// Bug: didn't match the beginning of the byte array
					if match[0] != 0 {
						panic("Assertion failure")
					}

					// Strip out the color code we just found
					readerBytes = readerBytes[match[1]+readerBytesIndex:]
					readerBytesLength = readerBytesLength - match[1] - readerBytesIndex
					readerBytesIndex = 0

					continue
				}
			}

			// Write processed byte
			writerStrippedBytes.WriteByte(readerBytes[readerBytesIndex])
			readerBytesIndex++
		}

		// Write stripped bytes
		writerStripped.Write(writerStrippedBytes.Bytes())
		writerStrippedBytes.Reset()
	}
}
