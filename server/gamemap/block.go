package gamemap

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/ehhio/ehhworldserver/server/utility"
)

// BlockMatrix A simple struct for a 2D matrix of Blocks and its size
type BlockMatrix struct {
	matrix [][]*Block    // 2D Matrix of Blocks
	size   *utility.Size // Size of the matrix
}

// Block A Block contains world cells. Useful for efficient updating of clients
type Block struct {
	size     *utility.Size     // The dimensions of the world block
	position *utility.Position // The position in the world ("World space coordinates")
	cells    [][]*Cell         // The 2D Cell matrix in the block
}

// NewBlock Creates a new Block of world Cells.
// w and h define the size of the block and its inner 2D matrix of world Cells
func NewBlock(w, h int) *Block {
	if w <= 0 || h <= 0 {
		log.Fatalf("Blocks must have dimensions greater than zero. Saw <%v, %v>.", w, h)
	}

	// Build CellMatrix
	matrix := make([][]*Cell, w)
	for i := range matrix {
		matrix[i] = make([]*Cell, h)
	}

	return &Block{
		size:  &utility.Size{Width: w, Height: w},
		cells: matrix,
	}
}

func (b *Block) String() string {
	return fmt.Sprintf("<Block>[%v, %v, Cells: %v]", b.position, b.size, b.cells)
}

// SetCells Set the Block's inner 2D Cell matrix to a new value
// New matrix must be the same size as the block
func (b *Block) SetCells(cells [][]*Cell) {
	newCellsWidth := len(cells)
	newCellsHeight := len(cells[0]) // note: assumes columns are all same height

	if newCellsWidth == b.size.Width && newCellsHeight == b.size.Height {
		b.cells = cells
	} else {
		log.Fatalf("Invalid Cell matrix shape for target block. Saw <%v, %v>. Expected %v.", newCellsWidth, newCellsHeight, b.size)
	}
}

// SetCellAt Set the Cell value within a world Block at a position
// x and y are Block space coordinates.
func (b *Block) SetCellAt(x, y int, cell *Cell) {
	if x >= 0 && x < b.size.Width && y >= 0 && y < b.size.Height {
		b.cells[x][y] = cell
	} else {
		log.Fatalf("Invalid Cell position in Block. Saw <%v, %v>. Block size was %v.", x, y, b.size)
	}
}

// GetCellAt Get the Cell from a Block at a position
// x and y are Block space coordinates.
func (b *Block) GetCellAt(x, y int) *Cell {
	targetCellPosition := utility.Position{X: x, Y: y}

	// Input validation
	if targetCellPosition.X < 0 || targetCellPosition.X >= b.size.Width || targetCellPosition.Y < 0 || targetCellPosition.Y >= b.size.Height {
		log.Fatalf("Invalid Cell position in Block. Saw %v. Block size was %v.", targetCellPosition, b.size)
	}

	return b.cells[targetCellPosition.X][targetCellPosition.Y]
}

// SetPosition Set the position of a world Block
// x and y are in World space coordinates.
func (b *Block) SetPosition(x, y int) {
	b.position = &utility.Position{X: x, Y: y}
}
