package anderson

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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

func (l PackageLister) listDeps(pkgs []*Package) ([]string, error) {
	var err error
	var path, seen []string

	dependencies := []string{}

	for _, p := range pkgs {
		if p.Standard {
			continue
		}

		if p.Error.Err != "" {
			err = errors.New("error loading packages")
			continue
		}

		path = append(path, p.Deps...)
	}

	if err != nil {
		return []string{}, err
	}

	var testImports []string
	for _, p := range pkgs {
		testImports = append(testImports, p.TestImports...)
		testImports = append(testImports, p.XTestImports...)
	}

	ps, err := l.loadPackages(testImports...)
	if err != nil {
		return []string{}, err
	}

	for _, p := range ps {
		if p.Standard {
			continue
		}

		if p.Error.Err != "" {
			log.Println(p.Error.Err)
			err = errors.New("error loading packages")
			continue
		}

		path = append(path, p.ImportPath)
		path = append(path, p.Deps...)
	}

	if err != nil {
		return []string{}, err
	}

	sort.Strings(path)
	path = uniq(path)

	ps, err = l.loadPackages(path...)
	if err != nil {
		return []string{}, err
	}

	for _, pkg := range ps {
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

		seen = append(seen, pkg.ImportPath)
		dependencies = append(dependencies, pkg.ImportPath)
	}

	if err != nil {
		return []string{}, err
	}

	return dependencies, nil
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
