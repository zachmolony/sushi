package main

import "os/exec"

// openFileManager opens the given directory in Finder on macOS.
func openFileManager(dir string) error {
	return exec.Command("open", dir).Start()
}
