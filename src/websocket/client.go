package websocket

// Note: Adapted Heavily from https://github.com/gorilla/websocket/blob/master/examples/chat

import (
	"log"
	"net/http"
	"time"

	"encoding/binary"

	"github.com/gorilla/websocket"
)

// Constants
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512

	// Size of a packed-message header, in bytes
	packedMessageHeaderSize = 2

	// Outbound message channel buffer size
	outboundMessageBuffer = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ClientMessage Structure to associate a client with a message
type ClientMessage struct {
	// The Client sending the message
	client *Client

	// The message
	message []byte
}

// Client Message handler structure around peers and their websocket connection
type Client struct {
	// Reference to parent Hub
	hub *Hub

	// Peer's websocket connection
	conn *websocket.Conn

	// Outbound message channel to peer
	outbound chan []byte
}

// handleWebsocketRequest Negotiates initial websocket request from the peer.
func handleWebsocketRequest(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, outbound: make(chan []byte, outboundMessageBuffer)}
	client.hub.register <- client

	// Start client message handler routines
	go client.outboundHandler()
	go client.inboundHandler()
}

// outboundHandler pumps messages from the parent hub to the peer websocket connection.
func (c *Client) outboundHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.outbound:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed our channel and wants us dead
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Printf("[Client Outbound][%s] Channel closed\n", c.conn.RemoteAddr().String())
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("[Client Outbound][%s] Error opening writer: %v\n", c.conn.RemoteAddr().String(), err)
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.outbound)
			for i := 0; i < n; i++ {
				w.Write(<-c.outbound)
			}

			if err := w.Close(); err != nil {
				log.Printf("[Client Outbound][%s] Error closing writer: %v\n", c.conn.RemoteAddr().String(), err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("[Client Outbound][%s] Error sending ping: %v\n", c.conn.RemoteAddr().String(), err)
				return
			}
		}
	}
}

// inboundHandler Pumps messages from the peer websocket connection to the parent hub.
func (c *Client) inboundHandler() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("[Client Inbound][%s] Read error: %v\n", c.conn.RemoteAddr().String(), err)
			}
			break
		}

		// Extract messages and pipe upstream to the hub
		offset := 0
		for offset < len(message) {
			// Determine location
			msglen := int(binary.BigEndian.Uint16(message[offset : offset+packedMessageHeaderSize]))

			// Extract
			payload := message[offset+packedMessageHeaderSize : offset+packedMessageHeaderSize+msglen]

			// Pipe
			c.hub.inbound <- ClientMessage{message: payload, client: c}

			// Update offset
			offset += packedMessageHeaderSize + msglen
		}
	}
}
