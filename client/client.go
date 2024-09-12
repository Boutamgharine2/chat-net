package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
)
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ./client <server-ip>:<port>")
        return
    }

    serverAddress := os.Args[1]
    conn, err := net.Dial("tcp", serverAddress)
    if err != nil {
        log.Fatalf("Erreur de connexion au serveur: %v", err)
    }
    defer conn.Close()

    go func() {
        reader := bufio.NewReader(conn)
        for {
            msg, err := reader.ReadString('\n')
            if err != nil {
                log.Printf("Erreur de lecture: %v", err)
                return
            }
            fmt.Print(msg)
        }
    }()

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        msg := scanner.Text()
        if msg == "" {
            continue
        }
        _, err := fmt.Fprintf(conn, "%s\n", msg)
        if err != nil {
            log.Printf("Erreur d'envoi: %v", err)
            return
        }
    }
}
