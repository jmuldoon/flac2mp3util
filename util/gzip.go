package util

import (
	"bufio"
	"compress/gzip"
	"os"
)

type Gzipper interface {
	Create(string) error
	Close()
	Write(*gzip.Reader) error
}

type GzipType struct {
	file       *os.File
	gzipWriter *gzip.Writer
	fileWriter *bufio.Writer
}

// Create sets up the GzipType struct to have all the writer properties
// needed to utilize the gzipping functionality
func (gz *GzipType) Create(path string) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	gzipw := gzip.NewWriter(f)
	filew := bufio.NewWriter(gzipw)
	gz.file = f
	gz.gzipWriter = gzipw
	gz.fileWriter = filew

	return nil
}

// Writer writes the data on gzip.reader
func (gz *GzipType) Write(gzipr *gzip.Reader) (err error) {
	buff := make([]byte, 0, 4096)

	if _, err := gzipr.Read(buff); err != nil {
		return err
	}
	for _, b := range buff {
		if err := gz.fileWriter.WriteByte(b); err != nil {
			return err
		}
	}
	return nil
}

// Close flushes the writer, closes the gzip.writer and then closes the file.
func (gz *GzipType) Close() {
	gz.fileWriter.Flush()
	gz.gzipWriter.Close() // close the gzip before the file
	gz.file.Close()
}
