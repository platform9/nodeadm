package utils

import (
	"io/ioutil"
	"log"
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
