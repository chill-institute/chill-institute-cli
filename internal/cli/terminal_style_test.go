package cli

import (
	"regexp"
	"testing"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func TestBoldWrapsWithANSI(t *testing.T) {
	t.Parallel()
	got := bold("hello")
	if stripANSI(got) != "hello" {
		t.Fatalf("bold() stripped = %q, want hello", stripANSI(got))
	}
	if got == "hello" {
		t.Fatal("bold() should contain ANSI codes")
	}
}

func TestDimWrapsWithANSI(t *testing.T) {
	t.Parallel()
	got := dim("label")
	if stripANSI(got) != "label" {
		t.Fatalf("dim() stripped = %q, want label", stripANSI(got))
	}
	if got == "label" {
		t.Fatal("dim() should contain ANSI codes")
	}
}
