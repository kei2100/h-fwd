package config

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestHeaders(t *testing.T) {
	header := http.Header{}
	header.Set("x-test", "test")
	h := Headers{
		Header:   header,
		Username: "user",
		Password: "pass",
	}
	if err := h.setup(); err != nil {
		t.Fatalf("failed to setup :%v", err)
	}

	tt := []struct {
		key  string
		want string
	}{
		{key: "X-Test", want: "test"},
		{key: "Authorization", want: "Basic dXNlcjpwYXNz"}, // dXNlcjpwYXNz = user:pass
	}
	for _, te := range tt {
		got := strings.Join(h.Header[te.key], "\t")
		if g, w := got, te.want; g != w {
			t.Errorf("%v got %v, want %v", te.key, g, w)
		}
	}
}

func TestHeaders_String(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer foo")
	h := &Headers{
		Header:   header,
		Username: "user",
		Password: "pass",
	}
	got := fmt.Sprintf("%v", h)
	want := `Username: user
Password: ****
Header: Authorization: **********
`

	if g, w := got, want; g != w {
		t.Errorf("String() got %v, want %v", g, w)
	}
}
