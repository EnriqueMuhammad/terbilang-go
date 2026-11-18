// Command terbilang converts numbers into Indonesian words.
package main

import "os"

// version is overridden by GoReleaser via -X main.version.
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
