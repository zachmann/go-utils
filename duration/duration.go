package duration

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var unitMap = map[string]int64{
	"ns": int64(time.Nanosecond),
	"us": int64(time.Microsecond),
	"µs": int64(time.Microsecond), // U+00B5 = micro symbol
	"μs": int64(time.Microsecond), // U+03BC = Greek letter mu
	"ms": int64(time.Millisecond),
	"s":  int64(time.Second),
	"m":  int64(time.Minute),
	"h":  int64(time.Hour),
	"d":  int64(24 * time.Hour),       // Approximation
	"w":  int64(7 * 24 * time.Hour),   // Approximation
	"y":  int64(365 * 24 * time.Hour), // Approximation
}

const invalidDurErr = "time: invalid duration "

var errLeadingInt = errors.New("duration: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d", "w", "y".
func ParseDuration(s string) (time.Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d int64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, errors.New(invalidDurErr + orig)
	}
	for s != "" {
		var (
			v, f  int64       // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, errors.New(invalidDurErr + orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New(invalidDurErr + orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, s, err = leadingInt(s)
			if err != nil {
				return 0, errors.New(invalidDurErr + orig)
			}
			for n := pl - len(s); n > 0; n-- {
				scale *= 10
			}
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, errors.New(invalidDurErr + orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, errors.New("time: missing unit in duration " + orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, errors.New("time: unknown unit " + u + " in duration " + orig)
		}
		if v > (1<<63-1)/unit {
			// overflow
			return 0, errors.New(invalidDurErr + orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += int64(float64(f) * (float64(unit) / scale))
			if v < 0 {
				// overflow
				return 0, errors.New(invalidDurErr + orig)
			}
		}
		d += v
		if d < 0 {
			// overflow
			return 0, errors.New(invalidDurErr + orig)
		}
	}

	if neg {
		d = -d
	}
	return time.Duration(d), nil
}

// DurationOption is a type for parsing a duration string passed through a
// config file or by other means into a time.Duration.
// It implements the json.Unmarshaler and yaml.Unmarshaler interfaces.
// It also implements marshaler interfaces,
// but that output is not well readable by humans.
type DurationOption time.Duration

// UnmarshalJSON implements the json.Unmarshaler interface
func (do *DurationOption) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		var i float64
		if e := json.Unmarshal(bytes, &i); e != nil {
			return err
		}
		*do = DurationOption(time.Duration(i * float64(time.Second)))
		return nil
	}
	d, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*do = DurationOption(d)
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (do DurationOption) MarshalJSON() (b []byte, err error) {
	ns := do.Duration().Nanoseconds()
	str := fmt.Sprintf("%dns", ns)
	return json.Marshal(str)
}

// Duration returns the time.Duration
func (do DurationOption) Duration() time.Duration {
	return time.Duration(do)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (do *DurationOption) UnmarshalYAML(n *yaml.Node) error {
	switch n.Tag {
	case "!!int", "!!float":
		// Number → seconds
		var f float64
		if err := n.Decode(&f); err != nil {
			return err
		}
		*do = DurationOption(time.Duration(f * float64(time.Second)))
		return nil

	case "!!str", "":
		// String → parse with ParseDuration
		var s string
		if err := n.Decode(&s); err != nil {
			return err
		}
		d, err := ParseDuration(s)
		if err != nil {
			return err
		}
		*do = DurationOption(d)
		return nil
	}
	return fmt.Errorf("DurationOption (yaml) must be a number (seconds) or a string, got tag=%s", n.Tag)
}

// MarshalYAML implements the yaml.Marshaler interface
func (do DurationOption) MarshalYAML() (interface{}, error) {
	ns := time.Duration(do).Nanoseconds()
	return fmt.Sprintf("%dns", ns), nil
}
