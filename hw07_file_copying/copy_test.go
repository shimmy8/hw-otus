package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Test input file not found", func(t *testing.T) {
		err := Copy("./testdata/123.txt", "/tmp/123.txt", 100, 10)
		_, ok := err.(*os.PathError)
		require.True(t, ok)
	})

	t.Run("Test unsupported file error", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/123.txt", 100, 10)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("Test offset more than file error", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "/tmp/123.txt", 7000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}
