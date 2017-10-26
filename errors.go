package main

import "fmt"

type (
	testError struct {
		code int
		desc string
	}
)

func newTestError(code int, desc string) error {
	return &testError{code, desc}
}

func (e *testError) Error() string {
	return fmt.Sprintf("%s (%d)", e.desc, e.code)
}
