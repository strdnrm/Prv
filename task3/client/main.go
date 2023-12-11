package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	ConnRetries := 3
	ConnInterval := 3 * time.Second

	var conn net.Conn
	var err error

	for i := 0; i < ConnRetries; i++ {
		if i == ConnRetries {
			log.Println("out of conn retreies")
			return
		}
		conn, err = net.Dial("tcp", "localhost:8083")
		if err != nil {
			log.Println("Error connecting:", err)
			time.Sleep(ConnInterval)
			continue
		}
		break
	}

	defer conn.Close()

	listenErr := make(chan error)

	go WriteMsg(conn, listenErr)
	go ReadMsg(conn, listenErr)

	for err := range listenErr {
		log.Println("Error", err)
	}

}

func WriteMsg(conn net.Conn, errChan chan<- error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your nickname: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		errChan <- err
	}
	fmt.Println("Connected to server. Start chatting!")
	_, err = fmt.Fprint(conn, nickname)
	if err != nil {
		errChan <- err
	}

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			continue
		}
		fmt.Fprint(conn, text)
	}
}

func ReadMsg(conn net.Conn, errChan chan<- error) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			errChan <- err
			continue
		}
		fmt.Printf("\n%s", message)
	}
}
