package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fraenkel/candiedyaml"
	"github.com/mitchellh/colorstring"

	"github.com/xoebus/anderson/anderson"
)

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

	for _, importPath := range dependencies {
		path, err := anderson.LookGopath(importPath)
		if err != nil {
			fatal(fmt.Errorf("Could not find %s in your GOPATH...", importPath))
		}

		licenseType, licenseName, err := classifier.Classify(path, importPath)
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

func fatal(err error) {
	say(fmt.Sprintf("[red]> %s", err))
	os.Exit(1)
}

func say(message string) {
	fmt.Println(colorstring.Color(message))
}
