package player

import (
	"bitbucket.org/ehhio/ehhworldserver/server/collision"
	"bitbucket.org/ehhio/ehhworldserver/server/game"
	"bitbucket.org/ehhio/ehhworldserver/server/object"
	"bitbucket.org/ehhio/ehhworldserver/server/player/minimap"
	"bitbucket.org/ehhio/ehhworldserver/server/utility"
)

const playerCollisionSize = 0.1

// Verify that Player implements Trackable
var _ collision.Trackable = &Player{}

// Verify that *Player implements Trackable
var _ collision.Trackable = (*Player)(nil)

// Verify that Player implements IGameObject
var _ game.IGameObject = &Player{}

// Verify that *Player implements IGameObject
var _ game.IGameObject = (*Player)(nil)

// Player is a struct defining a game object representing a remote game client
type Player struct {
	name       string
	position   *utility.PositionHighResolution
	Dirty      DirtyFlagsBitSet
	minimap    *minimap.Minimap
	objectType *object.TypeFlagsBitSet
}

// NewPlayer creates a new player object.
// Name names the player. Position sets the player's position. Size sets the size of the player's map.
func NewPlayer(name string, position *utility.PositionHighResolution, size *utility.Size) *Player {
	flags := object.NewTypeFlagsBitSet(object.FlagPlayer)
	minimap := minimap.NewMinimap(size)
	return &Player{
		name:       name,
		position:   position,
		Dirty:      DirtyFlagsBitSet{},
		minimap:    minimap,
		objectType: flags,
	}
}

// func (p *Player) IsDirty() bool {
// 	return p.dirty.position || p.dirty.orientation || p.dirty.health
// }

// func (p *Player) IsPositionDirty() bool {
// 	return p.dirty.position
// }

// func (p *Player) IsOrientationDirty() bool {
// 	return p.dirty.orientation
// }

// // IsHealthDirty returns true if the player health has channged since last update
// func (p *Player) IsHealthDirty() bool {
// 	return p.dirty.health
// }

// // SetMinimap sets the game minimap for a Player
// func (p *Player) SetMinimap(minimap *minimap.Minimap) {
// 	p.minimap = minimap
// }

// Update updates the player using a consistent time step
func (p *Player) Update(dt int64, g *game.Game) {

}

// Render renders the player to clients, interpolating state into the future
func (p *Player) Render(dt int64, g *game.Game) {

}

// GetAABBBottomLeftPoint returns the bottom left point of the Player collision AABB, implementing collision.Trackable
func (p *Player) GetAABBBottomLeftPoint() *utility.PositionHighResolution {
	return &utility.PositionHighResolution{X: p.position.X - playerCollisionSize/2.0, Y: p.position.Y - playerCollisionSize/2.0}
}

// GetAABBSize returns the width and height dimensions of the Player collision AABB, implementing collision.Trackable
func (p *Player) GetAABBSize() *utility.SizeHighResolution {
	return &utility.SizeHighResolution{Width: playerCollisionSize, Height: playerCollisionSize}
}

// GetObjectFlags returns the object type bit flags of the Player, implementing collision.Trackable
func (p *Player) GetObjectFlags() *object.TypeFlagsBitSet {
	return p.objectType
}
