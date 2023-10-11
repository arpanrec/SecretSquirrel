package iss

import (
	"log"
	"os"
	"path/filepath"
)

func GetData(p string) (string, error) {

	// Read the file and return the contents
	log.Println("Reading file: ", p)
	d, err := os.ReadFile(p)
	return string(d), err
}

func PutData(p string, data string) (bool, error) {
	dir := filepath.Dir(p)
	log.Println("Creating directory: ", dir)
	errMakeDir := os.MkdirAll(dir, 0755)
	if errMakeDir != nil {
		log.Fatal(errMakeDir)
		return false, errMakeDir
	}
	log.Printf("Writing file: %s", p)
	errWriteFile := os.WriteFile(p, []byte(data), 0644)
	if errWriteFile != nil {
		log.Fatal(errWriteFile)
		return false, errWriteFile
	}
	return true, nil
}
