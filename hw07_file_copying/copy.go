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

func copyNWithProgress(dst io.Writer, src io.Reader, nTotal int64) (int64, error) {
	var copied int64

	var chunkSize int64 = 16
	progreesBar := NewProgressBar(nTotal)

	for {
		remain := nTotal - copied
		if remain == 0 {
			break
		}
		if remain < chunkSize {
			chunkSize = remain
		}

		n, err := io.CopyN(dst, src, chunkSize)
		copied += n
		progreesBar.Print(copied)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return copied, err
		}
	}
	return copied, nil
}

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
