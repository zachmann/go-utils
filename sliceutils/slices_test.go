package sliceutils

import (
	"testing"
)

func sliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func containsSubset[T comparable](subs [][]T, target []T) bool {
	for _, s := range subs {
		if sliceEqual(s, target) {
			return true
		}
	}
	return false
}

func TestSubsets_NilInput(t *testing.T) {
	var in []int
	out := Subsets(in)
	if out != nil {
		t.Fatalf("expected nil for nil input, got %v", out)
	}
}

func TestSubsets_EmptyInput(t *testing.T) {
	in := []int{}
	out := Subsets(in)
	if out == nil {
		t.Fatalf("expected empty slice (not nil) for empty input")
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 subsets, got %d: %v", len(out), out)
	}
}

func TestSubsets_Singleton(t *testing.T) {
	in := []int{1}
	out := Subsets(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 subset, got %d: %v", len(out), out)
	}
	if !containsSubset(out, []int{1}) {
		t.Fatalf("expected subset [1] present, got %v", out)
	}
}

func TestSubsets_TwoElements(t *testing.T) {
	in := []int{1, 2}
	out := Subsets(in)
	// 2^2 - 1 = 3 subsets
	if len(out) != 3 {
		t.Fatalf("expected 3 subsets, got %d: %v", len(out), out)
	}
	expected := [][]int{{1}, {2}, {1, 2}}
	for _, e := range expected {
		if !containsSubset(out, e) {
			t.Fatalf("missing expected subset %v in %v", e, out)
		}
	}
	// Ensure empty subset not present
	if containsSubset(out, []int{}) {
		t.Fatalf("did not expect empty subset in %v", out)
	}
}

func TestSubsets_ThreeElements_OrderWithinSubsets(t *testing.T) {
	in := []int{1, 2, 3}
	out := Subsets(in)
	// 2^3 - 1 = 7 subsets
	if len(out) != 7 {
		t.Fatalf("expected 7 subsets, got %d: %v", len(out), out)
	}
	expected := [][]int{{1}, {2}, {3}, {1, 2}, {1, 3}, {2, 3}, {1, 2, 3}}
	for _, e := range expected {
		if !containsSubset(out, e) {
			t.Fatalf("missing expected subset %v in %v", e, out)
		}
	}
	// Check preservation of order within subsets: [2,1] should not be present
	if containsSubset(out, []int{2, 1}) || containsSubset(out, []int{3, 1}) || containsSubset(out, []int{3, 2}) {
		t.Fatalf("unexpected reversed-order subset found in %v", out)
	}
}

func TestSubsets_Strings(t *testing.T) {
	in := []string{"a", "b"}
	out := Subsets(in)
	if len(out) != 3 {
		t.Fatalf("expected 3 subsets, got %d: %v", len(out), out)
	}
	expected := [][]string{{"a"}, {"b"}, {"a", "b"}}
	for _, e := range expected {
		if !containsSubset(out, e) {
			t.Fatalf("missing expected subset %v in %v", e, out)
		}
	}
}
