package util

import (
	"os"
	"path/filepath"
)

// AbsolutePathHelper Takes a relative path and string parts and utilizes
// path.filepath.Abs and Join to return a valid UNC path.
func AbsolutePathHelper(rel string, parts ...string) (path string, err error) {
	abs, err := filepath.Abs(rel)
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{abs}, parts...)...), nil
}

// MakeDirectory creates the directory if it doesn't already exist
func MakeDirectory(path string) (err error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// if it exists we ignore the error, by not creating the dir
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		os.Mkdir(abs, 0777)
	}
	return nil
}
