package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

func main() {
	listener, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatal("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Server started. Listening on :8083")

	go handleMessages()

	listenErr := make(chan error)

	go AcceptConn(listener, listenErr)

	for err := range listenErr {
		log.Println("Error", err)
	}
}

func AcceptConn(listener net.Listener, errChan chan<- error) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
			continue
		}

		reader := bufio.NewReader(conn)
		nickname, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			continue
		}

		client := Client{
			conn:     conn,
			nickname: nickname,
		}
		register <- conn

		go handleConnection(client, errChan)
	}

}

func handleConnection(client Client, errChan chan<- error) {
	defer client.conn.Close()

	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			continue
		}
		nick := strings.TrimSuffix(client.nickname, "\r\n")
		log.Printf("Received message from %s: %s", nick, message)

		broadcast <- Message{client, message}
	}
}

func handleMessages() {
	for {
		select {
		case message := <-broadcast:
			for conn := range clients {
				if conn != message.client.conn {
					nick := strings.TrimSuffix(message.client.nickname, "\r\n")
					conn.Write([]byte(fmt.Sprintf("%s: %s", nick, message.message)))
				}
			}
		case conn := <-register:
			clients[conn] = Client{}
		case conn := <-unregister:
			delete(clients, conn)
		}
	}
}
