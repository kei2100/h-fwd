package errors

import (
	"fmt"
	"testing"
)

func TestDoOrSkip(t *testing.T) {
	door := new(DoOrSkip)

	i := 0

	err1 := fmt.Errorf("err1")
	err2 := fmt.Errorf("err2")

	fn1 := func() error {
		i++
		return nil
	}
	fn2 := func() error {
		i += 2
		return err1
	}
	fn3 := func() error {
		i += 4
		return err2
	}

	door.DoOrSkip(fn1)
	door.DoOrSkip(fn2)
	door.DoOrSkip(fn3)

	if g, w := i, 3; g != w {
		t.Errorf("'i' got %v, want %v", g, w)
	}
	if g, w := door.Err(), err1; g != w {
		t.Errorf("door.Err() got %v, want %v", g, w)
	}

}
