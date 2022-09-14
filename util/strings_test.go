package util

import (
	"testing"
)

func TestContains(t *testing.T) {
	s := []string{"a", "b", "c"}

	sContainsA := Contains(s, "a")
	sContainsD := Contains(s, "d")

	if !sContainsA {
		t.Errorf("s slice should contain a")
	}

	if sContainsD {
		t.Errorf("s slice should not contain d")
	}
}
