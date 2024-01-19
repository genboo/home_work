package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Panicln("not all parameters are passed")
	}
	envDir := flag.Args()[0]
	if fileInfo, err := os.Stat(envDir); err != nil || !fileInfo.IsDir() {
		log.Panicln("wrong env directory path")
	}
	cmd := flag.Args()[1:]
	if _, err := exec.LookPath(cmd[0]); err != nil {
		log.Panicln("wrong command")
	}

	env, err := ReadDir(envDir)
	if err != nil {
		log.Panicln(err)
	}

	os.Exit(RunCmd(cmd, env))
}
