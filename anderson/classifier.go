package anderson

import (
	"fmt"

	"github.com/ryanuber/go-license"
)

type LicenseClassifier struct {
	Config Config
}

func (c LicenseClassifier) Classify(path string, importPath string) (LicenseStatus, string, error) {
	l, err := license.NewFromDir(path)

	if err != nil {
		switch err.Error() {
		case license.ErrNoLicenseFile:
			if contains(c.Config.Exceptions, importPath) {
				return LicenseTypeAllowed, "Unknown", nil
			}

			return LicenseTypeNoLicense, "Unknown", nil
		case license.ErrUnrecognizedLicense:
			if contains(c.Config.Exceptions, importPath) {
				return LicenseTypeAllowed, "Unknown", nil
			}

			return LicenseTypeUnknown, "Unknown", nil
		default:
			return LicenseTypeUnknown, "Error", fmt.Errorf("Could not determine license for: %s", importPath)
		}
	}

	if contains(c.Config.Blacklist, l.Type) {
		return LicenseTypeBanned, l.Type, nil
	}

	if contains(c.Config.Whitelist, l.Type) {
		return LicenseTypeAllowed, l.Type, nil
	}

	if contains(c.Config.Exceptions, importPath) {
		return LicenseTypeAllowed, l.Type, nil
	} else {
		return LicenseTypeMarginal, l.Type, nil
	}
}

func contains(haystack []string, needle string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}
