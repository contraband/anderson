package anderson

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
