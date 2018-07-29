package config

import "testing"

func TestDestination(t *testing.T) {
	c := &Destination{
		To:       "https://example.com/",
		Username: "user",
		Password: "pass",
		RewritePaths: map[string]string{
			"/user":    "/users",
			"/company": "/companies",
		},
	}

	if err := c.load(); err != nil {
		t.Fatalf("failed to load Destination config: %v", err)
	}
	if g, w := c.Host(), "example.com"; g != w {
		t.Errorf("Host got %v, want %v", g, w)
	}
	if g, w := c.Scheme(), "https"; g != w {
		t.Errorf("Scheme got %v, want %v", g, w)
	}
	if g, w := c.UserInfo().String(), "user:pass"; g != w {
		t.Errorf("UserInfo string got %v, want %v", g, w)
	}
	if g, w := len(c.PathRewriters()), 2; g != w {
		t.Errorf("len(c.PathRewriters()) got %v, want %v", g, w)
	}
}
