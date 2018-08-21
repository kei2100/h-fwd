package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTLSClient(t *testing.T) {
	tlsp := TLSClient{
		CACertPath:     "testdata/cacert.pem",
		PKCS12Path:     "testdata/clicert.pfx",
		PKCS12Password: "pass",
	}
	if err := tlsp.setup(); err != nil {
		t.Fatal(err)
	}

	clientCfg := tlsp.TLSClientConfig()
	clientCfg.InsecureSkipVerify = true

	servCert, err := tls.LoadX509KeyPair("testdata/servcert.pem", "testdata/servkey-nopass.pem")
	if err != nil {
		t.Fatalf("failed to load servcert: %v", err)
	}
	servCfg := &tls.Config{
		Certificates: []tls.Certificate{servCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCfg.RootCAs,
	}

	ok := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }
	serv := httptest.NewUnstartedServer(http.HandlerFunc(ok))
	serv.TLS = servCfg
	serv.StartTLS()
	defer serv.Close()

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: clientCfg}}
	resp, err := client.Get(serv.URL)
	if err != nil {
		t.Fatalf("failed GET request: %v", err)
	}
	defer resp.Body.Close()
	if g, w := resp.StatusCode, 200; g != w {
		t.Errorf("resp.StatusCode got %v, want %v", g, w)
	}
}

func TestTLSClient_Format(t *testing.T) {
	tlsp := &TLSClient{
		CACertPath:     "testdata/cacert.pem",
		PKCS12Path:     "testdata/clicert.pfx",
		PKCS12Password: "pass",
	}
	got := fmt.Sprintf("%s", tlsp)
	want := `CACertPath: testdata/cacert.pem
PKCS12Path: testdata/clicert.pfx
PKCS12Password: ****
`
	if g, w := got, want; g != w {
		t.Errorf("Strings() got %v, want %v", g, w)
	}
}
