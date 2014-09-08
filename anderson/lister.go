package anderson

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Godeps struct {
	Deps []Dependency
}

type Dependency struct {
	ImportPath string
}

type GodepsLister struct{}

func (l GodepsLister) ListDependencies() ([]string, error) {
	out := new(bytes.Buffer)

	cmd := exec.Command(
		"bash",
		"-c",
		"go list -f $'{{range $dep := .Deps}}{{$dep}}\n{{end}}' ./... | "+
			"xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'",
	)

	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	deps := strings.Split(out.String(), "\n")

	// strip trailing linebreak causing extra split
	return deps[0 : len(deps)-1], nil
}
