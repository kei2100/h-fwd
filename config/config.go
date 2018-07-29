package config

import "github.com/kei2100/fwxy/errors"

// Parameters is the configuration parameters
type Parameters struct {
	Destination
}

// Load configuration given parameters
func (p *Parameters) Load() error {
	errs := errors.NewMultiLine()
	errs.AddIfErr(p.Destination.load())
	if errs.Len() > 0 {
		return errs
	}
	return nil
}
