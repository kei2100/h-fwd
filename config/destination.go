package config

import (
	"net/url"

	"fmt"

	"github.com/kei2100/fwxy/rewrite"
)

// Destination is configuration parameters for the destination
type Destination struct {
	To           string            // Destination proto://host[:port]
	Username     string            // Username or blank. for basic authN
	Password     string            // Password or blank. for basic authN
	RewritePaths map[string]string // map[oldPath]newPath

	host          string // host[:port]
	scheme        string // http or https
	userInfo      *url.Userinfo
	pathRewriters []rewrite.PathRewriter
}

// Host returns the host or host:port
func (dp *Destination) Host() string {
	return dp.host
}

// Scheme returns the scheme
func (dp *Destination) Scheme() string {
	return dp.scheme
}

// UserInfo returns the *url.UserInfo or nil
func (dp *Destination) UserInfo() *url.Userinfo {
	return dp.userInfo
}

// PathRewriters returns path rewriters
func (dp *Destination) PathRewriters() []rewrite.PathRewriter {
	return dp.pathRewriters
}

// load configuration given parameters
func (dp *Destination) load() error {
	if dp == nil {
		return nil
	}

	u, err := url.Parse(dp.To)
	if err != nil {
		return fmt.Errorf("config: failed to parse Destination.To %v  to URL: %v", dp.To, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("config: invalid scheme %v", u.Scheme)
	}
	dp.host = u.Host
	dp.scheme = u.Scheme

	if dp.Username != "" {
		dp.userInfo = url.UserPassword(dp.Username, dp.Password)
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
