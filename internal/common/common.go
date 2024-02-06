package common

import (
	"log"
	"os"
)

func IfFileExistsReplaceStringVal(l *string) error {
	log.Println("StringOrFile : ", *l)
	fileInfo, err := os.Stat(*l)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	if fileInfo.IsDir() {
		return nil
	}
	data, err := os.ReadFile(*l)
	if err != nil {
		return err
	}
	*l = string(data)
	return nil
}
