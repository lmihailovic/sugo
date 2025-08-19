package main

import (
	"os"
	"text/template"
)

type Blogpost struct {
	Title   string
	Content []string
}

func main() {
	tmpl, err := template.ParseFiles("templates/blogpost.html")
	if err != nil {
		panic(err)
	}

	blogpost := Blogpost{"Prvi clanak", []string{"lorem ipsum", "dolor sit amet"}}

	err = tmpl.Execute(os.Stdout, blogpost)
	if err != nil {
		panic(err)
	}
}
