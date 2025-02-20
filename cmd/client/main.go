package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/udpclient"
)

func main() {
	f, perr := os.Create("cpu-client.pprof")
	if perr != nil {
		log.Fatal(perr)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	defer fmt.Println()
	argsValues := os.Args

	args.HelpClient(argsValues)

	serverIp := args.ValueFromArg(argsValues, args.SERVER_ARG)
	mcastIp := args.ValueFromArg(argsValues, args.MCAST_ARG)
	ethName := args.ValueFromArg(argsValues, args.NET_INTER_ARG)
	saveFile := args.ValueFromArgFileSave(argsValues, args.FILE_SAVE)

	args.ValidateMandatoryClient(serverIp, mcastIp, ethName)

	timer := args.ValueFromArg(argsValues, args.TIMER_ARG)
	timerNumber := args.ConvertTimer(timer)

	mpegtsBuffer := args.ValueFromArg(argsValues, args.MPEGTS_BUF_ARG)
	mpegtsBufSize := args.ConvertMpegtsBuf(mpegtsBuffer)

	mpegtsPkt := args.ValueFromArg(argsValues, args.MPEGTS_PKT)
	mpegtsPktSize := args.ConvertMpegtsPktSize(mpegtsPkt)

	var file filehandler.FileHandler
	if saveFile {
		file = filehandler.NewFileHandler()
		if err := file.NewFile("client.bin"); err != nil {
			panic(err)
		}
	}

	tcpCon, errc := udpclient.NewTCPClient(serverIp)
	if errc != nil {
		panic(errc)
	}

	udpListener := udpclient.NewUDPListener(mpegtsBufSize, mpegtsPktSize, timerNumber, tcpCon, file)

	if err := udpListener.SetUpListener(ethName, mcastIp); err != nil {
		panic(err)
	}

	if errU := udpListener.Listen(); errU != nil {
		panic(errU)
	}
}
