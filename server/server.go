package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	defaultPort    = "8989"
	maxConnections = 10
)

var (
	clients  = make(map[net.Conn]string)
	mu       sync.Mutex
	messages []string
)

func main() {
	port := defaultPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
	defer listener.Close()

	messagesChan := make(chan string)

	go func() {
		for msg := range messagesChan {
			broadcast(msg)
		}
	}()

	log.Printf("Serveur à l'écoute sur le port %s", port)
	for {
		if len(clients) >= maxConnections {
			log.Println("Nombre maximal de connexions atteint, nouvelle connexion ignorée.")
			continue
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erreur lors de l'acceptation de la connexion: %v", err)
			continue
		}

		log.Print("Demande du nom du client")
		_, err = conn.Write([]byte("Veuillez entrer votre nom : "))
		if err != nil {
			log.Printf("Erreur lors de l'écriture au client: %v", err)
			conn.Close()
			continue
		}

		name, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Erreur lors de la lecture du nom du client: %v", err)
			conn.Close()
			continue
		}
		name = name[:len(name)-1] // Supprimer le caractère de nouvelle ligne

		mu.Lock()
		clients[conn] = name
		mu.Unlock()

		broadcast(fmt.Sprintf("Client %s a rejoint le chat", name))

		go handleConnection(conn, name, messagesChan)
	}
}
