package main

import (
	"archive/zip"
	"compress/flate"
	"io"
	"io/fs"
	"os"
)

type WriteCounter struct {
	writer io.Writer
	count  int64
	name   string
}

func NewWriteCounter(writer io.Writer, name string) *WriteCounter {
	w := &WriteCounter{
		writer: writer,
		name:   name,
	}

	return w
}

func (w *WriteCounter) Reset() {
	w.count = 0
}

func (w *WriteCounter) Count() int64 {
	return w.count
}

func (w *WriteCounter) Write(data []byte) (int, error) {
	// fmt.Printf("  - [%s] counter write: %d\n", w.name, len(data))
	c, err := w.writer.Write(data)
	w.count += int64(c)
	return c, err
}

type ZipFile struct {
	Fd       *os.File
	writer   *zip.Writer
	lcounter *WriteCounter
	counter  *WriteCounter
	Level    int
}

func NewZipFile(filename string) (*ZipFile, error) {
	fd, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	lcounter := NewWriteCounter(fd, "LowLevel")
	writer := zip.NewWriter(lcounter)

	file := &ZipFile{
		Fd:       fd,
		writer:   writer,
		lcounter: lcounter,
		Level:    flate.BestCompression,
	}

	writer.RegisterCompressor(zip.Deflate, file.Compressor)
	return file, nil
}

func (f *ZipFile) Close() error {
	if err := f.writer.Close(); err != nil {
		return err
	}

	if err := f.Fd.Close(); err != nil {
		return err
	}

	return nil
}

func (f *ZipFile) Compressor(w io.Writer) (io.WriteCloser, error) {
	f.counter = NewWriteCounter(w, "Compressed")
	return flate.NewWriter(f.counter, f.Level)
}

func (f *ZipFile) AddDirectory(filename string, info fs.FileInfo) (int64, int64, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return 0, 0, err
	}

	header.Method = zip.Store
	w, err := f.writer.CreateHeader(header)
	if err != nil {
		return 0, 0, err
	}

	_, err = w.Write([]byte{})
	if err != nil {
		return 0, 0, err
	}

	f.counter.Reset()
	f.writer.Flush()
	return 0, 0, nil
}

func (f *ZipFile) AddRegularFile(filename string, info fs.FileInfo) (int64, int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}

	defer fd.Close()

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return 0, 0, err
	}

	header.Name = filename
	header.Method = zip.Deflate

	w, err := f.writer.CreateHeader(header)
	if err != nil {
		return 0, 0, err
	}

	sizeOriginal, err := io.Copy(w, fd)
	if err != nil {
		return 0, 0, err
	}

	f.writer.Flush()
	sizeCompressed := f.lcounter.Count()
	f.counter.Reset()

	// sizeLowLevelSize := f.lcounter.Count()
	f.lcounter.Reset()

	// fmt.Printf("  + original: %d, compressed: %d, LSize: %d\n", sizeOriginal, sizeCompressed, sizeLowLevelSize)

	return sizeOriginal, sizeCompressed, nil
}
