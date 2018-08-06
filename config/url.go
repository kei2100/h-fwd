package config

import (
	"fmt"

	"github.com/kei2100/fwxy/rewrite"
)

// URL is configuration parameters for the destination
type URL struct {
	RewritePaths  map[string]string // map[oldPath]newPath
	pathRewriters []rewrite.PathRewriter
}

// PathRewriters returns path rewriters
func (dp *URL) PathRewriters() []rewrite.PathRewriter {
	return dp.pathRewriters
}

// load configuration given parameters
func (dp *URL) load() error {
	if dp == nil {
		return nil
	}

	for old, new := range dp.RewritePaths {
		rwr, err := rewrite.NewRewriter(old, new)
		if err != nil {
			return fmt.Errorf("config: failed to interpret the rewrite string %v to %v", old, new)
		}
		dp.pathRewriters = append(dp.pathRewriters, rwr)
	}
	return nil
}
