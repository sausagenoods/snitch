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

package snitch

import "net"

type SniProxy struct {
	// Bind address:port
	BindAddr string
	// HTTPS destination address:port
	DestAddr string
	// Function for conditional proxying
	AllowFunc func(serverName string) bool
}

// Starts SNI proxy server
func (s *SniProxy) Serve() error {
	l, err := net.Listen("tcp", s.BindAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}
