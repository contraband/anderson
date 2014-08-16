package anderson

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func findPackage(packagePath string) error {
	d, err := os.Stat(packagePath)
	if err != nil {
		return err
	}
	if m := d.Mode(); m.IsDir() {
		return nil
	}
	return os.ErrPermission
}

func LookGopath(packagePath string) (string, error) {
	gopathenv := os.Getenv("GOPATH")
	if gopathenv == "" {
		return "", errors.New("GOPATH not set")
	}

	for _, dir := range strings.Split(gopathenv, ":") {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, "src", packagePath)
		if err := findPackage(path); err == nil {
			return path, nil
		}
	}

	return "", errors.New("could not find package in GOPATH")
}
