package oerrs

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

func NewList() *ErrorList {
	return &ErrorList{}
}

func NewSafeList(includeCallers bool) *ErrorList {
	var mux sync.RWMutex
	return &ErrorList{
		mux: &mux,
	}
}

type ErrorList struct {
	mux  *sync.RWMutex
	errs []error
}

func (e *ErrorList) Safe() {
	if e.mux != nil {
		panic("race")
	}
	e.mux = new(sync.RWMutex)
}

func (e *ErrorList) Errorf(format string, args ...interface{}) {
	err := Errorf(format, args...)
	defer e.lock(true)()
	if AlwaysWithCaller {
		err = &wrapped{err, Caller(1)}
	}
	e.errs = append(e.errs, err)
}

func (e *ErrorList) ErrorCallerf(skip int, format string, args ...interface{}) {
	err := ErrorCallerf(skip+1, format, args...)
	defer e.lock(true)()
	e.errs = append(e.errs, err)
}

func (e *ErrorList) PushIf(err error) bool {
	if err == nil {
		return false
	}
	if AlwaysWithCaller {
		err = &wrapped{err, Caller(1)}
	}
	defer e.lock(true)()
	e.errs = append(e.errs, err)
	return true
}

func (e *ErrorList) PushCallerIf(err error, skip int) bool {
	if err == nil {
		return false
	}

	err = &wrapped{err, Caller(skip + 1)}
	defer e.lock(true)()
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

	if len(e.errs) == 0 {
		return nil
	}
	return e
}

// Error implements the error interface
func (e *ErrorList) Error() string {
	defer e.lock(false)()
	return fmtErrors(e.errs, AlwaysWithCaller)
}

// Unwrap implements errors.Unwrap by returning the next error in the
// chain or nil if there are no more errors.
func (e *ErrorList) Unwrap() error {
	defer e.lock(false)()
	if len(e.errs) == 1 {
		return nil
	}

	mux := e.mux
	if mux != nil {
		mux = new(sync.RWMutex)
	}
	return &ErrorList{errs: e.errs[1:], mux: mux}
}

// As implements errors.As by checking if target is ErrorList, otherwise will go through each error.
func (e *ErrorList) As(target interface{}) bool {
	defer e.lock(false)()
	if oe, ok := target.(*ErrorList); ok {
		defer oe.lock(true)()
		oe.errs = append([]error(nil), e.errs...)
		return true
	}

	for _, err := range e.errs {
		if errors.As(err, target) {
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
			if !Is(err, oe.errs[i]) {
				return false
			}
		}

		return true
	}

	for _, err := range e.errs {
		if Is(err, target) {
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

func (e *ErrorList) lock(rw bool) func() {
	if e == nil || e.mux == nil {
		return func() {}
	}
	mux := e.mux
	if rw {
		mux.Lock()
		return mux.Unlock
	}
	mux.RLock()
	return mux.RUnlock
}

func fmtErrors(es []error, withStack bool) string {
	if len(es) == 0 {
		return ""
	}

	verb := "%v"
	if withStack {
		verb = "%+v"
	}

	if len(es) == 1 {
		return fmt.Sprintf(verb, es[0])
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "%d errors occurred:", len(es))

	for _, err := range es {
		buf.WriteString("\n\t* ")
		fmt.Fprintf(&buf, verb, err)
	}

	return buf.String()
}
