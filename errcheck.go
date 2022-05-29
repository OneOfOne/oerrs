package oerrs

type tryErr struct {
	*wrappedError
}

var CatchAll func(err error, fr *Frame)

func Catch(err *error, handler func(err error, fr *Frame)) {
	switch pv := recover().(type) {
	case tryErr:
		if err != nil {
			*err = Errorf("%s: %w", pv.fr.String(), pv.err)
		}
		if handler != nil {
			handler(pv.err, pv.fr)
		}
		if CatchAll != nil {
			CatchAll(pv.err, pv.fr)
		}
	case nil:
	default:
		panic(pv)
	}
}

func Try(err error) {
	checkErr(err)
}

func Try1[T1 any](v1 T1, err error) T1 {
	checkErr(err)
	return v1
}

func Try2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	checkErr(err)
	return v1, v2
}

func Try3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	checkErr(err)
	return v1, v2, v3
}

func Try4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	checkErr(err)
	return v1, v2, v3, v4
}

func Try5[T1, T2, T3, T4, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, err error) (T1, T2, T3, T4, T5) {
	checkErr(err)
	return v1, v2, v3, v4, v5
}

func checkErr(err error) {
	if err != nil {
		c := Caller(3)
		c.trim = true
		panic(tryErr{&wrappedError{err, c}})
	}
}
