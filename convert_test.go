package terbilang_test

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/EnriqueMuhammad/terbilang-go"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n    int64
		want string
	}{
		{0, "nol"},
		{1, "satu"},
		{2, "dua"},
		{9, "sembilan"},
		{10, "sepuluh"},
		{11, "sebelas"},
		{12, "dua belas"},
		{15, "lima belas"},
		{19, "sembilan belas"},
		{20, "dua puluh"},
		{21, "dua puluh satu"},
		{30, "tiga puluh"},
		{99, "sembilan puluh sembilan"},
		{100, "seratus"},
		{101, "seratus satu"},
		{110, "seratus sepuluh"},
		{111, "seratus sebelas"},
		{115, "seratus lima belas"},
		{121, "seratus dua puluh satu"},
		{200, "dua ratus"},
		{999, "sembilan ratus sembilan puluh sembilan"},
		{1000, "seribu"},
		{1001, "seribu satu"},
		{1011, "seribu sebelas"},
		{1100, "seribu seratus"},
		{1111, "seribu seratus sebelas"},
		{2000, "dua ribu"},
		{2500, "dua ribu lima ratus"},
		{10000, "sepuluh ribu"},
		{10011, "sepuluh ribu sebelas"},
		{11000, "sebelas ribu"},
		{11111, "sebelas ribu seratus sebelas"},
		{100000, "seratus ribu"},
		{101000, "seratus satu ribu"},
		{110000, "seratus sepuluh ribu"},
		{111000, "seratus sebelas ribu"},
		{200000, "dua ratus ribu"},
		{1000000, "satu juta"},
		{1000001, "satu juta satu"},
		{1001000, "satu juta seribu"},
		{1011000, "satu juta sebelas ribu"},
		{1100100, "satu juta seratus ribu seratus"},
		{1250000, "satu juta dua ratus lima puluh ribu"},
		{2000000, "dua juta"},
		{123456789, "seratus dua puluh tiga juta empat ratus lima puluh enam ribu tujuh ratus delapan puluh sembilan"},
		{1000000000, "satu miliar"},
		{1000000000000, "satu triliun"},
		{1000000000000000, "satu kuadriliun"},
		{1000000000000000000, "satu kuintiliun"},
		{-1, "minus satu"},
		{-11, "minus sebelas"},
		{-100, "minus seratus"},
		{-1000, "minus seribu"},
		{-1250000, "minus satu juta dua ratus lima puluh ribu"},
		{math.MaxInt64, "sembilan kuintiliun dua ratus dua puluh tiga kuadriliun tiga ratus tujuh puluh dua triliun tiga puluh enam miliar delapan ratus lima puluh empat juta tujuh ratus tujuh puluh lima ribu delapan ratus tujuh"},
		{math.MinInt64, "minus sembilan kuintiliun dua ratus dua puluh tiga kuadriliun tiga ratus tujuh puluh dua triliun tiga puluh enam miliar delapan ratus lima puluh empat juta tujuh ratus tujuh puluh lima ribu delapan ratus delapan"},
	}

	for _, tt := range tests {
		got := terbilang.Convert(tt.n)
		if got != tt.want {
			t.Errorf("Convert(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestConvertSignAndSpacing(t *testing.T) {
	t.Parallel()

	for n := int64(0); n <= 2500; n++ {
		got := terbilang.Convert(n)
		assertWords(t, n, got)
		if n == 0 {
			continue
		}
		neg := terbilang.Convert(-n)
		if neg != "minus "+got {
			t.Fatalf("Convert(%d) = %q, want %q", -n, neg, "minus "+got)
		}
	}
}

func TestScales(t *testing.T) {
	t.Parallel()

	scales := []struct {
		zeros int
		want  string
	}{
		{3, "seribu"},
		{4, "sepuluh ribu"},
		{5, "seratus ribu"},
		{6, "satu juta"},
		{9, "satu miliar"},
		{12, "satu triliun"},
		{15, "satu kuadriliun"},
		{18, "satu kuintiliun"},
		{21, "satu sekstiliun"},
		{24, "satu septiliun"},
		{27, "satu oktiliun"},
		{30, "satu noniliun"},
		{33, "satu desiliun"},
	}

	for _, sc := range scales {
		in := "1" + strings.Repeat("0", sc.zeros)
		got, err := terbilang.Parse(in)
		if err != nil {
			t.Fatalf("Parse(10^%d): %v", sc.zeros, err)
		}
		if got != sc.want {
			t.Fatalf("Parse(10^%d) = %q, want %q", sc.zeros, got, sc.want)
		}
		if strings.HasPrefix(sc.want, "satu ") {
			in2 := "2" + strings.Repeat("0", sc.zeros)
			got2, err := terbilang.Parse(in2)
			if err != nil {
				t.Fatal(err)
			}
			want2 := "dua " + strings.TrimPrefix(sc.want, "satu ")
			if got2 != want2 {
				t.Fatalf("Parse(2*10^%d) = %q, want %q", sc.zeros, got2, want2)
			}
		}
	}

	if _, err := terbilang.Parse("2" + strings.Repeat("0", 3)); err != nil {
		t.Fatal(err)
	} else if got, _ := terbilang.Parse("2000"); got != "dua ribu" {
		t.Fatalf("dua ribu = %q", got)
	}

	biggest := strings.Repeat("9", 36)
	got, err := terbilang.Parse(biggest)
	if err != nil {
		t.Fatalf("36 nines: %v", err)
	}
	if !strings.Contains(got, "desiliun") || strings.Contains(got, "  ") {
		t.Fatalf("36 nines = %q", got)
	}

	_, err = terbilang.Parse("1" + strings.Repeat("0", 36))
	if !errors.Is(err, terbilang.ErrMagnitude) {
		t.Fatalf("10^36 error = %v, want ErrMagnitude", err)
	}
}

func assertWords(t *testing.T, n int64, got string) {
	t.Helper()
	if got == "" {
		t.Fatalf("Convert(%d) is empty", n)
	}
	if strings.TrimSpace(got) != got || strings.Contains(got, "  ") {
		t.Fatalf("Convert(%d) = %q has broken spacing", n, got)
	}
	if n == 0 && got != "nol" {
		t.Fatalf("zero = %q", got)
	}
	if n > 0 && strings.Contains(got, "nol") {
		t.Fatalf("Convert(%d) = %q contains nol", n, got)
	}
}
