package oerrs

import (
	"encoding/json"
	"io"
	"testing"

	"golang.org/x/xerrors"
)

func TestList(t *testing.T) {
	var errs ErrorList
	errs.saveCaller = true
	errs.PushIf(io.EOF)
	errs.PushIf(io.ErrClosedPipe)
	errs.PushIf(io.ErrClosedPipe)
	err := errs.Err()
	t.Logf("%v", err)
	var trg *ErrorList
	t.Log(xerrors.As(&errs, &trg))
	trg.saveCaller = true
	t.Log(trg)
	trg.SupportJSON(true, false)
	j, _ := json.MarshalIndent(trg, "", "\t")
	t.Logf("%s", j)
	trg.SupportJSON(true, true)
	j, _ = json.MarshalIndent(trg, "", "\t")
	t.Logf("%s", j)
}
