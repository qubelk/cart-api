package main

import (
	"cart-api/server"
	"log"
)

func main() {
	handlers := server.NewHTTPHandler()
	server := server.NewServer(handlers)

	log.Println("Server started on localhost:3000/carts")
	server.StartServer()
}
