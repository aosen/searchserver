package utils

import (
	"log"
	"os"
	"path/filepath"
)

func GetFilelist(path string) []string {
	pstr := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		pstr = append(pstr, path)
		return nil
	})
	PutError(err)
	return pstr
}

func PutError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
