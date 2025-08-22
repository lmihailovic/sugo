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

// Extracts the json formatted front matter from a content file. Returns
// the front matter of said file and the index at which Markdown content starts.
func GetFrontMatter(filePath string, delimiter string) (map[string]any, int, error) {
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
		return nil, -1, errors.New("invalid json in front matter: " + frontMatter)
	}

	data := make(map[string]any, 0)

	err = json.Unmarshal([]byte(frontMatter), &data)
	if err != nil {
		return nil, endIndex + len(delimiter), err
	}

	// for k, v := range data {
	// 	fmt.Printf("\nk: %v\tv: %v\n", k, v)
	// }

	return data, endIndex + len(delimiter), nil
}

func GetSpecificFrontMatter(filePath string, delimiter string, target string) (any, error) {
	frontMatter, _, err := GetFrontMatter(filePath, delimiter)
	if err != nil {
		return nil, err
	}

	targetData := frontMatter[target]
	return targetData, nil
}

// Writes the HTML to specified destination using the provided template and content.
func GenerateHtmlFile(templateFilePath string, contentFilePath string, destinationRootPath string) error {
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	frontMatter, fmEndIndex, err := GetFrontMatter(contentFilePath, "+++")
	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(contentFilePath)
	if err != nil {
		return err
	}

	fileContent := string(fileBytes)[fmEndIndex:]
	var buf bytes.Buffer

	err = goldmark.Convert([]byte(fileContent), &buf)
	if err != nil {
		return err
	}

	pageData := frontMatter
	pageData["Content"] = buf.String()

	_, baseName := filepath.Split(contentFilePath)
	htmlFileName := strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".html"

	contentRoot := "content"

	destinationPath, err := filepath.Rel(contentRoot, contentFilePath)
	if err != nil {
		return err
	}

	destinationSubDir, _ := filepath.Split(destinationPath)

	fullDestinationPath := filepath.Join(destinationRootPath, destinationSubDir)

	err = os.MkdirAll(fullDestinationPath, os.ModePerm)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Join(fullDestinationPath, htmlFileName))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, pageData)
	if err != nil {
		return err
	}

	return nil
}

// Returns all files inside of directory recursively.
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

	files := ListFiles(contentPath)

	for _, filename := range files {

		filenameParentDir, _ := filepath.Split(filename)
		templPath, err := filepath.Rel("content", filenameParentDir)
		if err != nil {
			panic(err)
		}

		customTempl, err := GetSpecificFrontMatter(filename, "+++", "Template")
		if err != nil {
			panic(err)
		}
		var customTemplPath string
		if customTempl != nil {
			fmt.Printf("\nCustom template specified: %v\n", customTempl)
			customTemplPath = filepath.Join("custom", customTempl.(string))
		}

		if strings.HasSuffix(filename, "index.md") {

			templateFullPath := ""
			if customTempl != nil {
				templateFullPath = filepath.Join("templates", customTemplPath)
			} else {
				templateFullPath = filepath.Join("templates", templPath, "section.html")
			}

			// log.Println("Reached file: " + filename)
			err := GenerateHtmlFile(templateFullPath, filename, *outputRootPath)
			if err != nil {
				panic(err)
			}
		} else {
			templateFullPath := ""
			if customTempl != nil {
				templateFullPath = filepath.Join("templates", customTemplPath)
			} else {
				templateFullPath = filepath.Join("templates", templPath, "single.html")
			}

			// log.Println("Reached file: " + filename)
			err := GenerateHtmlFile(templateFullPath, filename, *outputRootPath)
			if err != nil {
				panic(err)
			}
		}
	}
}
