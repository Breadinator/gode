package main

import (
	"os"
	"strings"
)

func GetPaths(variable string) []string {
	env := os.Getenv(variable)
	return strings.Split(env, ";")
}

func GetGoSrc() []string {
	paths := GetPaths("GOPATH")
	srcDirectories := make([]string, len(paths))

	for i, path := range paths {
		var slashes string
		if path[len(path):] != "\\" {
			slashes += "\\"
		}
		srcDirectories[i] = path + slashes + "src"
	}

	return srcDirectories
}
