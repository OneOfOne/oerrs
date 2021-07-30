package oerrs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/xerrors"
)

func NewList(includeCallers bool) *ErrorList {
	return &ErrorList{
		saveCaller: includeCallers,
	}
}

func NewListWithJSON(includeCallers bool) *ErrorList {
	return &ErrorList{
		saveCaller:  includeCallers,
		jsonExport:  true,
		jsonCallers: includeCallers,
	}
}

type ErrorList struct {
	errs        []error
	mux         *sync.RWMutex
	saveCaller  bool // if true, Error() will return the caller
	noFlatten   bool // if true, won't flatten child ErrorLists
	jsonExport  bool
	jsonCallers bool
}

func (e *ErrorList) IncludeCallers(v bool) *ErrorList {
	defer e.lock(true)()
	e.saveCaller = v
	return e
}

func (e *ErrorList) Flatten(v bool) *ErrorList {
	defer e.lock(true)()
	e.noFlatten = !v
	return e
}

func (e *ErrorList) SupportJSON(on, includeCallers bool) *ErrorList {
	defer e.lock(true)()
	e.jsonExport, e.jsonCallers = on, includeCallers
	if includeCallers {
		e.saveCaller = includeCallers
	}
	return e
}

func (e *ErrorList) ThreadSafe() {
	if e.mux != nil {
		panic("race")
	}
	e.mux = new(sync.RWMutex)
}

func (e *ErrorList) PushIf(err error) bool {
	return e.PushIfSkipCaller(err, 1)
}

func (e *ErrorList) PushIfSkipCaller(err error, skip int) bool {
	if err == nil {
		return false
	}
	defer e.lock(true)()

	if !e.noFlatten {
		if oe, ok := err.(*ErrorList); ok {
			if oe == e {
				return false
			}
			defer oe.lock(false)()
			e.errs = append(e.errs, oe.errs...)
			return true
		}
	}

	if e.saveCaller {
		err = WrapSkipCaller(err, e.jsonExport, skip+1)
	}

	e.errs = append(e.errs, err)
	return true
}

func (e *ErrorList) Reset() {
	if e == nil {
		return
	}

	defer e.lock(true)()
	e.errs = nil
}

func (e *ErrorList) Err() error {
	if e == nil {
		return nil
	}

	defer e.lock(false)()

	nilCount := 0
	for _, err := range e.errs {
		if err == nil {
			nilCount++
		}
	}
	if nilCount == len(e.errs) {
		return nil
	}
	return e
}

// Error implements the error interface
func (e *ErrorList) Error() string {
	defer e.lock(false)()
	return fmtErrors(e.errs, e.saveCaller)
}

// Unwrap implements errors.Unwrap by returning the next error in the
// chain or nil if there are no more errors.
func (e *ErrorList) Unwrap() error {
	defer e.lock(false)()
	if len(e.errs) == 1 {
		return nil
	}

	return &ErrorList{errs: e.errs[1:], saveCaller: e.saveCaller}
}

// As implements errors.As by checking if target is ErrorList, otherwise will go through each error.
func (e *ErrorList) As(target interface{}) bool {
	defer e.lock(false)()
	if oe, ok := target.(*ErrorList); ok {
		defer oe.lock(false)()
		e.errs, e.saveCaller = oe.errs, oe.saveCaller
		e.jsonExport, e.jsonCallers = oe.jsonExport, oe.jsonCallers
		return true
	}

	for _, err := range e.errs {
		if xerrors.As(err, target) {
			return true
		}
	}
	return false
}

// Is implements errors.Is by comparing if target is ErrorList, otherwise will go through each error.
func (e *ErrorList) Is(target error) bool {
	defer e.lock(false)()

	if oe, ok := target.(*ErrorList); ok {
		defer oe.lock(false)()

		if len(e.errs) != len(oe.errs) {
			return false
		}

		for i, err := range e.errs {
			if !xerrors.Is(err, oe.errs[i]) {
				return false
			}
		}

		return true
	}

	for _, err := range e.errs {
		if xerrors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *ErrorList) Errors() []error {
	defer e.lock(false)()
	errs := make([]error, 0, len(e.errs))
	return append(errs, e.errs...)
}

func (e *ErrorList) MarshalJSON() ([]byte, error) {
	defer e.lock(false)()
	if e == nil || !e.jsonExport || len(e.errs) == 0 {
		return []byte("nil"), nil
	}
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, err := range e.errs {
		if err == nil {
			continue
		}
		je := &jsonError{
			Msg: err.Error(),
		}
		if e.jsonCallers {
			if err, ok := err.(*wrappedError); ok && err.frame != nil {
				je.Func, je.File, je.Line = err.frame.Location()
			}
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		b, _ := json.Marshal(je)
		buf.Write(b)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (e *ErrorList) lock(rw bool) func() {
	if e == nil || e.mux == nil {
		return func() {}
	}
	if rw {
		e.mux.Lock()
		return e.mux.Unlock
	}
	e.mux.RLock()
	return e.mux.RUnlock
}

func fmtErrors(es []error, withStack bool) string {
	if len(es) == 0 {
		return ""
	}

	if len(es) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", es[0])
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d errors occurred:", len(es))

	verb := "%v"
	if withStack {
		verb = "%+v"
	}
	for _, err := range es {
		buf.WriteString("\n\t* ")
		fmt.Fprintf(&buf, verb, err)
	}

	return buf.String()
}
