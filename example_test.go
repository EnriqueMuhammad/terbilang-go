package terbilang_test

import (
	"fmt"

	"github.com/EnriqueMuhammad/terbilang-go"
)

func ExampleConvert() {
	fmt.Println(terbilang.Convert(1250000))
	fmt.Println(terbilang.Convert(11))
	fmt.Println(terbilang.Convert(-100))
	fmt.Println(terbilang.Convert(0))
	// Output:
	// satu juta dua ratus lima puluh ribu
	// sebelas
	// minus seratus
	// nol
}

func ExampleParse() {
	fmt.Println(must(terbilang.Parse("10,05")))
	fmt.Println(must(terbilang.Parse("Rp1.250.000")))
	// Output:
	// sepuluh koma nol lima
	// satu juta dua ratus lima puluh ribu
}

func ExampleRupiah() {
	fmt.Println(terbilang.Rupiah(1250000))
	fmt.Println(terbilang.Rupiah(-1500))
	// Output:
	// satu juta dua ratus lima puluh ribu rupiah
	// minus seribu lima ratus rupiah
}

func ExampleRupiahSen() {
	fmt.Println(must(terbilang.RupiahSen(1500, 50)))
	fmt.Println(must(terbilang.RupiahSen(0, 50)))
	// Output:
	// seribu lima ratus rupiah lima puluh sen
	// lima puluh sen
}

func ExampleFormatRupiah() {
	fmt.Println(terbilang.FormatRupiah(1250000))
	fmt.Println(terbilang.FormatRupiah(-1500))
	// Output:
	// Rp1.250.000,00
	// -Rp1.500,00
}

func ExampleFormatRupiahSen() {
	fmt.Println(must(terbilang.FormatRupiahSen(1500, 5)))
	// Output: Rp1.500,05
}

func ExampleParseRupiah() {
	fmt.Println(must(terbilang.ParseRupiah("1.250.000,5")))
	// Output: satu juta dua ratus lima puluh ribu rupiah lima puluh sen
}

func ExampleFormatRupiahString() {
	fmt.Println(must(terbilang.FormatRupiahString("1250000.5")))
	// Output: Rp1.250.000,50
}

func must(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}
