package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty command", func(t *testing.T) {
		code := RunCmd([]string{}, Environment{})
		require.EqualValues(t, code, codeWrongCmd)
	})
	t.Run("wrong command", func(t *testing.T) {
		code := RunCmd([]string{"/dev/test"}, Environment{})
		require.EqualValues(t, code, codeRunError)
	})
	t.Run("correct command run", func(t *testing.T) {
		code := RunCmd([]string{"/bin/bash"}, Environment{})
		require.EqualValues(t, code, 0)
	})
}
