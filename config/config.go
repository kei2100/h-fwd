package config

import (
	"fmt"

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

// String returns string representation of this configuration. useful for debugging.
func (p *Parameters) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s%s%s", p.URL.String(), p.Headers.String(), p.TLSClient.String())
}
