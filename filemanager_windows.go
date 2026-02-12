package main

import "os/exec"

// openFileManager opens the given directory in Explorer on Windows.
func openFileManager(dir string) error {
	return exec.Command("explorer", dir).Start()
}
