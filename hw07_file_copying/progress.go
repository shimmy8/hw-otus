package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type ProgressBar struct {
	total   int64
	curr    int64
	percent int
	rate    string
}

func (pb *ProgressBar) New(total int64) {
	pb.total = total
}

func (pb *ProgressBar) Print(curr int64) {
	pb.curr = curr
	prevPercent := pb.percent
	pb.percent = pb.getPercent()
	if pb.percent != prevPercent {
		pb.rate = strings.Repeat("#", pb.percent/2)
	}
	fmt.Printf("\rCopy [%-50s]%3d%% %8d/%d", pb.rate, pb.percent, pb.curr, pb.total)
	if pb.curr == pb.total {
		fmt.Println("\nDone!")
	}
}

func (pb *ProgressBar) getPercent() int {
	return int((float32(pb.curr) / float32(pb.total)) * 100)
}

func copyNWithProgress(dst io.Writer, src io.Reader, nTotal int64) (int64, error) {
	var copied int64

	var chunkSize int64 = 16
	var progreesBar ProgressBar
	progreesBar.New(nTotal)
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
