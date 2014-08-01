package license

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLicense(t *testing.T) {
	l := New("MyLicense", "Some license text.")
	if l.Type != "MyLicense" {
		t.Fatalf("bad license type: %s", l.Type)
	}
	if l.Text != "Some license text." {
		t.Fatalf("bad license text: %s", l.Text)
	}
}

func TestNewFromFile(t *testing.T) {
	f, err := ioutil.TempFile("", "go-license")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f.Name())

	licenseText := "The MIT License (MIT)"

	if _, err := f.WriteString(licenseText); err != nil {
		t.Fatalf("err: %s", err)
	}

	l, err := NewFromFile(f.Name())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if l.Type != "MIT" {
		t.Fatalf("unexpected license type: %s", l.Type)
	}

	if l.Text != licenseText {
		t.Fatalf("unexpected license text: %s", l.Text)
	}

	if l.File != f.Name() {
		t.Fatalf("unexpected file path: %s", l.File)
	}

	// Fails properly if the file doesn't exist
	if _, err := NewFromFile("/tmp/go-license-nonexistent"); err == nil {
		t.Fatalf("expected error loading non-existent file")
	}

	// Fails properly if license type from file is not guessable
	if err := os.Truncate(f.Name(), 0); err != nil {
		t.Fatalf("err: %s", err)
	}
	f.WriteString("No license data")
	if _, err := NewFromFile(f.Name()); err == nil {
		t.Fatalf("expected error guessing license type from non-license file")
	}
}

func TestNewFromDir(t *testing.T) {
	d, err := ioutil.TempDir("", "go-license")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(d)

	// Fails properly if the directory contains no license files
	if _, err := NewFromDir(d); err == nil {
		t.Fatalf("expected error loading empty directory")
	}

	fPath := filepath.Join(d, "LICENSE")
	f, err := os.Create(fPath)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	licenseText := "The MIT License (MIT)"

	if _, err := f.WriteString(licenseText); err != nil {
		t.Fatalf("err: %s", err)
	}

	l, err := NewFromDir(d)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if l.Type != "MIT" {
		t.Fatalf("unexpected license type: %s", l.Type)
	}

	if l.Text != licenseText {
		t.Fatalf("unexpected license text: %s", l.Text)
	}

	if l.File != fPath {
		t.Fatalf("unexpected file path: %s", l.File)
	}

	// Fails properly if the directory does not exist
	if _, err := NewFromDir("go-license-nonexistent"); err == nil {
		t.Fatalf("expected error loading non-existent directory")
	}

	// Fails properly if the directory specified is actually a file
	if _, err := NewFromDir(fPath); err == nil {
		t.Fatalf("expected error loading file as directory")
	}
}

func TestLicenseRecognized(t *testing.T) {
	// Known licenses are recognized
	l := New("MIT", "The MIT License (MIT)")
	if !l.Recognized() {
		t.Fatalf("license was not recognized")
	}

	// Unknown licenses are not recognized
	l = New("None", "No license text")
	if l.Recognized() {
		t.Fatalf("fake license was recognized")
	}
}

func TestLicenseTypes(t *testing.T) {
	licenseStrings := []string{
		"The MIT License (MIT)",
		"Apache License\nVersion 2.0, January 2004",
		"GNU GENERAL PUBLIC LICENSE\nVersion 2, June 1991",
		"GNU GENERAL PUBLIC LICENSE\nVersion 3, 29 June 2007",
		"GNU LESSER GENERAL PUBLIC LICENSE\nVersion 2.1, February 1999",
		"GNU LESSER GENERAL PUBLIC LICENSE\nVersion 3, 29 June 2007",
		"Mozilla Public License Version 2.0",
		"Redistribution and use in source and binary forms\n4. Neither",
		"Redistribution and use in source and binary forms\n* Redistributions",
		"Redistribution and use in source and binary forms\nFreeBSD Project.",
		"(CDDL)\nVersion 1.0",
		"Eclipse Public License - v 1.0",
	}

	for _, s := range licenseStrings {
		l := New("", s)
		if err := l.GuessType(); err != nil {
			t.Fatalf("failed to identify license: %s", s)
		}
	}
}
