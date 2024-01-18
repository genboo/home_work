package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("dir not exists", func(t *testing.T) {
		_, err := ReadDir("/dev/test")
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("correct read", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.Nil(t, err)
		require.EqualValues(t, len(env), 5)
		require.EqualValues(t, env["BAR"].Value, "bar")
		require.False(t, env["BAR"].NeedRemove)
		require.EqualValues(t, env["EMPTY"].Value, "")
		require.False(t, env["EMPTY"].NeedRemove)
		require.EqualValues(t, env["FOO"].Value, "   foo\nwith new line")
		require.False(t, env["FOO"].NeedRemove)
		require.EqualValues(t, env["HELLO"].Value, "\"hello\"")
		require.False(t, env["HELLO"].NeedRemove)
		require.EqualValues(t, env["UNSET"].Value, "")
		require.True(t, env["UNSET"].NeedRemove)
	})
}
