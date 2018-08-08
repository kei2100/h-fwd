package fwxy

import (
	"errors"
	"net/http"
	"net/url"

	"path"

	"io"
	"log"

	"github.com/kei2100/fwxy/config"
)

// NewHandler returns http.Handler which performs forward proxy.
// dst represents destination base URL. the format must be "http[s]://[user:pass@]host[:port][/base/path]"
// params are configuration parameters of the forward proxy.
func NewHandler(dst *url.URL, params *config.Parameters) (http.Handler, error) {
	if dst == nil {
		return nil, errors.New("fwxy: failed to create new handler. dst must be set")
	}
	// TODO dst validation
	forwarder := &http.Client{
		Transport: &http.Transport{TLSClientConfig: params.TLSClientConfig()},
	}
	return &server{dst: dst, params: params, forwarder: forwarder}, nil
}

type server struct {
	dst       *url.URL
	params    *config.Parameters
	forwarder *http.Client
}

func (s *server) ServeHTTP(w http.ResponseWriter, orig *http.Request) {
	req, err := http.NewRequest(orig.Method, orig.URL.String(), orig.Body)
	if err != nil {
		log.Printf("fwxy: failed to create a new request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.copyHeader(orig, req)
	s.rewriteHeader(req)
	s.rewriteURL(req.URL)

	res, err := s.forwarder.Do(req)
	if err != nil {
		log.Printf("fwxy: an error occurrd while forwarding the request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for h, vv := range res.Header {
		for _, v := range vv {
			w.Header().Add(h, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func (s *server) copyHeader(orig, req *http.Request) {
	req.Header = make(http.Header)
	for k, v := range orig.Header {
		if _, ok := hopByHopHeaders[k]; ok {
			continue
		}
		req.Header[k] = v
	}
}

func (s *server) rewriteHeader(req *http.Request) {
	for k := range s.params.Header {
		if k == "Host" {
			req.Host = s.params.Header.Get(k)
			continue
		}
		req.Header[k] = s.params.Header[k]
	}
}

func (s *server) rewriteURL(reqURL *url.URL) {
	for _, rewrite := range s.params.PathRewriters() {
		if ok := rewrite.Do(reqURL); ok {
			break
		}
	}

	dstURL := *s.dst
	if dstURL.User == nil {
		dstURL.User = reqURL.User
	}
	dstURL.Path = path.Join(dstURL.Path, reqURL.Path)
	dstURL.ForceQuery = reqURL.ForceQuery
	dstURL.RawQuery = reqURL.RawQuery
	dstURL.Fragment = reqURL.Fragment

	*reqURL = dstURL
}

// Hop-by-hop headers, which are meaningful only for a single
// transport-level connection, and are not stored by caches or
// forwarded by proxies.
//
// https://tools.ietf.org/html/rfc2616#section-13.5.1
var hopByHopHeaders = map[string]struct{}{
	// Header names are canonicalized (see http.Request or http.Response).
	"Connection":          {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"TE":                {},
	"Trailers":          {},
	"Transfer-Encoding": {},
	"Upgrade":           {},
}
