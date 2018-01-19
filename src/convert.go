package main

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"os"
	"regexp"
	"strings"

	"./myarr"
)

var reComment = regexp.MustCompile(`^#`)
var reHr = regexp.MustCompile(`^====`)
var reHeadLine = regexp.MustCompile(`^\*`)
var reNotation = regexp.MustCompile(`^{{}}`)
var reDivOpen = regexp.MustCompile(`^{`)
var reDivClose = regexp.MustCompile(`^}`)
var reBqOpen = regexp.MustCompile(`^>>`)
var reBqClose = regexp.MustCompile(`^<<`)
var rePre = regexp.MustCompile(`^\t`)
var reNotParagraph = regexp.MustCompile(`^([*#\t:\-\+]|====|\{|\}|>>|<<|$)`)

func convert(src string) {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	lines := myarr.NewMyArr()
	for s.Scan() {
		lines.Push(s.Text())
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}
	title := lines.Pop()

	buf := myarr.NewMyArr()
	buf.Push(`<!--`, title, `-->`)

	for lines.Size() > 0 {
		first := lines.First()
		switch {
		case first == "":
			lines.Pop()
		case reComment.MatchString(first):
			buf.Push(`<!--`).Concat(lines.TakeBlock(reComment)).Push(`-->`)
		case reHr.MatchString(first):
			lines.Pop()
			buf.Push(`<hr />`)
		case reHeadLine.MatchString(first):
			buf.Push(headLine(lines.Pop()))
		case reNotation.MatchString(first):
			buf.Push(lines.Pop()) //TODO
		case reDivOpen.MatchString(first):
			buf.Push(divOpen(lines.Pop()))
		case reDivClose.MatchString(first):
			lines.Pop()
			buf.Push(`</div>`)
		case reBqOpen.MatchString(first):
			lines.Pop()
			buf.Push(`<blockquote>`)
		case reBqClose.MatchString(first):
			lines.Pop()
			buf.Push(`</blockquote>`)
		case rePre.MatchString(first):
			buf.Push(`<pre><code>`).Concat(lines.TakeBlock(rePre).Map(html.EscapeString)).Push(`</pre></code>`)
		case !reNotParagraph.MatchString(first):
			buf.Push(`<p>`).Concat(lines.TakeBlockNot(reNotParagraph).Map(paragraph)).Push(`</p>`)
		default:
			buf.Push(lines.Pop())
		}
	}
	content := buf.Join("\n")
	execute(title, content)
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

func headLine(line string) string {
	level := min(2+strings.Count(line, "*"), 6)
	content := strings.Replace(line, "*", "", -1)
	return fmt.Sprintf("<h%d>%s</h%d>", level, content, level)
}

func paragraph(line string) string {
	return line + `<br />`
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
