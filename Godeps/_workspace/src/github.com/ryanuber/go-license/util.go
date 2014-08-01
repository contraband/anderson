package license

import "strings"

// scanLeft will scan through a block of text, delmited by newline characters,
// and check for lines matching the provided text on the left-hand side.
func scanLeft(text, match string) bool {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), match) {
			return true
		}
	}
	return false
}

// scanRight will scan through a block of text, delmited by newline characters,
// and check for lines matching the provided text on the right-hand side.
func scanRight(text, match string) bool {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasSuffix(strings.TrimSpace(line), match) {
			return true
		}
	}
	return false
}
