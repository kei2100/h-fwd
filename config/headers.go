package config

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// Headers is configuration parameters for the http header
type Headers struct {
	Header   http.Header
	Username string // Username or blank. for basic authN
	Password string // Password or blank. for basic authN
}

// setup configuration given parameters
func (h *Headers) setup() error {
	if len(h.Username) > 0 {
		src := []byte(h.Username + ":" + h.Password)
		dst := base64.StdEncoding.EncodeToString(src)
		//if h.Header == nil {
		//	h.Header = make(http.Header, 0)
		//}
		h.Header.Set("Authorization", "Basic "+string(dst))
	}
	return nil
}

// String returns string representation of this configuration. useful for debugging.
func (h *Headers) String() string {
	b := strings.Builder{}
	if h == nil {
		return b.String()
	}
	b.WriteString(fmt.Sprintf("Username: %s\n", h.Username))
	b.WriteString(fmt.Sprintf("Password: %s\n", strings.Repeat("*", len(h.Password))))
	for k := range h.Header {
		v := h.Header.Get(k)
		if k == "Authorization" {
			v = strings.Repeat("*", len(v))
		}
		b.WriteString(fmt.Sprintf("Header: %s: %s\n", k, v))
	}
	return b.String()
}
