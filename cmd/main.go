package main

import (
	"github.com/dollarkillerx/inventory/internal/server"
	"github.com/dollarkillerx/inventory/internal/utils"

	"log"
)

func main() {
	utils.InitJWT()

	server := server.NewServer()
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
