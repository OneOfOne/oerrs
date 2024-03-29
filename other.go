package oerrs

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
)

var (
	AlwaysWithCaller                  = true
	WrappedErrorTextIncludesFrameInfo = false
)

// type aliases from xerrors
type (
	Formatter = xerrors.Formatter
	Printer   = xerrors.Printer
	Wrapper   = xerrors.Wrapper
)

// HasIs is a handy interface to check if an error implements Is
type HasIs interface {
	Is(error) bool
}

// HasAs is a handy interface to check if an error implements As
type HasAs interface {
	As(interface{}) bool
}

// As is an alias to errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is is an alias to errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Unwrap is an alias to errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Opaque is an alias to xerrors.Opaque
func Opaque(err error) error {
	return xerrors.Opaque(err)
}

func Unframe(err error) (*Frame, error) {
	for err != nil {
		if fr, ok := err.(Framer); ok {
			return fr.Frame(), err
		}
		err = Unwrap(err)
	}
	return nil, nil
}

func Errorf(format string, args ...interface{}) error {
	if AlwaysWithCaller {
		return ErrorCallerf(1, format, args...)
	}
	return fmterrorf(format, args...)
}

func ErrorCallerf(skip int, format string, args ...interface{}) error {
	err := fmterrorf(format, args...)
	return &wrapped{err, Caller(skip + 1)}
}

func fmterrorf(format string, args ...interface{}) error {
	if len(args) == 0 {
		return String(format)
	}
	return fmt.Errorf(format, args...)
}

func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &ErrorList{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}
