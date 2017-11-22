package portforward

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// ConnectionForwarder handles all connections by forwarding it to the target.
// It implements the connectionHandler interface
type ConnectionForwarder struct {
	Target          string
	InboundLogPath  string
	OutboundLogPath string
}

func (forwarder ConnectionForwarder) handleConnection(l net.Listener) {
	log.Printf("About to forward traffic to %s", forwarder.Target)
	inboundConn, err := l.Accept()
	check(err, "could not accept client connection")
	log.Printf("A new client '%v' connected!\n", inboundConn.RemoteAddr())

	outboundConn, err := net.Dial("tcp", forwarder.Target)
	check(err, "could not connect to target")
	log.Printf("outbound connection to server %v established!\n", outboundConn.RemoteAddr())

	go func() {
		// send all inbound traffic to outbound connection and also tee to log file
		teedTraffic := io.TeeReader(inboundConn, outboundConn)
		appendStreamToFile(teedTraffic, forwarder.InboundLogPath)
		err := inboundConn.Close()
		check(err, fmt.Sprint("Could not close inbound connection ", inboundConn.RemoteAddr()))
	}()
	go func() {
		// send all outbound traffic to inbound connection and also tee to log file
		teedTraffic := io.TeeReader(outboundConn, inboundConn)
		appendStreamToFile(teedTraffic, forwarder.OutboundLogPath)
		err := outboundConn.Close()
		check(err, fmt.Sprint("Could not close outbound connection ", outboundConn.RemoteAddr()))
	}()
}

func appendStreamToFile(r io.Reader, s string) {
	f, err := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, fmt.Sprint("Could not open file ", s))
	w := bufio.NewWriter(f)
	_, err = io.Copy(w, r)
	check(err, fmt.Sprint(s, " errored before getting EOL"))
	err = w.Flush()
	check(err, fmt.Sprint("Could not flush file ", s))
	err = f.Close()
	check(err, fmt.Sprint("Could not close file ", s))
}
