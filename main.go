package main

import (
	"log"
)

func main() {

	server := NewTcpServer(":3000")
	log.Fatal(server.Start())

}
