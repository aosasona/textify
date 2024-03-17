package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	content   string
	extension = "kt"
	ignore    string
)

func main() {
	var targetFile string

	flag.StringVar(&targetFile, "file", "file.txt", "file to write to")
	flag.StringVar(&extension, "extension", "kt", "file extension to read")
	flag.StringVar(&ignore, "ignore", "", "directories to ignore, separated by comma")
	flag.Parse()

	extension = "." + extension
	sourceDir := flag.Arg(0)
	pathsToIgnore := strings.Split(ignore, ",")

	if sourceDir == "" {
		panic("source directory is required")
	}

	if targetFile == "" {
		panic("target file is required")
	}

	file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open file:" + err.Error())
	}
	defer file.Close()

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		panic("failed to read directory:" + err.Error())
	}

	readFiles(files, sourceDir, pathsToIgnore)

	_, err = file.WriteString(content)
	if err != nil {
		panic("failed to write to file to target file:" + err.Error())
	}
}

func readFiles(files []fs.DirEntry, root string, pathsToIgnore []string) {
	for _, file := range files {
		if file.IsDir() {
			if slices.Contains(pathsToIgnore, file.Name()) {
				continue
			}

			contents, err := os.ReadDir(filepath.Join(root, file.Name()))
			if err != nil {
				panic("failed to read directory:" + err.Error())
			}

			readFiles(contents, filepath.Join(root, file.Name()), pathsToIgnore)
		}

		if file.Type().IsRegular() {
			readFile(file.Name(), filepath.Join(root, file.Name()))
		}
	}
}

func readFile(filename string, path string) {
	if !strings.HasSuffix(path, extension) {
		return
	}

	c, err := os.ReadFile(path)
	if err != nil {
		panic("failed to read file:" + err.Error())
	}

	fmt.Println("Reading file:", filename)
	content += fmt.Sprintf("// %s\n", filename)
	content += string(c)
	content += "\n\n\n"
}
