# Contributing

Thanks for looking at terbilang-go. Bug reports, spelling corrections, and small patches are welcome.

## Develop

The module targets Go 1.22 or newer. From a clone:

```bash
go test -race -count=1 -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go vet ./...
golangci-lint run
```

`gofmt` is the formatting standard. Please keep new behavior covered by a table test, and add a godoc example when the public API changes.

## Spellings

Invoice readings are deliberate: `sebelas`, `seratus`, `seribu`, and `satu juta` (not `sejuta`). `miliar` follows the KBBI spelling. If you think a reading is wrong, open an issue with the number, the current output, and the form you expect, ideally with a source (KBBI, an invoice convention, or a statute).
