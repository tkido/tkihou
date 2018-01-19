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
var reDl = regexp.MustCompile(`^:`)
var reUl = regexp.MustCompile(`^-`)
var reOl = regexp.MustCompile(`^\+`)
var reTable = regexp.MustCompile(`^\|`)

var reNotP = regexp.MustCompile(`^([*#\t:\-\+]|====|\{|\}|>>|<<|$)`)
var reTableEnd = regexp.MustCompile(`\*$`)

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
		case reDl.MatchString(first):
			buf.Push(`<dl>`).Concat(lines.TakeBlock(reDl).Map(definition)).Push(`</dl>`)
		case reUl.MatchString(first):
			buf.Concat(list("ul", reUl, lines.TakeBlock(reUl)))
		case reOl.MatchString(first):
			buf.Concat(list("ol", reOl, lines.TakeBlock(reOl)))
		case reTable.MatchString(first):
			buf.Push(`<table border="1"><tbody align="center">`).Concat(lines.TakeBlock(reTable).Map(tr)).Push(`</tbody></table>`)
		case !reNotP.MatchString(first):
			buf.Push(`<p>`).Concat(lines.TakeBlockNot(reNotP).Map(paragraph)).Push(`</p>`)
		default:
			buf.Push(lines.Pop())
		}
	}
	content := buf.Join("\n")
	execute(title, content)
}

func tr(line string) string {
	buf := myarr.NewMyArr()
	var tag string
	if reTableEnd.MatchString(line) {
		tag = `th`
	} else {
		tag = `td`
	}
	buf.Push(`<tr>`)
	for _, col := range strings.Split(line, "|") {
		if col == "" || col == "*" {
			continue
		}
		buf.Push(fmt.Sprintf("<%s>", tag))
		buf.Push(inline(col))
		buf.Push(fmt.Sprintf("</%s>", tag))
	}
	buf.Push(`</tr>`)
	return buf.Join("")
}

func list(tag string, re *regexp.Regexp, lines *myarr.MyArr) *myarr.MyArr {
	buf := myarr.NewMyArr()
	buf.Push(fmt.Sprintf("<%s>", tag))
	close := false
	for lines.Size() > 0 {
		if re.MatchString(lines.First()) {
			buf.Concat(list(tag, re, lines.TakeBlock(re)))
		} else {
			if close {
				buf.Push(`</li>`)
			}
			close = true
			buf.Push(`<li>`, inline(lines.Pop()))
		}
	}
	if close {
		buf.Push(`</li>`)
	}
	buf.Push(fmt.Sprintf("</%s>", tag))
	return buf
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
	level := 2 + strings.Count(line, "*")
	if level > 6 {
		level = 6
	}
	content := strings.Replace(line, "*", "", -1)
	return fmt.Sprintf("<h%d>%s</h%d>", level, content, level)
}

func paragraph(line string) string {
	return line + `<br />`
}

func inline(line string) string {
	return line
}
func definition(line string) string {
	pair := strings.Split(line, ":")
	if len(pair) != 2 {
		log.Fatal("definition(): invalid argument")
	}
	return fmt.Sprintf("<dt>%s</dt><dd>%s</dd>", inline(pair[0]), inline(pair[1]))
}
