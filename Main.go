package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func main() {
	srcDirectories := GetGoSrc()
	var mods, workspaces []string
	var wg sync.WaitGroup

	for _, dir := range srcDirectories {
		wg.Add(1)
		go HandleDir(dir, &mods, &workspaces, &wg)
	}

	wg.Wait()

	for _, workspace := range workspaces {
		wg.Add(1)
		go AppendDirsToWorkspace(&workspace, &wg, &mods)
	}

	wg.Wait()

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
	} else if len(workspaces) == 1 {
		if err := exec.Command("code", workspaces[0]).Run(); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("More than one workspace found, not opening")
	}
}

func HandleDir(dir string, mods, workspaces *[]string, wg *sync.WaitGroup) {
	if newMods, err := GetMods(dir); err == nil {
		*mods = append(*mods, newMods...)
	}

	if newWorkspaces, err := GetWorkspaces(dir); err == nil {
		*workspaces = append(*workspaces, newWorkspaces...)
	}

	(*wg).Done()
}
