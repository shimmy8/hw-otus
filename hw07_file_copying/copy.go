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

	numBytesToCopy := fileSize - offset
	if limit > 0 && limit < numBytesToCopy {
		numBytesToCopy = limit
	}

	// seek to offset
	inpFile.Seek(offset, 0)

	reader := io.ReadCloser(inpFile)
	defer reader.Close()
	writer := io.WriteCloser(outpFile)
	defer writer.Close()

	// *2 to show progress for both reading and writing
	progressBar := NewProgressBar(numBytesToCopy * 2)
	_, copyErr := io.CopyN(
		NewProgressWriter(writer, progressBar),
		NewProgressReader(reader, progressBar),
		numBytesToCopy,
	)
	if copyErr != nil {
		return copyErr
	}

	return nil
}
