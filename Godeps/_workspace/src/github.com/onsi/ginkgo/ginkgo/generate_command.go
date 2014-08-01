package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func BuildGenerateCommand() *Command {
	var noDot bool
	flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
	flagSet.BoolVar(&noDot, "nodot", false, "If set, generate will generate a test file that does not . import ginkgo and gomega")

	return &Command{
		Name:         "generate",
		FlagSet:      flagSet,
		UsageCommand: "ginkgo generate <filename(s)>",
		Usage: []string{
			"Generate a test file named filename_test.go",
			"If the optional <filenames> argument is omitted, a file named after the package in the current directory will be created.",
			"Accepts the following flags:",
		},
		Command: func(args []string, additionalArgs []string) {
			generateSpec(args, noDot)
		},
	}
}

var specText = `package {{.Package}}_test

import (
	. "{{.PackageImportPath}}"

	{{if .IncludeImports}}. "github.com/onsi/ginkgo"{{end}}
	{{if .IncludeImports}}. "github.com/onsi/gomega"{{end}}
)

var _ = Describe("{{.Subject}}", func() {

})
`

type specData struct {
	Package           string
	Subject           string
	PackageImportPath string
	IncludeImports    bool
}

func generateSpec(args []string, noDot bool) {
	if len(args) == 0 {
		err := generateSpecForSubject("", noDot)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("")
			os.Exit(1)
		}
		fmt.Println("")
		return
	}

	var failed bool
	for _, arg := range args {
		err := generateSpecForSubject(arg, noDot)
		if err != nil {
			failed = true
			fmt.Println(err.Error())
		}
	}
	fmt.Println("")
	if failed {
		os.Exit(1)
	}
}

func generateSpecForSubject(subject string, noDot bool) error {
	packageName := getPackage()
	if subject == "" {
		subject = packageName
	} else {
		subject = strings.Split(subject, ".go")[0]
		subject = strings.Split(subject, "_test")[0]
	}

	formattedSubject := strings.Replace(strings.Title(strings.Replace(subject, "_", " ", -1)), " ", "", -1)

	data := specData{
		Package:           packageName,
		Subject:           formattedSubject,
		PackageImportPath: getPackageImportPath(),
		IncludeImports:    !noDot,
	}

	targetFile := fmt.Sprintf("%s_test.go", subject)
	if fileExists(targetFile) {
		return fmt.Errorf("%s already exists.", targetFile)
	} else {
		fmt.Printf("Generating ginkgo test for %s in:\n  %s\n", data.Subject, targetFile)
	}

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	specTemplate, err := template.New("spec").Parse(specText)
	if err != nil {
		return err
	}

	specTemplate.Execute(f, data)
	goFmt(targetFile)
	return nil
}

func getPackageImportPath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	sep := string(filepath.Separator)
	paths := strings.Split(workingDir, sep+"src"+sep)
	if len(paths) == 1 {
		fmt.Printf("\nCouldn't identify package import path.\n\n\tginkgo generate\n\nMust be run within a package directory under $GOPATH/src/...\nYou're going to have to change UNKNOWN_PACKAGE_PATH in the generated file...\n\n")
		return "UNKNOWN_PACKAGE_PATH"
	}
	return filepath.ToSlash(paths[len(paths)-1])
}
