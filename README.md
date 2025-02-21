# UDP to TCP to UDP : Multicast mpegts stream

This application has the intention to listen in a specific remote UDP IP and send its packets through TCP connection and convert back to UDP, all of this in the same order as received.

This may introduce some latency, but I'm trading latency for order.

But you may wonder _way it would be out of order?_, goroutines, the server will handle every connection in a goroutine as the client also sends the buffered data in a goroutine.

# Usage;

Compile server;

go build -o server cmd/server

Compile client;

go build -o client cmd/client

# CLI

## Server

Usage: ./server [options]

Options:

```
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
        Number of packets(1316b \* n) to send on a TCP connection.
        Default: 50
        Example: -mpegtsBuffer 60

    -packetSize <size>
        Packet size setting.
        Valid range: 188 (min) to 1316 (max)
        Default: 1316
        Example: -packetSize 376

  -save_file
      Boolean flag, on/off, if passed it will save a file with all mpegts packets = client.bin

    -h
        Help
        Example: -h
```

## Client

Usage: udp-tcp-server [options]

Options:

```
  -listen_ip <IP:PORT>
      IP:PORT to listen.
      Example: -listen_ip 0.0.0.0:0000

  -local_mcast <IP:PORT>
      Ethernet IP to send on.
      Example: -local_mcast 0.0.0.0:0000

  -remote_mcast <IP:PORT>
      Multicast IP to send to.
      Example: -remote_mcast 0.0.0.0:0000

  -q_window <number>
	  Number of items in the queue before processing.
	  Default: 4
      Example: -q_window 20

  -save_file
      Boolean flag, on/off, if passed it will save a file with all mpegts packets = server.bin

  -h
      Help
      Example: -h
```
