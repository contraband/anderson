package license

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// Recognized license types
	LicenseMIT      = "MIT"
	LicenseNewBSD   = "NewBSD"
	LicenseFreeBSD  = "FreeBSD"
	LicenseApache20 = "Apache-2.0"
	LicenseMPL20    = "MPL-2.0"
	LicenseGPL20    = "GPL-2.0"
	LicenseGPL30    = "GPL-3.0"
	LicenseLGPL21   = "LGPL-2.1"
	LicenseLGPL30   = "LGPL-3.0"
	LicenseCDDL10   = "CDDL-1.0"
	LicenseEPL10    = "EPL-1.0"

	// Various errors
	ErrNoLicenseFile       = "license: unable to find any license file"
	ErrUnrecognizedLicense = "license: could not guess license type"
)

// A set of reasonable license file names to use when guessing where the
// license may be.
var DefaultLicenseFiles = []string{
	"LICENSE", "LICENSE.txt", "LICENSE.md", "license.txt",
	"COPYING", "COPYING.txt", "COPYING.md", "copying.txt",
}

// A slice of standardized license abbreviations
var KnownLicenses = []string{
	LicenseMIT,
	LicenseNewBSD,
	LicenseFreeBSD,
	LicenseApache20,
	LicenseMPL20,
	LicenseGPL20,
	LicenseGPL30,
	LicenseLGPL21,
	LicenseLGPL30,
	LicenseCDDL10,
	LicenseEPL10,
}

// License describes a software license
type License struct {
	Type string // The type of license in use
	Text string // License text data
	File string // The path to the source file, if any
}

// New creates a new License from explicitly passed license type and data
func New(licenseType, licenseText string) *License {
	l := &License{
		Type: licenseType,
		Text: licenseText,
	}
	return l
}

// NewFromFile will attempt to load a license from a file on disk, and guess the
// type of license based on the bytes read.
func NewFromFile(path string) (*License, error) {
	licenseText, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	l := &License{
		Text: string(licenseText),
		File: path,
	}

	if err := l.GuessType(); err != nil {
		return nil, err
	}

	return l, nil
}

// NewFromDir will search a directory for well-known and accepted license file
// names, and if one is found, read in its content and guess the license type.
func NewFromDir(dir string) (*License, error) {
	l := new(License)

	if err := l.GuessFile(dir); err != nil {
		return nil, err
	}

	return NewFromFile(l.File)
}

// Recognized determines if the license is known to go-license.
func (l *License) Recognized() bool {
	for _, license := range KnownLicenses {
		if license == l.Type {
			return true
		}
	}
	return false
}

// GuessFile searches a given directory (non-recursively) for files with well-
// established names that indicate license content.
func (l *License) GuessFile(dir string) error {
	d, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !d.IsDir() {
		return fmt.Errorf("license: cannot search %s: not a directory", dir)
	}

	for _, file := range DefaultLicenseFiles {
		filePath := filepath.Join(dir, file)
		_, err := os.Stat(filePath)
		if err == nil {
			l.File = filePath
			return nil
		}
	}
	return errors.New(ErrNoLicenseFile)
}

// GuessType will scan license text and attempt to guess what license type it
// describes. It will return the license type on success, or an error if it
// cannot accurately guess the license type.
func (l *License) GuessType() error {
	switch {
	case scanLeft(l.Text, "The MIT License"):
		l.Type = LicenseMIT

	case scanLeft(l.Text, "Apache License"):
		switch {
		case scanLeft(l.Text, "Version 2.0"):
			l.Type = LicenseApache20
		}

	case scanLeft(l.Text, "GNU GENERAL PUBLIC LICENSE"):
		switch {
		case scanLeft(l.Text, "Version 2"):
			l.Type = LicenseGPL20
		case scanLeft(l.Text, "Version 3"):
			l.Type = LicenseGPL30
		}

	case scanLeft(l.Text, "GNU LESSER GENERAL PUBLIC LICENSE"):
		switch {
		case scanLeft(l.Text, "Version 2.1"):
			l.Type = LicenseLGPL21
		case scanLeft(l.Text, "Version 3"):
			l.Type = LicenseLGPL30
		}

	case scanLeft(l.Text, "Mozilla Public License Version 2.0"):
		l.Type = LicenseMPL20

	case scanLeft(l.Text, "Redistribution and use in source and binary forms"):
		switch {
		case scanLeft(l.Text, "* Redistribution"):
			l.Type = LicenseNewBSD
		case scanRight(l.Text, "FreeBSD Project."):
			l.Type = LicenseFreeBSD
		}

	case scanRight(l.Text, "(CDDL)"):
		switch {
		case scanRight(l.Text, "Version 1.0"):
			l.Type = LicenseCDDL10
		}

	case scanLeft(l.Text, "Eclipse Public License - v 1.0"):
		l.Type = LicenseEPL10

	default:
		return errors.New(ErrUnrecognizedLicense)
	}

	return nil
}
