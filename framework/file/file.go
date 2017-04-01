package file

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//get filename extensions
func GetFileExt(fileName string) string {
	if fileName == "" {
		return ""
	} else {
		index := strings.LastIndex(fileName, ".")
		if index < 0 {
			return ""
		} else {
			return string(fileName[index:])
		}
	}
}

//check filename is exist
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
