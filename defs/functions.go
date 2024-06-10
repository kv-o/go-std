package defs

import (
	"fmt"

	"git.sr.ht/~kvo/go-std/errors"
)

// Get returns the nth element of slice s. Returns error if slice s does not
// have an element at index n.
//
// For most access attempts on strings, []rune(str) will be a more appropriate
// choice than []byte(str) for the parameter s, as no individual byte in
// []byte(str) is guaranteed to hold a single Unicode code point.
func Get[T any](s []T, n int) (T, error) {
	var none T
	if n > len(s)-1 || n < 0 {
		return none, errors.New(
			fmt.Sprintf("index out of range [%d] with length %d", n, len(s)),
			nil,
		)
	}
	return s[n], nil
}

// Contains checks slice s for the existence of an element elem.
func Has[T comparable](s []T, elem T) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
}

// HasOnly checks if all elements of s are elem.
func HasOnly[T comparable](s []T, elem T) bool {
	for _, v := range s {
		if v != elem {
			return false
		}
	}
	return true
}

// Remove attempts to remove element elem from slice s and return the resulting
// slice. If elem is not present in s, s is returned unchanged.
func Remove[T comparable](s []T, elem T) []T {
	for i, v := range s {
		if v == elem {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
