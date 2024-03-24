package main

import (
	"log"
	"net"

	"github.com/achintya-7/ez-chat/server"
)

func main() {
	s := server.NewServer()
	go s.Run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		c := s.NewClient(conn)
		go c.ReadInput()
	}
}
