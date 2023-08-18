package main

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	total   int64
	curr    int64
	percent int
	rate    string
}

func NewProgressBar(total int64) (*ProgressBar, error) {
	pb := &ProgressBar{
		total: total,
	}
	return pb, nil
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
