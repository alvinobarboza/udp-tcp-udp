package main

import (
	"fmt"
	"os"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/tcpserver"
)

func main() {
	defer fmt.Println()
	argsValues := os.Args

	args.HelpServer(argsValues)

	listenIp := args.ValueFromArg(argsValues, args.LISTEN_IP)
	localMcastIp := args.ValueFromArg(argsValues, args.LOCAL_MCAST)
	remoteMcastIp := args.ValueFromArg(argsValues, args.REMOTE_MCAST)

	args.ValidateMandatoryServer(listenIp, localMcastIp, remoteMcastIp)

	// file := filehandler.NewFileHandler()
	// file.NewFile("teste.bin")

	updSender, err := tcpserver.NewUDPSender(localMcastIp, remoteMcastIp)
	if err != nil {
		panic(err)
	}

	tcpServer := tcpserver.NewTCPServer(listenIp, updSender, nil)

	if errL := tcpServer.Listen(); errL != nil {
		panic(errL)
	}
}
