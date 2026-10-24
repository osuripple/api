package common

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	for _, n := range []int{0, 1, 7, 32, 50} {
		s, err := RandomString(n)
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != n || strings.Trim(s, "0123456789abcdef") != "" {
			t.Errorf("RandomString(%d) = %q", n, s)
		}
	}
}

func TestRandomLetters(t *testing.T) {
	for _, n := range []int{0, 1, 7, 32, 50} {
		s, err := RandomLetters(n)
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != n || strings.Trim(s, LetterChars) != "" {
			t.Errorf("RandomLetters(%d) = %q", n, s)
		}
	}
}

func TestRandomStringFromChars(t *testing.T) {
	if _, err := RandomStringFromChars(5, ""); err == nil {
		t.Error("expected an error with empty chars")
	}

	const chars = "ABC123"
	const n = 60000
	s, err := RandomStringFromChars(n*len(chars), chars)
	if err != nil {
		t.Fatal(err)
	}
	counts := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		counts[s[i]]++
	}
	if len(counts) != len(chars) {
		t.Fatalf("expected %d distinct chars, got %v", len(chars), counts)
	}
	// each char is expected n times; allow a generous margin to keep the test stable
	for c, got := range counts {
		if got < n*95/100 || got > n*105/100 {
			t.Errorf("char %q appeared %d times, expected about %d", c, got, n)
		}
	}
}
