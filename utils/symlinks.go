package utils

import (
	"io/ioutil"
	"os"
	"log"
	"path/filepath"
)

// Create symlinks of all the files inside sourceDir to targetDir
func CreateSymLinks(sourceDir, targetDir string) {
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Print("Creating symlink for " + f.Name())
		err = os.Symlink(filepath.Join(sourceDir,f.Name()), filepath.Join(targetDir,f.Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
