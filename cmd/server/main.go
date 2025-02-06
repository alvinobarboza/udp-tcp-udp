package main

import (
	"fmt"

	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/tcp"
)

func main() {
	// local := "10.50.0.120:1234"
	// remote := "234.50.99.2:6000"
	defer fmt.Println()

	file := filehandler.NewFileHandler()
	file.NewFile("teste.bin")

	tcpServer := tcp.NewTCPServer("localhost:3002", file)

	if errL := tcpServer.Listen(); errL != nil {
		panic(errL)
	}
}
