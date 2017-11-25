package main

import (
	"flag"

	"github.com/conradho/gogogo/pkg/mynet"
)

func main() {
	var (
		portPtr   = flag.Int("port", 8000, "port to listen to")
		targetPtr = flag.String("target", "localhost:1234", "target to forward to")
	)
	flag.Parse()

	f := portforward.ConnectionForwarder{*targetPtr, "inbound.log", "outbound.log"}

	s := portforward.Server{make(chan bool, 1)}
	s.Loop(*portPtr, f)

}
