package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
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

func normaliseTargetFileName(targetFile string) (string, error) {
	// Expand tilda to home directory
	if strings.HasPrefix(targetFile, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %s", err.Error())
		}

		targetFile = strings.Replace(targetFile, "~", home, 1)
	}

	targetFile, err := filepath.Abs(targetFile)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %s", err.Error())
	}

	return targetFile, nil
}

func main() {
	var targetFile string

	flag.StringVar(&targetFile, "file", "./file.txt", "file to write to")
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

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		panic("failed to read directory:" + err.Error())
	}

	readFiles(files, sourceDir, pathsToIgnore)

	targetFile, err = normaliseTargetFileName(targetFile)
	if err != nil {
		panic(err.Error())
	}

	err = os.WriteFile(targetFile, []byte(content), 0644)
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

	log.Println("Reading file: ", filename)
	content += fmt.Sprintf("// File: %s\n", filename)
	content += string(c)
	content += "\n\n\n\n"
}
