package portforward

import (
	"fmt"
	"io"
	"log"
	"net"
)

func Forward(port int, target string) {
	log.Printf("About to forward traffic from port %d to %s", port, target)

	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("could not start server on %d: %v", port, err)
	}

	for {
		client, err := incoming.Accept()
		if err != nil {
			log.Fatal("could not accept client connection", err)
		}
		defer client.Close()
		log.Printf("client '%v' connected!\n", client.RemoteAddr())

		forwardConnection, err := net.Dial("tcp", target)
		if err != nil {
			log.Fatal("could not connect to target", err)
		}
		defer forwardConnection.Close()
		log.Printf("connection to server %v established!\n", forwardConnection.RemoteAddr())
		go func() { io.Copy(forwardConnection, client) }()
		go func() { io.Copy(client, forwardConnection) }()
	}
}
