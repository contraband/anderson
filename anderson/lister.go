package anderson

import (
	"encoding/json"
	"errors"
	"os"
)

type Godeps struct {
	Deps []Dependency
}

type Dependency struct {
	ImportPath string
}

type GodepsLister struct{}

func (l GodepsLister) ListDependencies() ([]string, error) {
	godepsFile, err := os.Open("Godeps/Godeps.json")
	if err != nil {
		return []string{}, errors.New("Couldn't find your Godeps.json file!")
	}
	defer godepsFile.Close()

	var godep Godeps
	if err := json.NewDecoder(godepsFile).Decode(&godep); err != nil {
		return []string{}, errors.New("Your Godeps file wasn't valid JSON!")
	}

	deps := []string{}

	for _, dep := range godep.Deps {
		deps = append(deps, dep.ImportPath)
	}

	return deps, nil
}
