package terbilang

import (
	"math/big"
	"strings"
)

// scales[i] is the name of 1000^i. The units group has no name.
// Indonesian uses the short scale: miliar is 10^9 (KBBI), not the older milyar spelling.
var scales = [...]string{
	"",
	"ribu",
	"juta",
	"miliar",
	"triliun",
	"kuadriliun",
	"kuintiliun",
	"sekstiliun",
	"septiliun",
	"oktiliun",
	"noniliun",
	"desiliun",
}

// small holds nol and the words for 1..11.
var small = [...]string{
	"nol",
	"satu",
	"dua",
	"tiga",
	"empat",
	"lima",
	"enam",
	"tujuh",
	"delapan",
	"sembilan",
	"sepuluh",
	"sebelas",
}

// Convert returns the Indonesian words for n.
//
// Zero is "nol". Negative numbers are prefixed with "minus".
// The conventional forms sebelas, seratus, and seribu are used for 11, 100,
// and 1000. Larger powers use satu, as in "satu juta", which is the form
// expected on invoices.
func Convert(n int64) string {
	v := new(big.Int).SetInt64(n)
	neg := v.Sign() < 0
	v.Abs(v)
	words, err := spellAbs(v)
	if err != nil {
		// int64 is at most kuintiliun, which is inside the supported scales.
		panic(err)
	}
	if neg {
		return "minus " + words
	}
	return words
}

// spellAbs spells a non-negative integer. n is not modified.
func spellAbs(n *big.Int) (string, error) {
	if n.Sign() < 0 {
		panic("terbilang: spellAbs called with a negative value")
	}
	if n.Sign() == 0 {
		return "nol", nil
	}

	base := big.NewInt(1000)
	remaining := new(big.Int).Set(n)
	quot := new(big.Int)
	mod := new(big.Int)
	parts := make([]string, 0, 8)

	for scale := 0; remaining.Sign() > 0; scale++ {
		if scale >= len(scales) {
			return "", ErrMagnitude
		}
		quot.QuoRem(remaining, base, mod)
		if mod.Sign() > 0 {
			parts = append(parts, spellGroup(int(mod.Int64()), scale))
		}
		remaining.Set(quot)
	}

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, " "), nil
}

// spellGroup spells one 1..999 group and, except for the units group, its scale.
// The ribu group is the only scale that uses the "se" form: 1000 is "seribu",
// while 1_000_000 is "satu juta".
func spellGroup(n, scale int) string {
	if scale == 1 && n == 1 {
		return "seribu"
	}
	words := spellBelow1000(n)
	if scale == 0 {
		return words
	}
	return words + " " + scales[scale]
}

func spellBelow1000(n int) string {
	if n < 12 {
		return small[n]
	}
	if n < 20 {
		return small[n-10] + " belas"
	}
	if n < 100 {
		words := small[n/10] + " puluh"
		if n%10 == 0 {
			return words
		}
		return words + " " + small[n%10]
	}

	var words string
	if n/100 == 1 {
		words = "seratus"
	} else {
		words = small[n/100] + " ratus"
	}
	if n%100 == 0 {
		return words
	}
	return words + " " + spellBelow1000(n%100)
}

func spellDigits(digits string) string {
	parts := make([]string, len(digits))
	for i := 0; i < len(digits); i++ {
		parts[i] = small[digits[i]-'0']
	}
	return strings.Join(parts, " ")
}
