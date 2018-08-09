package config

import (
	"fmt"

	"github.com/kei2100/h-fwd/rewrite"
)

// URL is configuration parameters for the destination
type URL struct {
	RewritePaths  map[string]string // map[oldPath]newPath
	pathRewriters []rewrite.PathRewriter
}

// PathRewriters returns path rewriters
func (u *URL) PathRewriters() []rewrite.PathRewriter {
	return u.pathRewriters
}

// load configuration given parameters
func (u *URL) load() error {
	if u == nil {
		return nil
	}

	for old, new := range u.RewritePaths {
		rwr, err := rewrite.NewRewriter(old, new)
		if err != nil {
			return fmt.Errorf("config: failed to interpret the rewrite string %v to %v", old, new)
		}
		u.pathRewriters = append(u.pathRewriters, rwr)
	}
	return nil
}
