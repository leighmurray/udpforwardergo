package main

import "fmt"
import "net"
import "encoding/binary"
import "bytes"

type JacktripHeader struct {
	TimeStamp                  uint64
	SeqNumber                  uint16
	BufferSize                 uint16
	SamplingRate               uint8
	BitResolution              uint8
	NumIncomingChannelsFromNet uint8
	NumOutgoingChannelsToNet   uint8
}

func main() {
	udpListener, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 4464,
		IP:   net.ParseIP("0.0.0.0"),
	})

	if err != nil {
		panic(err)
	}

	defer udpListener.Close()

	fmt.Printf("server listening %s\n", udpListener.LocalAddr().String())
	remoteIPAddresses := make(map[string]net.UDPAddr)

	for {
		message := make([]byte, 1024)
		rlen, sender, err := udpListener.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		fmt.Printf("received from %s length %d\n", sender, rlen)

		jacktripHeaderBytes := message[:16]
		fmt.Printf("Jacktrip Header Bytes:\n%x\n", jacktripHeaderBytes)
		jacktripHeader := JacktripHeader{}
		err = binary.Read(bytes.NewReader(jacktripHeaderBytes), binary.LittleEndian, &jacktripHeader)

		if err != nil {
			fmt.Printf("Didn't decode properly\n")
			panic(err)
		}

		fmt.Printf("%+v\n", jacktripHeader)
		senderIPString := sender.String()
		if _, exists := remoteIPAddresses[senderIPString]; !exists {
			fmt.Println("Adding Remote Address")
			remoteIPAddresses[senderIPString] = *sender
		} else {
			fmt.Println("Already exists")
		}

		for currentRemoteIPString, remoteUDPAddr := range remoteIPAddresses {
			if currentRemoteIPString == senderIPString {
				fmt.Printf("Not forwarding from %s to %s\n", senderIPString, currentRemoteIPString)
				continue
			}

			fmt.Printf("Forwarding message to %s\n", currentRemoteIPString)
			_, err := udpListener.WriteToUDP(message[:rlen], &remoteUDPAddr)
			if err != nil {
				panic(err)
			}
		}
	}
}
