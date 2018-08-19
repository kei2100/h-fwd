package config

import (
	"encoding/base64"
	"net/http"
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
