package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
func GenerateHtmlFile(templateFilePath string, contentFilePath string, destinationRootPath string) (bool, error) {
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return false, err
	}

	frontMatter, fmEndIndex, err := GetFrontMatter(contentFilePath, "+++")
	if err != nil {
		return false, err
	}

	fileBytes, err := os.ReadFile(contentFilePath)
	if err != nil {
		return false, err
	}

	fileContent := string(fileBytes)[fmEndIndex:]
	var buf bytes.Buffer

	err = goldmark.Convert([]byte(fileContent), &buf)
	if err != nil {
		return false, err
	}

	htmlContent := Webpage{frontMatter["Title"], buf.String()}

	_, baseName := filepath.Split(contentFilePath)
	htmlFileName := strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".html"

	contentRoot := "content"

	destinationPath, err := filepath.Rel(contentRoot, contentFilePath)
	if err != nil {
		panic(err)
	}

	fullDestinationPath := filepath.Join(destinationRootPath, destinationPath)

	err = os.MkdirAll(fullDestinationPath, os.ModePerm)
	if err != nil {
		return false, err
	}

	outputFile, err := os.Create(filepath.Join(fullDestinationPath, htmlFileName))
	if err != nil {
		return false, err
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, htmlContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ListFiles(dir string) []string {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".md" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func main() {
	sitePath := flag.String("p", ".", "path to website directory")
	outputRootPath := flag.String("o", filepath.Join(*sitePath, "website"), "path for generated static web files")

	flag.Parse()

	contentPath := filepath.Join(*sitePath, "content")
	templatesPath := filepath.Join(*sitePath, "templates")

	files := ListFiles(contentPath)

	for _, filename := range files {

		fmt.Println(filename)
	}

	_, err := GenerateHtmlFile(templatesPath+"/blog/single.html", contentPath+"/blog/darkmode-difficulties.md", *outputRootPath)
	if err != nil {
		panic(err)
	}
}
