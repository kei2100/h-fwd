package config

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"log"

	"github.com/kei2100/h-fwd/errors"
	"golang.org/x/crypto/pkcs12"
)

// TLSClient is configuration parameters for the tls client certification
type TLSClient struct {
	CACertPath string

	PKCS12Path     string
	PKCS12Password string

	caCertPEM []byte
	certPEM   []byte
	keyPEM    []byte

	tlsConfig *tls.Config
}

// TLSClientConfig returns *tls.Config for the tls client certification
func (t *TLSClient) TLSClientConfig() *tls.Config {
	return t.tlsConfig
}

// load configuration given parameters
func (t *TLSClient) load() error {
	d := new(errors.DoOrSkip)
	d.DoOrSkip(t.loadCACert)
	d.DoOrSkip(t.loadPKCS12)
	d.DoOrSkip(t.loadTLSConfig)
	return d.Err()
}

// loadCACert normally called from the load() method
func (t *TLSClient) loadCACert() error {
	if t.CACertPath == "" {
		return nil
	}
	b, err := ioutil.ReadFile(t.CACertPath)
	if err != nil {
		return fmt.Errorf("config: failed to load ca cert file %v : %v", t.CACertPath, err)
	}
	t.caCertPEM = b
	return nil
}

// loadPKCS12 normally called from the load() method
func (t *TLSClient) loadPKCS12() error {
	if t.PKCS12Path == "" {
		return nil
	}
	b, err := ioutil.ReadFile(t.PKCS12Path)
	if err != nil {
		return fmt.Errorf("config: failed to load pkcs12 file %v : %v", t.PKCS12Path, err)
	}
	key, cert, err := pkcs12.Decode(b, t.PKCS12Password)
	if err != nil {
		return fmt.Errorf("config: failed to decode pkcs12 data: %v", err)
	}
	kp, err := encodePrivateKeyPEMToMemory(key)
	if err != nil {
		return err
	}
	t.keyPEM = kp
	t.certPEM = encodeCertPEMToMemory(cert)
	return nil
}

// loadTLSConfig normally called from the load() method
func (t *TLSClient) loadTLSConfig() error {
	cfg := tls.Config{}

	if certPEM := t.certPEM; certPEM != nil {
		cert, err := tls.X509KeyPair(certPEM, t.keyPEM)
		if err != nil {
			return fmt.Errorf("config: failed to create x509 keypair: %v", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}
	if certPEM := t.caCertPEM; certPEM != nil {
		p, err := x509.SystemCertPool()
		if err != nil {
			log.Printf("config: system cert pool is not available. creates a new cert pool: %v", err)
			p = x509.NewCertPool()
		}
		if ok := p.AppendCertsFromPEM(certPEM); !ok {
			return fmt.Errorf("config: failed to append ca cert file %v (%v bytes)", t.CACertPath, len(certPEM))
		}
		cfg.RootCAs = p
	}

	t.tlsConfig = &cfg
	return nil
}

func encodePrivateKeyPEMToMemory(key interface{}) ([]byte, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		kb := x509.MarshalPKCS1PrivateKey(k)
		pemb := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: kb})
		return pemb, nil
	case *ecdsa.PrivateKey:
		kb, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("config: failed to marshal ecdsa private key: %v", err)
		}
		pemb := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: kb})
		return pemb, nil
	default:
		return nil, fmt.Errorf("config: unknown private key type %T", key)
	}
}

func encodeCertPEMToMemory(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE", Bytes: cert.Raw,
	})
}
