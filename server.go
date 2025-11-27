package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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

func (t *TcpServer) Start() error {
	var err error
	t.Listener, err = net.Listen("tcp", t.Addr)

	if err != nil {
		return err
	}

	for {
		conn, err := t.Listener.Accept()

		if err != nil {
			return err
		}

		go HandleConnection(conn)

	}

}

func HandleConnection(conn net.Conn) {

	buf := make([]byte, 1024)
	defer conn.Close()
	go func() {

		for {
			n, err := conn.Read(buf)

			if err != nil {
				log.Println("client desconected : ", err)
				return
			}
			fmt.Println(string(buf[:n]))

		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Println("write error:", err)
			return
		}

	}

}
