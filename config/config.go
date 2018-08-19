package config

import (
	"github.com/kei2100/h-fwd/errors"
)

// Parameters is the configuration parameters for the hfwd proxy server
type Parameters struct {
	URL
	Headers
	TLSClient
	Verbose bool
}

// Setup configuration given parameters
func (p *Parameters) Setup() error {
	errs := errors.NewMultiLine()
	errs.AddIfErr(p.URL.setup())
	errs.AddIfErr(p.Headers.setup())
	errs.AddIfErr(p.TLSClient.setup())
	if errs.Len() > 0 {
		return errs
	}
	return nil
}
