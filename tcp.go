package main

/**
 * @Author: hsm
 * @Email: hsmcool@163.com
 * @Date: 2023/5/16 15:09
 * @Desc:
 */

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"tunvitor/socks"
)

// Create a SOCKS server listening on addr and proxy to server.
func socksLocal(addr, server string) {
	//logf("listen socks proxy %s <-> %s", addr, server)
	log.Printf("listen tcp socks proxy %s <-> %s", addr, server)

	tcpLocal(addr, server, func(c net.Conn) (socks.Addr, error) { return socks.Handshake(c) })
}

// Listen on addr and proxy to server to reach target from getAddr.
func tcpLocal(addr, server string, getAddr func(net.Conn) (socks.Addr, error)) {
	//log.Printf("listen tcp proxy %s <-> %s:", addr, server)
	l, err := net.Listen("tcp", addr)

	if err != nil {
		logf("failed to listen on %s: %v", addr, err)
		return
	}

	for {

		c, err := l.Accept()
		log.Println("accept address: ", c.LocalAddr())
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}

		go func() {
			defer c.Close()
			tgt, err := getAddr(c)
			if err != nil {

				// UDP: keep the connection until disconnect then free the UDP socket
				if err == socks.InfoUDPAssociate {
					buf := make([]byte, 1)
					// block here
					for {
						_, err := c.Read(buf)
						if err, ok := err.(net.Error); ok && err.Timeout() {
							continue
						}
						log.Printf("UDP Associate End.")
						return
					}
				}

				log.Printf("failed to get target address: %v", err)
				return
			}

			rc, err := net.Dial("tcp", server)
			if err != nil {
				log.Printf("failed to connect to server %v: %v", server, err)
				return
			}
			defer rc.Close()

			//rc = shadow(rc)

			if _, err = rc.Write(tgt); err != nil {
				logf("failed to send target address: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, tgt)
			if err = relay(rc, c); err != nil {
				logf("relay error: %v", err)
			}
		}()
	}
}

// Listen on addr for incoming connections.
func tcpRemote(addr string) {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen on %s: %v", addr, err)
		return
	}

	log.Printf("listening TCP on %s", addr)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept: %v", err)
			continue
		}

		go func() {
			defer c.Close()

			tgt, err := socks.ReadAddr(c)
			if err != nil {
				// drain c to avoid leaking server behavioral features
				// see https://www.ndss-symposium.org/ndss-paper/detecting-probe-resistant-proxies/
				_, err = io.Copy(ioutil.Discard, c)
				if err != nil {
					log.Printf("discard error: %v", err)
				}
				return
			}

			rc, err := net.Dial("tcp", tgt.String())
			if err != nil {
				log.Printf("failed to connect to target: %v", err)
				return
			}
			defer rc.Close()
			//log.Printf("proxy %s <-> %s", c.RemoteAddr(), tgt)

			if err = relay(c, rc); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
