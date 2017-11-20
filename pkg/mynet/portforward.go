package portforward

import (
	"fmt"
	"io"
	"log"
	"net"
)

func Forward(port int, target string) {
	log.Printf("About to forward traffic from port %d to %s", port, target)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("could not start server on %d: %v", port, err)
	}

	for {
		ForwardIncomingConnection(l, target)
	}
}

func ForwardIncomingConnection(l net.Listener, target string) {
	inboundConn, err := l.Accept()
	if err != nil {
		log.Fatal("could not accept client connection", err)
	}
	log.Printf("client '%v' connected!\n", inboundConn.RemoteAddr())

	outboundConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Fatal("could not connect to target", err)
	}
	log.Printf("outbound connection to server %v established!\n", outboundConn.RemoteAddr())

	go func() {
		io.Copy(inboundConn, outboundConn)
		inboundConn.Close()
	}()
	go func() {
		io.Copy(outboundConn, inboundConn)
		outboundConn.Close()
	}()
}
