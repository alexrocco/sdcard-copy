package slice

import (
	"sort"
)

// diff gets the difference between a and b.
func Diff(a, b []string) []string {
	// Sort both slices
	sort.Strings(a)
	sort.Strings(b)

	var diff []string
	for _, x := range a {
		// sort.SearchStrings returns the index where the element should be added or where it exists in the slice
		i := sort.SearchStrings(b, x)

		// if the index is the length of the slice or the value of the index is different from the element, the element
		// does not exist in the slice
		if len(b) == i || b[i] != x {
			diff = append(diff, x)
		}
	}

	return diff
}