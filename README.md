# Snitch
[![GoDoc](https://godoc.org/gitlab.com/sausagenoods/snitch?status.svg)](https://godoc.org/gitlab.com/sausagenoods/snitch)

Snitch is a SNI proxy library for Golang that allows for conditional routing through user supplied allow function.

## Examples
### HTTPS Proxy Authentication
Proxies request to the HTTPS proxy server only if the server name contains the correct credentials. In this example only the requests to `s3cret.proxy.digilol.net` are proxied. The HTTPS proxy server is given the wildcard certificate for `*.proxy.digilol.net`, and is bound to `127.0.0.1:4443`.

This is useful when your browser/client does not support HTTP(S) proxy authentication.

```go
package main

import (
	"log"
	"strings"

	"gitlab.com/sausagenoods/snitch"
)

var proxyDomain = ".proxy.digilol.net"

func main() {
	s := &snitch.SniProxy{
		BindAddr: "0.0.0.0:443",
		DestAddr: "127.0.0.1:4443",
		AllowFunc: customAuth,
	}
	log.Fatal(s.Serve())
}

func customAuth(serverName string) bool {
	if !strings.HasSuffix(serverName, proxyDomain) {
		log.Println("Blocking connection to unauthorized backend")
		return false
	}
	creds := strings.ReplaceAll(serverName, proxyDomain, "")
	if creds == "s3cret" {
		return true
	}
	log.Println("Wrong credentials supplied:", creds)
	return false
}
```
