package oerrs

import (
	"fmt"
	"log"
	"testing"
)

func maybeerror(err bool) error {
	log.Printf("%+v", New("x"))
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
