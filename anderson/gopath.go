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
	paths, err := Gopaths()
	if err != nil {
		return "", err
	}

	for _, dir := range paths {
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

func Gopaths() ([]string, error) {
	gopathenv := os.Getenv("GOPATH")
	if gopathenv == "" {
		return []string{}, errors.New("GOPATH not set")
	}

	return strings.Split(gopathenv, ":"), nil
}
