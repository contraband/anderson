package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fraenkel/candiedyaml"
	"github.com/mitchellh/colorstring"
	"github.com/ryanuber/go-license"

	"github.com/xoebus/anderson/anderson"
)

type Config struct {
	Whitelist  []string `yaml:"whitelist"`
	Blacklist  []string `yaml:"blacklist"`
	Exceptions []string `yaml:"exceptions"`
}

type LicenseClassifier struct {
	Config Config
}

func (c LicenseClassifier) Classify(path string, importPath string) (LicenseStatus, string) {
	l, err := license.NewFromDir(path)

	if err != nil {
		switch err.Error() {
		case license.ErrNoLicenseFile:
			if contains(c.Config.Exceptions, importPath) {
				return LicenseTypeAllowed, "Unknown"
			}

			return LicenseTypeNoLicense, "Unknown"
		case license.ErrUnrecognizedLicense:
			if contains(c.Config.Exceptions, importPath) {
				return LicenseTypeAllowed, "Unknown"
			}

			return LicenseTypeUnknown, "Unknown"
		default:
			fatal(fmt.Sprintf("Could not determine license for: %s", importPath))
		}
	}

	if contains(c.Config.Blacklist, l.Type) {
		return LicenseTypeBanned, l.Type
	}

	if contains(c.Config.Whitelist, l.Type) {
		return LicenseTypeAllowed, l.Type
	}

	if contains(c.Config.Exceptions, importPath) {
		return LicenseTypeAllowed, l.Type
	} else {
		return LicenseTypeMarginal, l.Type
	}
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
		return "red"
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
	say("[blue]> Hold still citizen, scanning dependencies for contraband...")

	emptyConfig := true
	config := Config{}
	configFile, err := os.Open(".anderson.yml")
	if err == nil {
		if err := candiedyaml.NewDecoder(configFile).Decode(&config); err != nil {
			fatal("Looks like your .anderson.yml file is invalid YAML!")
		}

		emptyConfig = false
	}

	lister := anderson.GodepsLister{}
	classifier := LicenseClassifier{
		Config: config,
	}

	failed := false

	dependencies, err := lister.ListDependencies()
	if err != nil {
		fatal(err.Error())
	}

	for _, importPath := range dependencies {
		path, err := anderson.LookGopath(importPath)
		if err != nil {
			fatal(fmt.Sprintf("Could not find %s in your GOPATH...", importPath))
		}

		licenseType, licenseName := classifier.Classify(path, importPath)
		failed = failed || licenseType.FailsBuild()

		var message string
		var color string

		if emptyConfig {
			message = licenseName
			color = "white"
		} else {
			message = licenseType.Message()
			color = licenseType.Color()
		}

		whitespace := strings.Repeat(" ", 80-len(message)-len(importPath))
		say(fmt.Sprintf("[white]%s%s[%s]%s", importPath, whitespace, color, message))
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
