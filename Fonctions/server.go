package netcat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	clientsMutex sync.Mutex

	port       string
	clients    = make(map[net.Conn]string)
	messageLog []string
	logMutex   sync.Mutex
)

const maxCliens = 10

func Connection(con net.Conn) {
	defer con.Close()

	scaner := bufio.NewScanner(con)
	scaner.Scan()
	name := scaner.Text()
	if name == "" {
		name = "Anonyme"
	}
	clientsMutex.Lock()
	clients[con] = name
	clientsMutex.Unlock()
	broadcast(fmt.Sprintf("[%s][%s] a rejoint le chat", time.Now().Format(time.RFC3339), name)) // Inform other clients about the new connection

	sendHistory(con) // Send message history to the new client

	// Handle incoming messages
	for scaner.Scan() {
		message := scaner.Text()
		if strings.TrimSpace(message) != "" {
			logMutex.Lock()
			messageLog = append(messageLog, fmt.Sprintf("[%s][%s]: %s", time.Now().Format(time.RFC3339), name, message))
			logMutex.Unlock()
			broadcast(fmt.Sprintf("[%s][%s]: %s", time.Now().Format(time.RFC3339), name, message))
		}
	}
	// Handle client disconnect
	clientsMutex.Lock()
	delete(clients, con)
	clientsMutex.Unlock()
	broadcast(fmt.Sprintf("[%s][%s] has left the chat", time.Now().Format(time.RFC3339), name))
}

func broadcast(message string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for conn := range clients {
		fmt.Fprintln(conn, message)
	}
}

func sendHistory(conn net.Conn) {
	logMutex.Lock()
	defer logMutex.Unlock()

	for _, msg := range messageLog {
		fmt.Fprintln(conn, msg)
	}
}
