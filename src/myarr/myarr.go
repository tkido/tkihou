package myarr

import (
	"log"
	"regexp"
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
		log.Fatal("index out of bound")
	}
	return p.sl[0]
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

// Unshift is Unshift
func (p *MyArr) Unshift(s string) *MyArr {
	newSl := make([]string, 0, len(p.sl)+1)
	newSl = append(newSl, s)
	newSl = append(newSl, p.sl...)
	p.sl = newSl
	return p
}