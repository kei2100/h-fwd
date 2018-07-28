package rewrite

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
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
		return false
	}

	b := new(strings.Builder)
	b.WriteString(u.Scheme)
	b.WriteString("://")
	if ui := u.User.String(); len(ui) > 0 {
		b.WriteString(ui)
		b.WriteString("@")
	}
	b.WriteString(u.Host)
	b.WriteString(replaced)
	if len(u.RawQuery) > 0 {
		b.WriteString("?")
		b.WriteString(u.RawQuery)
	}
	if len(u.Fragment) > 0 {
		b.WriteString("#")
		b.WriteString(u.Fragment)
	}

	newu, err := url.Parse(b.String())
	if err != nil {
		log.Printf("rewrite: failed to regexp path rewrite: %v", err)
		return false
	}
	*u = *newu
	return true
}
