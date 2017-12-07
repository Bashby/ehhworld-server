package world

import (
	"image"
	"image/png"
	"log"
	"os"

	"bitbucket.org/ashbyb/go-map-generation/noise"

	"image/color"
)

func Generate(w int, h int) {
	// Generate Noise
	noiseGeneratorA := noise.NewFBMNoiseGenerator2D(16, 0.5, 2.0, 0.007, 0.5, 0.5, 0.0, 1.0, true, 3.0, false, 25.0)
	noiseGeneratorB := noise.NewFBMNoiseGenerator2D(16, 0.5, 2.0, 0.007, 0.5, 0.5, 0.0, 1.0, true, 1.0, true, 10.0)
	elevationNoise := noiseGeneratorA.BuildNoiseMatrix(w, h, 0.0, 255.0)
	moistureNoise := noiseGeneratorB.BuildNoiseMatrix(w, h, 0.0, 255.0)

	// Create an image from target size
	img := image.NewRGBA(image.Rect(0, 0, w*2, h*2))

	// Load in biome data
	reader, err := os.Open("world/biomes.v2.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	biomeData, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Render out the biome debug square
	debugOffset := 50
	for x := 0; x < 256; x++ {
		// Horizontal Lines
		img.Set(w+x+debugOffset, 0+debugOffset, color.RGBA{255, 0, 0, 255})
		img.Set(w+x+debugOffset, 255+debugOffset, color.RGBA{255, 0, 0, 255})
	}
	for y := 0; y < 256; y++ {
		// Vertical Lines
		img.Set(w+debugOffset, y+debugOffset, color.RGBA{255, 0, 0, 255})
		img.Set(w+255+debugOffset, y+debugOffset, color.RGBA{255, 0, 0, 255})
	}

	// Generate output
	for i := range elevationNoise {
		for j := range elevationNoise[i] {
			elevation := RoundToInt(elevationNoise[i][j])
			moisture := RoundToInt(moistureNoise[i][j])
			biome := DetermineBiome(biomeData, elevation, moisture)
			img.Set(i, j, biome)
			img.Set(i, j+h, color.RGBA{uint8(elevation), uint8(elevation), uint8(elevation), 255})
			img.Set(i+w, j+h, color.RGBA{uint8(moisture), uint8(moisture), uint8(moisture), 255})
			// fmt.Println("M: ", moisture, "E: ", elevation, "C: ", biome)
			// Debug Biome
			img.Set(w+moisture+debugOffset, 255-elevation+debugOffset, color.White)
		}
	}

	// Save to out.png
	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	png.Encode(f, img)
}

func DetermineBiome(biomes image.Image, elevation int, moisture int) color.Color {
	y := biomes.Bounds().Max.Y - elevation - 1
	if y < 0 {
		y = 0
	}
	x := moisture
	if x > biomes.Bounds().Max.X {
		x = biomes.Bounds().Max.X
	}
	return biomes.At(x, y)
}

// RoundToInt8 rounds 64-bit floats into integer numbers
func RoundToInt(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}
