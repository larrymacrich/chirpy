package main

import "github.com/larrymacrich/chirpy/internal/server"

func main() {
	port := "8080"
	server.StartServer(port)
}
