package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fraenkel/candiedyaml"
	"github.com/mitchellh/colorstring"

	"github.com/xoebus/anderson/anderson"
)

type License struct {
	Type anderson.LicenseStatus
	Name string
}

func main() {
	say("[blue]> Hold still citizen, scanning dependencies for contraband...")

	emptyConfig := true
	config := anderson.Config{}
	configFile, err := os.Open(".anderson.yml")
	if err == nil {
		if err := candiedyaml.NewDecoder(configFile).Decode(&config); err != nil {
			fatal(errors.New("Looks like your .anderson.yml file is invalid YAML!"))
		}

		emptyConfig = false
	}

	lister := anderson.GodepsLister{}
	classifier := anderson.LicenseClassifier{
		Config: config,
	}

	failed := false

	dependencies, err := lister.ListDependencies()
	if err != nil {
		fatal(err)
	}

	classified := map[string]License{}
	for _, importPath := range dependencies {
		path, err := anderson.LookGopath(importPath)
		if err != nil {
			fatal(fmt.Errorf("Could not find %s in your GOPATH...", importPath))
		}

		licenseType, licenseDeclarationPath, licenseName, err := classifier.Classify(path, importPath)
		failed = failed || licenseType.FailsBuild()

		containingGopath, err := anderson.ContainingGopath(importPath)
		if err != nil {
			fatal(fmt.Errorf("Unable to find containing GOPATH for %s: %s", licenseDeclarationPath, err))
		}

		relPath, err := filepath.Rel(filepath.Join(containingGopath, "src"), licenseDeclarationPath)
		if err != nil {
			fatal(fmt.Errorf("Unable to create relative path for %s: %s", licenseDeclarationPath, err))
		}

		classified[relPath] = License{Type: licenseType, Name: licenseName}
	}

	for relPath, license := range classified {
		var message string
		var color string

		if emptyConfig {
			message = license.Name
			color = "white"
		} else {
			message = license.Type.Message()
			color = license.Type.Color()
		}

		totalSize := len(message) + len(relPath)
		whitespace := " "
		if totalSize < 80 {
			whitespace = strings.Repeat(" ", 80-totalSize)
		}

		say(fmt.Sprintf("[white]%s%s[%s]%s", relPath, whitespace, color, message))
	}

	if failed {
		os.Exit(1)
	}
}

func fatal(err error) {
	say(fmt.Sprintf("[red]> %s", err))
	os.Exit(1)
}

func say(message string) {
	fmt.Println(colorstring.Color(message))
}
