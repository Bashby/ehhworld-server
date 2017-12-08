// Package game is where the core game loop lives.
// A game tracks a world state, receives updates, and send updates to clients as needed.
package game

import (
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/ehhio/ehhworldserver/server/collision"
	"bitbucket.org/ehhio/ehhworldserver/server/gamemap"
)

// Constants
const millisecondPerUpdate = time.Millisecond * 8

// IGameObject interface deinfes what a struct must implemented for the game to simulate it in the game loop
type IGameObject interface {
	// Update computes the object state change from the current state to the next state, over a provided delta time value
	Update(int64, *Game)

	// Render interpolates a partial state change from the current state to the next state, over a provided delta time value, and emits the changes
	Render(int64, *Game)
}

// Game is a struct the represents a game
type Game struct {
	// Game loop
	running      bool
	frame        uint16
	startTime    time.Time
	previousTick time.Time
	accumulator  time.Duration
	ticker       *time.Ticker

	// Data managers
	gamemap   *gamemap.GameMap
	collision *collision.Collision
	objects   map[IGameObject]IGameObject
	// playerBlocks map[*gamemap.Block]*player.Player
}

// NewGame creates a new game
func NewGame(gamemap *gamemap.GameMap) *Game {
	return &Game{
		running:   false,
		frame:     0,
		gamemap:   gamemap,
		collision: collision.NewCollision(),
		objects:   make(map[IGameObject]IGameObject),
		// network: network,
	}
}

// Start begins the game loop
func (g *Game) Start(tickRate int) {
	log.Info("Starting game.")
	g.running = true
	g.startTime = time.Now()
	g.accumulator = time.Duration(0)
	g.previousTick = g.startTime
	g.ticker = time.NewTicker(time.Second / time.Duration(tickRate))

	go func() {
		for g.running {
			select {
			case t := <-g.ticker.C:
				// Increment frame number
				if g.frame == math.MaxUint16 {
					g.frame = 0
				} else {
					g.frame++
				}

				// Tick game
				g.tick(t)
			}
		}
	}()
}

// Stop ends the game loop
func (g *Game) Stop() {
	log.WithFields(log.Fields{
		"uptime":     time.Now().Sub(g.startTime),
		"totalFames": g.frame,
	}).Info("Stopping game.")

	g.running = false
	g.ticker.Stop()
}

// GetGameMap returns the game map
func (g *Game) GetGameMap() *gamemap.GameMap {
	return g.gamemap
}

// GetCollision returns the game collision system
func (g *Game) GetCollision() *collision.Collision {
	return g.collision
}

// // GetMapBlocksSize returns the size of the game map in blocks
// func (g *Game) GetMapBlocksSize() *utility.Size {
// 	return g.gamemap.GetBlocksSize()
// }

// AddObject adds an object to the game to be tracked / simulated.
func (g *Game) AddObject(object IGameObject) {
	g.objects[object] = object

	// If Trackable, track collision
	if trackableObject, ok := object.(collision.Trackable); ok {
		g.collision.AddObject(trackableObject)
	}
}

// RemoveObject removes an object from the game.
func (g *Game) RemoveObject(object IGameObject) {
	delete(g.objects, object)

	// If Trackable, remove from collision tracking
	if trackedObject, ok := object.(collision.Trackable); ok {
		g.collision.DeleteObject(trackedObject)
	}
	// for i := range g.players {
	// 	if g.players[i] == player {
	// 		// Remove the player from slice by index
	// 		// Does not preserve order of the slice. Is memory leak safe.
	// 		g.players[i] = g.players[len(g.players)-1]
	// 		g.players[len(g.players)-1] = nil
	// 		g.players = g.players[:len(g.players)-1]
	// 		break
	// 	}
	// }
}

func (g *Game) tick(currentTime time.Time) {
	elapsed := currentTime.Sub(g.previousTick)
	g.previousTick = currentTime
	g.accumulator += elapsed

	log.WithFields(log.Fields{
		"dt from last tick": elapsed,
		"accumulated dt":    g.accumulator,
	}).Debug("Game tick call.")

	for g.accumulator >= millisecondPerUpdate {
		g.update()
		g.accumulator -= millisecondPerUpdate
	}

	frac := float64(g.accumulator) / float64(millisecondPerUpdate)

	log.WithFields(log.Fields{
		"fraction of normal dt": frac,
		"accumulated dt":        g.accumulator,
	}).Debug("Game about to render.")

	g.render(time.Duration(float64(millisecondPerUpdate) * frac))
}

// update simulates the game using a consistent time step
func (g *Game) update() {
	log.WithFields(log.Fields{
		"dt to simulate": millisecondPerUpdate,
		"accumulated dt": g.accumulator,
	}).Debug("Game update call.")

	for _, o := range g.objects {
		o.Update(int64(millisecondPerUpdate/time.Millisecond), g)
	}
}

// render updates clients of current game state, interpolated dt into the future
func (g *Game) render(dt time.Duration) {
	log.WithFields(log.Fields{
		"dt to interpolate": dt,
	}).Debug("Game render call.")

	for _, o := range g.objects {
		o.Render(int64(dt/time.Millisecond), g)
	}
}
