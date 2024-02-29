package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
)

func unzipListFiles(filename string) error {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		fmt.Printf("%s\n", file.Name)
	}

	return nil
}

func checkExtractDir(dirExtract string) error {
	if len(dirExtract) <= 0 {
		return nil
	}

	info, err := os.Stat(dirExtract)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dirExtract, 0755)
		}

		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dirExtract)
	}

	return nil
}

func unzipWriteFile(file *zip.File, targetFilename string) error {
	fd, err := os.Create(targetFilename)
	if err != nil {
		return err
	}

	defer fd.Close()

	reader, err := file.Open()
	if err != nil {
		return err
	}

	defer reader.Close()

	_, err = io.Copy(fd, reader)
	return err
}

func unzipExtractFiles(filename string, dirExtract string) error {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}

	defer reader.Close()

	if len(dirExtract) > 0 && dirExtract[len(dirExtract)-1] != '/' {
		dirExtract += "/"
	}

	if err := checkExtractDir(dirExtract); err != nil {
		return err
	}

	for _, file := range reader.File {
		var err error
		targetFilename := dirExtract + file.Name

		if file.FileInfo().IsDir() {
			fmt.Printf("   creating: %s\n", targetFilename)
			err = os.MkdirAll(targetFilename, file.Mode())
		} else {
			fmt.Printf(" extracting: %s\n", targetFilename)
			err = unzipWriteFile(file, targetFilename)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	listFiles := flag.Bool("l", false, "list files")
	dirExtract := flag.String("d", "", "directory to extract filess")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	filename := flag.Arg(0)
	var err error
	if *listFiles {
		err = unzipListFiles(filename)
	} else {
		err = unzipExtractFiles(filename, *dirExtract)
	}

	if err != nil {
		panic(err)
	}
}
