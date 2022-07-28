package utils

import (
	"os"
	"strings"

	"github.com/pkg/errors"
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

func AssertError(wantError bool, err error) bool {
	if wantError {
		if err == nil {
			return false
		}
	} else {
		if err != nil {
			return false
		}
	}
	return true
}

func SShHostSplit(hostname string) (user, host string) {
	var parts []string
	switch parts = strings.Split(hostname, "@"); {
	case len(parts) == 0:
		return "", parts[0]
	case len(parts) == 1:
		return parts[0], parts[1]
	default:
		return strings.Join(parts[0:len(parts)-1], "@"), parts[len(parts)-1]
	}
}
