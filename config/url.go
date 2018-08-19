package config

import (
	"fmt"

	"net/url"
	"regexp"
)

// URL is configuration parameters for the destination
type URL struct {
	RewritePaths  map[string]string // map[oldPath]newPath
	pathRewriters []PathRewriter
}

// PathRewriters returns path rewriters
func (u *URL) PathRewriters() []PathRewriter {
	return u.pathRewriters
}

// load configuration given parameters
func (u *URL) load() error {
	if u == nil {
		return nil
	}

	for old, new := range u.RewritePaths {
		rwr, err := newRewriter(old, new)
		if err != nil {
			return err
		}
		u.pathRewriters = append(u.pathRewriters, rwr)
	}
	return nil
}

// PathRewriter is an interface to path rewrite
type PathRewriter interface {
	// Do rewrites the URL
	Do(*url.URL) (rewrited bool)
}

// newRewriter creates a PathRewriter
func newRewriter(old, new string) (PathRewriter, error) {
	return newRegexpPathRewriter(old, new)
}

// regexpPathRewriter is an implementation of the PathRewriter using regexp
type regexpPathRewriter struct {
	rex  *regexp.Regexp
	repl string
}

func newRegexpPathRewriter(re, repl string) (*regexpPathRewriter, error) {
	rex, err := regexp.Compile(re)
	if err != nil {
		return nil, fmt.Errorf("config: failed to compile regexp for path rewriter %v: %v", re, err)
	}
	return &regexpPathRewriter{rex: rex, repl: repl}, nil
}

func (r *regexpPathRewriter) Do(u *url.URL) bool {
	var orig string
	if len(u.RawPath) != 0 {
		orig = u.RawPath
	} else {
		orig = u.Path
	}
	replaced := r.rex.ReplaceAllString(orig, r.repl)
	if replaced == orig {
		return false
	}
	u.Path = replaced
	return true
}
