package utils

import (
	"errors"
	"os"
)

func IsFileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
func GetCWD() string {
	path, err := os.Getwd()
	if err != nil {
	}
	return path
}

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
	}
	return homeDir
}
