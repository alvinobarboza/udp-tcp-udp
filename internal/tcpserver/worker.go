package tcpserver

import (
	"log"

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
		if w.oq.Length() < 10 {
			continue
		}

		data := w.oq.Pop()

		if data == nil {
			continue
		}

		log.Println("POP:", data.Counter, len(data.Data))
		if w.file != nil {
			w.file.Write(data.Data)
		}
		errCounter := 0
		rounds := len(data.Data) / args.MPEGTS_PKT_DEFAULT

		start := 0
		end := args.MPEGTS_PKT_DEFAULT
		for range rounds {
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
	}
}

func (w *worker) Enqueue(data *utils.TCPBuffData) {
	w.oq.Add(data)
}
