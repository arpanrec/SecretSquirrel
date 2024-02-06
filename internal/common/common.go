package common

import (
	"log"
	"os"
)

func ReadFileSureOrStop(l *string) []byte {
	log.Println("Reading file: ", *l)
	_, err := os.Stat(*l)
	if os.IsNotExist(err) {
		log.Fatalln("File does not exist: ", *l)
	}
	b, err := os.ReadFile(*l)
	if err != nil {
		log.Fatalln("Error reading file: ", err)
	}
	return b
}

func ReadFileStringSureOrStop(l *string) string {
	return string(ReadFileSureOrStop(l))
}
