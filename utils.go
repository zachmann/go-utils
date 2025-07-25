package zutils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

// NewInt returns a *int
func NewInt(i int) *int {
	return &i
}

// FirstNonEmpty is a utility function returning the first of the passed values that is not empty/zero
func FirstNonEmpty[C comparable](possibleValues ...C) C {
	var nullValue C
	for _, v := range possibleValues {
		if v != nullValue {
			return v
		}
	}
	return nullValue
}

// FirstNonEmptyFnc is a utility function returning the first of the passed
// values that is not empty/zero.
// In this function the values are not passed directly but a function that
// returns the value is passed instead. This enables lazy evaluation
func FirstNonEmptyFnc[C comparable](possibleValues ...func() C) C {
	var nullValue C
	for _, fnc := range possibleValues {
		if v := fnc(); v != nullValue {
			return v
		}
	}
	return nullValue
}

func RandomString(n int) (string, error) {
	byteLen := n * 3 / 4 // base64 expands by 4/3
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n], nil
}
