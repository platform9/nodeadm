package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/platform9/nodeadm/pkg/logrus"
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
	log.Infof("Downloading %s to location %s", url, fileName)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		log.Infof("\nFile already exists %s", fileName)
		if err := os.Chmod(fileName, mode); err != nil {
			log.Fatalf("\nFailed to set permissions for file %s, with error %v", fileName, err)
		}
	} else {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("\nFailed to create file %s with err %v", fileName, err)
		}
		defer file.Close()
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("\nFailed to download %s with error %v", url, err)
		}
		defer response.Body.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("\nFailed to save file %s with error %v", fileName, err)
		}
	}
	if err := os.Chmod(fileName, mode); err != nil {
		log.Fatalf("\nFailed to set permissions for file %s, with error %v", fileName, err)
	}
}
