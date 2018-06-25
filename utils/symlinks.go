package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Create symlinks of all the files inside sourceDir to targetDir
func CreateSymLinks(sourceDir, targetDir string, overwriteSymlinks bool) {
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}
	_, parentDir := filepath.Split(sourceDir)

	for _, f := range files {
		log.Print("Creating symlink for " + f.Name())

		symlinkPath := filepath.Join(targetDir, f.Name())

		if overwriteSymlinks {
			if _, err := os.Lstat(symlinkPath); err == nil {
				os.Remove(symlinkPath)
			}
		}

		err = os.Symlink(filepath.Join(parentDir, f.Name()), symlinkPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
