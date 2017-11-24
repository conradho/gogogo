package portforward

import (
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
	// block until there is an inbound connection we can accept
	h.t.Log("waiting to accept inbound connection", l)
	l.Accept()
	h.c <- true
}

func TestLoopListensOnPortAndHandlesConnection(t *testing.T) {
	c := make(chan bool, 5)
	h := mockHandler{c, t}

	if len(c) != 0 {
		t.Errorf("should not have called handler yet. called %d times", len(c))
	}

	go Loop(9268, h)
	waitForServerLoopToStartUp(t)

	if len(c) != 1 {
		t.Errorf("called handler %d times. expected 1", len(c))
	}
}

func waitForServerLoopToStartUp(t *testing.T) {
	// our mockHandler implementation could cause Get() to hang waiting for a response
	client := http.Client{
		Timeout: time.Duration(1 * time.Second),
	}
	for {
		t.Log("about to hit")
		time.Sleep(100 * time.Millisecond)
		_, err := client.Get("http://127.0.0.1:9268/")
		t.Log("err was", err)
		if opError, ok := err.(*net.OpError); !(ok && opError.Op == "read") {
			// error is a connection refused error because server is not listening yet
			break
		}
	}
}
