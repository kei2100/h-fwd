package hfwd

import (
	"fmt"
	"net/url"
	"testing"

	"net"
	"net/http"

	"github.com/kei2100/h-fwd/config"
)

func mustURL(u string) *url.URL {
	p, err := url.Parse(u)
	if err != nil {
		panic(fmt.Sprintf("failed to parse URL: %v", err))
	}
	return p
}

func configParam(subconfig ...interface{}) *config.Parameters {
	c := config.Parameters{}

	for _, sc := range subconfig {
		switch sc := sc.(type) {
		case config.URL:
			c.URL = sc
		case config.Headers:
			c.Headers = sc
		case config.TLSClient:
			c.TLSClient = sc
		}
	}

	if err := c.Load(); err != nil {
		panic(fmt.Sprintf("failed load configuration: %v", err))
	}
	return &c
}

func TestServer_rewriteURL(t *testing.T) {
	t.Run("rewrite path", func(t *testing.T) {
		tt := []struct {
			orig   *url.URL
			dst    *url.URL
			params *config.Parameters
			want   *url.URL
		}{
			{
				orig: mustURL("http://localhost:18000/foo/path?q=qv#frag"),
				dst:  mustURL("https://www.example.com/base"),
				params: configParam(config.URL{
					RewritePaths: map[string]string{"/foo/": "/bar/"},
				}),
				want: mustURL("https://www.example.com/base/bar/path?q=qv#frag"),
			},
			{
				orig: mustURL("http://ou:op@localhost:18000/foo"),
				dst:  mustURL("https://www.example.com"),
				want: mustURL("https://ou:op@www.example.com/foo"),
			},
			{
				orig: mustURL("http://ou:op@localhost:18000/foo"),
				dst:  mustURL("https://u:p@www.example.com"),
				want: mustURL("https://u:p@www.example.com/foo"),
			},
		}

		for _, te := range tt {
			s := &server{dst: te.dst, params: te.params}
			if s.params == nil {
				s.params = configParam()
			}
			s.rewriteURL(te.orig)
			if g, w := te.orig.String(), te.want.String(); g != w {
				t.Errorf("url got %v, want %v", g, w)
			}
		}
	})
}

// TODO delete
func TestCoreLogic(t *testing.T) {
	t.SkipNow()
	params := configParam()
	s := &server{dst: mustURL("https://www.google.com"), params: params, forwarder: http.DefaultClient}
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	http.Serve(ln, s)
}
