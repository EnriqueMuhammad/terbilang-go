package terbilang

import "errors"

var (
	// ErrEmpty is returned when the input is empty or only whitespace.
	ErrEmpty = errors.New("terbilang: empty number")

	// ErrSyntax is returned when the input is not a valid number.
	ErrSyntax = errors.New("terbilang: invalid number")

	// ErrMagnitude is returned when the integer part exceeds the largest
	// supported scale, desiliun (10^33, 36 digits).
	ErrMagnitude = errors.New("terbilang: magnitude exceeds supported range")

	// ErrSen is returned when sen is outside 0..99.
	ErrSen = errors.New("terbilang: sen out of range")
)

// errBadNumber is the internal sentinel for a numeric body that fails the grammar.
// Callers replace it with ErrSyntax and the original input.
var errBadNumber = errors.New("terbilang: bad number")
