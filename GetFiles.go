package main

import (
	"os"
	"path/filepath"
	"strings"
)

func HasExtension(path, ending string) bool {
	return path[len(path)-len(ending):] == ending
}

func GetMods(path string) ([]string, error) {
	var output []string
	err := filepath.Walk( // apparently WalkDir is better but I have no idea how it works
		path,
		func(path string, _ os.FileInfo, err error) error {
			if (err == nil) && HasExtension(path, "go.mod") {
				output = append(output, strings.Replace(filepath.Dir(path), "\\", "/", -1))
			}
			return nil
		},
	)

	return output, err
}

func GetWorkspaces(path string) ([]string, error) {
	var output []string

	err := filepath.Walk(
		path,
		func(path string, _ os.FileInfo, err error) error {
			if err == nil && HasExtension(path, ".code-workspace") {
				output = append(output, strings.Replace(path, "\\", "/", -1))
			}
			return nil
		},
	)

	return output, err
}
