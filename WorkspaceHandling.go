package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Workspace struct {
	Folders  []Folder
	Settings Settings
}
type Folder struct {
	Path string
}
type Settings struct{}

func AppendDirsToWorkspace(path *string, wg *sync.WaitGroup, dirs *[]string) {
	if content, err := ioutil.ReadFile(*path); err == nil {
		var workspace Workspace
		if err := json.Unmarshal(content, &workspace); err == nil {
			for _, dir := range *dirs {
				if !WorkspaceContainsDir(&workspace, path, &dir) {
					workspace.Folders = append(workspace.Folders, Folder{dir})
				}
			}
		}

		newJSONBytes, err := json.MarshalIndent(workspace, "", "    ")
		if err == nil {
			SaveJSON(path, &newJSONBytes)
		}
	}

	(*wg).Done()
}

func WorkspaceContainsDir(workspace *Workspace, workspacePath, dir *string) bool {
	for _, folder := range workspace.Folders {
		if folder.Path == *dir {
			return true
		}
		p, _ := filepath.Abs(path.Join(filepath.Dir(*workspacePath), folder.Path))
		if strings.Replace(p, "\\", "/", -1) == *dir {
			return true
		}
	}
	return false
}

func SaveJSON(path *string, jsonBytes *[]byte) {
	jsonString := string(*jsonBytes)

	jsonString = strings.ReplaceAll(jsonString, "\"Path\":", "\"path\":")
	jsonString = strings.ReplaceAll(jsonString, "\"Folders\":", "\"folders\":")
	jsonString = strings.ReplaceAll(jsonString, "\"Settings\":", "\"settings\":")

	if file, err := os.Create(*path); err == nil {
		defer file.Close()
		if len, err := file.WriteString(jsonString); err == nil {
			fmt.Printf("Wrote %v characters to %s\n", len, *path)
		} else {
			fmt.Printf("Failed to write to %s", *path)
		}
	}
}
