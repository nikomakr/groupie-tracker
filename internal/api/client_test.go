package api

import (
	"testing"
)

func TestContains(t *testing.T) {
	cases := []struct {
		s, substr string
		want      bool
	}{
		{"Queen", "queen", true},
		{"Queen", "Queen", true},
		{"Queen", "xyz", false},
		{"Freddie Mercury", "freddie", true},
		{"", "queen", false},
	}

	for _, c := range cases {
		got := Contains(c.s, c.substr)
		if got != c.want {
			t.Errorf("contains(%q, %q) = %v, want %v", c.s, c.substr, got, c.want)
		}
	}
}

func TestContainsEmpty(t *testing.T) {
	if !Contains("Queen", "") {
		t.Error("expected true when substr is empty")
	}
}