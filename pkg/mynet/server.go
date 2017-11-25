package portforward

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type connectionHandler interface {
	handleConnection(net.Listener)
}

// Server with a Loop() and a quit channel
type Server struct {
	Quit chan bool
}

// Loop starts listening on port and calls ConnectionHandler when a client connects in
func (s Server) Loop(port int, h connectionHandler) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	go func() {
		<-s.Quit
		log.Println("received termination signal from go channel")
		if err = l.Close(); err != nil {
			log.Fatal("could not close listener", err)
		}
		log.Println("terminated")
		s.Quit <- true
	}()

	check(err, fmt.Sprintf("could not start server on %d: %v", port, err))

	for {
		h.handleConnection(l)
	}
}

func check(err error, msg string) {
	if err != nil {
		log.Fatal(strings.Repeat("*", 80), "\n", msg, "\nError: ", err)
	}
}
