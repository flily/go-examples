package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func usage() {
	fmt.Printf("Usage: %s [file...]\n", os.Args[0])
	fmt.Printf("Display file checksum and block count, like `sum` in Linux and `cksum` in macOS\n")
}

func openFile(name string) (io.ReadCloser, error) {
	if name == "-" {
		return os.Stdin, nil
	}

	return os.Open(name)
}

func main() {
	algo := flag.String("a", "net", "algorithm to use")
	expect := flag.Uint64("check", 0, "checksum to check")
	flag.Usage = usage
	flag.Parse()

	toCheck := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "check" {
			toCheck = true
		}
	})

	checksum := NewChecksum(*algo)

	fileList := []string{"-"}
	if flag.NArg() > 0 {
		fileList = flag.Args()
	}

	for _, arg := range fileList {
		file, err := openFile(arg)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		defer file.Close()
		blocks := 0

		checksum.Reset()
		bufSize := checksum.BlockSize()
		buf := make([]byte, bufSize)
		for {
			length, err := file.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Printf("Error: %s\n", err)
				}
				break
			}

			blocks += 1
			_ = checksum.Update(buf[:length])
		}

		sum := checksum.Checksum()
		checkResult := ""
		if toCheck {
			if sum == *expect {
				checkResult = " [correct]"
			} else {
				checkResult = " [wrong]"
			}
		}

		switch checksum.Size() {
		case 2:
			fmt.Printf("%d 0x%04x %d %s%s\n", sum, sum, blocks, arg, checkResult)
		case 4:
			fmt.Printf("%d 0x%08x %d %s%s\n", sum, sum, blocks, arg, checkResult)
		case 8:
			fmt.Printf("%d 0x%016x %d %s%s\n", sum, sum, blocks, arg, checkResult)
		}
	}
}
