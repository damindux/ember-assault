package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
)

type Vector2 struct {
	X, Y float32
}

type Player struct {
    ID                    string
	Pos                   Vector2
	Height, Width, Health int32
}

type Client struct {
	clientAddr net.UDPAddr
	player     Player
}

func main() {
	addr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 10000,
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Server is listening on:", addr)

	clients := make(map[string]*Client)
    clientID := 1

	gob.Register(Vector2{})

	for {
		buf := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		var player Player
		buffer := bytes.NewBuffer(buf[:n])
		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&player); err != nil {
			fmt.Println("Error decoding player data:", err)
			continue
		}

        playerID := strconv.Itoa(clientID)
        if _, exists := clients[clientAddr.String()]; !exists {
            clientID++
            player.ID = playerID
        }

        fmt.Println("Received player data from:", clientAddr, player)

		// Update or add client
		clients[clientAddr.String()] = &Client{clientAddr: *clientAddr, player: player}

		// Broadcast player state to all clients
		for _, client := range clients {
			if client.clientAddr.String() != clientAddr.String() {
				var sendBuffer bytes.Buffer
				encoder := gob.NewEncoder(&sendBuffer)
				if err := encoder.Encode(player); err != nil {
					fmt.Println("Error encoding player:", err)
					continue
				}
				_, err := conn.WriteToUDP(sendBuffer.Bytes(), &client.clientAddr)
				if err != nil {
					fmt.Println("Error sending to client:", err)
				}
			}
		}
	}
}
