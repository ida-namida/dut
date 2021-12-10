# dang

This is a simple and (currently) deterministic decision rule execution service written in [Go](https://go.dev). Inspired by [Jboss' KIE]https://docs.jboss.org/drools/release/7.49.0.Final/drools-docs/html_single/index.html#_kiechapter)

## Limitations

- Currently accepted form field value types are limited to numbers (covering both integer and float) and boolean. Also string, but it hasn't been supported in any of the currently supported decision rules
- Currently available rules are only mapping numeric interval to output, mapping boolean condition to output, mapping threshold-comparison to output, and sum of numeric arguments
- For persistence, this initial implementation relies on the [Kopuro](https://github.com/ida-namida/kopuro) project, which serves solely as an auxiliary for this and the [Dut](https://github.com/ida-namida/dut) project

## Contributions

Contributions via pull requests are **very** welcome, especially to:
- support more form field value types
- add more decision rules
- add more implementation of configurable persistence options

The implementation is based on the usage of Golang's [text/template](https://pkg.go.dev/text/template) package