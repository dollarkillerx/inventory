package main

import (
	"github.com/dollarkillerx/inventory/internal/server"

	"log"
)

func main() {
	server := server.NewServer()
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
