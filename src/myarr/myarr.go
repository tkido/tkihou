package myarr

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// MyArr can push unshift etc.
type MyArr struct {
	sl []string
}

// NewMyArr is constructor of MyArr
func NewMyArr(arr ...string) *MyArr {
	a := MyArr{arr}
	return &a
}

// Concat is Concat
func (p *MyArr) Concat(o *MyArr) *MyArr {
	p.sl = append(p.sl, o.sl...)
	return p
}

// First is First
func (p *MyArr) First() string {
	if len(p.sl) == 0 {
		log.Fatal("First(): index out of bound")
	}
	return p.sl[0]
}

// Join is Join
func (p *MyArr) Join(seq string) string {
	return strings.Join(p.sl, seq)
}

// Map is Map
func (p *MyArr) Map(f func(string) string) *MyArr {
	newSl := make([]string, len(p.sl))
	for i, v := range p.sl {
		newSl[i] = f(v)
	}
	p.sl = newSl
	return p
}

// Pop is Pop
func (p *MyArr) Pop() string {
	if len(p.sl) == 0 {
		log.Fatal("index out of bound")
	}
	popped := p.sl[0]
	p.sl = p.sl[1:]
	return popped
}

// Push is Push
func (p *MyArr) Push(arr ...string) *MyArr {
	p.sl = append(p.sl, arr...)
	return p
}

// Size is Size
func (p *MyArr) Size() int {
	return len(p.sl)
}

// ReadLines is ReadLines
func ReadLines(path string) *MyArr {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	arr := NewMyArr()
	for s.Scan() {
		arr.Push(s.Text())
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}
	return arr
}

// TakeBlock is TakeBlock
func (p *MyArr) TakeBlock(re *regexp.Regexp) *MyArr {
	buf := NewMyArr()
	for _, line := range p.sl {
		if re.MatchString(line) {
			buf.Push(re.ReplaceAllString(line, ""))
		} else {
			break
		}
	}
	p.sl = p.sl[buf.Size():]
	return buf
}

// TakeBlockNot is TakeBlockNot
func (p *MyArr) TakeBlockNot(re *regexp.Regexp) *MyArr {
	buf := NewMyArr()
	for _, line := range p.sl {
		if !re.MatchString(line) {
			buf.Push(line)
		} else {
			break
		}
	}
	p.sl = p.sl[buf.Size():]
	return buf
}

// Unshift is Unshift
func (p *MyArr) Unshift(s string) *MyArr {
	newSl := make([]string, 0, len(p.sl)+1)
	newSl = append(newSl, s)
	newSl = append(newSl, p.sl...)
	p.sl = newSl
	return p
}
