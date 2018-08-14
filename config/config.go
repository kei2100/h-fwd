package config

import (
	"github.com/kei2100/h-fwd/errors"
)

// Parameters is the configuration parameters
type Parameters struct {
	URL
	Headers
	TLSClient
	Verbose bool
}

// Load configuration given parameters
func (p *Parameters) Load() error {
	errs := errors.NewMultiLine()
	errs.AddIfErr(p.URL.load())
	errs.AddIfErr(p.Headers.load())
	errs.AddIfErr(p.TLSClient.load())
	if errs.Len() > 0 {
		return errs
	}
	return nil
}
