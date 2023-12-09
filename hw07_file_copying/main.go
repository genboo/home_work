package main

import (
	"flag"
	"log"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	var err error
	var fileFrom *os.File
	fileFrom, err = os.Open(from)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileFrom)

	fileTo, err := os.Create(to)
	if err != nil {
		return
	}
	defer func(fileTo *os.File) {
		err = fileTo.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileTo)

	buff := make([]byte, limit)
	_, err = fileFrom.ReadAt(buff, offset)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(to, buff, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}
