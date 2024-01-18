package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

const (
	stopLetter = "="
)

var ErrorStopLetter = errors.New("имя не должно содержать =")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, e := range entries {
		if !e.IsDir() {
			if strings.Contains(e.Name(), stopLetter) {
				return nil, ErrorStopLetter
			}
			var buf []byte
			buf, err = os.ReadFile(fmt.Sprintf("%s/%s", dir, e.Name()))
			if err != nil {
				return nil, err
			}
			value := EnvValue{
				NeedRemove: len(buf) == 0,
			}
			if !value.NeedRemove {
				val := string(buf)
				val = strings.TrimRight(strings.TrimRight(strings.Split(val, "\n")[0], " "), "\t")
				val = string(bytes.ReplaceAll([]byte(val), []byte{0x00}, []byte("\n")))
				value.Value = val
			}
			env[e.Name()] = value
		}
	}
	return env, nil
}
