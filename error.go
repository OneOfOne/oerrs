package oerrs

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/xerrors"
)

// New returns an error that formats as the given text.
//
// The returned error contains a Frame set to the caller's location and
// implements Formatter to show this information when printed with details.
func New(text string, withJSON bool) error {
	return NewSkipCaller(text, withJSON, 1)
}

// NewSkipCaller returns an error that formats as the given text.
// and skips X frames.
// The returned error contains a Frame set to the caller's location and
// implements Formatter to show this information when printed with details.
func NewSkipCaller(text string, withJSON bool, skip int) error {
	return WrapSkipCaller(stringError(text), withJSON, skip+1)
}

func Wrap(err error, withJSON bool) error {
	return WrapSkipCaller(err, withJSON, 1)
}

func WrapSkipCaller(err error, withJSON bool, skip int) error {
	werr := wrappedError{err, Caller(skip + 1)}
	if withJSON {
		return &jsonError{wrappedError: werr}
	}
	return &werr
}

type stringError string

func (e stringError) Error() string { return string(e) }

// wrappedError is a trivial implementation of error with frame information
type wrappedError struct {
	s     error
	frame *Frame
}

func (e *wrappedError) Error() string {
	log.Printf("%#+v", e.s)
	return e.s.Error()
}

func (e *wrappedError) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *wrappedError) FormatError(p Printer) (next error) {
	p.Print(e.s)
	e.frame.Format(p)
	return nil
}

type jsonError struct {
	wrappedError `json:"-"`

	Msg  string `json:"error,omitempty"`
	Func string `json:"func,omitempty"`
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

func (e *jsonError) MarshalJSON() ([]byte, error) {
	type jsonError_ jsonError
	if e.wrappedError.s == nil && e.Msg != "" {
		return json.Marshal((*jsonError_)(e))
	}
	je := &jsonError_{
		Msg: e.Error(),
	}
	je.Func, je.File, je.Line = e.frame.Location()
	return json.Marshal(je)
}
