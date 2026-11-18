package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EnriqueMuhammad/terbilang-go"
)

const usageText = `terbilang converts numbers into Indonesian words.

Usage:
  terbilang [flags] [number ...]

If no numbers are given, numbers are read from standard input.
A line may hold one number (spaces allowed, as in "Rp 1.250.000,00")
or several numbers separated by spaces. Empty lines and lines starting
with # are ignored.

Flags:
  -format
    	print Indonesian Rupiah format (Rp1.250.000,00)
  -rupiah
    	spell the amount with a "rupiah" suffix and "sen" for fractions
  -version
    	print the version and exit

Use both -format and -rupiah to print the two forms together.
Negative numbers can be written directly: terbilang -11
`

func run(args []string, in io.Reader, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("terbilang", flag.ContinueOnError)
	fs.SetOutput(errOut)
	format := fs.Bool("format", false, "print Indonesian Rupiah format (Rp1.250.000,00)")
	rupiah := fs.Bool("rupiah", false, "spell with a rupiah suffix")
	showVersion := fs.Bool("version", false, "print version")
	fs.Usage = func() {
		_, _ = fmt.Fprint(errOut, usageText)
	}

	flagArgs, nums := splitArgs(args)
	if err := fs.Parse(flagArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if *showVersion {
		_, _ = fmt.Fprintf(out, "terbilang %s\n", version)
		return 0
	}
	nums = append(nums, fs.Args()...)

	if len(nums) == 0 {
		var err error
		nums, err = readNumbers(in)
		if err != nil {
			_, _ = fmt.Fprintf(errOut, "terbilang: read input: %v\n", err)
			return 1
		}
	}
	if len(nums) == 0 {
		fs.Usage()
		return 2
	}

	both := *format && *rupiah
	for i, raw := range nums {
		if err := writeOne(out, raw, *format, *rupiah); err != nil {
			_, _ = fmt.Fprintf(errOut, "terbilang: %v\n", err)
			return 1
		}
		if both && i != len(nums)-1 {
			if _, err := fmt.Fprintln(out); err != nil {
				_, _ = fmt.Fprintf(errOut, "terbilang: %v\n", err)
				return 1
			}
		}
	}
	return 0
}

func writeOne(out io.Writer, raw string, format, rupiah bool) error {
	if !format && !rupiah {
		words, err := terbilang.Parse(raw)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, words)
		return err
	}
	if format {
		formatted, err := terbilang.FormatRupiahString(raw)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, formatted); err != nil {
			return err
		}
	}
	if rupiah {
		words, err := terbilang.ParseRupiah(raw)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, words); err != nil {
			return err
		}
	}
	return nil
}

func splitArgs(args []string) (flags, nums []string) {
	for i, a := range args {
		if a == "--" {
			return flags, append(nums, args[i+1:]...)
		}
		if looksLikeNumber(a) {
			nums = append(nums, a)
			continue
		}
		flags = append(flags, a)
	}
	return flags, nums
}

// looksLikeNumber reports whether a should be parsed as a number instead of a flag.
// Valid numbers qualify, and so do broken ones that start like a number, so
// "1.2.3" is reported as invalid input rather than as an unknown flag.
func looksLikeNumber(s string) bool {
	if _, err := terbilang.Parse(s); err == nil {
		return true
	}
	rest := s
	if rest == "" {
		return false
	}
	if rest[0] == '+' || rest[0] == '-' {
		rest = strings.TrimSpace(rest[1:])
	}
	if rest == "" {
		return false
	}
	c := rest[0]
	if c == '.' || c == ',' || (c >= '0' && c <= '9') {
		return true
	}
	return len(rest) >= 2 && (rest[0] == 'R' || rest[0] == 'r') && (rest[1] == 'P' || rest[1] == 'p')
}

func readNumbers(in io.Reader) ([]string, error) {
	if f, ok := in.(*os.File); ok {
		fi, err := f.Stat()
		if err == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return nil, nil
		}
	}

	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var nums []string
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		nums = append(nums, tokensFromLine(line)...)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return nums, nil
}

func tokensFromLine(line string) []string {
	if _, err := terbilang.Parse(line); err == nil {
		return []string{line}
	}
	fields := strings.Fields(line)
	if len(fields) <= 1 {
		return []string{line}
	}
	return fields
}
