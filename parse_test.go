package terbilang_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/EnriqueMuhammad/terbilang-go"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"0", "nol"},
		{"000", "nol"},
		{"  1000  ", "seribu"},
		{"+11", "sebelas"},
		{"-0", "nol"},
		{"-11", "minus sebelas"},
		{"1250000", "satu juta dua ratus lima puluh ribu"},
		{"1.250.000", "satu juta dua ratus lima puluh ribu"},
		{"1,250,000", "satu juta dua ratus lima puluh ribu"},
		{"10.000", "sepuluh ribu"},
		{"1.250", "seribu dua ratus lima puluh"},
		{"3.141", "tiga ribu seratus empat puluh satu"},
		{"1,5", "satu koma lima"},
		{"1,25", "satu koma dua lima"},
		{"1,250", "satu koma dua lima nol"},
		{"1.5", "satu koma lima"},
		{"1.50", "satu koma lima nol"},
		{"1.05", "satu koma nol lima"},
		{"0.5", "nol koma lima"},
		{"0.05", "nol koma nol lima"},
		{"0.250", "nol koma dua lima nol"},
		{"0,250", "nol koma dua lima nol"},
		{".5", "nol koma lima"},
		{",5", "nol koma lima"},
		{"-0,5", "minus nol koma lima"},
		{"-1.05", "minus satu koma nol lima"},
		{"3,141", "tiga koma satu empat satu"},
		{"3.14159", "tiga koma satu empat satu lima sembilan"},
		{"12.3456", "dua belas koma tiga empat lima enam"},
		{"1.250.000,50", "satu juta dua ratus lima puluh ribu koma lima nol"},
		{"1,250,000.50", "satu juta dua ratus lima puluh ribu koma lima nol"},
		{"Rp1.250.000,00", "satu juta dua ratus lima puluh ribu koma nol nol"},
		{"Rp 1.250.000", "satu juta dua ratus lima puluh ribu"},
		{"rp10.000", "sepuluh ribu"},
		{"- Rp 11", "minus sebelas"},
		{"+Rp 100", "seratus"},
		{"010.000", "sepuluh ribu"},
	}

	for _, tt := range tests {
		got, err := terbilang.Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Parse(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	empty := []string{"", "   ", "\t"}
	for _, in := range empty {
		_, err := terbilang.Parse(in)
		if !errors.Is(err, terbilang.ErrEmpty) {
			t.Errorf("Parse(%q) error = %v, want ErrEmpty", in, err)
		}
	}

	invalid := []string{
		"abc",
		"-",
		"+",
		"Rp",
		"Rp ",
		"--1",
		"1.",
		"1,",
		"1.2.3",
		"12.34.567",
		"1,2,3",
		"12,34.567,89",
		"Rp-1",
		"1 000",
		"1.250.",
	}
	for _, in := range invalid {
		_, err := terbilang.Parse(in)
		if !errors.Is(err, terbilang.ErrSyntax) {
			t.Errorf("Parse(%q) error = %v, want ErrSyntax", in, err)
		}
	}

	longFrac := "0," + strings.Repeat("1", 65)
	if _, err := terbilang.Parse(longFrac); !errors.Is(err, terbilang.ErrSyntax) {
		t.Fatalf("long fraction error = %v", err)
	}
	okFrac := "0," + strings.Repeat("1", 64)
	got, err := terbilang.Parse(okFrac)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "nol koma ") || strings.Count(got, "satu") != 64 {
		t.Fatalf("64 digits = %q", got)
	}
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		"0",
		"1250000",
		"-11",
		"Rp1.250.000,50",
		"1,5",
		"10.000",
		"",
		"3.14159",
		"-0,05",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		words, err := terbilang.Parse(s)
		if err != nil {
			if words != "" {
				t.Fatalf("Parse error returned %q", words)
			}
			return
		}
		if words == "" || strings.TrimSpace(words) != words || strings.Contains(words, "  ") {
			t.Fatalf("Parse(%q) = %q", s, words)
		}
		cur, curErr := terbilang.ParseRupiah(s)
		if curErr == nil && (cur == "" || strings.Contains(cur, "  ")) {
			t.Fatalf("ParseRupiah(%q) = %q", s, cur)
		}
		formatted, fmtErr := terbilang.FormatRupiahString(s)
		if fmtErr == nil && !strings.Contains(formatted, "Rp") {
			t.Fatalf("FormatRupiahString(%q) = %q", s, formatted)
		}
	})
}
