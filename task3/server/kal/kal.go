package kal

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn     net.Conn
	nickname string
}

type Message struct {
	client  Client
	message string
}

var (
	clients    = make(map[net.Conn]Client)
	broadcast  = make(chan Message)
	register   = make(chan net.Conn)
	unregister = make(chan net.Conn)
)

func handleConnection(client Client) {
	defer client.conn.Close()

	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		fmt.Println("Received message from", client.nickname+":", message)

		broadcast <- Message{client, message}
	}
}

func handleMessages() {
	for {
		select {
		case message := <-broadcast:
			for conn := range clients {
				if conn != message.client.conn {
					conn.Write([]byte(message.client.nickname + ": " + message.message))
				} else {
					conn.Write([]byte("Message received\n"))
				}
			}
		case conn := <-register:
			clients[conn] = Client{}
		case conn := <-unregister:
			delete(clients, conn)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on :8080")

	go handleMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		reader := bufio.NewReader(conn)
		nickname, _ := reader.ReadString('\n')

		client := Client{
			conn:     conn,
			nickname: nickname,
		}
		register <- conn

		go handleConnection(client)
	}
}
