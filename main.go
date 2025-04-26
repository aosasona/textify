package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	sourceDir  string
	content    string
	extensions = []string{"go"}
	ignore     string
	filesCount int
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
	extensionsStr := flag.String("extensions", "go", "file extension to read")
	flag.StringVar(&ignore, "ignore", "node_modules", "directories to ignore, separated by comma")
	flag.Parse()

	extensions = strings.Split(*extensionsStr, ",")
	for i := range extensions {
		extensions[i] = "." + strings.TrimSpace(extensions[i])
	}

	sourceDir = flag.Arg(0)
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

	err = os.WriteFile(targetFile, []byte(strings.TrimSpace(content)), 0644)
	if err != nil {
		panic("failed to write to file to target file:" + err.Error())
	}

	slog.Info("done", slog.String("targetFile", targetFile), slog.Int("filesCount", filesCount))
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
	ext := filepath.Ext(filename)
	if !slices.Contains(extensions, ext) {
		slog.Debug("skipping file", slog.String("file", filename), slog.String("ext", ext))
		return
	}

	c, err := os.ReadFile(path)
	if err != nil {
		panic("failed to read file:" + err.Error())
	}

	slog.Info("reading file", slog.String("file", filename), slog.String("ext", ext))

	// For example, if the source directory is "src" and the file is "src/main.go", the nameWithDirectory will be "main.go"
	nameWithDirectory := strings.TrimPrefix(strings.TrimPrefix(path, sourceDir), "/")
	content += fmt.Sprintf("// File: %s\n", nameWithDirectory)
	content += strings.TrimSpace(string(c))
	content += "\n\n\n\n"
	filesCount++
}
