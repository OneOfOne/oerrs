package oerrs

import (
	"encoding/json"
	"fmt"

	"golang.org/x/xerrors"
)

func New(s string) error {
	err := String(s)
	if AlwaysWithCaller {
		return err.WithFrame(1)
	}
	return err
}

func Wrap(err error, skip int) error {
	return &wrappedError{err, Caller(skip + 1)}
}

// String is a plain string error, it can be converted to string and compared
type String string

func (e String) Error() string { return string(e) }
func (e String) WithFrame(skip int) error {
	return &wrappedError{
		err: e,
		fr:  Caller(skip + 1),
	}
}

func (e String) Is(o error) bool {
	return e.Error() == string(e)
}

// wrappedError is a trivial implementation of error with frame information
type wrappedError struct {
	err error
	fr  *Frame
}

func (e *wrappedError) Error() string {
	return e.err.Error()
}

func (e *wrappedError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *wrappedError) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *wrappedError) FormatError(p Printer) (next error) {
	p.Print(e.err)
	e.fr.Format(p)
	return nil
}

func (e *wrappedError) Is(target error) bool {
	return Is(e.err, target)
}

func (e *wrappedError) Frame() *Frame { return e.fr }

type jsonError struct {
	wrappedError `json:"-"`

	Msg  string `json:"error,omitempty"`
	Func string `json:"func,omitempty"`
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

func (e *jsonError) MarshalJSON() ([]byte, error) {
	type jsonError_ jsonError
	if e.err == nil && e.Msg != "" {
		return json.Marshal((*jsonError_)(e))
	}
	je := &jsonError_{
		Msg: e.Error(),
	}
	je.Func, je.File, je.Line = e.fr.Location()
	return json.Marshal(je)
}
