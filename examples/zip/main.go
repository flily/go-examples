package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
)

type ZipMode int

const (
	ZipModeAdd ZipMode = iota
	ZipModeUpdate
	ZipModeFreshen
	ZipModeDelete
	ZipModeCopy
)

type ZipConfigure struct {
	Mode        ZipMode
	Recursive   bool
	ZipFilename string
	Files       []string
}

func zipAddRegularFile(conf *ZipConfigure, zipFile *ZipFile, filename string, info fs.FileInfo) error {
	sizeOriginal, sizeCompressed, err := zipFile.AddRegularFile(filename, info)
	if err != nil {
		return err
	}

	sizeReduced := sizeOriginal - sizeCompressed
	ratio := 100.0 * float64(sizeReduced) / float64(sizeOriginal)
	fmt.Printf("  adding: %s (deflated %2.2f%%) %d/%d\n", filename, ratio, sizeCompressed, sizeOriginal)
	return nil
}

func zipAddDirectory(conf *ZipConfigure, zipFile *ZipFile, filename string, info fs.FileInfo) error {
	if filename[len(filename)-1] != '/' {
		filename += "/"
	}

	_, _, err := zipFile.AddDirectory(filename, info)
	if err != nil {
		return err
	}

	fmt.Printf("  adding: %s (stored 0.00%%)\n", filename)

	if conf.Recursive {
		entries, err := os.ReadDir(filename)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullname := path.Join(filename, entry.Name())
			err := zipAddFile(conf, zipFile, fullname)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func zipAddFile(conf *ZipConfigure, zipFile *ZipFile, filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return zipAddDirectory(conf, zipFile, filename, info)

	} else {
		return zipAddRegularFile(conf, zipFile, filename, info)
	}
}

func zipAdd(conf *ZipConfigure) error {
	zipFile, err := NewZipFile(conf.ZipFilename)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, filename := range conf.Files {
		err := zipAddFile(conf, zipFile, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func initFlags(conf *ZipConfigure) {
	flag.BoolFunc("u", "Update existing entries if newer on the file system and add new files. "+
		"If the archive does not exist issue warning then create a new archive.",
		func(string) error {
			conf.Mode = ZipModeUpdate
			return nil
		},
	)
	flag.BoolFunc("f", "Update existing entries of an archive if newer on the file system. "+
		"Does not add new files to the archive.",
		func(string) error {
			conf.Mode = ZipModeFreshen
			return nil
		},
	)
	flag.BoolFunc("d", "Select entries in an existing archive and delete them.",
		func(string) error {
			conf.Mode = ZipModeDelete
			return nil
		},
	)
	flag.BoolVar(&conf.Recursive, "r", false, "Travel the directory structure recursively.")
}

func main() {
	conf := &ZipConfigure{
		Mode: ZipModeAdd,
	}
	initFlags(conf)
	flag.Parse()

	if flag.NArg() <= 0 {
		flag.Usage()
		return
	}

	conf.ZipFilename = flag.Arg(0)
	conf.Files = flag.Args()[1:]
	var err error
	switch conf.Mode {
	case ZipModeAdd:
		err = zipAdd(conf)
	}

	if err != nil {
		panic(err)
	}
}
