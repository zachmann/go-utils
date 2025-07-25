package structutils

import (
	"strings"

	"github.com/fatih/structs"
)

// FieldTagNames returns a slice of the tag names for a []*structs.Field and the given tag
func FieldTagNames(s any, tag string) (names []string) {
	fields := structs.New(s).Fields()
	for _, f := range fields {
		if f == nil {
			continue
		}
		t := f.Tag(tag)
		if i := strings.IndexRune(t, ','); i > 0 {
			t = t[:i]
		}
		if t != "" && t != "-" {
			names = append(names, t)
		}
	}
	return
}
