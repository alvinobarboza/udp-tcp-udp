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
	saveFile := args.ValueFromArgFileSave(argsValues, args.FILE_SAVE)

	args.ValidateMandatoryServer(listenIp, localMcastIp, remoteMcastIp)

	qWindow := args.ValueFromArg(argsValues, args.Q_WINDOW)
	qWindowSize := args.ConvertQWindow(qWindow)

	log.SetFlags(log.Lshortfile)

	var file filehandler.FileHandler
	if saveFile {
		file = filehandler.NewFileHandler()
		file.NewFile("server.bin")
	}

	updSender, err := tcpserver.NewUDPSender(localMcastIp, remoteMcastIp)
	if err != nil {
		panic(err)
	}

	worker := tcpserver.NewWorker(file, updSender, utils.NewQueue(), qWindowSize)

	tcpServer := tcpserver.NewTCPServer(listenIp, worker)

	if errL := tcpServer.Listen(); errL != nil {
		panic(errL)
	}
}
