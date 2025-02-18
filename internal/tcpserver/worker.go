package tcpserver

import (
	"fmt"

	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

const MAGIC_NUMBER = 3000
const ERR_BEFORE_RETURN = 30

type Worker interface {
	Start(bool)
	Enqueue(*utils.TCPBuffData)
}

func NewWorker(file filehandler.FileHandler,
	udp UDPSender, oq utils.OrderedQueue) Worker {
	return &worker{
		file: file,
		udp:  udp,
		oq:   oq,
	}
}

type worker struct {
	udp  UDPSender
	file filehandler.FileHandler
	oq   utils.OrderedQueue
}

func (w *worker) Start(saveToFile bool) {
	defer w.udp.CloseConn()
	defer w.file.CloseConn()

	for {
		data := w.oq.Pop()

		if data == nil {
			continue
		}

		fmt.Println("POP:", data.Counter, data.MS)
		if saveToFile {
			w.file.Write(data.Data)
		}
		// errCounter := 0
		// pktToSend := make([]byte, 0)
		// for i, pkt := range data.Data {
		// 	if i%args.MPEGTS_PKT_DEFAULT == 0 {
		// 		time.Sleep(time.Microsecond * time.Duration(data.MS))
		// 		err := w.udp.Write(pktToSend)
		// 		if err != nil {
		// 			errCounter++
		// 			if errCounter > ERR_BEFORE_RETURN {
		// 				break
		// 			}
		// 		}
		// 		pktToSend = make([]byte, 0)
		// 	}
		// 	pktToSend = append(pktToSend, pkt)
		// }
		// if len(pktToSend) > 0 &&
		// 	len(pktToSend) < args.MPEGTS_PKT_DEFAULT &&
		// 	errCounter < ERR_BEFORE_RETURN {

		// 	time.Sleep(time.Microsecond * time.Duration(data.MS))
		// 	err := w.udp.Write(pktToSend)
		// 	if err != nil {
		// 		continue
		// 	}
		// }
	}
}

func (w *worker) Enqueue(data *utils.TCPBuffData) {
	w.oq.Add(data)
}
