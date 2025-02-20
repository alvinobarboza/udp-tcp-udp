package args

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type ArgValue string

const (
	// Args for both
	FILE_SAVE ArgValue = "-save_file" // file name to save, if specified
	//TODO: FILE_SIZE      ArgValue = "-save_file" // how long in bytes should it be saved, : default 10000000 bytes

	// Args for client
	SERVER_ARG     ArgValue = "-server_ip"    // server to send data
	MCAST_ARG      ArgValue = "-mcast"        // mcast to listen to
	NET_INTER_ARG  ArgValue = "-eth"          // from which interface to listen on
	TIMER_ARG      ArgValue = "-timer"        // how long should it run : default 0
	MPEGTS_BUF_ARG ArgValue = "-mpegtsBuffer" // pkt count to send on a TCP conn: default 50
	MPEGTS_PKT     ArgValue = "-packetSize"   // max 1316 min 188 default 1316

	// Args for server
	LISTEN_IP    ArgValue = "-listen_ip"    // TCP IP to listen
	LOCAL_MCAST  ArgValue = "-local_mcast"  // IP from local interface to listen on
	REMOTE_MCAST ArgValue = "-remote_mcast" // UDP IP to listen to

	// Default values
	TS_PACKET_DEFAULT  int = 188
	MPEGTS_PKT_DEFAULT int = 1316
	TIMER_DEFAULT      int = 0
	DEFAULT_TCP_PKT    int = 50

	USAGE_MESSAGE_CLIENT string = `
Usage: udp-tcp-client [options]

Options:
  -server_ip <IP:PORT>
      Server IP:PORT to send data.
      Example: -server_ip 0.0.0.0:0000

  -mcast <IP:PORT>
      Multicast IP to listen to.
      Example: -mcast 0.0.0.0:0000

  -eth <interface>
      Network interface to listen on (e.g., eth0).
      Example: -eth eth0

  -timer <seconds>
      Duration (in seconds) for which the app should run.
      Default: 0 (forever)
      Example: -timer 600

  -mpegtsBuffer <number> (>50)
      Number of packets to send on a TCP connection.
      Default: 50
      Example: -mpegtsBuffer 60

  -packetSize <size>
      Packet size setting.
      Valid range: 188 (min) to 1316 (max)
      Default: 1316
      Example: -packetSize 376
	
  -h 
      Help
      Example: -h


	`
	USAGE_MESSAGE_SERVER string = `
Usage: udp-tcp-server [options]

Options:
  -listen_ip <IP:PORT>
      IP:PORT to listen.
      Example: -listen_ip 0.0.0.0:0000

  -local_mcast <IP:PORT>
      Ethernet IP to send on.
      Example: -local_mcast 0.0.0.0:0000

  -remote_mcast <IP:PORT>
      Multicast IP to send to.
      Example: -remote_mcast 0.0.0.0:0000

  -h 
      Help
      Example: -h


	`
)

var ErrMandatoryArg = errors.New("missing mandatory argmunt")

func HelpClient(args []string) {
	if len(args) > 1 && args[1] == "-h" {
		messageAndExit(USAGE_MESSAGE_CLIENT)
	}
}

func HelpServer(args []string) {
	if len(args) > 1 && args[1] == "-h" {
		messageAndExit(USAGE_MESSAGE_SERVER)
	}
}

func ValueFromArg(args []string, key ArgValue) string {
	for i, arg := range args {
		if arg == string(key) {
			valueIndex := i + 1
			if valueIndex < len(args) {
				return args[valueIndex]
			}
		}
	}
	return ""
}

func ValueFromArgFileSave(args []string, key ArgValue) bool {
	for _, arg := range args {
		if arg == string(key) {
			return true
		}
	}
	return false
}

func ValidateMandatoryClient(s, m, e string) {
	if s == "" || m == "" || e == "" {
		fmt.Printf(
			"Error: %v\n\t-server_ip: %v\n\t-mcast: %v\n\t-eth: %v\n\n%v\n",
			ErrMandatoryArg,
			s, m, e, USAGE_MESSAGE_CLIENT)
		os.Exit(1)
	}
}

func ValidateMandatoryServer(s, l, r string) {
	if s == "" || l == "" || r == "" {
		fmt.Printf(
			"Error: %v\n\t-listen_ip: %v\n\t-local_mcast: %v\n\t-remote_mcast: %v\n\n%v\n",
			ErrMandatoryArg,
			s, l, r, USAGE_MESSAGE_SERVER)
		os.Exit(1)
	}
}

func ConvertMpegtsPktSize(pktSize string) int {
	defaultValue := MPEGTS_PKT_DEFAULT
	if pktSize != "" {
		shouldReturn := false
		message := "-packetSize should be a number multiplier of 188, max 1316 min 188 default 1316!"
		var err error
		defaultValue, err = strconv.Atoi(pktSize)
		if err != nil {
			shouldReturn = true
		}
		if defaultValue > MPEGTS_PKT_DEFAULT || defaultValue < TS_PACKET_DEFAULT {
			shouldReturn = true
		}
		if defaultValue%TS_PACKET_DEFAULT != 0 {
			shouldReturn = true
		}
		if shouldReturn {
			messageAndExit(message)
		}
		return defaultValue
	}

	return defaultValue
}

func ConvertTimer(timer string) int {
	return convertGeneric(timer, TIMER_DEFAULT, "-timer should be a number, how long should it run : default 0(forever)!")
}

func ConvertMpegtsBuf(buff string) int {
	return convertGeneric(buff, DEFAULT_TCP_PKT, "-mpegtsBuffer should be a number, pkt count to send on a TCP conn: default 50!")
}

func convertGeneric(value string, defaultValue int, message string) int {
	if value != "" {
		var err error
		number, err := strconv.Atoi(value)
		if err != nil {
			messageAndExit(message)
		}
		if number < defaultValue {
			return defaultValue
		}
		return number
	}
	return defaultValue
}

func messageAndExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}
