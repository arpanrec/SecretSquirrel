package common

import (
	"log"
	"os"
	"path/filepath"
)

func DeleteFileSureOrStop(l string) {
	log.Println("Deleting file: ", l)
	_, err := os.Stat(l)
	if os.IsNotExist(err) {
		log.Println("File does not exist: ", l)
	} else {
		log.Println("Deleting file: ", l)
		err := os.Remove(l)
		if err != nil {
			log.Fatalln("Error deleting file: ", err)
		}
	}
}

func ReadFileSureOrStop(l string) []byte {
	if filepath, err := filepath.Abs(l); err != nil {
		log.Fatalln("Error getting absolute path: ", err)
	} else {
		l = filepath
	}
	log.Println("Reading file: ", l)
	_, err := os.Stat(l)
	if os.IsNotExist(err) {
		log.Fatalln("File does not exist: ", l)
	}
	b, err := os.ReadFile(l)
	if err != nil {
		log.Fatalln("Error reading file: ", err)
	}
	return b
}

func ReadFileStringSureOrStop(l string) string {
	return string(ReadFileSureOrStop(l))
}
