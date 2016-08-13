package util

import (
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
