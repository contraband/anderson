package anderson

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Package struct {
	Dir        string
	Root       string
	ImportPath string
	Deps       []string
	Standard   bool

	TestGoFiles  []string
	TestImports  []string
	XTestGoFiles []string
	XTestImports []string

	Error struct {
		Err string
	}
}

type PackageLister struct{}

func (l PackageLister) ListDependencies() ([]string, error) {
	packages, err := l.loadPackages("./...")
	if err != nil {
		return []string{}, err
	}

	return l.listDeps(packages)
}

func (l PackageLister) loadPackages(name ...string) (packages []*Package, err error) {
	if len(name) == 0 {
		return nil, nil
	}

	args := []string{"list", "-e", "-json"}
	cmd := exec.Command("go", append(args, name...)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(stdout)
	for {
		info := new(Package)
		err = decoder.Decode(info)
		if err == io.EOF {
			break
		}
		if err != nil {
			info.Error.Err = err.Error()
		}
		packages = append(packages, info)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func (l PackageLister) listDeps(packages []*Package) ([]string, error) {
	var err error
	var path, seen []string

	dependencies := []string{}

	for _, pkg := range packages {
		if pkg.Standard {
			continue
		}

		if pkg.Error.Err != "" {
			err = errors.New("error loading packages")
			continue
		}

		path = append(path, pkg.Deps...)
	}

	if err != nil {
		return []string{}, err
	}

	var testImports []string
	for _, pkg := range packages {
		testImports = append(testImports, pkg.TestImports...)
		testImports = append(testImports, pkg.XTestImports...)
	}

	testPackages, err := l.loadPackages(testImports...)
	if err != nil {
		return []string{}, err
	}

	for _, pkg := range testPackages {
		if pkg.Standard {
			continue
		}

		if pkg.Error.Err != "" {
			err = errors.New("error loading packages")
			continue
		}

		path = append(path, pkg.ImportPath)
		path = append(path, pkg.Deps...)
	}

	if err != nil {
		return []string{}, err
	}

	sort.Strings(path)
	path = uniq(path)

	allPackages, err := l.loadPackages(path...)
	if err != nil {
		return []string{}, err
	}

	currentName, err := currentPackageName()
	if err != nil {
		return []string{}, err
	}

	for _, pkg := range allPackages {
		if pkg.Error.Err != "" {
			err = errors.New("error loading dependencies")
			continue
		}

		if pkg.Standard {
			continue
		}

		if containsPathPrefix(seen, pkg.ImportPath) {
			continue
		}

		if strings.HasPrefix(pkg.ImportPath, currentName) {
			continue
		}

		seen = append(seen, pkg.ImportPath)
		dependencies = append(dependencies, pkg.ImportPath)
	}

	if err != nil {
		return []string{}, err
	}

	return dependencies, nil
}

func currentPackageName() (string, error) {
	cmd := exec.Command("go", "list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func uniq(a []string) []string {
	i := 0
	s := ""
	for _, t := range a {
		if t != s {
			a[i] = t
			i++
			s = t
		}
	}
	return a[:i]
}

func containsPathPrefix(pats []string, s string) bool {
	for _, pat := range pats {
		if pat == s || strings.HasPrefix(s, pat+"/") {
			return true
		}
	}
	return false
}
