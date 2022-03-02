Jacktrip UDP forwarder
----------------------

This program can run on a server and accepts UDP packets on port 4464.
It stores all clients in a list and each time it receives a new packet it sends it straight on to the other clients.
The only thing jacktrip-specific is the JackTripHeader Struct which I use to debug the packets I am sending between
devices. If that is removed it's just a udpforwarder.

For Jacktrip it will either work with 2 sending and receiving clients, or one sending client and any number of receiving clients.
That's because jacktrip clients only expect a single audio stream from one other client.

```
go run .
```
