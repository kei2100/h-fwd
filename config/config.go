package config

import (
	"net/http"

	"github.com/kei2100/fwxy/errors"
)

// Parameters is the configuration parameters
type Parameters struct {
	Destination
	Header http.Header
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
