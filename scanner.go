package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScanFolder recursively walks a directory and upserts all .glb/.gltf files into the database.
// Returns the number of assets found.
func ScanFolder(db *Database, folder WatchFolder) (int, error) {
	count := 0

	err := filepath.WalkDir(folder.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Skip directories we can't read
			return nil
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".glb" && ext != ".gltf" {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil // skip files we can't stat
		}

		_, err = db.UpsertAsset(path, folder.ID, info.Size(), info.ModTime())
		if err != nil {
			fmt.Printf("warn: failed to upsert %s: %v\n", path, err)
			return nil
		}
		count++
		return nil
	})

	if err != nil {
		return count, fmt.Errorf("walk %s: %w", folder.Path, err)
	}

	// Prune assets that no longer exist on disk
	pruned, _ := db.PruneAssetsForFolder(folder.ID)
	if pruned > 0 {
		fmt.Printf("pruned %d missing assets from %s\n", pruned, folder.Path)
	}

	return count, nil
}

// ScanAllFolders scans every registered watch folder.
func ScanAllFolders(db *Database) error {
	folders, err := db.ListWatchFolders()
	if err != nil {
		return err
	}
	for _, f := range folders {
		count, err := ScanFolder(db, f)
		if err != nil {
			fmt.Printf("error scanning %s: %v\n", f.Path, err)
			continue
		}
		fmt.Printf("scanned %s: %d assets\n", f.Path, count)
	}
	return nil
}
