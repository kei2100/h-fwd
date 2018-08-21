package hfwd

import (
	"errors"
	"net/http"
	"net/url"

	"path"

	"io"
	"log"

	"net/http/httputil"

	"fmt"

	"bytes"
	"compress/flate"
	"compress/gzip"
	"io/ioutil"

	"github.com/dsnet/compress/brotli"
	"github.com/kei2100/h-fwd/config"
)

// NewHandler returns http.Handler which performs forward proxy.
func NewHandler(dst *url.URL, params *config.Parameters) (http.Handler, error) {
	if err := validateDestinatin(dst); err != nil {
		return nil, err
	}
	var tran http.RoundTripper
	tran = &http.Transport{TLSClientConfig: params.TLSClientConfig()}
	if params.Verbose {
		log.Printf("hfwd destination is %v", dst.String())
		log.Printf("hfwd configuration parameters are\n%s", params)
		tran = &verboseRoundTripper{
			chain: &http.Transport{TLSClientConfig: params.TLSClientConfig()},
		}
	}
	forwarder := &http.Client{
		Transport: tran,
	}
	return &server{dst: dst, params: params, forwarder: forwarder}, nil
}

func validateDestinatin(dst *url.URL) error {
	switch {
	case dst == nil:
	case dst.Scheme != "http" && dst.Scheme != "https":
	case len(dst.Opaque) > 0:
	case len(dst.RawQuery) > 0:
	case len(dst.Fragment) > 0:
		return errors.New("hfwd: destination URL format must be 'http[s]://[user:pass@]host[:port][/base/path]'")
	}
	return nil
}

type server struct {
	dst       *url.URL
	params    *config.Parameters
	forwarder *http.Client
}

func (s *server) ServeHTTP(w http.ResponseWriter, orig *http.Request) {
	req, err := http.NewRequest(orig.Method, orig.URL.String(), orig.Body)
	if err != nil {
		log.Printf("hfwd: failed to create a new request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.copyHeader(orig, req)
	s.rewriteHeader(req)
	s.rewriteURL(req.URL)

	res, err := s.forwarder.Do(req)
	if err != nil {
		log.Printf("hfwd: an error occurrd while forwarding the request: %v", err)
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

type verboseRoundTripper struct {
	chain http.RoundTripper
}

func (rt *verboseRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	var reqDump, resDump []byte
	defer func() {
		b := bytes.Buffer{}
		if len(reqDump) > 0 {
			b.WriteString("\n>>> hfwd send a request\n")
			b.Write(reqDump)
		}
		if len(resDump) > 0 {
			b.WriteString("<<< hfwd receive a response\n")
			b.Write(resDump)
		}
		log.Print(b.String())
	}()

	// dump request
	reqDump, err = httputil.DumpRequest(req, true)
	if err != nil {
		return nil, fmt.Errorf("hfwd: failed to dump the request: %v", err)
	}

	res, err = rt.chain.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// dump response
	resDump, err = httputil.DumpResponse(res, false)
	if err != nil {
		return nil, fmt.Errorf("hfwd: failed to dump the response header: %v", err)
	}
	// TODO configurable
	if false {
		resDumpBody := make([]byte, 0)
		body := res.Body
		cl := res.ContentLength
		defer body.Close()

		raw, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("hfwd: failed to read the response body: %v", err)
		}
		res.Body = ioutil.NopCloser(bytes.NewReader(raw))
		res.ContentLength = cl

		rd := ioutil.NopCloser(bytes.NewReader(raw))
		switch res.Header.Get("Content-Encoding") {
		case "gzip":
			if rd, err = gzip.NewReader(rd); err != nil {
				return nil, fmt.Errorf("hfwd: failed to create gzip reader for read the response body: %v", err)
			}
		case "deflate":
			rd = flate.NewReader(rd)
		case "br":
			if rd, err = brotli.NewReader(rd, nil); err != nil {
				return nil, fmt.Errorf("hfwd: failed to create brotli reader for read the response body: %v", err)
			}
		}

		resDumpBody, err = ioutil.ReadAll(rd)
		if err != nil {
			return nil, fmt.Errorf("hfwd: failed to dump the response body: %v", err)
		}
		if len(resDumpBody) > 0 {
			resDump = bytes.Join([][]byte{resDump, resDumpBody}, []byte("\r\n"))
		}
	}

	return res, nil
}
