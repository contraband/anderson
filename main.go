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

func main() {
	say("[blue]> Hold still citizen, scanning dependencies for contraband...")

	configFile, err := os.Open(".anderson.yml")
	if err != nil {
		panic(err)
	}

	var config Config
	if err := candiedyaml.NewDecoder(configFile).Decode(&config); err != nil {
		panic(err)
	}

	godepsFile, err := os.Open("Godeps/Godeps.json")
	if err != nil {
		panic(err)
	}

	var godep Godeps
	if err := json.NewDecoder(godepsFile).Decode(&godep); err != nil {
		panic(err)
	}

	for _, dependency := range godep.Deps {
		path, err := LookGopath(dependency.ImportPath)
		if err != nil {
			continue
		}

		l, err := license.NewFromDir(path)
		whitespace := strings.Repeat(" ", 80-10-len(dependency.ImportPath))
		if err != nil {
			if err.Error() == "license: unable to find any license file" {
				say(fmt.Sprintf("[white]%s%s[magenta]NO LICENSE", dependency.ImportPath, whitespace))
			} else {
				panic(err)
			}
			continue
		}

		if contains(config.Blacklist, l.Type) {
			say(fmt.Sprintf("[white]%s%s[red]CONTRABAND", dependency.ImportPath, whitespace))
			continue
		}

		if contains(config.Whitelist, l.Type) {
			say(fmt.Sprintf("[white]%s%s[green]CHECKS OUT", dependency.ImportPath, whitespace))
			continue
		}

		if contains(config.Greylist, l.Type) {
			if contains(config.Exceptions, dependency.ImportPath) {
				say(fmt.Sprintf("[white]%s%s[green]CHECKS OUT", dependency.ImportPath, whitespace))
			} else {
				say(fmt.Sprintf("[white]%s%s[yellow]BORDERLINE", dependency.ImportPath, whitespace))
			}
			continue
		}

		fmt.Println(l.Type)
	}
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
