package utility

import (
	"fmt"
	"math/rand"
)

// Position is a point in a 2D coordinate system, using ints
// TODO: Replace with native golang/geo lib equivalents
type Position struct {
	X, Y int
}

func (p *Position) String() string {
	return fmt.Sprintf("<Position>[x: %v, y: %v]", p.X, p.Y)
}

// IsWithinBounds returns true if the position is contained by the bounds
func (p *Position) IsWithinBounds(bounds *Size) bool {
	if p.X < 0 || p.X >= bounds.Width || p.Y < 0 || p.Y >= bounds.Height {
		return false
	}
	return true
}

// PositionHighResolution is a point in a 2D coordinate system, using floats
type PositionHighResolution struct {
	X, Y float64
}

func (p *PositionHighResolution) String() string {
	return fmt.Sprintf("<PositionHighResolution>[x: %v, y: %v]", p.X, p.Y)
}

// IsWithinBounds returns true if the position is contained by the bounds; uses HighResolution float64s
func (p *PositionHighResolution) IsWithinBounds(bounds *SizeHighResolution) bool {
	if p.X < 0 || p.X >= bounds.Width || p.Y < 0 || p.Y >= bounds.Height {
		return false
	}
	return true
}

// Size describes an object's 2D width and height, using ints
type Size struct {
	Width, Height int
}

func (s *Size) String() string {
	return fmt.Sprintf("<Size>[w: %v, h: %v]", s.Width, s.Height)
}

// SizeHighResolution describes an object's 2D width and height, using float64s
type SizeHighResolution struct {
	Width, Height float64
}

func (s *SizeHighResolution) String() string {
	return fmt.Sprintf("<SizeHighResolution>[w: %v, h: %v]", s.Width, s.Height)
}

// BooleanMatrix is a wrapper for a 2D matrix of bools and its size
type BooleanMatrix struct {
	Matrix [][]bool
	Size   *Size
}

// NewBooleanMatrix constructs an empty 2D matrix of requested size
func NewBooleanMatrix(size *Size) *BooleanMatrix {
	matrix := make([][]bool, size.Width)
	for i := range matrix {
		matrix[i] = make([]bool, size.Height)
	}

	return &BooleanMatrix{
		Matrix: matrix,
		Size:   size,
	}
}

func (m *BooleanMatrix) String() string {
	return fmt.Sprintf("<BooleanMatrix>[%v, %v]", m.Size, m.Matrix)
}

// RandomFloat64InRange returns a float64 randomly within the min and max values passed
func RandomFloat64InRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// RandomIntInRange returns a int randomly within the min and max values passed
func RandomIntInRange(min, max int) int {
	return rand.Int()*(max-min) + min
}

// GeneratePlayerName creates a placeholder name for a player
func GeneratePlayerName() string {
	return fmt.Sprintf("anon%v", rand.Intn(999))
}

// Clamp restricts an int to a specified int range, inclusive
func Clamp(val, min, max int) int {
	if val > max {
		return max
	} else if val < min {
		return min
	}

	return val
}

// ClampHighResolution restricts a float64 to a specified float64 range, inclusive
func ClampHighResolution(val, min, max float64) float64 {
	if val > max {
		return max
	} else if val < min {
		return min
	}

	return val
}
