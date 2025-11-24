package main

import (
	"log"

	"github.com/ghosind/antdb/server"
)

func main() {
	s := server.NewServer()

	err := s.Listen()
	if err != nil {
		log.Fatalf("Failed to start AntDB: %v", err)
	}
}
