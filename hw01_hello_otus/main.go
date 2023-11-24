package main

import (
	"fmt"
	"os"

	"golang.org/x/example/hello/reverse"
)

func main() {
	_, err := fmt.Fprintln(os.Stdout, reverse.String("Hello, OTUS!"))
	if err != nil {
		return
	}
}
