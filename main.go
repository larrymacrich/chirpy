package main

import (
	"log"

	"github.com/larrymacrich/chirpy/internal/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	port := "8080"
	server.StartServer(port)
}
