package terbilang

import (
	"fmt"
	"math/big"
	"strings"
)

// Rupiah spells amount and appends "rupiah".
// 1250000 becomes "satu juta dua ratus lima puluh ribu rupiah".
func Rupiah(amount int64) string {
	words, err := spellRupiahAmount(amount, 0)
	if err != nil {
		// int64 is inside the supported scales, and sen 0 needs no range check.
		panic(err)
	}
	return words
}

// RupiahSen spells a rupiah amount with an optional sen part (0–99).
// A zero sen part is omitted. When rupiah is zero and sen is not, only the
// sen is spelled: RupiahSen(0, 50) is "lima puluh sen".
func RupiahSen(rupiah int64, sen int) (string, error) {
	if sen < 0 || sen > 99 {
		return "", ErrSen
	}
	return spellRupiahAmount(rupiah, sen)
}

// FormatRupiah formats amount as Indonesian Rupiah with two decimal places.
// 1250000 becomes "Rp1.250.000,00". Negative amounts get a leading minus:
// "-Rp1.250.000,00".
func FormatRupiah(amount int64) string {
	v := new(big.Int).SetInt64(amount)
	neg := v.Sign() < 0
	v.Abs(v)
	return formatRupiah(neg, v, 0)
}

// FormatRupiahSen formats rupiah and sen (0–99) as Indonesian currency.
// FormatRupiahSen(1250000, 50) is "Rp1.250.000,50".
func FormatRupiahSen(rupiah int64, sen int) (string, error) {
	if sen < 0 || sen > 99 {
		return "", ErrSen
	}
	v := new(big.Int).SetInt64(rupiah)
	neg := v.Sign() < 0
	v.Abs(v)
	return formatRupiah(neg, v, sen), nil
}

// ParseRupiah parses an amount with the same syntax as Parse and spells it
// for an invoice. Whole amounts end in "rupiah". One or two fractional digits
// are sen; a single digit is a tenth, so "1,5" is "satu rupiah lima puluh sen".
// More than two fractional digits is an error.
func ParseRupiah(input string) (string, error) {
	p, neg, sen, err := parseCurrency(input)
	if err != nil {
		return "", err
	}
	words, err := spellRupiah(neg, p.whole, sen)
	if err != nil {
		return "", fmt.Errorf("%w: %q", err, strings.TrimSpace(input))
	}
	return words, nil
}

// FormatRupiahString parses an amount with the same syntax as Parse and
// formats it as Indonesian Rupiah. "1250000.5" and "1.250.000,5" both become
// "Rp1.250.000,50". More than two fractional digits is an error.
func FormatRupiahString(input string) (string, error) {
	p, neg, sen, err := parseCurrency(input)
	if err != nil {
		return "", err
	}
	return formatRupiah(neg, p.whole, sen), nil
}

func parseCurrency(input string) (parsed, bool, int, error) {
	p, err := parseNumber(input)
	if err != nil {
		return parsed{}, false, 0, err
	}
	neg, sen, err := p.sen()
	if err != nil {
		return parsed{}, false, 0, fmt.Errorf("%w: %q", err, strings.TrimSpace(input))
	}
	return p, neg, sen, nil
}

func (p parsed) sen() (bool, int, error) {
	sen := 0
	switch len(p.frac) {
	case 0:
	case 1:
		sen = int(p.frac[0]-'0') * 10
	case 2:
		sen = int(p.frac[0]-'0')*10 + int(p.frac[1]-'0')
	default:
		return false, 0, fmt.Errorf("%w: at most two decimal places for rupiah", ErrSyntax)
	}
	showNeg := p.neg && (p.whole.Sign() != 0 || sen != 0)
	return showNeg, sen, nil
}

func spellRupiahAmount(rupiah int64, sen int) (string, error) {
	v := new(big.Int).SetInt64(rupiah)
	neg := v.Sign() < 0
	v.Abs(v)
	return spellRupiah(neg, v, sen)
}

func spellRupiah(neg bool, whole *big.Int, sen int) (string, error) {
	if whole.Sign() == 0 && sen > 0 {
		words := Convert(int64(sen))
		if neg {
			return "minus " + words + " sen", nil
		}
		return words + " sen", nil
	}
	words, err := spellAbs(whole)
	if err != nil {
		return "", err
	}
	if neg {
		words = "minus " + words
	}
	words += " rupiah"
	if sen > 0 {
		words += " " + Convert(int64(sen)) + " sen"
	}
	return words, nil
}

func formatRupiah(neg bool, whole *big.Int, sen int) string {
	digits := whole.String()
	var b strings.Builder
	b.Grow(len(digits) + len(digits)/3 + 8)
	if neg {
		b.WriteByte('-')
	}
	b.WriteString("Rp")
	writeGrouped(&b, digits)
	b.WriteByte(',')
	b.WriteByte(byte('0' + sen/10))
	b.WriteByte(byte('0' + sen%10))
	return b.String()
}

func writeGrouped(b *strings.Builder, digits string) {
	lead := len(digits) % 3
	if lead == 0 {
		lead = 3
	}
	b.WriteString(digits[:lead])
	for i := lead; i < len(digits); i += 3 {
		b.WriteByte('.')
		b.WriteString(digits[i : i+3])
	}
}
