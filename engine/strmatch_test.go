package engine

import "testing"

func test_strMatch(t *testing.T) {
	tests := []struct {
		pattern string
		str     string
		nocase  bool
		match   bool
	}{
		// Basic tests
		{"*", "hello", false, true},
		{"h*", "hello", false, true},
		{"*o", "hello", false, true},
		{"he?lo", "hello", false, true},
		{"he?lo", "healo", false, true},
		{"he?lo", "heaoo", false, false},

		// Character class
		{"h[aeiou]llo", "hello", false, true},
		{"h[aeiou]llo", "hillo", false, true},
		{"h[aeiou]llo", "hollo", false, true},
		{"h[!aeiou]llo", "hallo", false, false},

		// Range
		{"[a-z]ello", "hello", false, true},
		{"[a-z]ello", "aello", false, true},
		{"[a-z]ello", "Aello", false, false},
		{"[A-Z]ello", "Aello", false, true},

		// Case insensitive
		{"h*", "HELLO", true, true},
		{"H*o", "hello", true, true},
		{"he?lo", "HELLO", true, true},
		{"he?lo", "heALO", true, true},
		{"h[aeiou]llo", "HELLO", true, true},

		// Edge cases
		{"*", "", false, true},
		{"?", "", false, false},
		{"*", "anystring", false, true},
		{"[a-c]*", "bcdef", false, true},
		{"*xyz", "abcxyz", false, true},
		{"*", "abc", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			if got := strMatch(tt.pattern, tt.str, tt.nocase); got != tt.match {
				t.Errorf("strMatch(%q, %q, %v) = %v, want %v", tt.pattern, tt.str, tt.nocase, got, tt.match)
			}
		})
	}
}
