package main

import (
	"fmt"
	"net"
	"strings"
)

type TcpServer struct {
	Addr     string
	Listener net.Listener
}

func NewTcpServer(addr string) *TcpServer {
	return &TcpServer{
		Addr: addr,
	}
}

var clients = make(map[net.Conn]string)

var broadcast = make(chan string)

func (t *TcpServer) Start() error {
	var err error
	t.Listener, err = net.Listen("tcp", t.Addr)
	defer t.Listener.Close()

	if err != nil {
		return err
	}
	fmt.Println("Server running on :3000")

	go handleBroadcast()
	for {
		conn, err := t.Listener.Accept()

		if err != nil {
			return err
		}
		conn.Write([]byte("Enter your username :"))
		buf := make([]byte, 1024)

		n, _ := conn.Read(buf)
		username := strings.TrimSpace(string(buf[:n]))

		clients[conn] = username
		msg := fmt.Sprintf("%s join the chat !! ", clients[conn])
		fmt.Println(msg)
		broadcast <- msg
		go handleClient(conn)

	}

}

func handleClient(conn net.Conn) {
	defer func() {
		msg := fmt.Sprintf("%s left the chat !!", clients[conn])
		fmt.Println(msg)
		broadcast <- msg
		delete(clients, conn)
		conn.Close()
	}()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		msg := fmt.Sprintf("%s:%s", clients[conn], string(buf[:n]))

		broadcast <- msg
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		for conn := range clients {
			fmt.Fprintln(conn, msg)
		}
	}
}
