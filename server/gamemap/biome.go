package gamemap

import (
	"fmt"
	"image/color"
)

// BiomeDefinition A simple struct for a Biome and its color.
type BiomeDefinition struct {
	biome Biome
	color color.Color
}

func (bd *BiomeDefinition) String() string {
	return fmt.Sprintf("<BiomeDefinition>[%v, color: %v]", bd.biome, bd.color)
}

// Biome An enum for the type of biomes in the world
type Biome int

//go:generate stringer -type=Biome

// World Biome Types enum
const (
	SaltwaterDeep Biome = iota
	SaltwaterShallow
	FreshwaterDeep
	FreshwaterShallow
	Shore
	SubtropicalDesert
	Grassland
	TropicalSeasonalForest
	TropicalRainForest
	TemperateDesert
	TemperateDeciduousForest
	TemperateRainForest
	Shrubland
	Taiga
	Scorched
	Bare
	Tundra
	Snow
)

// BiomePalette World biome color palette
var BiomePalette color.Palette = []color.Color{
	color.RGBA{0xFF, 0x00, 0xFF, 0x00}, // SaltwaterDeep UNUSED
	color.RGBA{0xFF, 0x00, 0xFF, 0x00}, // SaltwaterShallow UNUSED
	color.RGBA{59, 58, 105, 0xFF},      // FreshwaterDeep
	color.RGBA{66, 105, 184, 0xFF},     // FreshwaterShallow
	color.RGBA{160, 144, 119, 0xFF},    // Shore
	color.RGBA{210, 185, 139, 0xff},    // SubtropicalDesert
	color.RGBA{136, 170, 85, 0xff},     // Grassland
	color.RGBA{85, 153, 68, 0xff},      // TropicalSeasonalForest
	color.RGBA{51, 119, 85, 0xff},      // TropicalRainForest
	color.RGBA{201, 210, 155, 0xFF},    // TemperateDesert
	color.RGBA{103, 148, 89, 0xFF},     // TemperateDeciduousForest
	color.RGBA{68, 136, 85, 0xFF},      // TemperateRainForest
	color.RGBA{136, 153, 119, 0xFF},    // Shrubland
	color.RGBA{153, 170, 119, 0xFF},    // Taiga
	color.RGBA{85, 85, 85, 0xFF},       // Scorched
	color.RGBA{136, 136, 136, 0xFF},    // Bare
	color.RGBA{187, 187, 170, 0xFF},    // Tundra
	color.RGBA{221, 221, 228, 0xFF},    // Snow
	// color.RGBA{0x00, 0x00, 0x00, 0xFF}, // Black
	// color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White
}
