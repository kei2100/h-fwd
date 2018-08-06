package config

import "testing"

func TestDestination(t *testing.T) {
	c := &URL{
		RewritePaths: map[string]string{
			"/user":    "/users",
			"/company": "/companies",
		},
	}

	if err := c.load(); err != nil {
		t.Fatalf("failed to load URL config: %v", err)
	}
	if g, w := len(c.PathRewriters()), 2; g != w {
		t.Errorf("len(c.PathRewriters()) got %v, want %v", g, w)
	}
}
