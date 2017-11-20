package portforward

import (
	"io"
	"log"
	"net"
)

// ConnectionForwarder handles all connections by forwarding it to the target.
// It implements the connectionHandler interface
type ConnectionForwarder struct {
	Target string
}

func (f ConnectionForwarder) handleConnection(l net.Listener) {
	log.Printf("About to forward traffic to %s", f.Target)
	inboundConn, err := l.Accept()
	checkError(err, "could not accept client connection")
	log.Printf("A new client '%v' connected!\n", inboundConn.RemoteAddr())

	outboundConn, err := net.Dial("tcp", f.Target)
	checkError(err, "could not connect to target")
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
