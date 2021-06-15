package oerrs

import "golang.org/x/xerrors"

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

// As is an alias to xerrors.As
func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

// Is is an alias to xerrors.Is
func Is(err, target error) bool {
	return xerrors.Is(err, target)
}

// Unwrap is an alias to xerrors.Unwrap
func Unwrap(err error) error {
	return xerrors.Unwrap(err)
}

// Opaque is an alias to xerrors.Opaque
func Opaque(err error) error {
	return xerrors.Opaque(err)
}

// Errorf is an alias to xerrors.Errorf
func Errorf(format string, a ...interface{}) error {
	return xerrors.Errorf(format, a...)
}
