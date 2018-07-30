package errors

// DoOrSkip executes given functions in order.
// If an error occurrs during execution, it skips
type DoOrSkip struct {
	err error
}

// Err returns the error or nil.
func (d *DoOrSkip) Err() error {
	return d.err
}

// DoOrSkip executes given functions in order.
// If an error occurrs during execution, it skips
func (d *DoOrSkip) DoOrSkip(do func() error) {
	if d.err != nil {
		return
	}
	d.err = do()
}
