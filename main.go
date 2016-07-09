package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/mitchellh/colorstring"

	"github.com/contraband/anderson/anderson"
)

type License struct {
	Type anderson.LicenseStatus
	Name string
}

type Lister interface {
	ListDependencies() ([]string, error)
}

func main() {
	config, missingConfig := loadConfig()
	lister := lister()
	classifier := anderson.LicenseClassifier{
		Config: config,
	}

	info("Hold still citizen, scanning dependencies for contraband...")
	dependencies, err := lister.ListDependencies()
	if err != nil {
		fatalf(err.Error())
	}

	failed := false
	classified := map[string]License{}
	for _, importPath := range dependencies {
		path, err := anderson.LookGopath(importPath)
		if err != nil {
			fatalf("Could not find %s in your GOPATH...", importPath)
		}

		licenseType, licenseDeclarationPath, licenseName, err := classifier.Classify(path, importPath)
		failed = failed || licenseType.FailsBuild()

		containingGopath, err := anderson.ContainingGopath(importPath)
		if err != nil {
			fatalf("Unable to find containing GOPATH for %s: %s", licenseDeclarationPath, err)
		}

		relPath, err := filepath.Rel(filepath.Join(containingGopath, "src"), licenseDeclarationPath)
		if err != nil {
			fatalf("Unable to create relative path for %s: %s", licenseDeclarationPath, err)
		}

		classified[relPath] = License{
			Type: licenseType,
			Name: licenseName,
		}
	}

	for relPath, license := range classified {
		var message string
		var messageLen int

		if missingConfig {
			message = fmt.Sprintf("[white]%s", license.Name)
			messageLen = len(license.Name)
		} else {
			message = fmt.Sprintf("(%s) [%s]%10s", license.Name, license.Type.Color(), license.Type.Message())
			messageLen = len(license.Name) + len("() ") + 9 // length of all messages
		}

		totalSize := messageLen + len(relPath)
		whitespace := " "
		if totalSize < 80 {
			whitespace = strings.Repeat(" ", 80-totalSize)
		}

		say(fmt.Sprintf("[white]%s%s%s", relPath, whitespace, message))
	}

	if failed {
		os.Exit(1)
	}
}

func loadConfig() (config anderson.Config, missing bool) {
	configFile, err := os.Open(".anderson.yml")
	if err != nil {
		return config, true
	}

	if err := candiedyaml.NewDecoder(configFile).Decode(&config); err != nil {
		fatalf("Looks like your .anderson.yml file is invalid YAML!")
	}

	return config, false
}

func lister() Lister {
	if isStdinPipe() {
		return anderson.StdinLister{}
	}

	return anderson.PackageLister{}
}

func isStdinPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func fatalf(err string, args ...interface{}) {
	message := fmt.Sprintf(err, args)
	say(fmt.Sprintf("[red]> %s", message))
	os.Exit(1)
}

func info(message string) {
	say(fmt.Sprintf("[blue]> %s", message))
}

func say(message string) {
	fmt.Println(colorstring.Color(message))
}
