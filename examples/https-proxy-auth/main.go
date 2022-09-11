/*
 * Snitch is a SNI proxy library.
 * Copyright (C) 2022 İrem Kuyucu <siren@kernal.eu>

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// This example showcases how to use Snitch for authenticating HTTPS proxy users
// through domain name. The HTTPS proxy server is given a wildcard certificate 
// and is bound to localhost. DestAddr field is set to the HTTPS proxy address.

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
