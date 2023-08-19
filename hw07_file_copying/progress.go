package main

import (
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

func NewProgressBar(total int64) *ProgressBar {
	pb := &ProgressBar{
		total: total,
	}
	return pb
}

func (pb *ProgressBar) print() {
	prevPercent := pb.percent
	pb.percent = pb.getPercent()
	if pb.percent != prevPercent {
		pb.rate = strings.Repeat("#", pb.percent/2)
	}
	fmt.Printf("\rCopy [%-50s]%3d%%", pb.rate, pb.percent)
	if pb.curr == pb.total {
		fmt.Println("\nDone!")
	}
}

func (pb *ProgressBar) Inc(n int) {
	pb.curr += int64(n)
	pb.print()
}

func (pb *ProgressBar) getPercent() int {
	return int((float32(pb.curr) / float32(pb.total)) * 100)
}

type ProgressReader struct {
	Reader io.Reader
	pb     *ProgressBar
}

func NewProgressReader(reader io.Reader, pb *ProgressBar) *ProgressReader {
	return &ProgressReader{
		Reader: reader,
		pb:     pb,
	}
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)

	if err == nil || err == io.EOF {
		pr.pb.Inc(n)
	}
	return n, err
}

type ProgressWriter struct {
	Writer io.Writer
	pb     *ProgressBar
}

func NewProgressWriter(writer io.Writer, pb *ProgressBar) *ProgressWriter {
	return &ProgressWriter{
		Writer: writer,
		pb:     pb,
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.Writer.Write(p)

	if err == nil {
		pw.pb.Inc(n)
	}
	return n, err
}
