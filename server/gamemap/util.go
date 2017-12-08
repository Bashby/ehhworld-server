package gamemap

import (
	"image"
	"image/color"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/llgcode/draw2d"
	"github.com/pzsz/voronoi"
)

// GenerationMode defines the type of world generation method to use
type GenerationMode int

//go:generate stringer -type=GenerationMode

const (
	// Noise defines a generation mode using noise augmented with fractal brownian motion (fbm)
	Noise GenerationMode = iota
	// Voronoi defines a generation mode using a voronoi diagram augmented with noise
	Voronoi
)

// NewGenerationMode Creates a new generation mode based on an integer type
func NewGenerationMode(value int) (mode GenerationMode) {
	switch value {
	case 0:
		mode = Noise
	case 1:
		mode = Voronoi
	default:
		log.Fatalln("Unrecognized world generation mode, saw: ", value)
	}

	return
}

func determineBiome(biomes image.Image, elevation int, moisture int, fuzz bool, fuzzFactor float64) (x, y int, biome color.Color) {
	y = biomes.Bounds().Max.Y - elevation - 1
	if fuzz {
		y += roundToInt(rand.Float64()*(2*fuzzFactor) - fuzzFactor)
	}
	if y < 0 {
		y = 0
	} else if y > biomes.Bounds().Max.Y-1 {
		y = biomes.Bounds().Max.Y - 1
	}
	x = moisture
	if fuzz {
		x += roundToInt(rand.Float64()*(2*fuzzFactor) - fuzzFactor)
	}
	if x < 0 {
		x = 0
	} else if x > biomes.Bounds().Max.X-1 {
		x = biomes.Bounds().Max.X - 1
	}
	biome = biomes.At(x, y)
	return
}

// RoundToInt rounds 64-bit floats into integer numbers
func roundToInt(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

func fontNamer(fontData draw2d.FontData) string {
	fontFileName := fontData.Name + "/" + fontData.Name
	if fontData.Style&draw2d.FontStyleBold != 0 {
		fontFileName += "-Bold"
	} else if fontData.Style&draw2d.FontStyleItalic != 0 {
		fontFileName += "-Italic"
	} else {
		fontFileName += "-Regular"
	}
	fontFileName += ".ttf"
	return fontFileName
}

func averageColors(c1, c2 color.RGBA) color.RGBA {
	return color.RGBA{
		R: (c1.R + c2.R) / 2,
		G: (c1.G + c2.G) / 2,
		B: (c1.B + c2.B) / 2,
		A: 0xFF,
	}
}

type colorSampling struct {
	color color.Color
	point image.Point
}

type diagram struct {
	*voronoi.Diagram
	Center voronoi.Vertex
}
