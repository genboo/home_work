package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("wrong file", func(t *testing.T) {
		err := Copy("/dev/urandom", "out.txt", 0, 0)
		_ = os.Remove("out.txt")
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("offset greater than file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 100000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("not existed from file", func(t *testing.T) {
		err := Copy("testdata/input_1.txt", "out.txt", 0, 0)
		_ = os.Remove("out.txt")
		require.NotNil(t, err)
	})

	t.Run("cant create to file", func(t *testing.T) {
		out := "/out.txt"
		err := Copy("testdata/input.txt", out, 0, 0)
		_ = os.Remove(out)
		require.NotNil(t, err)
	})

	t.Run("negative limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 0, -1)
		require.ErrorIs(t, err, ErrNegativeLimit)
	})

	t.Run("negative offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", -1, 0)
		require.ErrorIs(t, err, ErrNegativeOffset)
	})
}
