# h-fwd
Simple HTTP request forwarder written in Golang

## Usage
```
$ hfwd https://example.com
2018/08/14 11:00:00 hfwd listening on 127.0.0.1:8080
```

Sending a request to http://127.0.0.1:8080, will be forwarded to https://example.com 

## Downloads
[releases](../../releases)

## Features
Regexp based path rewriting
```
$ hfwd https://example.com --rewrite='^/status_(.+)$:/status/$1'

# http://127.0.0.1:8080/status_200 => https://example.com/status/200
```

Additional or overwrite request headers
```
$ ./bin/hfwd https://example.com --header="User-Agent: MyAgent"

# User-Agent:  MyAgent
```

Basic authentication
```
$ hfwd https://example.com --username=user --password=pass

# Authorization: Basic dXNlcjpwYXNz
```

Add SSL/TLS certificate to Root CAs
```
$ hfwd https://example.com --ca-cert=/path/to/cert
```

SSL/TLS client authentication
```
$ hfwd https://example.com --pkcs12=/path/to/pkcs12 --pkcs12-password=pass
```

More info
```
$ ./bin/hfwd -h
hfwd is a simple HTTP forward proxy

Usage:
  hfwd <destination URL> [flags]

Flags:
      --ca-cert string           path of the additional CA certificate PEM
  -H, --header strings           list for the additional http headers (-H Host:https://custom.example.com -H 'User-Agent:My Agent'
  -h, --help                     help for hfwd
  -l, --listen string            listen addr:port (default "127.0.0.1:8080")
  -p, --password string          password for the basic authentication
      --pkcs12 string            path of the PKCS12 encoded file for the client certification
      --pkcs12-password string   password for the PKCS12 file
  -r, --rewrite strings          list for path rewrite (-r /old:/new -r /o:/n OR -r /old:/new,/o:/n)
  -u, --username string          username for the basic authentication
      --verbose                  verbose output
```
