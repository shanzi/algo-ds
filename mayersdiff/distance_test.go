package mayersdiff

import "testing"

func TestIdenticalStrings(t *testing.T) {
	if EditingDistanceStrings("", "") != 0 {
		t.Error("Editing distance of two empty strings should be 0")
	}

	if EditingDistanceStrings("abc", "abc") != 0 {
		t.Error("Editing distance of 'abc' and 'abc' should be 0")
	}

	if EditingDistanceStrings(
		"The quick brown fox jumps over the lazy dog",
		"The quick brown fox jumps over the lazy dog",
	) != 0 {
		t.Errorf(
			"Editing distance of '%s' and '%s' should be 0",
			"The quick brown fox jumps over the lazy dog",
			"The quick brown fox jumps over the lazy dog",
		)
	}
}

func TestTotallyDifferentStrings(t *testing.T) {
	if EditingDistanceStrings("abc", "def") != 6 {
		t.Error(
			"Editing distance of '%s' and '%s' strings should be %d",
			"abc", "def", 6,
		)
	}

	if EditingDistanceStrings("abcdef", "xyz") != 9 {
		t.Error(
			"Editing distance of '%s' and '%s' strings should be %d",
			"abcdef", "xyz", 9,
		)
	}

	if EditingDistanceStrings("abc", "uvwxyz") != 9 {
		t.Error(
			"Editing distance of '%s' and '%s' strings should be %d",
			"abc", "uvwxyz", 9,
		)
	}
}

func TestEditingDistance(t *testing.T) {
	if EditingDistanceStrings("abc", "cde") != 4 {
		t.Error(
			"Editing distance of '%s' and '%s' strings should be %d",
			"abc", "cde", 4,
		)
	}

	if d := EditingDistanceStrings("abcabba", "cbabac"); d != 5 {
		t.Errorf(
			"Editing distance of '%s' and '%s' strings should be %d, but got %d",
			"abcabba", "cbabac", 5, d,
		)
	}

	if d := EditingDistanceStrings("bjaposdif", "bawjpoid"); d != 7 {
		t.Errorf(
			"Editing distance of '%s' and '%s' strings should be %d, but got %d",
			"bjaposdif", "bawjpoid", 7, d,
		)
	}

	if d := EditingDistanceStrings("fajsodivjpwevasdbug", "ajpivaspvheuvbaosidufh"); d != 19 {
		t.Errorf(
			"Editing distance of '%s' and '%s' strings should be %d, but got %d",
			"fajsodivjpwevasdbug", "ajpivaspvheuvbaosidufh", 19, d,
		)
	}
}
