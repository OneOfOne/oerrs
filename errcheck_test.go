package oerrs

import (
	"fmt"
	"testing"
)

func maybeerror(err bool) error {
	if err {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

func doErrorChecks() (err error) {
	defer Catch(&err, nil)

	Try(maybeerror(true))
	return
}

func TestCatcher(t *testing.T) {
	t.Log(doErrorChecks())
	defer Catch(nil, func(err error, fr *Frame) { t.Log(err, fr) })

	Try(doErrorChecks())
}
