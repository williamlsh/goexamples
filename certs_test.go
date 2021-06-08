package x

import (
	"testing"
)

func TestGenerateCerts(t *testing.T) {
	if err := GenerateCerts(); err != nil {
		t.Fatal(err)
	}
	if err := removeCerts(); err != nil {
		t.Fatal(err)
	}
}
