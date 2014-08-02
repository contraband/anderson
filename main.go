package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fraenkel/candiedyaml"
	"github.com/mitchellh/colorstring"
	"github.com/ryanuber/go-license"
)

type Config struct {
	Whitelist  []string `yaml:"whitelist"`
	Greylist   []string `yaml:"greylist"`
	Blacklist  []string `yaml:"blacklist"`
	Exceptions []string `yaml:"exceptions"`
}

type Godeps struct {
	Deps []Dependency
}

type Dependency struct {
	ImportPath string
}

type GodepsLister struct{}

func (l GodepsLister) ListDependencies() []string {
	godepsFile, err := os.Open("Godeps/Godeps.json")
	if err != nil {
		fatal("Couldn't find your Godeps.json file!")
	}
	defer godepsFile.Close()

	var godep Godeps
	if err := json.NewDecoder(godepsFile).Decode(&godep); err != nil {
		fatal("Your Godeps file wasn't valid JSON!")
	}

	deps := []string{}

	for _, dep := range godep.Deps {
		deps = append(deps, dep.ImportPath)
	}

	return deps
}

type LicenseClassifier struct {
	Config Config
}

func (c LicenseClassifier) Classify(path string, importPath string) LicenseStatus {
	l, err := license.NewFromDir(path)

	if err != nil {
		switch err.Error() {
		case license.ErrNoLicenseFile:
			return LicenseTypeNoLicense
		case license.ErrUnrecognizedLicense:
			return LicenseTypeUnknown
		default:
			fatal(fmt.Sprintf("Could not determine license for: %s", importPath))
		}
	}

	if contains(c.Config.Blacklist, l.Type) {
		return LicenseTypeBanned
	}

	if contains(c.Config.Whitelist, l.Type) {
		return LicenseTypeAllowed
	}

	if contains(c.Config.Greylist, l.Type) {
		if contains(c.Config.Exceptions, importPath) {
			return LicenseTypeAllowed
		} else {
			return LicenseTypeMarginal
		}
	}

	panic("uh oh")
}

type LicenseStatus int

func (s LicenseStatus) Color() string {
	switch s {
	case LicenseTypeUnknown:
		return "magenta"
	case LicenseTypeNoLicense:
		return "cyan"
	case LicenseTypeAllowed:
		return "green"
	case LicenseTypeBanned:
		return "read"
	case LicenseTypeMarginal:
		return "yellow"
	default:
		return "red"
	}
}

func (s LicenseStatus) Message() string {
	switch s {
	case LicenseTypeUnknown:
		return "UNKNOWN"
	case LicenseTypeNoLicense:
		return "NO LICENSE"
	case LicenseTypeAllowed:
		return "CHECKS OUT"
	case LicenseTypeBanned:
		return "CONTRABAND"
	case LicenseTypeMarginal:
		return "BORDERLINE"
	default:
		return "ERROR"
	}
}

func (s LicenseStatus) FailsBuild() bool {
	switch s {
	case LicenseTypeUnknown:
		return true
	case LicenseTypeNoLicense:
		return true
	case LicenseTypeAllowed:
		return false
	case LicenseTypeBanned:
		return true
	case LicenseTypeMarginal:
		return true
	default:
		return true
	}
}

const (
	LicenseTypeUnknown LicenseStatus = iota
	LicenseTypeNoLicense
	LicenseTypeBanned
	LicenseTypeAllowed
	LicenseTypeMarginal
)

func main() {
	license.DefaultLicenseFiles = []string{
		"LICENSE", "LICENSE.txt", "LICENSE.md", "license.txt",
		"COPYING", "COPYING.txt", "COPYING.md", "copying.txt",
		"MIT.LICENSE",
	}

	say("[blue]> Hold still citizen, scanning dependencies for contraband...")

	configFile, err := os.Open(".anderson.yml")
	if err != nil {
		fatal("You seem to be missing your .anderson.yml...")
	}

	var config Config
	if err := candiedyaml.NewDecoder(configFile).Decode(&config); err != nil {
		panic(err)
	}

	lister := GodepsLister{}
	classifier := LicenseClassifier{
		Config: config,
	}

	failed := false
	for _, importPath := range lister.ListDependencies() {
		path, err := LookGopath(importPath)
		if err != nil {
			fatal(fmt.Sprintf("Could not find %s in your GOPATH...", importPath))
		}

		licenseType := classifier.Classify(path, importPath)
		failed = failed || licenseType.FailsBuild()
		message := licenseType.Message()

		whitespace := strings.Repeat(" ", 80-len(message)-len(importPath))
		say(fmt.Sprintf("[white]%s%s[%s]%s", importPath, whitespace, licenseType.Color(), message))
	}

	if failed {
		os.Exit(1)
	}
}

func fatal(message string) {
	say(fmt.Sprintf("[red]> %s", message))
	os.Exit(1)
}

func say(message string) {
	fmt.Println(colorstring.Color(message))
}

func contains(haystack []string, needle string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}
