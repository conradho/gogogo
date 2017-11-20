package portforward

import (
	"fmt"
	"log"
	"net"
)

type connectionHandler interface {
	handleConnection(net.Listener)
}

// Loop starts listening on port and calls ConnectionHandler when a client connects in
func Loop(port int, h connectionHandler) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	checkError(err, fmt.Sprintf("could not start server on %d: %v", port, err))

	for {
		h.handleConnection(l)
	}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
