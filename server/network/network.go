package network

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/ehhio/ehhworldserver/server/game"

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

	// Maximum message size allowed from peer, in bytes
	maxMessageSize = 512

	// Size of a packed-message header, in bytes
	packedMessageHeaderSize = 2

	// Outbound message channel buffer size
	outboundMessageBuffer = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOrigin,
}

// // Network is a top-level networking object containing a websocket communication hub
// type Network struct {
// 	address string
// 	hub     *Hub
// }

// // NewNetwork creates a new network object for serving websocket requests from clients
// func NewNetwork(address string) *Network {
// 	return &Network{
// 		address: address,
// 		hub:     NewHub(),
// 	}
// }

// CheckOrigin validates incoming websocket connection origins
func CheckOrigin(r *http.Request) bool {
	return r.Header.Get("Origin") == "http://localhost:8080"
}

// Serve starts running a hub to handle clients via websocket connections
func Serve(address string, game *game.Game) *Hub {
	// Start a hub
	hub := NewHub(game)
	go hub.Start()

	// Configure the webserver to point to the hub
	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocketRequest(hub, w, r)
	})

	// Listen
	go func() {
		log.WithFields(log.Fields{
			"address": address,
		}).Info("Starting webserver")

		err := http.ListenAndServe(address, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	return hub
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"url": r.URL.String(),
	}).Info("Webserver serving request")

	if r.URL.Path != "/" {
		http.Error(w, "Not Found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	http.ServeFile(w, r, "./network/websocket_test.html")
}

// handleWebsocketRequest negotiates initial websocket requests from peers
func handleWebsocketRequest(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(hub, conn)
	client.hub.register <- client

	// Start client message handler routines
	go client.outboundHandler()
	go client.inboundHandler()
}
