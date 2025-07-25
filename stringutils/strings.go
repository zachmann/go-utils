package utils

// EqualIfSet checks if two strings are equal if they are both not empty
func EqualIfSet(a, b string) bool {
	return a == "" || b == "" || a == b
}
