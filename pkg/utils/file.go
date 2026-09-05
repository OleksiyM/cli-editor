package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadFile reads the content of a file and returns it as a string
func ReadFile(filename string) (string, error) {
	content, err := ioutil.ReadAll(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file
func WriteFile(filename string, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

// GetFileExtension returns the extension of a file
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// GetFileName returns the name of a file without the path
func GetFileName(filename string) string {
	return filepath.Base(filename)
}

// GetFileDir returns the directory of a file
func GetFileDir(filename string) string {
	return filepath.Dir(filename)
}