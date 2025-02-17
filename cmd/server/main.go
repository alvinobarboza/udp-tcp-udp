package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/tcpserver"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

func main() {
	defer fmt.Println()
	argsValues := os.Args

	args.HelpServer(argsValues)

	listenIp := args.ValueFromArg(argsValues, args.LISTEN_IP)
	localMcastIp := args.ValueFromArg(argsValues, args.LOCAL_MCAST)
	remoteMcastIp := args.ValueFromArg(argsValues, args.REMOTE_MCAST)

	args.ValidateMandatoryServer(listenIp, localMcastIp, remoteMcastIp)

	log.SetFlags(log.Lshortfile)

	// TODO: Adjust filehandling
	file := filehandler.NewFileHandler()
	file.NewFile("teste.bin")
	// END

	updSender, err := tcpserver.NewUDPSender(localMcastIp, remoteMcastIp)
	if err != nil {
		panic(err)
	}

	worker := tcpserver.NewWorker(file, updSender, utils.NewQueue())

	tcpServer := tcpserver.NewTCPServer(listenIp, worker)

	if errL := tcpServer.Listen(); errL != nil {
		panic(errL)
	}
}
