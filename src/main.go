package main

import (
	"io/ioutil"
	"math/rand"
	"time"

	"log"

	"flag"

	"bitbucket.org/ashbyb/go-map-generation/websocket"
	"bitbucket.org/ashbyb/go-map-generation/world"
)

var width int
var height int
var address string
var debug bool

func init() {
	// Define input parameters
	flag.BoolVar(&debug, "debug", true, "whether to print debug statements during execution")
	flag.IntVar(&width, "width", 512, "width of the game world")
	flag.IntVar(&height, "height", 512, "height of the game world")
	flag.StringVar(&address, "address", ":8080", "http service address to listen on")

	// Seed pRNG
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// Parse inputs
	flag.Parse()

	// Configure logging
	if debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	// Generate the world
	world.Generate(width, height)

	// Serve
	websocket.Serve(address)
}
