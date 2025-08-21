package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"slices"
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

	contentFileParentDirs := strings.Split(contentFilePath, "/")
	contentFileName := contentFileParentDirs[len(contentFileParentDirs)-1]
	baseName := strings.TrimSuffix(contentFileName, filepath.Ext(contentFileName))
	htmlFileName := baseName + ".html"

	destinationPath := ""
	for i := slices.Index(contentFileParentDirs, "content") + 1; i < len(contentFileParentDirs)-1; i++ {
		destinationPath = filepath.Join(destinationPath, contentFileParentDirs[i])
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

// For each file in content, find the template with the same name and apply it.

func main() {
	sitePath := flag.String("p", ".", "path to website directory")
	outputRootPath := flag.String("o", *sitePath+"/website", "path for outputted static web files")

	flag.Parse()

	contentPath := *sitePath + "/content"

	_, err := GenerateHtmlFile("templates/blog.html", contentPath+"/blog/darkmode-difficulties.md", *outputRootPath)
	if err != nil {
		panic(err)
	}
}
