package utils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func ReplaceString(file string, from string, to string) {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s with error %v", file, err)
	}
	newContents := strings.Replace(string(read), from, to, -1)
	err = ioutil.WriteFile(file, []byte(newContents), 0)
	if err != nil {
		log.Fatalf("Failed to write file %s with error %v", file, err)
	}
}

func Download(fileName string, url string, mode os.FileMode) {
	log.Printf("Downloading %s to location %s", url, fileName)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		log.Printf("File already exists %s\n", fileName)
		if err := os.Chmod(fileName, mode); err != nil {
			log.Fatalf("Failed to set permissions for file %s, with error %v\n", fileName, err)
		}
	} else {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Failed to create file %s with err %v\n", fileName, err)
		}
		defer file.Close()
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to download %s with error %v\n", url, err)
		}
		defer response.Body.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("Failed to save file %s with error %v\n", fileName, err)
		}
	}
	if err := os.Chmod(fileName, mode); err != nil {
		log.Fatalf("Failed to set permissions for file %s, with error %v\n", fileName, err)
	}
}
