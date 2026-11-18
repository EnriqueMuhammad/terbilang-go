package terbilang

import (
	"fmt"
	"math/big"
	"strings"
)

// maxFracDigits limits how long a fractional part may be.
// Invoice amounts only need a couple of digits; the cap keeps a hostile
// input from turning into an enormous word list.
const maxFracDigits = 64

// parsed is a decimal number exactly as written.
// frac keeps leading zeros, which matter for "koma nol lima".
type parsed struct {
	neg   bool
	whole *big.Int
	frac  string
}

// Parse spells the number in input.
//
// The integer part may use Indonesian grouping (dots) or US grouping (commas).
// A comma is the decimal separator, except when commas are used only as
// thousands separators (1,250,000) or when a dot is the rightmost separator
// (1,250.50). A single dot followed by exactly three digits is a thousands
// separator (10.000), matching Indonesian notation. Any other single dot is
// a decimal point (1.5, 3.14159).
//
// Fractional digits are read one by one after "koma", so zeros in the
// fraction are preserved: 1,05 is "satu koma nol lima".
//
// An optional leading sign and an "Rp" prefix are accepted.
// The integer part must fit in the supported scales, through desiliun (10^33).
func Parse(input string) (string, error) {
	p, err := parseNumber(input)
	if err != nil {
		return "", err
	}
	words, err := p.words()
	if err != nil {
		return "", fmt.Errorf("%w: %q", err, strings.TrimSpace(input))
	}
	return words, nil
}

func (p parsed) words() (string, error) {
	whole, err := spellAbs(p.whole)
	if err != nil {
		return "", err
	}
	// "-0" is just zero. A written fraction keeps the sign: "-0,5".
	if p.neg && (p.whole.Sign() != 0 || p.frac != "") {
		whole = "minus " + whole
	}
	if p.frac == "" {
		return whole, nil
	}
	return whole + " koma " + spellDigits(p.frac), nil
}

func parseNumber(input string) (parsed, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return parsed{}, ErrEmpty
	}

	neg := false
	if s[0] == '+' || s[0] == '-' {
		neg = s[0] == '-'
		s = strings.TrimSpace(s[1:])
		if s == "" {
			return parsed{}, syntax(input)
		}
	}
	if hasRpPrefix(s) {
		s = strings.TrimSpace(s[2:])
		if s == "" {
			return parsed{}, syntax(input)
		}
	}

	whole, frac, err := splitNumeric(s)
	if err != nil {
		return parsed{}, syntax(input)
	}
	if len(frac) > maxFracDigits {
		return parsed{}, syntax(input)
	}

	whole = strings.TrimLeft(whole, "0")
	if whole == "" {
		whole = "0"
	}
	w, ok := new(big.Int).SetString(whole, 10)
	if !ok {
		return parsed{}, syntax(input)
	}
	return parsed{neg: neg, whole: w, frac: frac}, nil
}

func syntax(input string) error {
	return fmt.Errorf("%w: %q", ErrSyntax, strings.TrimSpace(input))
}

func hasRpPrefix(s string) bool {
	return len(s) >= 2 && (s[0] == 'R' || s[0] == 'r') && (s[1] == 'P' || s[1] == 'p')
}

// splitNumeric applies the separator rules documented on Parse.
func splitNumeric(body string) (string, string, error) {
	if !validNumeric(body) {
		return "", "", errBadNumber
	}
	dotN := strings.Count(body, ".")
	commaN := strings.Count(body, ",")
	switch {
	case dotN == 0 && commaN == 0:
		return body, "", nil
	case dotN > 0 && commaN > 0:
		return splitMixed(body)
	case commaN > 0:
		return splitSingle(body, ',', commaN)
	default:
		return splitSingle(body, '.', dotN)
	}
}

func validNumeric(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || c == '.' || c == ',' {
			continue
		}
		return false
	}
	return true
}

func splitMixed(body string) (string, string, error) {
	lastDot := strings.LastIndexByte(body, '.')
	lastComma := strings.LastIndexByte(body, ',')
	dec := byte('.')
	group := byte(',')
	decAt := lastDot
	if lastComma > lastDot {
		dec = ','
		group = '.'
		decAt = lastComma
	}

	intPart := body[:decAt]
	frac := body[decAt+1:]
	if !digitsOnly(frac) {
		return "", "", errBadNumber
	}
	if strings.IndexByte(intPart, dec) >= 0 {
		return "", "", errBadNumber
	}
	if intPart == "" {
		return "0", frac, nil
	}
	if strings.IndexByte(intPart, group) < 0 {
		if !digitsOnly(intPart) {
			return "", "", errBadNumber
		}
		return intPart, frac, nil
	}
	whole, err := stripGroups(intPart, group)
	if err != nil {
		return "", "", err
	}
	return whole, frac, nil
}

func splitSingle(body string, sep byte, count int) (string, string, error) {
	if count > 1 {
		whole, err := stripGroups(body, sep)
		if err != nil {
			return "", "", err
		}
		return whole, "", nil
	}

	i := strings.IndexByte(body, sep)
	left, right := body[:i], body[i+1:]
	if !digitsOnly(right) {
		return "", "", errBadNumber
	}
	if left == "" {
		left = "0"
	}
	if !digitsOnly(left) {
		return "", "", errBadNumber
	}
	// 10.000 is sepuluh ribu. 1.5 and 3.14159 stay decimal.
	// A zero integer part is always decimal, so 0.250 is not 250.
	if sep == '.' && len(right) == 3 && !allZeroDigits(left) {
		return left + right, "", nil
	}
	return left, right, nil
}

func stripGroups(s string, sep byte) (string, error) {
	parts := strings.Split(s, string(sep))
	if len(parts) < 2 {
		return "", errBadNumber
	}
	for i, p := range parts {
		if !digitsOnly(p) {
			return "", errBadNumber
		}
		if i == 0 {
			if len(p) > 3 {
				return "", errBadNumber
			}
			continue
		}
		if len(p) != 3 {
			return "", errBadNumber
		}
	}
	return strings.Join(parts, ""), nil
}

func digitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func allZeroDigits(s string) bool {
	if s == "" {
		return true
	}
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			return false
		}
	}
	return true
}
