package sliceutils

import (
	"slices"
)

// RemoveFromSlice removes an element from the slice
func RemoveFromSlice[C comparable](slice []C, v C) []C {
	result := make([]C, 0, len(slice))
	for _, vv := range slice {
		if vv != v {
			result = append(result, vv)
		}
	}
	return result
}

// SliceContains checks if a slice contains a value
func SliceContains[C comparable](v C, slice []C) bool {
	return slices.Contains(slice, v)
}

// Reverse returns the slice in reverse order.
func Reverse[V any](ivs []V) []V {
	if ivs == nil {
		return nil
	}
	l := len(ivs)
	ovs := make([]V, l)
	for i := range ivs {
		l--
		ovs[i] = ivs[l]
	}
	return ovs
}

// IsSubsetOf checks if a slice `subset` is a subset of the slice `of`,
// i.e. it is verified that `of` contains all the elements of `subset`
func IsSubsetOf[T comparable](subset []T, of []T) bool {
	for _, sub := range subset {
		if !slices.Contains(of, sub) {
			return false
		}
	}
	return true
}

// Subsets returns all non-empty subsets (power set without empty set)
// of the given slice. It preserves element order within subsets.
// If the input slice is nil, it returns nil.
func Subsets[T any](vs []T) [][]T {
	if vs == nil {
		return nil
	}

	// Build subsets without including the empty subset
	result := make([][]T, 0)
	for _, v := range vs {
		n := len(result)
		// Add the single-element subset
		result = append(result, []T{v})
		// Extend existing subsets with the new element
		for i := 0; i < n; i++ {
			subset := result[i]
			ns := make([]T, len(subset)+1)
			copy(ns, subset)
			ns[len(subset)] = v
			result = append(result, ns)
		}
	}

	return result
}
