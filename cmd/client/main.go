package main

import (
	"fmt"
	"os"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/udp"
)

func main() {

	defer fmt.Println()
	argsValues := os.Args

	args.HelpClient(argsValues)

	serverIp := args.ValueFromArg(argsValues, args.SERVER_ARG)
	mcastIp := args.ValueFromArg(argsValues, args.MCAST_ARG)
	ethName := args.ValueFromArg(argsValues, args.NET_INTER_ARG)

	args.ValidateMandatoryClient(serverIp, mcastIp, ethName)

	timer := args.ValueFromArg(argsValues, args.TIMER_ARG)
	timerNumber := args.ConvertTimer(timer)

	mpegtsBuffer := args.ValueFromArg(argsValues, args.MPEGTS_BUF_ARG)
	mpegtsBufSize := args.ConvertMpegtsBuf(mpegtsBuffer)

	mpegtsPkt := args.ValueFromArg(argsValues, args.MPEGTS_PKT)
	mpegtsPktSize := args.ConvertMpegtsPktSize(mpegtsPkt)

	// file := filehandler.NewFileHandler()
	// if err := file.NewFile("teste.bin"); err != nil {
	// 	panic(err)
	// }

	tcpCon, errc := tcp.NewTCPClient("localhost:3002")
	if errc != nil {
		panic(errc)
	}

	udpListener := udp.NewUDPListener(tcpMultiplierBuf, packetSize, timerSeconds, tcpCon)

	if err := udpListener.SetUpListener("Ethernet", "234.50.99.2:6000"); err != nil {
		panic(err)
	}

	if errU := udpListener.Listen(); errU != nil {
		panic(errU)
	}
}
