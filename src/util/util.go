package util

import (
	"bytes"
	"regexp"
)

// ReplaceAllStringFuncSubmatches is ReplaceAllStringFuncSubmatches
func ReplaceAllStringFuncSubmatches(re *regexp.Regexp, src string, repl func([]string) string) string {
	buf := bytes.Buffer{}
	anchor := 0
	sms := re.FindAllStringSubmatchIndex(src, -1)
	for _, sm := range sms {
		if sm[0] == -1 {
			continue
		}
		buf.WriteString(src[anchor:sm[0]])
		anchor = sm[1]
		br := []string{}
		for i := 0; i < len(sm)/2; i++ {
			if sm[2*i] != -1 {
				br = append(br, src[sm[2*i]:sm[2*i+1]])
			} else {
				br = append(br, "")
			}
		}
		buf.WriteString(repl(br))
	}
	buf.WriteString(src[anchor:len(src)])
	return buf.String()
}
