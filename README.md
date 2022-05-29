# oerrs

[![Go Reference](https://pkg.go.dev/badge/go.oneofone.dev/oerrs.svg)](https://pkg.go.dev/go.oneofone.dev/oerrs)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/OneOfOne/oerrs)

[godocs]: https://pkg.go.dev/go.oneofone.dev/oerrs

`oerrs` is a package for Go that builds on top of [`golang.org/x/xerrors`](https://golang.org/x/xerrors).

Adds an ErrorList with optional stack traces.

## Installation and Docs

Install using `go get go.oneofone.dev/oerrs`.

Package documentation: <https://pkg.go.dev/go.oneofone.dev/oerrs>

### Requires go version 1.18 or newer

## Features

* Passes through most [`golang.org/x/xerrors`](https://golang.org/x/xerrors) functions and interfaces to make life easier.
* A complete drop-in replacement for `xerrors`
* A complete drop-in replacement for `errgroup`
* `ErrorList` handles multiple errors in a sane way.
* `ErrorList` can optionally be toggled to enable thread safety.
* `Try/Catch` *needs doc*.
* All errors can support be toggled to support JSON output.

## Usage

`oerrs` was made to make error handling easier

## Examples

```go
var errs oerrs.ErrorList

errs.PushIf(step1())
errs.PushIf(step2())
if errs.PushIf(step3()) {
	// do something
}
return errs.Err()

```

`oerrs.ErrorList` implements `error`

```go

if err := something(); err != nil {
	if el, ok := err.(*oerrs.ErrorList); ok {
		// Use merr.Errors
	}
	// or
	var el *oerrs.ErrorList
	if errors.As(err, &el) {
		// use el.Errors
	}
}


```

```go
// try / catch

func doStuff() (res Response, err error) {
	defer oerrs.Catch(&err, nil)

	a := Try1(someFunThatReturnValAndErr())
	// do something with a
	b := Try1(someOtherFunc(a))

	return a + b, nil
}

```

## License

The code of errgroup is mostly copied from  [`golang.org/x/sync`](https://golang.org/x/sync) for personal convience, all rights reserved to go devs.

MIT
