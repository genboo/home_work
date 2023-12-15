package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	BuffSize            = 256
	ProgressBarSections = 20
	ProgressBarTemplate = "\r[%s%s] %v%%"
	ProgressBarEmpty    = "_"
	ProgressBarFill     = "#"
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var err error
	var fileFrom *os.File
	fileFrom, err = os.Open(fromPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileFrom)

	var fi os.FileInfo
	fi, err = fileFrom.Stat()
	if err != nil {
		return err
	}
	if fi.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	var fileTo *os.File
	fileTo, err = os.Create(toPath)
	if err != nil {
		return err
	}
	defer func(fileTo *os.File) {
		err = fileTo.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileTo)

	if offset > 0 {
		_, err = io.CopyN(io.Discard, fileFrom, offset)
		if err != nil {
			return err
		}
	}
	if limit == 0 || limit+offset > fi.Size() {
		limit = fi.Size() - offset
	}
	done := make(chan struct{})
	progress := make(chan int64)
	go func() {
		for {
			select {
			case <-done:
				fmt.Printf(ProgressBarTemplate+"\n", strings.Repeat(ProgressBarFill, ProgressBarSections), "", "100")
				return
			case v := <-progress:
				percent := float64(v) / float64(limit)
				section := int(float64(ProgressBarSections) * percent)
				filled := strings.Repeat(ProgressBarFill, section)
				empty := strings.Repeat(ProgressBarEmpty, ProgressBarSections-section)
				fmt.Printf(ProgressBarTemplate, filled, empty, math.Floor(percent*100))
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()
	var bytesCopied int64
	var num int64
	buff := int64(BuffSize)
	for {
		if limit-bytesCopied < BuffSize {
			buff = limit - bytesCopied
		}
		num, err = io.CopyN(fileTo, fileFrom, buff)
		if err != nil || num == 0 {
			break
		}
		bytesCopied += num
		progress <- bytesCopied
		if buff == 0 {
			break
		}
	}
	done <- struct{}{}
	return nil
}
