package udpclient

import (
	"log"
	"net"
	"time"

	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

type UDPListener interface {
	SetUpListener(ethName, udpIp string) error
	Listen() error
}

type udpListener struct {
	tcpMultiplierBuf int
	packetSize       int
	timerSeconds     int
	eth              *net.Interface
	udpAddr          *net.UDPAddr
	tcpHandler       TCPClient
	fileHandler      filehandler.FileHandler
}

func NewUDPListener(
	tcpMultiplierBuf,
	packetSize,
	timerSeconds int,
	tcp TCPClient,
	file filehandler.FileHandler,
) UDPListener {

	return &udpListener{
		tcpMultiplierBuf: tcpMultiplierBuf,
		packetSize:       packetSize,
		timerSeconds:     timerSeconds,
		tcpHandler:       tcp,
		fileHandler:      file,
	}
}

func (ul *udpListener) SetUpListener(ethName, udpIp string) error {
	inter, errI := net.InterfaceByName(ethName)
	if errI != nil {
		return errI
	}

	udpAddr, err := net.ResolveUDPAddr("udp", udpIp)
	if err != nil {
		return err
	}

	ul.eth = inter
	ul.udpAddr = udpAddr

	return nil
}

func (ul *udpListener) Listen() error {
	conn, err := net.ListenMulticastUDP("udp", ul.eth, ul.udpAddr)

	if err != nil {
		return err
	}
	defer conn.Close()

	done := ul.setTimeout()
	errChan := make(chan error)
	connSignal := make(chan string)

	for {
		tcpCon, errCon := ul.tcpHandler.GetConn()
		if errCon != nil {
			return errCon
		}
		select {
		case res := <-done:
			close(connSignal)
			ul.tcpHandler.Close(tcpCon)
			log.Println(res)
			return nil
		case err := <-errChan:
			close(connSignal)
			ul.tcpHandler.Close(tcpCon)
			return err
		default:
			tcpBuffer := make([]byte, 0)

			now := time.Now()
			for i := 0; i < ul.tcpMultiplierBuf; i++ {
				buf := make([]byte, ul.packetSize)
				countBytes, _, errC := conn.ReadFrom(buf)

				if errC != nil {
					return errC
				}
				tcpBuffer = append(tcpBuffer, buf[:countBytes]...)
			}
			tcpBuff := &utils.TCPBuffData{
				Data:    tcpBuffer,
				MS:      uint32(time.Since(now).Abs().Microseconds()),
				Counter: uint64(now.UnixMilli()),
			}

			go ul.tcpHandler.Write(
				tcpCon, tcpBuff,
				errChan, connSignal,
			)
		}
	}
}

func (ul *udpListener) setTimeout() chan string {
	done := make(chan string)
	if ul.timerSeconds > 0 {
		go func() {
			time.Sleep(time.Duration(ul.timerSeconds) * time.Second)
			done <- "Time ended!"
		}()
	}
	return done
}
