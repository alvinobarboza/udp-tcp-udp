package udpclient

import (
	"fmt"
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

	end := ""
	var errReturn error

	go func() {
		for {
			select {
			case val := <-done:
				close(connSignal)
				end = val
				return
			case val := <-errChan:
				close(connSignal)
				end = val.Error()
				errReturn = val
				return
			}

		}
	}()

	bufCh := make(chan []byte, 100)

	go func() {
		now := time.Now()
		for data := range bufCh {
			// if le > 20 {
			// fmt.Println(len(bufCh), len(data))
			// }
			tcpBuff := &utils.TCPBuffData{
				Data:    data,
				MS:      uint32(time.Duration(time.Since(now).Microseconds()) / time.Duration(ul.tcpMultiplierBuf)),
				Counter: uint64(time.Now().UnixMilli()),
			}
			fmt.Println(time.Since(now).Microseconds())
			now = time.Now()
			go ul.tcpHandler.Write(tcpBuff, errChan, connSignal)
			// errChan <- fmt.Errorf("ended")

		}
	}()

	buf := make([]byte, ul.packetSize)
	for {
		if end != "" {
			log.Println(end)
			return errReturn
		}

		tcpBuffer := make([]byte, 0)
		for range ul.tcpMultiplierBuf {
			read, _, errC := conn.ReadFromUDP(buf)
			if errC != nil {
				return errC
			}
			tcpBuffer = append(tcpBuffer, buf[:read]...)
		}
		bufCh <- tcpBuffer
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
