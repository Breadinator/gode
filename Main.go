package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func main() {
	// Get and parse GOPATH src directories into a slice of strings
	srcDirectories := GetGoSrc()

	var mods, workspaces []string
	var wg sync.WaitGroup

	// populate earlier-declared `mods` and `workspaces` variables with all go modules and workspaces in the source directories
	for _, dir := range srcDirectories {
		wg.Add(1)
		go HandleDir(dir, &mods, &workspaces, &wg)
	}

	wg.Wait()

	// for each workspace found, update it to include every module found
	for _, workspace := range workspaces {
		wg.Add(1)
		go AppendDirsToWorkspace(&workspace, &wg, &mods)
	}

	wg.Wait()

	// launches vscode if only 1 workspace was found
	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
	} else if len(workspaces) == 1 {
		if err := exec.Command("code", workspaces[0]).Run(); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("More than one workspace found, not opening.")
	}
}

// Looks through a given directory then appends all go modules and vscode workspaces to provided slices
func HandleDir(dir string, mods, workspaces *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	if newMods, err := GetMods(dir); err == nil {
		*mods = append(*mods, newMods...)
	}

	if newWorkspaces, err := GetWorkspaces(dir); err == nil {
		*workspaces = append(*workspaces, newWorkspaces...)
	}
}
