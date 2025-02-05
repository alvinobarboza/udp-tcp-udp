package udp

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

type PktConsumer interface {
	Write([]byte, chan error)
}

type UPDListener interface {
	SetUpListener(ethName, udpIp string) error
	Listen() error
}

type udpListener struct {
	tcpMultiplierBuf int
	packetSize       int
	timerSeconds     int
	eth              *net.Interface
	udpAddr          *net.UDPAddr
	writer           PktConsumer
}

func NewUDPListener(tcpMultiplierBuf, packetSize, timerSeconds int, writer PktConsumer) UPDListener {
	return &udpListener{
		tcpMultiplierBuf: tcpMultiplierBuf,
		packetSize:       packetSize,
		timerSeconds:     timerSeconds,
		writer:           writer,
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

	errorTerminator := make(chan error)
	for {
		select {
		case res := <-done:
			log.Println(res)
			return nil
		case err := <-errorTerminator:
			log.Println(err)
			return err
		default:
			tcpBuffer := make([]byte, 0)

			for i := 0; i < ul.tcpMultiplierBuf; i++ {
				buf := make([]byte, ul.packetSize)
				countBytes, _, errC := conn.ReadFrom(buf)

				if errC != nil {
					return errC
				}
				tcpBuffer = append(tcpBuffer, buf[:countBytes]...)
				fmt.Printf("%04d %08d\r", countBytes, len(tcpBuffer))
			}

			go ul.writer.Write(tcpBuffer, errorTerminator)
		}
	}
}

func (ul *udpListener) setTimeout() chan string {
	done := make(chan string)
	if ul.timerSeconds > 0 {
		go func() {
			time.Sleep(time.Duration(rand.IntN(ul.timerSeconds)) * time.Second)
			done <- "Time ended!"
		}()
	}
	return done
}
