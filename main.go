package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

type Webpage struct {
	Title   string
	Content string
}

// Extracts the json formatted front matter from a content file. Returns
// the front matter of said file and the index at which Markdown content starts.
func GetFrontMatter(filePath string, delimiter string) (map[string]string, int, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, -1, err
	}

	fileContent := string(fileBytes)

	startIndex := strings.Index(fileContent, delimiter)
	if startIndex == -1 {
		return nil, -1, errors.New("no front matter delimiter found in content file")
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

// Writes the HTML to specified destination using the provided template and content.
func GenerateHtmlFile(templateFilePath string, contentFilePath string, destinationPath string) {
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		panic(err)
	}

	frontMatter, fmEndIndex, err := GetFrontMatter(contentFilePath, "+++")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	fileBytes, err := os.ReadFile(contentFilePath)
	if err != nil {
		panic(err)
	}

	fileContent := string(fileBytes)[fmEndIndex:]

	err = goldmark.Convert([]byte(fileContent), &buf)
	if err != nil {
		panic(err)
	}

	htmlContent := Webpage{frontMatter["Title"], buf.String()}

	err = tmpl.Execute(os.Stdout, htmlContent)
	if err != nil {
		panic(err)
	}
}

// For each file in content, find the template with the same name and apply it.

func main() {
	sitePath := flag.String("p", ".", "path to website directory")

	flag.Parse()

	contentPath := *sitePath + "/content"

	GenerateHtmlFile("templates/blog.html", contentPath+"/blog/smth.md", *sitePath+"/website")
}
