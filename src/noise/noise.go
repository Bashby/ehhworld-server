package noise

import (
	"math"
	"math/rand"

	"github.com/tbogdala/noisey"
)

// type FBMOctave struct {
// 	frequency float64
// 	amplitude float64
// }

// type FBMNoiseGenerator struct {
// 	generator      *opensimplex.Noise
// 	seed           int64
// 	octaves        []*FBMOctave
// 	baseFrequency  float64
// 	baseAmplitude  float64
// 	lacunarity     float64
// 	gain           float64
// 	totalAmplitude float64
// }

// func NewFBMOctave(frequency float64, amplitude float64) *FBMOctave {
// 	return &FBMOctave{frequency, amplitude}
// }

// func NewFBMNoiseGenerator(octaveCount int, baseFrequency float64, baseAmplitude float64, lacunarity float64, gain float64) *FBMNoiseGenerator {
// 	// Parse inputs
// 	f := baseFrequency
// 	a := baseAmplitude
// 	l := lacunarity
// 	g := gain
// 	o := make([]*FBMOctave, octaveCount)

// 	// Seed the noise generator
// 	seed := time.Now().Unix()
// 	rand.Seed(seed)

// 	// Build octaves
// 	totalAmplitude := 0.0
// 	for i := range o {
// 		o[i] = NewFBMOctave(f, a)
// 		totalAmplitude += a
// 		f *= l
// 		a *= g
// 	}

// 	return &FBMNoiseGenerator{opensimplex.NewWithSeed(rand.Int63()), seed, o, baseFrequency, baseAmplitude, l, g, totalAmplitude}
// }

// // ReSeed Seed the generator with a new random seed
// func (n *FBMNoiseGenerator) ReSeed() {
// 	newSeed := rand.Int63()
// 	n.generator = opensimplex.NewWithSeed(newSeed)
// 	n.seed = newSeed
// }

// // GenerateFBMNoise Generate a noise value using FBM
// func (n *FBMNoiseGenerator) GenerateFBMNoise(x int, y int) float64 {
// 	noise := 0.0
// 	for _, octave := range n.octaves {
// 		noise += MapNoise(n.generator.Eval2(float64(x)*octave.frequency, float64(y)*octave.frequency)) * octave.amplitude
// 	}

// 	// Average noise values across octaves
// 	noise /= n.totalAmplitude

// 	return noise
// }

// // BuildNoiseMatrix Builds a 2d matrix of given size filled with simplex noise and Fractal Brownian Motion
// func (n *FBMNoiseGenerator) BuildNoiseMatrix(w int, h int, min int, max int) [][]float64 {
// 	data := make([][]float64, w)
// 	for x := 0; x < w; x++ {
// 		data[x] = make([]float64, h)
// 		for y := 0; y < h; y++ {
// 			noise := n.GenerateFBMNoise(x, y)

// 			// fmt.Println(noise)

// 			// Normalize noise [-1.0, 1.0] to requested range [min, max]
// 			noise = float64((max-min))*((noise+0.7)/1.4) + float64(min)

// 			//noise = noise*float64(max-min)/2 + float64(max+min)/2

// 			data[x][y] = noise
// 		}
// 	}

// 	return data
// }

// func MapNoise(noise float64) float64 {
// 	return noise/2.0 + 0.5
// }

// Octaves        int                // the number of octaves in each fbm calculation
// Frequency      float64            // the starting frequency, decays along lacunarity
// lacunarity     float64            // used to decay the frequency across each octave
// gain           float64            // used to decay the amplitude across each octave

type FBMNoiseGenerator2D struct {
	Generator          *noisey.Scale2D // the interface to generate basic noise
	Redistributed      bool
	RedistributeFactor float64
	Terraced           bool
	TerraceFactor      float64
}

func NewFBMNoiseGenerator2D(octaveCount int, persistence float64, lacunarity float64, frequency float64, scale float64, bias float64, min float64, max float64, isRedistributed bool, redistributeFactor float64, isTerraced bool, terraceFactor float64) *FBMNoiseGenerator2D {
	randGen := rand.New(rand.NewSource(rand.Int63()))
	nGen := noisey.NewOpenSimplexGenerator(randGen)
	fbm := noisey.NewFBMGenerator2D(&nGen, octaveCount, persistence, lacunarity, frequency)
	sFbm := noisey.NewScale2D(&fbm, scale, bias, min, max)
	return &FBMNoiseGenerator2D{&sFbm, isRedistributed, redistributeFactor, isTerraced, terraceFactor}
}

// BuildNoiseMatrix Builds a 2d matrix of given size filled with simplex noise and Fractal Brownian Motion
func (n *FBMNoiseGenerator2D) BuildNoiseMatrix(w int, h int, min float64, max float64) [][]float64 {
	b := noisey.NewBuilder2D(n.Generator, w, h)
	b.Bounds = noisey.Builder2DBounds{0.0, 0.0, float64(w), float64(h)}
	b.Build()

	data := make([][]float64, w)
	for x := 0; x < w; x++ {
		data[x] = make([]float64, h)
		for y := 0; y < h; y++ {
			noise := n.Generator.Get2D(float64(x), float64(y))

			// Redistribute output
			if n.Redistributed {
				noise = math.Pow(noise, n.RedistributeFactor)
			}

			// Terrace output
			if n.Terraced {
				noise = float64(RoundToInt(noise*n.TerraceFactor)) / n.TerraceFactor
			} else {

			}

			// Normalize noise [0.0, 1.0] to requested range [min, max]
			noise = (max-min)*noise + min

			data[x][y] = noise
		}
	}

	// fmt.Println(data)

	return data
}

// RoundToInt8 rounds 64-bit floats into integer numbers
func RoundToInt(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}
