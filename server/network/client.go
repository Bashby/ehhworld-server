package network

// Note: Adapted Heavily from https://github.com/gorilla/websocket/blob/master/examples/chat

import (
	"encoding/binary"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// ClientMessage Structure to associate a client with a message
type ClientMessage struct {
	// The Client sending/receiving the message
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

// NewClient constructs an object to represent a remote peer that will communicate with us over a websocket
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		outbound: make(chan []byte, outboundMessageBuffer)}
}

// outboundHandler pumps messages from the parent hub to the peer websocket connection.
func (c *Client) outboundHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.WithFields(log.Fields{
			"client": c.conn.RemoteAddr().String(),
		}).Info("Client outbound handler stopping")

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
				log.WithFields(log.Fields{
					"client": c.conn.RemoteAddr().String(),
				}).Info("Client outbound handler; Channel closed by hub request")

				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.WithFields(log.Fields{
					"client": c.conn.RemoteAddr().String(),
					"error":  err,
				}).Info("Client outbound handler; Writer open error")

				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.outbound)
			for i := 0; i < n; i++ {
				w.Write(<-c.outbound)
			}

			if err := w.Close(); err != nil {
				log.WithFields(log.Fields{
					"client": c.conn.RemoteAddr().String(),
					"error":  err,
				}).Info("Client outbound handler; Writer close error")

				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.WithFields(log.Fields{
					"client": c.conn.RemoteAddr().String(),
					"error":  err,
				}).Info("Client outbound handler; Error sending ping")

				return
			}
		}
	}
}

// inboundHandler Pumps messages from the peer websocket connection to the parent hub.
func (c *Client) inboundHandler() {
	defer func() {
		log.WithFields(log.Fields{
			"client": c.conn.RemoteAddr().String(),
		}).Info("Client inbound handler stopping")

		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{
				"client": c.conn.RemoteAddr().String(),
				"error":  err,
			}).Info("Client inbound handler; Read error")

			break
		}

		// Only support binary messages
		if messageType != websocket.BinaryMessage {
			continue
		}

		// Handle blank or too short messages
		offset := 0
		totalMessageLength := len(message)
		if totalMessageLength == 0 || totalMessageLength <= offset+packedMessageHeaderSize {
			continue
		}

		// Extract messages and pipe upstream to the hub
		for offset < totalMessageLength {
			// Header read bounds check
			if offset+packedMessageHeaderSize >= totalMessageLength {
				break
			}

			// Determine location
			msglen := int(binary.BigEndian.Uint16(message[offset : offset+packedMessageHeaderSize]))

			// Payload read bounds check
			if offset+packedMessageHeaderSize+msglen >= totalMessageLength {
				break
			}

			// Extract
			payload := message[offset+packedMessageHeaderSize : offset+packedMessageHeaderSize+msglen]

			// Pipe
			c.hub.inbound <- &ClientMessage{message: payload, client: c}

			// Update offset
			offset += packedMessageHeaderSize + msglen
		}
	}
}