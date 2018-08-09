package cli

import (
	"fmt"
	"net/http"
	"strings"

	"log"

	"github.com/kei2100/h-fwd/config"
	"github.com/spf13/cobra"
)

var (
	// option parameters for the url configuration
	rewritePaths []string
)

var (
	// option parameters for the headers configuration
	headers  []string
	username string
	password string
)

var (
	// options parameters for the client certification
	caCertPath     string
	pkcs12Path     string
	pkcs12Password string
)

func init() {
	flags := RootCmd.PersistentFlags()

	flags.StringVarP(&username, "username", "u", "", "username for the basic authentication")
	flags.StringVarP(&password, "password", "p", "", "password for the basic authentication")
	flags.StringSliceVarP(&rewritePaths, "rewrite", "r", []string{}, "list for path rewrite (-r /old:/new -r /o:/n OR -r /old:/new,/o:/n)")

	flags.StringSliceVarP(&headers, "header", "H", []string{}, "list for the additional http headers (-H Host:https://example.com -H 'User-Agent:My Agent'")

	flags.StringVar(&caCertPath, "ca-cert", "", "path of the PEM encoded CA certificate")
	flags.StringVar(&pkcs12Path, "pkcs12", "", "path of the PKCS12 encoded file for the client certification")
	flags.StringVar(&pkcs12Password, "pkcs12-password", "", "password for the PKCS12 file")
}

// RootCmd for CLI
var RootCmd = &cobra.Command{
	Use:   "hfwd",
	Short: "hfwd is a simple HTTP forward proxy",
	Run: func(cmd *cobra.Command, args []string) {
		param := config.Parameters{}
		param.RewritePaths = parseRewritePaths(rewritePaths)

		param.Header = parseHeaders(headers)
		param.Username = username
		param.Password = password

		param.CACertPath = caCertPath
		param.PKCS12Path = pkcs12Path
		param.PKCS12Password = pkcs12Password

		if err := param.Load(); err != nil {
			log.Fatalf("failed to load configuration: %v", err)
		}
		fmt.Printf("%+v", param)
	},
}

func parseRewritePaths(rewritePaths []string) map[string]string {
	m := make(map[string]string, len(rewritePaths))
	for _, p := range rewritePaths {
		sp := strings.SplitAfterN(p, ":", 2)
		m[sp[0]] = sp[1]
	}
	return m
}

func parseHeaders(headers []string) http.Header {
	hh := make(http.Header, len(headers))
	for _, h := range headers {
		sp := strings.SplitAfterN(h, ":", 2)
		hh.Add(sp[0], sp[1])
	}
	return hh
}
