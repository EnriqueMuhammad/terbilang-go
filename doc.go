// Package terbilang converts numbers into Indonesian words and formats
// Indonesian Rupiah amounts for invoices and receipts.
//
//	fmt.Println(terbilang.Convert(1250000))
//	// satu juta dua ratus lima puluh ribu
//
//	fmt.Println(terbilang.Rupiah(1250000))
//	// satu juta dua ratus lima puluh ribu rupiah
//
//	fmt.Println(terbilang.FormatRupiah(1250000))
//	// Rp1.250.000,00
package terbilang
