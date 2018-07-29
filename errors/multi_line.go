package errors

import (
	"fmt"
	"strings"
)

// NewMultiLine creates new MultiLine error by given error(s)
func NewMultiLine(errs ...error) MultiLine {
	e := new(multiLine)
	for _, err := range errs {
		e.Add(err)
	}
	return e
}

// MultiLine represents multiple errors, prints line by line
type MultiLine interface {
	error
	// Add an error.
	// If the error is nil, Add panics.
	Add(error)

	// Len returns the count of errors
	Len() int
}

type multiLine struct {
	errs []error
}

func (e *multiLine) Add(err error) {
	if err == nil {
		panic("errors: err is nil")
	}
	e.errs = append(e.errs, err)
}

func (e *multiLine) Len() int {
	return len(e.errs)
}

func (e *multiLine) Error() string {
	b := new(strings.Builder)
	for _, e := range e.errs {
		if b.Len() > 0 {
			fmt.Fprintln(b, "")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}
