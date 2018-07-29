package errors

import (
	"fmt"
	"testing"
)

func TestMultiLine_Error(t *testing.T) {
	e := NewMultiLine(fmt.Errorf("foo"))
	want := "foo"
	if g, w := e.Error(), want; g != w {
		t.Errorf("got %v, want %v", g, w)
	}

	e = NewMultiLine(fmt.Errorf("foo"))
	e.Add(fmt.Errorf("bar"))

	want = fmt.Sprintln("foo")
	want += "bar"

	if g, w := e.Error(), want; g != w {
		t.Errorf("got %v, want %v", g, w)
	}
}
