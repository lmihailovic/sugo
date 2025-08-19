package main

import (
	"bytes"
	"os"
	"text/template"

	"github.com/yuin/goldmark"
)

type Blogpost struct {
	Title   string
	Content string
}

func main() {
	tmpl, err := template.ParseFiles("templates/blogpost.html")
	if err != nil {
		panic(err)
	}

	blogpost := Blogpost{"Prvi clanak", ""}

	// err = tmpl.Execute(os.Stdout, blogpost)
	// if err != nil {
	// 	panic(err)
	// }

	var buf bytes.Buffer

	file, err := os.ReadFile("content/blog/darkmode-difficulties.md")
	if err != nil {
		panic(err)
	}

	err = goldmark.Convert(file, &buf)
	if err != nil {
		panic(err)
	}

	blogpost.Content = buf.String()

	err = tmpl.Execute(os.Stdout, blogpost)
	if err != nil {
		panic(err)
	}
}
