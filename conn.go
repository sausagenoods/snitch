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
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func readClientHello(reader io.Reader) (*tls.ClientHelloInfo, error) {
	var hello *tls.ClientHelloInfo
	err := tls.Server(mockConn{reader: reader}, &tls.Config{
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			hello = new(tls.ClientHelloInfo)
			*hello = *argHello
			return nil, nil
		},
	}).Handshake()
	if hello == nil {
		return nil, err
	}
	return hello, nil
}

func peekClientHello(reader io.Reader) (*tls.ClientHelloInfo, io.Reader, error) {
	peekedBytes := new(bytes.Buffer)
	hello, err := readClientHello(io.TeeReader(reader, peekedBytes))
	if err != nil {
		return nil, nil, err
	}
	return hello, io.MultiReader(peekedBytes, reader), nil
}

func (s *SniProxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	if err := clientConn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		log.Print(err)
		return
	}
	clientHello, clientReader, err := peekClientHello(clientConn)
	if err != nil {
	log.Print(err)
		return
	}
	if err := clientConn.SetReadDeadline(time.Time{}); err != nil {
		log.Println(err)
		return
	}
	if ok := s.AllowFunc(clientHello.ServerName); !ok {
		return
	}
	backendConn, err := net.DialTimeout("tcp", s.DestAddr, 60 * time.Second)
	if err != nil {
		log.Print(err)
		return
	}
	defer backendConn.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(clientConn, backendConn)
		clientConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()
	go func() {
		io.Copy(backendConn, clientReader)
		backendConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()
	wg.Wait()
}
