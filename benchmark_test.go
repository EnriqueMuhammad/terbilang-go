package terbilang_test

import (
	"testing"

	"github.com/EnriqueMuhammad/terbilang-go"
)

var blackhole string

func BenchmarkConvert(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		blackhole = terbilang.Convert(9223372036854775807)
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		blackhole, _ = terbilang.Parse("Rp1.250.000,50")
	}
}

func BenchmarkFormatRupiah(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		blackhole = terbilang.FormatRupiah(1250000)
	}
}

func BenchmarkRupiah(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		blackhole = terbilang.Rupiah(1250000)
	}
}
