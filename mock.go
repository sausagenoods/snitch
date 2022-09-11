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

import (
	"io"
	"net"
	"time"
)

// Mock io.Reader to satisfy the net.Conn interface
// https://www.agwa.name/blog/post/writing_an_sni_proxy_in_go
type mockConn struct {
	reader io.Reader
}

func (conn mockConn) Read(p []byte) (int, error) { return conn.reader.Read(p) }
func (conn mockConn) Write(p []byte) (int, error) { return 0, nil }
func (conn mockConn) Close() error { return nil }
func (conn mockConn) LocalAddr() net.Addr { return nil }
func (conn mockConn) RemoteAddr() net.Addr { return nil }
func (conn mockConn) SetDeadline(t time.Time) error { return nil }
func (conn mockConn) SetReadDeadline(t time.Time) error { return nil }
func (conn mockConn) SetWriteDeadline(t time.Time) error { return nil }
