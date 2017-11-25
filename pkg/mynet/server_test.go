package portforward

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

type mockHandler struct {
	c chan bool
	t *testing.T
}

func (h mockHandler) handleConnection(l net.Listener) {
	// blocks until there is an inbound connection we can accept
	h.t.Log("Waiting to process a new inbound connection", l)
	l.Accept()
	h.c <- true
}

func TestLoopListensOnPortAndHandlesConnection(t *testing.T) {
	processed := make(chan bool, 3)
	h := mockHandler{processed, t}

	if len(processed) != 0 {
		t.Errorf("should not have called handler yet. called %d times", len(processed))
	}

	s := Server{make(chan bool, 1)}
	go s.Loop(9268, h)
	makeSuccessfulConnection(t, 9268)

	if len(processed) != 1 {
		t.Errorf("called handler %d times. expected 1", len(processed))
	}
}

func TestLoopHandlesMultipleConnections(t *testing.T) {
	processed := make(chan bool, 3)
	h := mockHandler{processed, t}

	s := Server{make(chan bool, 1)}
	go s.Loop(8989, h)

	for i := 0; i < 5; i++ {
		makeSuccessfulConnection(t, 8989)
	}
}

func makeSuccessfulConnection(t *testing.T, port int) {
	// our mockHandler implementation could cause Get() to hang waiting for a response
	client := http.Client{
		Timeout: time.Duration(500 * time.Millisecond),
	}
	t.Log("Attempting to connect to the server")
	for i := 0; i < 10; i++ {
		_, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
		t.Logf("response err for %d-th request was %v", i, err)
		if err == nil {
			return
		}
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("Unable to make a successful connection")
}
