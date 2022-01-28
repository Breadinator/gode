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

// Data structure for the top-level workspace json files
type Workspace struct {
	Folders  []Folder `json:"folders"`
	Settings Settings `json:"settings"`
}
type Folder struct {
	Path string `json:"path"`
}
type Settings struct{}

// Appends folders to a workspace
// Not sure if channels are really necessary
func AppendDirsToWorkspace(path *string, wg *sync.WaitGroup, dirs *[]string) {
	defer wg.Done()

	content, err := ioutil.ReadFile(*path)

	if err != nil {
		return
	}

	// number of folders that aren't already in the workspace
	newFolders := false

	// declares a workspace, loads json to it from the earlier file read then continues if no err
	var workspace Workspace
	if err := json.Unmarshal(content, &workspace); err == nil {
		// channel used to send folders to add to the workspace defined 2 lines ago
		ch := make(chan *Folder, 5)

		// these next vars only ever on a single gorouting so shouldn't race
		todo := 0 // number of folders to check
		done := 0 // number of folders checked

		// for each dir start goroutine that checks if
		for _, dir := range *dirs {
			go func(ch chan<- *Folder, workspace Workspace, path string, dir string) {
				if !WorkspaceContainsDir(&workspace, &path, &dir) {
					f := Folder{dir}
					ch <- &f
				} else {
					ch <- nil
				}
			}(ch, workspace, *path, dir)
			todo++
		}

		// loop for receiving on the channel
		for {
			// breaks if all folders received
			if done >= todo {
				close(ch)
				break
			}

			// await folder from channel
			folder := <-ch

			// appends folder if not nil to workspace
			if folder != nil {
				workspace.Folders = append(workspace.Folders, *folder)
				newFolders = true
			}

			done++
		}

	}

	// exits out of program if there weren't any new folders located
	if !newFolders {
		fmt.Printf("No new folders found, no write performed to %s\n", *path)
		return
	}

	// saves changes
	newJSONBytes, err := json.MarshalIndent(workspace, "", "    ")
	if err == nil {
		SaveJSON(*path, &newJSONBytes)
	}
}

// Checks if a Workspace struct already contains a given directory
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

// Writes to path from a slice of bytes then logs the result of that write
func SaveJSON(path string, jsonBytes *[]byte) {
	jsonString := string(*jsonBytes)

	if file, err := os.Create(path); err == nil {
		defer file.Close()
		if len, err := file.WriteString(jsonString); err == nil {
			fmt.Printf("Wrote %v characters to %s\n", len, path)
		} else {
			fmt.Printf("Failed to write to %s\n", path)
		}
	}
}
