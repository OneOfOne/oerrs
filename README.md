# oerrs

[![Go Reference](https://pkg.go.dev/badge/go.oneofone.dev/oerrs.svg)](https://pkg.go.dev/go.oneofone.dev/oerrs)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/go.oneofone.dev/oerrs)

[godocs]: https://pkg.go.dev/go.oneofone.dev/oerrs

`oerrs` is a package for Go that builds on top of [`golang.org/x/xerrors`](https://golang.org/x/xerrors)

Adds an ErrorList with optional stack traces.

## Installation and Docs

Install using `go get go.oneofone.dev/oerrs`.

Package documentation: https://pkg.go.dev/go.oneofone.dev/oerrs

### Requires go version 1.15 or newer


## Features

* Passes through most [`golang.org/x/xerrors`](https://golang.org/x/xerrors) functions and interfaces to make life easier.
* A complete drop-in replacement for `xerrors`
* `ErrorList` handles multiple errors in a sane way.
* `ErrorList` can optionally be toggled to enable thread safety.
* All errors can support be toggled to support JSON output.

## Usage

`oerrs` was made to make error handling easier

**Example**

```go
var errs oerrs.ErrorList

// if you want know the caller:
errs.RecorderCaller = true

errs.PushIf(step1())
errs.PushIf(step2())
if errs.PushIf(step3()) {
	// do something
}
return errs.Err()


**Accessing the list of errors**

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

## License

MIT
