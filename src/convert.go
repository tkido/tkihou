package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^([*#\t:\-\+]|====|\{|\}|>>|<<|$)`)

func convert(src string) {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	lines := make([]string, 0, 128)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}

	title, lines := lines[0], lines[1:]

	buf := make([]string, 0, 128)
	buf = append(buf, "<!--", title, "-->")

	for len(lines) != 0 {
		line := lines[0]
		switch {
		case line == "":
			lines = lines[1:]
		case strings.HasPrefix(line, "===="):
			buf = append(buf, "<hr />")
			lines = lines[1:]
		case strings.HasPrefix(line, "#"):
			tmp := take(lines, "#")
			buf = append(buf, "<!--")
			buf = append(buf, tmp...)
			buf = append(buf, "-->")
			lines = lines[len(tmp):]
		default:
			tmp := takeParagraph(lines)
			buf = append(buf, "<p>")
			buf = append(buf, tmp...)
			buf = append(buf, "</p>")
			lines = lines[len(tmp):]
		}
	}

	content := strings.Join(buf, "\n")

	execute(title, content)
}

func take(lines []string, mark string) []string {
	buf := make([]string, 0, 8)
	for len(lines) != 0 && strings.HasPrefix(lines[0], mark) {
		buf = append(buf, strings.Replace(lines[0], "#", "", 1))
		lines = lines[1:]
	}
	return buf
}

func takeParagraph(lines []string) []string {
	buf := make([]string, 0, 8)
	for len(lines) != 0 && !re.MatchString(lines[0]) {
		buf = append(buf, lines[0])
		lines = lines[1:]
	}
	return buf
}
