package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Create symlinks of all the files inside sourceDir to targetDir
func CreateSymLinks(sourceDir, targetDir string) {
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}
	_, parentDir := filepath.Split(sourceDir)

	for _, f := range files {
		log.Print("Creating symlink for " + f.Name())
		log.Println()

		err = os.Symlink(filepath.Join(parentDir, f.Name()), filepath.Join(targetDir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
