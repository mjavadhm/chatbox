package main

import (
	"strings"
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "217.196.107.82:8087")
	if err != nil {
        fmt.Println("Error connecting to server:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Connected to server")

    go readMessages(conn)

    reader := bufio.NewReader(os.Stdin)
    for {
        //fmt.Print("Enter message: ")
        msg, _ := reader.ReadString('\n')
        if msg != "" {
            fmt.Fprint(conn, msg)
        }
    }
}

func readMessages(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading message:", err)
            break
        }
        msg = strings.TrimSpace(msg)
        if msg != "" {
            fmt.Println("\nMessage from server:", msg)
        }
    }
}
