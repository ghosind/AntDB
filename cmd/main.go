package main

import (
	"log"
	"os"

	"github.com/ghosind/antdb/config"
	"github.com/ghosind/antdb/server"
)

func main() {
	cfg, err := config.ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	options := config.BuildOptionsByConfig(cfg)
	s := server.NewServer(options...)

	err = s.Listen()
	if err != nil {
		log.Fatalf("Failed to start AntDB: %v", err)
	}
}
