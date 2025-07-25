package maputils

import (
	"github.com/zachmann/go-utils/sliceutils"
)

// Keys returns the keys of the map m.
// The order of the keys is not deterministic
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// MergeMaps merges two or more maps into one; overwrite determines if values
// are overwritten if already set or not
func MergeMaps(overwrite bool, mm ...map[string]any) map[string]any {
	if !overwrite {
		return MergeMaps(true, sliceutils.Reverse(mm)...)
	}
	all := make(map[string]any)
	for _, m := range mm {
		for k, v := range m {
			all[k] = v
		}
	}
	return all
}
