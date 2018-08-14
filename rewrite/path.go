package rewrite

import (
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/kei2100/h-fwd/env"
)

// PathRewriter is an interface to path rewrite
type PathRewriter interface {
	// Do rewrites the URL
	Do(*url.URL) (rewrited bool)
}

// NewRewriter creates a PathRewriter
func NewRewriter(old, new string) (PathRewriter, error) {
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
		return nil, fmt.Errorf("rewrite: failed to compile regexp %v: %v", re, err)
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
		if env.Verbose {
			log.Printf("rewrite: path not rewrite\nexpr: %v\nrequested: %v", r.rex.String(), orig)
		}
		return false
	}
	u.Path = replaced
	if env.Verbose {
		log.Printf("rewrite: path rewrited\nexpr: %v\nrequested: %v\nrewrited: %v", r.rex.String(), orig, replaced)
	}
	return true
}
