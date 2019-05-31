package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Server struct {
	Router *mux.Router
	Port   int
}

func (s *Server) Run() {
	// Initializing the server
	addr := fmt.Sprintf(":%d", s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("Error initializing the server : %s", err)
	}

	done := make(chan bool)
	// Start the server
	go http.Serve(listener, s.Router)

	logrus.Infof("Successfully started the Milobella::Oratio server on port %d !", s.Port)
	<-done
}
