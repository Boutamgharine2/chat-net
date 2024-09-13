package netcat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	
)

func StartClient(port string, server string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server, port)) // crier un connection tcp
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	go receiveMessages(conn)

	for {
		fmt.Print("Enter message: ")
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		if len(message) > 1 {
			fmt.Fprintln(conn, message)
		}
	}
}

func receiveMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
