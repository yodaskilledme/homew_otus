package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	// fileInfo validation
	if !file.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	// offset validation
	if offset > file.Size() {
		return ErrOffsetExceedsFileSize
	}
	// limit adjusting
	if limit == 0 || limit > file.Size()-offset {
		limit = file.Size() - offset
	}

	inFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer closeOrLogFatal(*inFile)

	outFIle, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeOrLogFatal(*outFIle)

	_, err = inFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(inFile)

	_, err = io.CopyN(outFIle, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()

	return nil
}

func closeOrLogFatal(file os.File) {
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
