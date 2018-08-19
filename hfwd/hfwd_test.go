package hfwd

import (
	"fmt"
	"net/url"
	"testing"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/kei2100/h-fwd/config"
)

var dstMux = http.NewServeMux()

func init() {
	dstMux.Handle("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Echo", r.Header.Get("X-Echo"))
		w.WriteHeader(200)
		w.Write([]byte("foo body"))
	}))
	dstMux.Handle("/dumpHeaders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(r.Header)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}))
}

func assertOKResponse(t *testing.T, res *http.Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("response returns err: %v", err)
	}
	if g, w := res.StatusCode, 200; g != w {
		t.Errorf("res.StatusCode got %v, want %v", g, w)
	}
}

func withRunProxy(dstURL string, params *config.Parameters, test func(proxyURL string)) {
	h, err := NewHandler(mustURL(dstURL), params)
	if err != nil {
		panic(fmt.Sprintf("hfwd: failed to create hfwd Handler for a test: %v", err))
	}
	proxyServer := httptest.NewServer(h)
	defer proxyServer.Close()
	test(proxyServer.URL)
}

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

	if err := c.Setup(); err != nil {
		panic(fmt.Sprintf("failed setup configuration: %v", err))
	}
	return &c
}

func TestServer_ServeHTTP(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		dstServer := httptest.NewServer(dstMux)
		defer dstServer.Close()
		withRunProxy(dstServer.URL, &config.Parameters{}, func(proxyURL string) {
			res, err := http.Get(proxyURL + "/foo")
			assertOKResponse(t, res, err)

			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			if g, w := string(b), "foo body"; g != w {
				t.Errorf("res.Body got %v, want %v", g, w)
			}
		})
	})

	t.Run("copy and rewriteHeader", func(t *testing.T) {
		dstServer := httptest.NewServer(dstMux)
		defer dstServer.Close()
		rwHeader := make(http.Header)
		rwHeader.Set("User-Agent", "my agent")
		rwHeader.Set("Host", "http://example.com:8888")
		params := configParam(config.Headers{Header: rwHeader})

		withRunProxy(dstServer.URL, params, func(proxyURL string) {
			req, err := http.NewRequest("GET", proxyURL+"/dumpHeaders", nil)
			if err != nil {
				t.Fatalf("failed to create a new request: %v", err)
			}
			req.Header.Set("X-Echo", "echo")
			res, err := http.DefaultClient.Do(req)
			assertOKResponse(t, res, err)

			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			dumpHeaders := make(http.Header)
			json.Unmarshal(b, &dumpHeaders)

			if g, w := dumpHeaders.Get("User-Agent"), "my agent"; g != w {
				t.Errorf("res.Header[User-Agent] got %v, want %v", g, w)
			}
			if g, w := dumpHeaders.Get("X-Echo"), "echo"; g != w {
				t.Errorf("res.Header[X-Echo] got %v, want %v", g, w)
			}
		})
	})

	t.Run("rewriteURL", func(t *testing.T) {
		dstServer := httptest.NewServer(dstMux)
		defer dstServer.Close()
		params := configParam(config.URL{RewritePaths: map[string]string{"/bar": "/foo"}})

		withRunProxy(dstServer.URL, params, func(proxyURL string) {
			res, err := http.Get(proxyURL + "/bar")
			assertOKResponse(t, res, err)

			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			if g, w := string(b), "foo body"; g != w {
				t.Errorf("res.Body got %v, want %v", g, w)
			}
		})
	})
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
