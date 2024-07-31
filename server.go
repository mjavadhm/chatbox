package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "sync"
)

var (
    clients   = make(map[net.Conn]bool)
    broadcast = make(chan string)
    mutex     = &sync.Mutex{}
)

func main() {
    ln, err := net.Listen("tcp", ":8087")
    if err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }
    defer ln.Close()

    go handleMessages()

    fmt.Println("Server started on :8087")
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        mutex.Lock()
        clients[conn] = true
        mutex.Unlock()

        go handleConnection(conn)
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        mutex.Lock()
        for conn := range clients {
            _, err := fmt.Fprintln(conn, msg)
            if err != nil {
                fmt.Println("Error sending message:", err)
                conn.Close()
                delete(clients, conn)
            }
        }
        mutex.Unlock()
    }
}

func handleConnection(conn net.Conn) {
    defer func() {
        mutex.Lock()
        delete(clients, conn)
        mutex.Unlock()
        conn.Close()
    }()

    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Client disconnected:", conn.RemoteAddr().String())
            break
        }
        broadcast <- msg
    }
}
