package oerrs

import (
	"errors"
	"fmt"
	"io"
	"testing"
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
	t.Log(errors.As(&errs, &trg))
	t.Log(trg)
}

func TestErrIs(t *testing.T) {
	ErrNoCredentials := String("no credentials")
	err := fmt.Errorf("%w: %s", ErrNoCredentials, errors.New("some error"))
	if errors.Is(err, errors.New("other error")) {
		t.Error("should not be equal")
	}

	if !errors.Is(err, ErrNoCredentials) {
		t.Error("should be equal")
	}

	err = Wrap(err, 1)
	t.Logf("%s", err)
	t.Logf("%+v", err)
	t.Logf("%#+v", err)
}
