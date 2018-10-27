package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: website-status-checker links.txt")
		os.Exit(1)
	}

	links, err := readLinks(os.Args[1])

	if err != nil {
		panic("Failed to open file")
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	for l := range c {
		go func(link string, channel chan string) {
			time.Sleep(5 * time.Second)
			checkLink(link, channel)
		}(l, c)
	}
}

func checkLink(link string, c chan string) {
	start := time.Now()
	_, err := http.Get(link)
	t := time.Now()
	elapsed := t.Sub(start).String()

	if err != nil {
		fmt.Println(link + " is down (" + elapsed + ")")
	} else {
		fmt.Println(link + " is up (" + elapsed + ")")
	}

	c <- link
}

func readLinks(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	if file, err = os.Open(path); err != nil {
		return
	}

	defer func() {
		err = file.Close()
	}()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
