package main

import "os/exec"

// openFileManager opens the given directory in the default file manager on Linux.
func openFileManager(dir string) error {
	return exec.Command("xdg-open", dir).Start()
}
