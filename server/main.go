package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"time"

	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"

	"bitbucket.org/ehhio/ehhworldserver/server/game"
	"bitbucket.org/ehhio/ehhworldserver/server/gamemap"
	"bitbucket.org/ehhio/ehhworldserver/server/network"
)

var quiet bool
var level string
var width int
var height int
var mode int
var seed int64
var mapBlockSize int
var address string
var serveGame bool
var tick int

func init() {
	// Define input parameters
	flag.BoolVar(&quiet, "quiet", false, "whether to print any log statements during execution")
	flag.StringVar(&level, "level", "info", "the log level to output during execution. (e.g. 'panic', 'fatal', 'error', 'warn', 'info', or 'debug'")
	flag.IntVar(&width, "width", 512, "width of the game world")
	flag.IntVar(&height, "height", 512, "height of the game world")
	flag.IntVar(&mode, "mode", 0, "The map generator mode to use. 0 = 'noise', 1 = 'voronoi'. (default: 0)")
	flag.Int64Var(&seed, "seed", time.Now().UTC().UnixNano(), "World generation seed, defaults to random seed.")
	flag.IntVar(&mapBlockSize, "mapBlockSize", 4, "The size of blocks to break the game map into for transport.")
	flag.StringVar(&address, "address", ":8080", "The webserver address to listen on.")
	flag.BoolVar(&serveGame, "serve", false, "Start a game loop and run a webserver to serve the game world.")
	flag.IntVar(&tick, "tickrate", 60, "Times per second the game ticks and then updates players.")
}

func main() {
	// Parse inputs
	flag.Parse()
	generationMode := gamemap.NewGenerationMode(mode)

	// Seed pRNG
	rand.Seed(seed)

	// Configure logging
	if quiet {
		log.SetLevel(log.FatalLevel)
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(colorable.NewColorableStdout())
		level, err := log.ParseLevel(level)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)
		handler := func() {
			log.Error("Shit is fucked for some reason.")
		}
		log.RegisterExitHandler(handler)
	}

	// Generate the world
	gameMap := gamemap.NewGameMap(width, height, mapBlockSize, mapBlockSize)
	gameMap.Generate(generationMode, seed)

	// Serve, when prompted
	if serveGame {
		exitChan := make(chan int)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		go func() {
			<-signalChan
			exitChan <- 1
		}()

		// Start game and serve
		game := game.NewGame(gameMap)
		game.Start(tick)
		hub := network.Serve(address, game)

		// Wait for kill signal
		<-exitChan

		// Stop
		hub.Stop()
		game.Stop()
	}
}
