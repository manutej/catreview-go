package wiringmap

import (
	"testing"
)

func TestSanitizeID(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"pkg/analysis/foo.go", "pkg.analysis.foo.go"},
		{"encoding/json", "encoding.json"},
		{"github.com/foo/bar", "github.com.foo.bar"},
		{"a-b_c.d", "a-b_c.d"},
		{"weird@name", "weird_name"},
	}
	for _, tc := range tests {
		if got := SanitizeID(tc.in); got != tc.want {
			t.Errorf("SanitizeID(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestCheckSanitizeCollisions(t *testing.T) {
	err := CheckSanitizeCollisions(map[string]string{
		"a/b": SanitizeID("a/b"),
		"c":   SanitizeID("c"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = CheckSanitizeCollisions(map[string]string{
		"a/b": SanitizeID("a/b"),
		"a.b": SanitizeID("a.b"),
	})
	if err == nil {
		t.Fatal("expected collision error")
	}
}
