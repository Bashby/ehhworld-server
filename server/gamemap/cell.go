package gamemap

import (
	"fmt"

	"bitbucket.org/ehhio/ehhworldserver/server/utility"
)

// CellShape The shape of a world Cell.
// Used for rendering.
type CellShape int

//go:generate stringer -type=CellShape

// CellShape enum
const (
	Full CellShape = iota
	TopHalf
	BottomHalf
	LeftHalf
	RightHalf
	TopLeftQuarter
	TopRightQuarter
	BottomLeftQuarter
	BottomRightQuarter
)

// Cell Smallest unit in the game map.
// Each cell has a biome, a shape, a parent Block that contains it, and a position in their parent.
type Cell struct {
	block    *Block            // The parent Block
	position *utility.Position // The position in the parent Block ("Block space coordinates")
	biome    *BiomeDefinition  // The biome of the cell
	shape    CellShape         // The shape of the cell
}

// NewCell Create a new game map Cell.
// x, y are Block space coordinates.
// Will automatically set itself in its parent as part of instantiation.
func NewCell(position utility.Position, biomeDef *BiomeDefinition, parent *Block, shape CellShape) (cell *Cell) {
	// Instantiate
	cell = &Cell{
		block:    parent,
		position: &position,
		biome:    biomeDef,
		shape:    shape,
	}

	// Register with parent block
	cell.block.SetCellAt(position.X, position.Y, cell)

	return
}

func (c *Cell) String() string {
	return fmt.Sprintf("<Cell>[%v, %v, %v]", c.position, c.biome, c.shape)
}
