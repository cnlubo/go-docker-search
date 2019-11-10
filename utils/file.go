package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// check file is or is not exist
func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// check directory is or is not empty
func IsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// return relative file path if filepath is absolute
func ParseAbsPath(filepath string, homeDir string) string {

	var result string
	if ok := path.IsAbs(filepath); ok {
		if strings.HasPrefix(filepath, homeDir) {
			result = strings.Replace(filepath, homeDir, "~", 1)
		} else {
			result = filepath
		}
	} else {
		result = filepath
	}
	return result

}

// return absolute file path if filepath is relative
func ParseRelPath(filepath string, homeDir string) string {

	var result string
	if ok := path.IsAbs(filepath); !ok {
		result = strings.Replace(filepath, "~", homeDir, 1)
	} else {
		result = filepath
	}
	return result

}

func IsSymLink(filepath string) (error, bool) {

	fileInfo, err := os.Lstat(filepath)

	if err != nil {
		return err, false
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return nil, true
	} else {
		return nil, false
	}

}

// return the original file path if file  is a symbol link
func ParseOriginalFilePath(filepath string) string {

	err, link := IsSymLink(filepath)
	if err != nil {
		return ""
	}
	if link {
		originFile, err := os.Readlink(filepath)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return originFile
	}
	return filepath
}
