package license

import "testing"

func TestScan(t *testing.T) {
	text := "she sells sea shells by the sea shore"

	// scanLeft matches properly
	scanLeftPass := []string{
		"she sells sea shells",
		"she sells sea shells by the sea shore",
	}
	for _, s := range scanLeftPass {
		if !scanLeft(text, s) {
			t.Fatalf("%s did not match during scanLeft", s)
		}
	}

	// scanLeft rejects messages that shouldn't match
	scanLeftFail := []string{
		"by the sea shore",
		" she sells sea shells by the sea shore",
	}
	for _, s := range scanLeftFail {
		if scanLeft(text, s) {
			t.Fatalf("%s matched during scanLeft", s)
		}
	}

	// scanRight matches properly
	scanRightPass := []string{
		"by the sea shore",
		"she sells sea shells by the sea shore",
	}
	for _, s := range scanRightPass {
		if !scanRight(text, s) {
			t.Fatalf("%s did not match during scanRight", s)
		}
	}

	// scanRight rejects messages that shouldn't match
	scanRightFail := []string{
		"she sells sea shells",
		"she sells sea shells by the sea shore ",
	}
	for _, s := range scanRightFail {
		if scanRight(text, s) {
			t.Fatalf("%s matched during scanRight", s)
		}
	}
}
