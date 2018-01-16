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

func execute(title, content string) {
	ioutil.WriteFile(rstTxt, []byte(content), os.ModePerm)
	tmpl, _ := template.ParseFiles(tmplTxt)
	rst, _ := os.Create(rstHTML)
	defer rst.Close()
	data := Data{title, content}
	tmpl.Execute(rst, data)
}