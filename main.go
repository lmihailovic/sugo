package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/yuin/goldmark/extension"

	"github.com/yuin/goldmark"
)

// GetFrontMatter extracts the json formatted front matter from a content file. Returns
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

	return data, endIndex + len(delimiter), nil
}

// GetSpecificFrontMatter returns the value of targeted key in front matter map.
func GetSpecificFrontMatter(filePath string, delimiter string, target string) (any, error) {
	frontMatter, _, err := GetFrontMatter(filePath, delimiter)
	if err != nil {
		return nil, err
	}

	targetData := frontMatter[target]
	return targetData, nil
}

// GetTextContent reads a Markdown file, processes its content from a given offset, and converts it to HTML text.
func GetTextContent(filePath string, offset int) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	fileContent := fileBytes[offset:]
	var buf bytes.Buffer

	mdParser := goldmark.New(goldmark.WithExtensions(extension.GFM))

	err = mdParser.Convert(fileContent, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ListFiles returns all files inside of directory recursively.
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

// GetChildPages scans the specified content directory for markdown files and returns a map
// of their HTML file names mapped to their titles.
func GetChildPages(url string, indexesOnly bool) map[string]map[string]any {
	var pages = make(map[string]map[string]any)
	root := filepath.Join("content", url)
	rootDepth := strings.Count(root, string(os.PathSeparator))

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".md" {
			htmlFilePath := strings.TrimSuffix(path, filepath.Ext(path))
			htmlFilePath += ".html"

			if !indexesOnly && strings.Contains(htmlFilePath, "index.html") {
				return nil
			} else if indexesOnly && !strings.Contains(htmlFilePath, "index.html") {
				return nil
			}

			relPath, err := filepath.Rel("content", htmlFilePath)
			if err != nil {
				log.Fatal(err)
			}

			relDepth := strings.Count(relPath, string(os.PathSeparator))

			if d.IsDir() && (indexesOnly && relDepth-rootDepth == 0) || relDepth-rootDepth > 1 {
				return filepath.SkipDir
			}

			if indexesOnly {
				relDepth := strings.Count(relPath, string(os.PathSeparator))
				if relDepth-rootDepth == 0 {
					return nil
				}
			}

			fullPath := filepath.Join("/", relPath)

			pageData, _, err := GetFrontMatter(path, "+++")
			if err != nil {
				log.Fatal(err)
			}

			pages[fullPath] = pageData
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return pages
}

// SortPages sorts a map of pages by the specified key in ascending or descending order, returning a slice of pages.
func SortPages(pages map[string]map[string]any, sortKey string, desc bool) (sortedPages []map[string]any) {
	out := make([]map[string]any, 0)
	for link, data := range pages {
		data["Link"] = link
		out = append(out, data)
	}

	if sortKey == "Date" {
		sort.Slice(out, func(i, j int) bool {
			t1, err := time.Parse("2-1-2006", out[i][sortKey].(string))
			if err != nil {
				log.Fatal(err)
			}

			t2, err := time.Parse("2-1-2006", out[j][sortKey].(string))
			if err != nil {
				log.Fatal(err)
			}

			if desc {
				return t1.After(t2)
			} else {
				return t1.Before(t2)
			}
		})

		return out
	}

	sort.Slice(out, func(i, j int) bool {
		if desc {
			return out[i][sortKey].(string) > out[j][sortKey].(string)
		}
		return out[i][sortKey].(string) < out[j][sortKey].(string)
	})

	return out
}

// CopyStaticDir copies all files from the source directory to the destination directory.
// src specifies the source directory containing the files to be copied.
// dst specifies the destination directory where the files will be placed.
func CopyStaticDir(src string, dst string) error {
	destFS := os.DirFS(src)

	err := os.CopyFS(dst, destFS)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	sitePath := flag.String("i", ".", "path to website directory")
	outputRootPath := flag.String("o", filepath.Join(*sitePath, "website"), "path for generated static web files")
	devServer := flag.Bool("d", false, "run dev server")

	flag.Parse()

	if *devServer {
		fileServer := http.FileServer(http.Dir(*outputRootPath))
		http.Handle("/", fileServer)

		fmt.Printf("Starting server at http://localhost:8080/\n")

		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}

	contentPath := filepath.Join(*sitePath, "content")
	staticPath := filepath.Join(*sitePath, "static")

	files := ListFiles(contentPath)

	for _, file := range files {
		relPath, err := filepath.Rel(contentPath, file)
		if err != nil {
			panic(err)
		}
		relPathDirs, fileName := filepath.Split(relPath)
		baseFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		templFile := "single.gohtml"
		if baseFileName == "index" {
			templFile = "section.gohtml"
		}

		tmpl := template.New("").Funcs(template.FuncMap{
			"GetChildPages": GetChildPages,
			"SortPages":     SortPages,
		})

		tmpl = template.Must(tmpl.ParseFiles(
			"templates/_layouts/base.gohtml",
			"templates/_layouts/head.gohtml",
			"templates/_layouts/header.gohtml",
			"templates/_layouts/footer.gohtml",
			filepath.Join("templates", relPathDirs, templFile), // specific page
		))

		frontmatter, offset, err := GetFrontMatter(file, "+++")
		if err != nil {
			panic(err)
		}

		pageData := frontmatter
		pageData["Content"], err = GetTextContent(file, offset)
		if err != nil {
			panic(err)
		}

		htmlFileName := baseFileName + ".html"

		destPath := filepath.Join(*outputRootPath, relPathDirs)
		err = os.MkdirAll(destPath, os.ModePerm)
		if err != nil {
			panic(err)
		}

		outputFile, err := os.Create(filepath.Join(destPath, htmlFileName))
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		err = tmpl.ExecuteTemplate(outputFile, "base.gohtml", pageData)
		if err != nil {
			return
		}
	}

	err := CopyStaticDir(staticPath, *outputRootPath)
	if err != nil {
		panic(err)
	}

}
