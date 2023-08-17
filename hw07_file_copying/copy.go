package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// get file info
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	// file is empty or size unknown
	if fileSize == 0 {
		return ErrUnsupportedFile
	}
	// check offset
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	inpFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer inpFile.Close()

	outpFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outpFile.Close()

	numRemainFileBytes := fileSize - offset
	// set limit to remain filesize
	if limit == 0 || limit > numRemainFileBytes {
		limit = numRemainFileBytes
	}

	// seek to offset
	inpFile.Seek(offset, 0)

	reader := io.ReadCloser(inpFile)
	defer reader.Close()
	writer := io.WriteCloser(outpFile)
	defer writer.Close()

	_, copyErr := copyNWithProgress(writer, reader, limit)
	if copyErr != nil {
		return copyErr
	}

	return nil
}
