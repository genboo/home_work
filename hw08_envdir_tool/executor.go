package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	codeWrongCmd = 15
	codeRunError = 16
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return codeWrongCmd
	}
	c := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	for key, e := range env {
		if i := IndexOf(c.Env, key); i > -1 {
			c.Env = append(c.Env[0:i], c.Env[i+1:]...)
		}
		if !e.NeedRemove {
			c.Env = append(c.Env, fmt.Sprintf("%s=%s", key, e.Value))
		}
	}
	if err := c.Run(); err != nil {
		log.Println(err)
		return codeRunError
	}
	return c.ProcessState.ExitCode()
}

func IndexOf(slice []string, val string) int {
	if len(slice) == 0 {
		return -1
	}
	for i, v := range slice {
		if strings.HasPrefix(v, val+"=") {
			return i
		}
	}
	return -1
}
