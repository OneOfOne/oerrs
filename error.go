package oerrs

import (
	"encoding/json"
	"fmt"

	"golang.org/x/xerrors"
)

func New(s string) error {
	err := String(s)
	if AlwaysWithCaller {
		return err.WithFrame(2)
	}
	return err
}

// String is a plain string error, it can be converted to string and compared
type String string

func (e String) Error() string { return string(e) }
func (e String) WithFrame(skip int) error {
	return Wrap(e, skip+1)
}

func (e String) Is(o error) bool {
	return e.Error() == string(e)
}

func Wrap(err error, skip int) error {
	return &wrapped{err, Caller(skip + 1)}
}

// wrapped is a trivial implementation of error with frame information
type wrapped struct {
	err error
	fr  *Frame
}

func (e *wrapped) Error() string {
	if WrappedErrorTextIncludesFrameInfo {
		return fmt.Sprintf("%s: %v", e.fr.String(), e.err)
	}
	return e.err.Error()
}

func (e *wrapped) Unwrap() error {
	return e.err
}

func (e *wrapped) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *wrapped) FormatError(p Printer) (next error) {
	p.Print(e.err)
	e.fr.Format(p)
	return nil
}

func (e *wrapped) Is(target error) bool {
	return Is(e.err, target)
}

func (e *wrapped) Frame() *Frame { return e.fr }

type JSONError struct {
	wrapped `json:"-"`

	Err  string `json:"error,omitempty"`
	Func string `json:"func,omitempty"`
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

func (e *JSONError) MarshalJSON() ([]byte, error) {
	type jsonError_ JSONError
	if e.err == nil && e.Err != "" {
		return json.Marshal((*jsonError_)(e))
	}
	je := &jsonError_{
		Err: e.err.Error(),
	}
	je.Func, je.File, je.Line = e.fr.Location()
	return json.Marshal(je)
}
