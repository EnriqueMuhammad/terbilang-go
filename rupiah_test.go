package terbilang_test

import (
	"errors"
	"math"
	"testing"

	"github.com/EnriqueMuhammad/terbilang-go"
)

func TestRupiah(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n    int64
		want string
	}{
		{0, "nol rupiah"},
		{1, "satu rupiah"},
		{11, "sebelas rupiah"},
		{1000, "seribu rupiah"},
		{1250000, "satu juta dua ratus lima puluh ribu rupiah"},
		{-1250000, "minus satu juta dua ratus lima puluh ribu rupiah"},
		{-1, "minus satu rupiah"},
	}
	for _, tt := range tests {
		if got := terbilang.Rupiah(tt.n); got != tt.want {
			t.Errorf("Rupiah(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestRupiahSen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		rupiah int64
		sen    int
		want   string
	}{
		{0, 0, "nol rupiah"},
		{0, 5, "lima sen"},
		{0, 50, "lima puluh sen"},
		{1, 0, "satu rupiah"},
		{1500, 50, "seribu lima ratus rupiah lima puluh sen"},
		{100, 11, "seratus rupiah sebelas sen"},
		{-20, 5, "minus dua puluh rupiah lima sen"},
		{-20, 0, "minus dua puluh rupiah"},
	}
	for _, tt := range tests {
		got, err := terbilang.RupiahSen(tt.rupiah, tt.sen)
		if err != nil {
			t.Errorf("RupiahSen(%d, %d) error = %v", tt.rupiah, tt.sen, err)
			continue
		}
		if got != tt.want {
			t.Errorf("RupiahSen(%d, %d) = %q, want %q", tt.rupiah, tt.sen, got, tt.want)
		}
	}

	for _, sen := range []int{-1, 100, 150} {
		if _, err := terbilang.RupiahSen(1, sen); !errors.Is(err, terbilang.ErrSen) {
			t.Errorf("RupiahSen sen %d error = %v, want ErrSen", sen, err)
		}
		if _, err := terbilang.FormatRupiahSen(1, sen); !errors.Is(err, terbilang.ErrSen) {
			t.Errorf("FormatRupiahSen sen %d error = %v, want ErrSen", sen, err)
		}
	}
}

func TestFormatRupiah(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n    int64
		want string
	}{
		{0, "Rp0,00"},
		{1, "Rp1,00"},
		{12, "Rp12,00"},
		{100, "Rp100,00"},
		{999, "Rp999,00"},
		{1000, "Rp1.000,00"},
		{1250000, "Rp1.250.000,00"},
		{1000000000, "Rp1.000.000.000,00"},
		{-1250000, "-Rp1.250.000,00"},
		{-1, "-Rp1,00"},
		{math.MaxInt64, "Rp9.223.372.036.854.775.807,00"},
		{math.MinInt64, "-Rp9.223.372.036.854.775.808,00"},
	}
	for _, tt := range tests {
		if got := terbilang.FormatRupiah(tt.n); got != tt.want {
			t.Errorf("FormatRupiah(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestFormatRupiahSen(t *testing.T) {
	t.Parallel()

	got, err := terbilang.FormatRupiahSen(1250000, 50)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Rp1.250.000,50" {
		t.Fatalf("got %q", got)
	}
	got, err = terbilang.FormatRupiahSen(1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Rp1,05" {
		t.Fatalf("got %q", got)
	}
	got, err = terbilang.FormatRupiahSen(0, 50)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Rp0,50" {
		t.Fatalf("got %q", got)
	}
	got, err = terbilang.FormatRupiahSen(-1250000, 50)
	if err != nil {
		t.Fatal(err)
	}
	if got != "-Rp1.250.000,50" {
		t.Fatalf("got %q", got)
	}
}

func TestCurrencyStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in    string
		words string
		fmt   string
	}{
		{"1250000", "satu juta dua ratus lima puluh ribu rupiah", "Rp1.250.000,00"},
		{"1250000,5", "satu juta dua ratus lima puluh ribu rupiah lima puluh sen", "Rp1.250.000,50"},
		{"1250000.50", "satu juta dua ratus lima puluh ribu rupiah lima puluh sen", "Rp1.250.000,50"},
		{"0,05", "lima sen", "Rp0,05"},
		{"-0,50", "minus lima puluh sen", "-Rp0,50"},
		{"-0", "nol rupiah", "Rp0,00"},
		{"Rp1.250.000", "satu juta dua ratus lima puluh ribu rupiah", "Rp1.250.000,00"},
		{"Rp 1.250,5", "seribu dua ratus lima puluh rupiah lima puluh sen", "Rp1.250,50"},
		{"10.000", "sepuluh ribu rupiah", "Rp10.000,00"},
		{"1,5", "satu rupiah lima puluh sen", "Rp1,50"},
		{"-1.05", "minus satu rupiah lima sen", "-Rp1,05"},
	}
	for _, tt := range tests {
		words, err := terbilang.ParseRupiah(tt.in)
		if err != nil {
			t.Errorf("ParseRupiah(%q) error = %v", tt.in, err)
			continue
		}
		if words != tt.words {
			t.Errorf("ParseRupiah(%q) = %q, want %q", tt.in, words, tt.words)
		}
		formatted, err := terbilang.FormatRupiahString(tt.in)
		if err != nil {
			t.Errorf("FormatRupiahString(%q) error = %v", tt.in, err)
			continue
		}
		if formatted != tt.fmt {
			t.Errorf("FormatRupiahString(%q) = %q, want %q", tt.in, formatted, tt.fmt)
		}
	}

	if _, err := terbilang.ParseRupiah("1.2345"); !errors.Is(err, terbilang.ErrSyntax) {
		t.Fatalf("too many places: %v", err)
	}
	if _, err := terbilang.FormatRupiahString("3,141"); !errors.Is(err, terbilang.ErrSyntax) {
		t.Fatalf("three decimal digits: %v", err)
	}
	if _, err := terbilang.ParseRupiah("nope"); !errors.Is(err, terbilang.ErrSyntax) {
		t.Fatalf("syntax: %v", err)
	}
	if _, err := terbilang.ParseRupiah("   "); !errors.Is(err, terbilang.ErrEmpty) {
		t.Fatalf("empty: %v", err)
	}
}

func TestFormatRoundTrip(t *testing.T) {
	t.Parallel()

	samples := []int64{0, 1, 11, 100, 1000, 1250000, -1250000, math.MaxInt64, math.MinInt64}
	for _, n := range samples {
		formatted := terbilang.FormatRupiah(n)
		got, err := terbilang.ParseRupiah(formatted)
		if err != nil {
			t.Fatalf("round trip %d via %s: %v", n, formatted, err)
		}
		if got != terbilang.Rupiah(n) {
			t.Fatalf("round trip %d: %q via %s, want %q", n, got, formatted, terbilang.Rupiah(n))
		}
	}

	formatted, err := terbilang.FormatRupiahSen(1500, 5)
	if err != nil {
		t.Fatal(err)
	}
	got, err := terbilang.ParseRupiah(formatted)
	if err != nil {
		t.Fatal(err)
	}
	want, err := terbilang.RupiahSen(1500, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("sen round trip %q -> %q, want %q", formatted, got, want)
	}
}
