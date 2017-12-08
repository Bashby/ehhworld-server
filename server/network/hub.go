package network

import (
	"encoding/binary"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/ehhio/ehhworldserver/server/utility"

	"bitbucket.org/ehhio/ehhworldserver/server/game"
	"bitbucket.org/ehhio/ehhworldserver/server/player"
	"bitbucket.org/ehhio/ehhworldserver/server/protobuf"

	"github.com/golang/protobuf/proto"
)

// Hub is a networking message manager between clients and the game, using websockets.
type Hub struct {
	// Serving is true if the hub is serving clients
	serving bool

	// Game object the hub will pipe data to
	game *game.Game

	// Clients maps network clients to players in the connected game
	clients map[*Client]*player.Player

	// Inbound messages channel from Clients
	inbound chan *ClientMessage

	// Outbound messages channel to Clients
	outbound chan *ClientMessage

	// Register requests channel from new Clients
	register chan *Client

	// Unregister requests channel from existing Clients
	unregister chan *Client
}

// NewHub constructs a websocket Hub to manage clients and messages to and from them
func NewHub(game *game.Game) *Hub {
	return &Hub{
		clients:    make(map[*Client]*player.Player),
		inbound:    make(chan *ClientMessage),
		outbound:   make(chan *ClientMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		game:       game,
	}
}

// Start starts the serving loop for the hub to handle websocket messages from clients
func (h *Hub) Start() {
	defer func() {
		log.Printf("Hub is stopping.")
	}()

	h.serving = true

	// Core serving loop
	for h.serving {
		select {
		case client := <-h.register:
			h.handleClientConnect(client)
		case client := <-h.unregister:
			h.handleClientDisconnect(client)
		case message := <-h.inbound:
			h.handleInboundMessage(message)
		case message := <-h.outbound:
			h.handleOutboundMessage(message)
		}
	}

	// Drop all players
	for client, player := range h.clients {
		log.WithFields(log.Fields{
			"client address": client.conn.RemoteAddr().String(),
		}).Info("Client Disconnected; Hub is stopping")

		h.game.RemoveObject(player)
		delete(h.clients, client)
		close(client.outbound)
	}
}

// Stop stops the hub serving clients
func (h *Hub) Stop() {
	log.Info("Stopping webserver")
	h.serving = false
}

func (h *Hub) handleClientConnect(client *Client) {
	log.WithFields(log.Fields{
		"client address": client.conn.RemoteAddr().String(),
	}).Info("New Client Connected")

	player := player.NewPlayer(utility.GeneratePlayerName(), h.game.GetGameMap().RandomPositionHighResolution(), h.game.GetGameMap().GetBlocksSize())
	h.game.AddObject(player)
	h.clients[client] = player
}

func (h *Hub) handleClientDisconnect(client *Client) {
	if player, exists := h.clients[client]; exists {
		log.WithFields(log.Fields{
			"client address": client.conn.RemoteAddr().String(),
		}).Info("Client Disconnected; Own volition")

		h.game.RemoveObject(player)
		delete(h.clients, client)
		close(client.outbound)
	}
}

func (h *Hub) handleInboundMessage(msg *ClientMessage) {
	go processMessage(msg)
}

func (h *Hub) handleOutboundMessage(msg *ClientMessage) {
	if msg.client != nil {
		// Message a client
		if player, exists := h.clients[msg.client]; exists {
			select {
			case msg.client.outbound <- packMessage(msg.message):
			default:
				log.WithFields(log.Fields{
					"client address": msg.client.conn.RemoteAddr().String(),
				}).Info("Client Disconnected; Buffer full")

				h.game.RemoveObject(player)
				delete(h.clients, msg.client)
				close(msg.client.outbound)
			}
		}
	} else {
		// Broadcast message
		for client, player := range h.clients {
			select {
			case client.outbound <- packMessage(msg.message):
			default:
				log.WithFields(log.Fields{
					"client address": client.conn.RemoteAddr().String(),
				}).Info("Client Disconnected; Buffer full")

				h.game.RemoveObject(player)
				delete(h.clients, client)
				close(client.outbound)
			}
		}
	}
}

func packMessage(msg []byte) []byte {
	// Build a header for the outgoing message
	messageHeader := make([]byte, packedMessageHeaderSize)
	binary.BigEndian.PutUint16(messageHeader, uint16(len(msg)))
	return append(messageHeader, msg...)
}

func processMessage(message *ClientMessage) {
	// Decode message
	wrapper := &protobuf.Message{}
	err := proto.Unmarshal(message.message, wrapper)
	if err != nil {
		log.Fatal("Unmarshaling: ", err)
	}

	// Process payload
	switch msg := wrapper.Payload.(type) {
	case *protobuf.Message_Move:
		handleMove(msg, message.client.hub)
	case *protobuf.Message_Attack:
		handleAttack(msg, message.client.hub)
	}
}

func handleMove(msg *protobuf.Message_Move, hub *Hub) {
	log.Println("Twas a Move message: ", msg.Move.Direction)
}

func handleAttack(msg *protobuf.Message_Attack, hub *Hub) {
	log.Println("Twas a Attack message: ", msg.Attack.Target)
}

// test := &websocket.Message{
// 	Type: 1,
// 	Payload: &websocket.Message_Move{
// 		Move: &websocket.Move{
// 			Direction: "Left",
// 		},
// 	},
// }
// data, err := proto.Marshal(test)
// if err != nil {
// 	log.Fatal("marshaling error: ", err)
// }
