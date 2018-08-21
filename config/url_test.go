package config

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	c := &URL{
		RewritePaths: map[string]string{
			"/user":    "/users",
			"/company": "/companies",
		},
	}

	if err := c.setup(); err != nil {
		t.Fatalf("failed to setup URL config: %v", err)
	}
	if g, w := len(c.PathRewriters()), 2; g != w {
		t.Errorf("len(c.PathRewriters()) got %v, want %v", g, w)
	}
}

func TestURL_String(t *testing.T) {
	c := &URL{
		RewritePaths: map[string]string{
			"/user": "/users",
		},
	}
	got := fmt.Sprintf("%v", c)
	want := "RewritePath: /user: /users\n"
	if g, w := got, want; g != w {
		t.Errorf("String() got %v, want %v", g, w)
	}

}

func mustURL(u string) *url.URL {
	uu, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return uu
}

func TestRegexpPathRewriter_Do(t *testing.T) {
	tt := []struct {
		re   string
		repl string
		url  *url.URL
		want *url.URL
	}{
		{
			re:   "^/bar/",
			repl: "/foo/",
			url:  mustURL("https://example.com/bar/bar"),
			want: mustURL("https://example.com/foo/bar"),
		},
		{
			re:   "^/bar/",
			repl: "/foo/",
			url:  mustURL("https://user:pass@example.com/bar/bar?x=y#here"),
			want: mustURL("https://user:pass@example.com/foo/bar?x=y#here"),
		},
		{
			re:   "/いいい/",
			repl: "/あああ/",
			url:  mustURL("https://example.com/あああ/いいい/あああ?x=ワイ#ここ"),
			want: mustURL("https://example.com/あああ/あああ/あああ?x=ワイ#ここ"),
		},
	}

	for _, te := range tt {
		r, err := newRegexpPathRewriter(te.re, te.repl)
		if err != nil {
			t.Errorf("failed to create rewriter. re: %v, repl: %v, url: %v, msg: %v", te.re, te.repl, te.url.String(), err)
			continue
		}

		replaced := *te.url
		ret := r.Do(&replaced)

		if replaced.String() != te.want.String() {
			t.Errorf("got %v, want %v. re: %v, repl: %v", replaced.String(), te.want.String(), te.re, te.repl)
		}
		if (replaced.String() == te.url.String()) == ret {
			t.Errorf("ret %v, got %v, before %v", ret, replaced.String(), te.url.String())
		}
	}
}
