package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type NetworkClient struct {
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
	player     *Player
}

// Initialize the network client
func InitNetworkClient(serverIP string, serverPort int, player *Player) (*NetworkClient, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}

	return &NetworkClient{
		conn:       conn,
		serverAddr: serverAddr,
		player:     player,
	}, nil
}

// SendPlayerData sends the player's current state to the server
func (nc *NetworkClient) SendPlayerData() {
	ticker := time.NewTicker(2 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		if err := encoder.Encode(nc.player); err != nil {
			fmt.Println("Error encoding player data:", err)
			continue
		}

		_, err := nc.conn.Write(buffer.Bytes())
		if err != nil {
			fmt.Println("Error sending player data to server:", err)
			continue
		}
	}
}

// ReceiveUpdates listens for updates from the server
func (nc *NetworkClient) ReceiveUpdates() {
	for {
		buf := make([]byte, 1024)
		n, _, err := nc.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			continue
		}

		var players []Player
		buffer := bytes.NewBuffer(buf[:n])
		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&players); err != nil {
			fmt.Println("Error decoding data from server:", err)
			continue
		}
	}
}

// Close closes the network connection
func (nc *NetworkClient) Close() {
	nc.conn.Close()
}
