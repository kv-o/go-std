// Package std collects implements an alternative standard library.
//
// While Go's standard library and built-in identifiers typically do well in
// catering to the needs of the programmer, some functions that could assist
// the programmer in general programming cases do not exist in the standard Go
// distribution. Package std attempts to resolve this issue by providing a set
// of useful identifiers for general programming.
package std

import "fmt"

// Access attempts to access the nth element of slice s.
func Access[T comparable](s []T, n int) error {
	if n > len(s)-1 || n < 0 {
		return fmt.Errorf("index out of range [%d] with length %d", n, len(s))
	}
	return nil
}

// Contains checks slice s for the existence of an element elem.
func Contains[T comparable](s []T, elem T) bool {
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
