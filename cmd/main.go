package main

import (
	"github.com/dollarkillerx/inventory/internal/server"
	"github.com/dollarkillerx/inventory/internal/utils"

	"log"
)

func main() {
	utils.InitJWT()

	ser := server.NewServer()
	if err := ser.Run(); err != nil {
		log.Fatalln(err)
	}
}
