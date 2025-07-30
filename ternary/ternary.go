package ternary

// If returns a if con is true, b otherwise
func If(con bool, a, b interface{}) interface{} {
	if con {
		return a
	}
	return b
}
