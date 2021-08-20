package oerrs

import (
	"io"
	"testing"

	"golang.org/x/xerrors"
)

func TestList(t *testing.T) {
	var errs ErrorList
	errs.PushIf(io.EOF)
	err := errs.Err()
	t.Logf("%v", err)
	errs.PushIf(io.ErrClosedPipe)
	errs.PushIf(io.ErrClosedPipe)
	err = errs.Err()
	t.Logf("%v", err)
	var trg *ErrorList
	t.Log(xerrors.As(&errs, &trg))
	t.Log(trg)
}
