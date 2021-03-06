// Copyright 2018 gopcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Command server provides a connection establishment of OPC UA Secure Conversation as a server.

XXX - Currently this command just handles the UACP connection from any client.
*/
package main

import (
	"context"
	"flag"
	"log"
	"math/rand"

	"github.com/wmnsk/gopcua/services"

	"github.com/wmnsk/gopcua/uacp"
	"github.com/wmnsk/gopcua/uasc"
	"github.com/wmnsk/gopcua/utils"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://example.com/foo/bar", "OPC UA Endpoint URL")
		bufsize  = flag.Int("bufsize", 0xffff, "Receive Buffer Size")
	)
	flag.Parse()

	listener, err := uacp.Listen(*endpoint, uint32(*bufsize))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started listening on %s.", listener.Endpoint())

	cfg := uasc.NewConfig(
		1, "http://opcfoundation.org/UA/SecurityPolicy#None", nil, nil, 1, rand.Uint32(),
	)

	for {
		func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			conn, err := listener.Accept(ctx)
			if err != nil {
				log.Print(err)
				return
			}
			defer func() {
				conn.Close()
				log.Printf("Successfully closed connection with %v", conn.RemoteAddr())
			}()
			log.Printf("Successfully established connection with %v", conn.RemoteAddr())

			secChan, err := uasc.ListenAndAcceptSecureChannel(ctx, conn, cfg)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				secChan.Close()
				log.Printf("Successfully closed secure channel with %v", conn.RemoteAddr())
			}()
			log.Printf("Successfully opened secure channel with %v", conn.RemoteAddr())

			buf := make([]byte, 1024)
			/*
				n, err := secChan.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Successfully received message: %x\n%s", buf[:n], utils.Wireshark(0, buf[:n]))

				sc, err := uasc.Decode(buf[:n])
				if err != nil {
					log.Println("Couldn't decode received bytes as UASC")
					return
				}
				log.Printf("Successfully decoded as UASC: %v", sc)
			*/

			n, err := secChan.ReadService(buf)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Successfully received message: %x\n%s", buf[:n], utils.Wireshark(0, buf[:n]))

			srv, err := services.Decode(buf[:n])
			if err != nil {
				log.Println("Couldn't decode received bytes as Service")
				return
			}
			log.Printf("Successfully decoded as Service: %v", srv)
		}()
	}
}
