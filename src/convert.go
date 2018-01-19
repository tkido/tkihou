package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"log"
	"net/url"
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

var r = `
  \*\*(.+?)\*\*                  # $1: em
| \*(.+?)\*                      # $2: strong
| \\\-(.+?)\-                    # $3: del
| \\_(.+?)_                      # $4: u
| \\(.+?)\\                      # $5: i
| >>(.+?)<<                      # $6: q
| \{(.+?)\}                      # $7: notation
| \[([^;]+?);w\]                 # $8: wikipedia
| \[([^;]+?);g\]                 # $9: google
| \[([^;]+?);nd\]                # $10: nico_dic
| \[([^;]+?);ej\]                # $11: weblio
| \[([0-9^;]+?);y\]              # $12: yahoo finance Japan
| \[([A-Z^;]+?);y\]              # $13: yahoo finance America
| \[([^;]+?);(https?://.+?)\]    # $14:label, $15: URI
`
var reReComment = regexp.MustCompile(`(?m)(\s+)|(\#.*$)`)
var reInlineCheck = regexp.MustCompile(reReComment.ReplaceAllString(r, ""))
var reInlineConvert = regexp.MustCompile(`《START》.+?《END》`)

func inlineConvert(line string) string {
	br := strings.Split(line, "《SEP》")
	if em := br[1]; em != "" {
		return fmt.Sprintf(`<em>%s</em>`, inline(em))
	} else if strong := br[2]; strong != "" {
		return fmt.Sprintf(`<strong>%s</strong>`, inline(strong))
	} else if del := br[3]; del != "" {
		return fmt.Sprintf(`<del>%s</del>`, inline(del))
	} else if u := br[4]; u != "" {
		return fmt.Sprintf(`<u>%s</u>`, inline(u))
	} else if i := br[5]; i != "" {
		return fmt.Sprintf(`<i>%s</i>`, inline(i))
	} else if q := br[6]; q != "" {
		return fmt.Sprintf(`<q>%s</q>`, inline(q))
	} else if nota := br[7]; nota != "" {
		return notation(nota)
	} else if wikipedia := br[8]; wikipedia != "" {
		return fmt.Sprintf(`<a href="http://ja.wikipedia.org/wiki/%s" target="_blank">%s</a>`, url.PathEscape(wikipedia), html.EscapeString(wikipedia))
	} else if google := br[9]; google != "" {
		return fmt.Sprintf(`<a href="http://www.google.com/search?num=50&hl=ja&q=%s&lr=lang_ja" target="_blank">%s</a>`, url.PathEscape(google), html.EscapeString(google))
	} else if nicodic := br[10]; nicodic != "" {
		return fmt.Sprintf(`<a href="http://dic.nicovideo.jp/a/%s" target="_blank">%s</a>`, url.PathEscape(nicodic), html.EscapeString(nicodic))
	} else if weblio := br[11]; weblio != "" {
		return fmt.Sprintf(`<a href="http://ejje.weblio.jp/content/%s" target="_blank">%s</a>`, url.PathEscape(weblio), html.EscapeString(weblio))
	} else if codeJp := br[12]; codeJp != "" {
		return fmt.Sprintf(`<a href="http://stocks.finance.yahoo.co.jp/stocks/detail/?code=%s" target="_blank">%s</a>`, url.PathEscape(codeJp), html.EscapeString(codeJp))
	} else if codeUs := br[13]; codeUs != "" {
		return fmt.Sprintf(`<a href="http://finance.yahoo.com/q?s=%s" target="_blank">%s</a>`, url.PathEscape(codeUs), html.EscapeString(codeUs))
	} else if label, uri := br[14], br[15]; label != "" {
		return fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, uri, html.EscapeString(label))
	}
	return line
}
func inline(line string) string {
	line = reInlineCheck.ReplaceAllString(line, `《START》《SEP》$1《SEP》$2《SEP》$3《SEP》$4《SEP》$5《SEP》$6《SEP》$7《SEP》$8《SEP》$9《SEP》$10《SEP》$11《SEP》$12《SEP》$13《SEP》$14《SEP》$15《SEP》《END》`)
	line = reInlineConvert.ReplaceAllStringFunc(line, inlineConvert)
	return line
}

//TODO
func notation(line string) string {
	return line
}

func tr(line string) string {
	buf := bytes.Buffer{}
	var tag string
	if reTableEnd.MatchString(line) {
		tag = `th`
	} else {
		tag = `td`
	}
	buf.WriteString(`<tr>`)
	for _, col := range strings.Split(line, "|") {
		if col == "" || col == "*" {
			continue
		}
		buf.WriteString(fmt.Sprintf("<%s>", tag))
		buf.WriteString(inline(col))
		buf.WriteString(fmt.Sprintf("</%s>", tag))
	}
	buf.WriteString(`</tr>`)
	return buf.String()
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
	return inline(line) + `<br />`
}

func definition(line string) string {
	pair := strings.Split(line, ":")
	if len(pair) != 2 {
		log.Fatal("definition(): invalid argument")
	}
	return fmt.Sprintf("<dt>%s</dt><dd>%s</dd>", inline(pair[0]), inline(pair[1]))
}
