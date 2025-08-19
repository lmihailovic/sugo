package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

type Blogpost struct {
	Title   string
	Content string
}

func GetFrontMatter(filePath string, delimiter string) (map[string]string, int, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, -1, err
	}

	fileContent := string(fileBytes)

	startIndex := strings.Index(fileContent, delimiter)
	if startIndex == -1 {
		return nil, 1, errors.New("no front matter delimiter found in content file")
	}
	startIndex += len(delimiter)

	endIndex := strings.Index(fileContent[startIndex:], delimiter)
	endIndex += len(delimiter)

	frontMatter := "{" + fileContent[startIndex:endIndex] + "}\n"

	if !json.Valid([]byte(frontMatter)) {
		return nil, -1, errors.New("invalid json in front matter")
	}

	data := make(map[string]string, 0)

	err = json.Unmarshal([]byte(frontMatter), &data)
	if err != nil {
		return nil, endIndex + len(delimiter), err
	}

	return data, endIndex + len(delimiter), nil
}

func main() {
	contentPath := "content/blog/smth.md"

	tmpl, err := template.ParseFiles("templates/blogpost.html")
	if err != nil {
		panic(err)
	}

	frontMatter, fmEndIndex, err := GetFrontMatter(contentPath, "+++")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	fileBytes, err := os.ReadFile(contentPath)
	if err != nil {
		panic(err)
	}

	fileContent := string(fileBytes)[fmEndIndex:]

	err = goldmark.Convert([]byte(fileContent), &buf)
	if err != nil {
		panic(err)
	}

	blogpost := Blogpost{frontMatter["Title"], buf.String()}

	err = tmpl.Execute(os.Stdout, blogpost)
	if err != nil {
		panic(err)
	}
}
