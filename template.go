package main

import (
	"io/ioutil"
	"os"
	"text/template"
)

// Data data for template
type Data struct {
	Title   string
	Content string
}

var tmpl *template.Template

func execute(title, content string) {
	if tmpl == nil {
		tmplTxt := tmplTkido
		tmpl, _ = template.ParseFiles(tmplTxt)
	}
	ioutil.WriteFile(rstTxt, []byte(content), os.ModePerm)
	rst, _ := os.Create(rstHTML)
	defer rst.Close()
	data := Data{title, content}
	tmpl.Execute(rst, data)
}
