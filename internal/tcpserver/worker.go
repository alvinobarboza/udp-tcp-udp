package tcpserver

import (
	"fmt"

	"github.com/alvinobarboza/udp-tcp-udp/internal/args"
	"github.com/alvinobarboza/udp-tcp-udp/internal/filehandler"
	"github.com/alvinobarboza/udp-tcp-udp/internal/utils"
)

const MAGIC_NUMBER = 3000
const ERR_BEFORE_RETURN = 30

type Worker interface {
	Start()
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

func (w *worker) Start() {
	defer w.udp.CloseConn()
	if w.file != nil {
		defer w.file.CloseConn()
	}

	for {
		data := w.oq.Pop()

		if data == nil {
			continue
		}

		fmt.Println("POP:", data.Counter, data.MS)
		if w.file != nil {
			w.file.Write(data.Data)
		}
		errCounter := 0
		// pktToSend := make([]byte, args.MPEGTS_PKT_DEFAULT)
		rounds := len(data.Data) / args.MPEGTS_PKT_DEFAULT

		start := 0
		end := args.MPEGTS_PKT_DEFAULT
		for range rounds {
			// copy(pktToSend, data.Data[start:end])
			err := w.udp.Write(data.Data[start:end])
			if err != nil {
				errCounter++
				if errCounter > ERR_BEFORE_RETURN {
					break
				}
			}
			start = end
			end += args.MPEGTS_PKT_DEFAULT
		}
		// if len(pktToSend) > 0 &&
		// 	len(pktToSend) < args.MPEGTS_PKT_DEFAULT &&
		// 	errCounter < ERR_BEFORE_RETURN {

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
