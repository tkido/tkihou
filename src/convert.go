package main

import (
	"bufio"
	"fmt"
	"html"
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

	tmp := make([]string, 0)

	for len(lines) != 0 {
		line := lines[0]
		switch {
		case line == "":
			lines = lines[1:]
		case strings.HasPrefix(line, "#"):
			tmp, lines = split(lines, "#")
			buf = append(buf, "<!--")
			buf = append(buf, mapString(tmp, pass)...)
			buf = append(buf, "-->")
		case strings.HasPrefix(line, "===="):
			tmp, lines = split(lines, "====")
			buf = append(buf, mapString(tmp, hr)...)
		case strings.HasPrefix(line, "*"):
			tmp, lines = split(lines, "*")
			buf = append(buf, mapString(tmp, headline)...)
		case strings.HasPrefix(line, "{{}}"): // notation
			tmp, lines = split(lines, "{{}}")
			buf = append(buf, mapString(tmp, pass)...)
		case strings.HasPrefix(line, "{"):
			tmp, lines = split(lines, "{")
			buf = append(buf, mapString(tmp, divOpen)...)
		case strings.HasPrefix(line, "}"):
			tmp, lines = split(lines, "}")
			buf = append(buf, mapString(tmp, divClose)...)
		case strings.HasPrefix(line, ">>"):
			tmp, lines = split(lines, ">>")
			buf = append(buf, mapString(tmp, bqOpen)...)
		case strings.HasPrefix(line, "<<"):
			tmp, lines = split(lines, "<<")
			buf = append(buf, mapString(tmp, bqClose)...)
		case strings.HasPrefix(line, "\t"):
			tmp, lines = split(lines, "\t")
			buf = append(buf, "<pre><code>")
			buf = append(buf, mapString(tmp, html.EscapeString)...)
			buf = append(buf, "</pre></code>")
		default:
			tmp, lines = splitParagraph(lines)
			buf = append(buf, "<p>")
			buf = append(buf, mapString(tmp, paragraph)...)
			buf = append(buf, "</p>")
		}
	}

	content := strings.Join(buf, "\n")

	execute(title, content)
}

func split(lines []string, mark string) ([]string, []string) {
	buf := make([]string, 0, 8)
	for i, line := range lines {
		if strings.HasPrefix(line, mark) {
			buf = append(buf, strings.Replace(lines[i], mark, "", 1))
		} else {
			break
		}
	}
	return buf, lines[len(buf):]
}

func splitParagraph(lines []string) ([]string, []string) {
	buf := make([]string, 0, 8)
	for _, line := range lines {
		if !re.MatchString(line) {
			buf = append(buf, line)
		} else {
			break
		}
	}
	return buf, lines[len(buf):]
}

func divOpen(line string) string {
	switch {
	case strings.Contains(line, "aa"):
		return `<div class="ascii-art">`
	case strings.Contains(line, "ep"):
		return `<div class="epigraph">`
	default:
		return `<div>`
	}
}
func divClose(line string) string {
	return `</div>`
}
func bqOpen(line string) string {
	return `<blockquote>`
}
func bqClose(line string) string {
	return `</blockquote>`
}
func headline(line string) string {
	level := min(3+strings.Count(line, "*"), 6)
	content := strings.Replace(line, "*", "", -1)
	return fmt.Sprintf("<h%d>%s</h%d>", level, content, level)
}
func hr(line string) string {
	return "<hr />"
}
func paragraph(line string) string {
	return line + `<br />`
}
func pass(line string) string {
	return line
}

func inline(line string) string {
	return line
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func mapString(x []string, f func(string) string) []string {
	r := make([]string, len(x))
	for i, e := range x {
		r[i] = f(e)
	}
	return r
}
