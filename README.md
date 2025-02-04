# UDP to TCP to UDP

This application has the intention to listen in a specific remote UDP IP and send its packets through TCP connection and convert back to UDP, all of this in the same order as received.

This may introduce some latency, but I'm trading latency for order.

But you may wonder _way it would be out of order?_, goroutines, the server will handle every connection in a goroutine as the client also sends the buffered data in a goroutine.
