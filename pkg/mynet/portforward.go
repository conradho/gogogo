package portforward

import (
	"fmt"
	"io"
	"log"
	"net"
)

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func Forward(port int, target string) {
	log.Printf("About to forward traffic from port %d to %s", port, target)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	CheckError(err, fmt.Sprintf("could not start server on %d: %v", port, err))

	for {
		handleConnection(l, target)
	}
}

func handleConnection(l net.Listener, target string) {
	inboundConn, err := l.Accept()
	CheckError(err, "could not accept client connection")
	log.Printf("A new client '%v' connected!\n", inboundConn.RemoteAddr())

	outboundConn, err := net.Dial("tcp", target)
	CheckError(err, "could not connect to target")
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
