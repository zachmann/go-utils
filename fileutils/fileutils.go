package fileutils

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func evalSymlink(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '~' {
		path = os.Getenv("HOME") + path[1:]
	}
	evalPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return evalPath
}

// FileExists checks if a given file exists.
func FileExists(path string) bool {
	if _, err := os.Stat(evalSymlink(path)); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.WithError(err).Error()
		return false
	}
}

// ReadFile reads a given file and returns the content.
func ReadFile(filename string) ([]byte, error) {
	filename = evalSymlink(filename)
	log.WithField("filepath", filename).Trace("Reading file...")
	return os.ReadFile(filename)
}

// ReadFileFromLocations checks if a file exists in one of the passed
// directories and returns the content. If no file is found, nil is returned
func ReadFileFromLocations(filename string, locations []string) ([]byte, string) {
	for _, dir := range locations {
		if strings.HasPrefix(dir, "~") {
			homeDir := os.Getenv("HOME")
			dir = filepath.Join(homeDir, dir[1:])
		}
		filep := filepath.Join(dir, filename)
		log.WithField("filepath", filep).Debug("Try to read file")
		data, err := ReadFile(filep)
		if err == nil {
			return data, dir
		}
	}
	return nil, ""
}
