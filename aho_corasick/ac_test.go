package aho_corasick

import (
	"testing"
)

var dict = []string{
	"abc",
	"bcd",
	"cde",
	"abcdefg",
	"123",
	"456",
	"789",
	"56789",
	"not found",
}

var str = "abcdefg123456789"

func TestSearch(t *testing.T) {
	ac := New(dict)
	expect := map[string]bool{}
	found := map[string]bool{}

	for _, w := range dict {
		expect[w] = true
	}

	expect["not found"] = false

	for m := range ac.Search(str) {
		w := str[m.Start():m.End()]
		if expect[w] != true {
			t.Errorf("%s should not have been found\n", w)
		}
		found[w] = true
	}

	for w, v := range expect {
		if v && found[w] != true {
			t.Errorf("%s should have been found\n", w)
		}
	}
}

func TestSearchCallback(t *testing.T) {
	ac := New(dict)
	expect := map[string]bool{}
	found := map[string]bool{}

	for _, w := range dict {
		expect[w] = true
	}

	expect["not found"] = false

	ac.SearchCallback(str, func(start, end int) {
		w := str[start:end]
		if expect[w] != true {
			t.Errorf("%s should not have been found\n", w)
		}
		found[w] = true
	})

	for w, v := range expect {
		if v && found[w] != true {
			t.Errorf("%s should have been found\n", w)
		}
	}
}
