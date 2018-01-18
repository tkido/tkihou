package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func convert(src string) {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	if !s.Scan() {
		log.Fatal("No Title!!")
	}
	title := s.Text()
	buf := bytes.NewBufferString(fmt.Sprintf("<!--\n%s\n-->\n", title))
	for s.Scan() {
		buf.WriteString(s.Text())
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}
	content := buf.String()

	execute(title, content)
}
