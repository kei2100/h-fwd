package config

import (
	"net/url"

	"fmt"

	"github.com/kei2100/fwxy/rewrite"
)

// Destination configuration
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
func (c *Destination) Host() string {
	return c.host
}

// Scheme returns the scheme
func (c *Destination) Scheme() string {
	return c.scheme
}

// UserInfo returns the *url.UserInfo or nil
func (c *Destination) UserInfo() *url.Userinfo {
	return c.userInfo
}

// PathRewriters returns path rewriters
func (c *Destination) PathRewriters() []rewrite.PathRewriter {
	return c.pathRewriters
}

// load url infos given configuration
func (c *Destination) load() error {
	if c == nil {
		return nil
	}

	u, err := url.Parse(c.To)
	if err != nil {
		return fmt.Errorf("config: failed to parse Destination.To %v  to URL: %v", c.To, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("config: invalid scheme %v", u.Scheme)
	}
	c.host = u.Host
	c.scheme = u.Scheme

	if c.Username != "" {
		c.userInfo = url.UserPassword(c.Username, c.Password)
	}

	for old, new := range c.RewritePaths {
		rwr, err := rewrite.NewRewriter(old, new)
		if err != nil {
			return fmt.Errorf("config: failed to interpret the rewrite string %v to %v", old, new)
		}
		c.pathRewriters = append(c.pathRewriters, rwr)
	}
	return nil
}
