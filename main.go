package main

import (
	"fmt"
	"net"
	"os"

	netcat "netcat/Fonctions"
)

var (
	con  net.Conn
	port string
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		port = ":8989"
		
		netcat.StartClient(port, os.Args[1])

	} else if len(os.Args) == 2 {
		port = os.Args[1]
		netcat.StartClient(port, os.Args[1])
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port $server")
		os.Exit(1)
	}
}
