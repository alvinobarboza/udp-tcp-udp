package udp

import (
	"fmt"
	"net"
)

type PktConsumer interface {
	Write([]byte) ([]byte, error)
}

type UPDListener interface {
	SetUpListener(ethName, udpIp string) error
	Listen() error
}

type udpListener struct {
	tcpMultiplierBuf int
	packetSize       int
	tcpPktSend       int
	eth              *net.Interface
	udpAddr          *net.UDPAddr
	writer           PktConsumer
}

func NewUDPListener(tcpMultiplierBuf, packetSize, tcpPktSend int, writer PktConsumer) UPDListener {
	return &udpListener{
		tcpMultiplierBuf: tcpMultiplierBuf,
		packetSize:       packetSize,
		tcpPktSend:       tcpPktSend,
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

	counter := 0
	for {
		if ul.tcpPktSend > 0 {
			counter++

		}
		if counter == ul.tcpPktSend {
			break
		}

		tcpBuffer := make([]byte, 0)

		for i := 0; i < ul.tcpMultiplierBuf; i++ {
			buf := make([]byte, ul.packetSize)
			countBytes, _, errC := conn.ReadFrom(buf)

			if errC != nil {
				return errC
			}
			tcpBuffer = append(tcpBuffer, buf[:countBytes]...)
			fmt.Printf("%v %v %v %v\r", countBytes, i, len(tcpBuffer), counter)
		}

		// TODO:Migrate to goroutine
		// channels to break loop in case of error
		_, err := ul.writer.Write(tcpBuffer)
		if err != nil {
			return err
		}

	}

	return nil
}
