package gamemap

import (
	"image"
	"image/color"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/ehhio/ehhworldserver/server/noise"
	"bitbucket.org/ehhio/ehhworldserver/server/utility"

	"golang.org/x/image/draw"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
)

// GameMap A map made up of blocks containing cells of biome data.
type GameMap struct {
	size      *utility.Size // Size of the game map in cells
	seed      int64         // pRNG seed used to generate the game map
	blocks    *BlockMatrix  // 2D matrix containing Blocks
	blockSize *utility.Size // Size of Blocks in the map
}

// NewGameMap creates a new empty game map.
// w and h define the size of the game map in cells.
// bW and bH define the size of blocks that contain the cells. The block size must be a multiple of the game map size.
// To populate the game map with data, use the Generate() method.
func NewGameMap(w, h, bW, bH int) *GameMap {
	mapSize := utility.Size{Width: w, Height: h}
	blockSize := utility.Size{Width: bW, Height: bH}
	matrixSize := utility.Size{Width: w / bW, Height: h / bH}

	// Input validation
	if mapSize.Width < blockSize.Width || mapSize.Height < blockSize.Height {
		log.Fatalf("Block size may not be larger than map size. Block size: %v, Map size: %v", blockSize, mapSize)
	} else if mapSize.Width%blockSize.Width != 0 || mapSize.Height%blockSize.Height != 0 {
		log.Fatalf("Block size must divide evenly into map size. Block size: %v, Map size: %v.", blockSize, mapSize)
	}

	// Build matrix struct
	matrix := make([][]*Block, matrixSize.Width)
	for i := range matrix {
		matrix[i] = make([]*Block, matrixSize.Height)
	}

	return &GameMap{
		size:      &mapSize,
		blocks:    &BlockMatrix{matrix: matrix, size: &matrixSize},
		blockSize: &blockSize,
	}
}

// GetSize returns the size of the game map in cells
func (m *GameMap) GetSize() *utility.Size {
	return m.size
}

// GetBlocksSize returns the size of the game map in blocks
func (m *GameMap) GetBlocksSize() *utility.Size {
	return m.blocks.size
}

// GetSeed returns the seed used to generate the game map
func (m *GameMap) GetSeed() int64 {
	return m.seed
}

// GetBlockSize returns the size of a single block in the game map
func (m *GameMap) GetBlockSize() *utility.Size {
	return m.blockSize
}

// SetBlocks Set the 2D Blocks matrix in the map
// New matrix must be same dimensions as existing matrix
func (m *GameMap) SetBlocks(blocks [][]*Block) {
	newBlocksSize := utility.Size{Width: len(blocks), Height: len(blocks[0])} // note: assumes columns are all same height

	if newBlocksSize.Width == m.blocks.size.Width && newBlocksSize.Height == m.blocks.size.Height {
		m.blocks.matrix = blocks
	} else {
		log.Fatalf("Invalid Block matrix shape. Saw %v. Expected %v.", newBlocksSize, m.blocks.size)
	}
}

// SetBlockAt Set the Block in the game map at a position
// x and y are map space coordinates.
func (m *GameMap) SetBlockAt(x, y int, block *Block) {
	targetWorldPosition := utility.Position{X: x, Y: y}
	newBlockPosition := m.mapToBlockCoordinates(targetWorldPosition)

	// Input Validation
	if newBlockPosition.X < 0 || newBlockPosition.X >= m.blocks.size.Width || newBlockPosition.Y < 0 || newBlockPosition.Y >= m.blocks.size.Height {
		log.Fatalf("Invalid Block position in Map. Tried %v (parsed from map coordinates %v), Block matrix size was %v.", newBlockPosition, targetWorldPosition, m.blocks.size)
	}

	// Set Block
	m.blocks.matrix[newBlockPosition.X][newBlockPosition.Y] = block
}

// GetBlockAt Get the Block in the game map at a position
// x and y are map space coordinates
func (m *GameMap) GetBlockAt(x, y int) *Block {
	targetWorldPosition := utility.Position{X: x, Y: y}
	targetBlockPosition := m.mapToBlockCoordinates(targetWorldPosition)

	if targetBlockPosition.X < 0 || targetBlockPosition.X >= m.blocks.size.Width || targetBlockPosition.Y < 0 || targetBlockPosition.Y >= m.blocks.size.Height {
		log.Fatalf("Invalid Block position in Map. Tried %v (parsed from Map coordinates %v), Block matrix size was %v.", targetBlockPosition, targetWorldPosition, m.blocks.size)
	}

	return m.blocks.matrix[targetBlockPosition.X][targetBlockPosition.Y]
}

// GetCellAt Get the Cell from the Block in the game Map at a position
// x and y are Map space coordinates
func (m *GameMap) GetCellAt(x, y int) *Cell {
	targetWorldPosition := utility.Position{X: x, Y: y}
	targetCellPosition := m.mapToCellCoordinates(targetWorldPosition)

	targetBlock := m.GetBlockAt(targetWorldPosition.X, targetWorldPosition.Y)

	if targetBlock == nil {
		return nil
	}

	return targetBlock.GetCellAt(targetCellPosition.X, targetCellPosition.Y)
}

// Generate populates the game map blocks with biome data.
// mode determines the generation mode to use.
// seed is used to initialize any pRNG functionality during generation.
func (m *GameMap) Generate(mode GenerationMode, seed int64) {
	log.WithFields(log.Fields{
		"mode": mode,
		"seed": seed,
		"size": m.size,
	}).Info("Generating world.")

	// Record seed
	m.seed = seed

	// Init biome data
	biomeImg := prepareBiomeData(*m.size, "../assets/image/biomes.v5.png")

	// Init noise data
	elevationNoise, moistureNoise, elevationImg, moistureImg := prepareNoiseData(*m.size)

	// Generate World using requested method
	var worldImg *image.RGBA
	switch mode {
	case Noise:
		worldImg = noiseWorldGeneration(*m.size, elevationNoise, moistureNoise, biomeImg)
	case Voronoi:
		worldImg = voronoiWorldGeneration(*m.size, elevationNoise, moistureNoise, biomeImg)
	default:
		log.Fatalf("Invalid Map generation mode. Saw %v.", mode)
	}

	// Populate map blocks
	m.populateBlocksFromImage(worldImg)

	// Generate output image
	prepareOutputImage(*m.size, worldImg, biomeImg, elevationImg, moistureImg)
}

func (m *GameMap) populateBlocksFromImage(image *image.RGBA) {
	// TODO(bashby) Is parsing a 2D grid into a isometric grid enough of a problem to be worth fixing?

	imageSize := image.Rect.Size()

	// Input Validation
	if imageSize.X > m.size.Width || imageSize.Y > m.size.Height {
		log.Fatalf("Insufficient room to populate. Image is larger than game map. Image %v. Map %v.", imageSize, m.size)
	}

	// Partition
	for x := 0; x < imageSize.X; x++ {
		for y := 0; y < imageSize.Y; y++ {
			// Determine the target block
			block := m.GetBlockAt(x, y)
			if block == nil {
				block = NewBlock(m.blockSize.Width, m.blockSize.Height)
				block.SetPosition(x, y)
				m.SetBlockAt(x, y, block)
			}

			// Determine biome for cell
			pixel := image.At(x, y)
			biome := &BiomeDefinition{
				biome: Biome(BiomePalette.Index(pixel)),
				color: BiomePalette.Convert(pixel),
			}

			// Create cell
			cellPosition := m.mapToCellCoordinates(utility.Position{X: x, Y: y})
			NewCell(cellPosition, biome, block, Full)
		}
	}
}

// ToImage outputs game map data to a png image
func (m *GameMap) ToImage() {
	// Create an output target
	img := image.NewRGBA(image.Rect(0, 0, m.size.Width, m.size.Height))

	for x := 0; x < m.size.Width; x++ {
		for y := 0; y < m.size.Height; y++ {
			pixel := m.GetCellAt(x, y)
			img.Set(x, y, pixel.biome.color)
		}
	}

	// Save image
	draw2dimg.SaveToPngFile("../assets/image/generated/out_world.png", img)
}

func (m *GameMap) mapToBlockCoordinates(mapPos utility.Position) utility.Position {
	if mapPos.X < 0 || mapPos.X >= m.size.Width || mapPos.Y < 0 || mapPos.Y >= m.size.Height {
		log.Fatalf("Map position is out of bounds. Tried %v, Map size was %v.", mapPos, m.size)
	}
	return utility.Position{X: int(mapPos.X / m.blockSize.Width), Y: int(mapPos.Y / m.blockSize.Height)}
}

func (m *GameMap) mapToCellCoordinates(mapPos utility.Position) utility.Position {
	if mapPos.X < 0 || mapPos.X >= m.size.Width || mapPos.Y < 0 || mapPos.Y >= m.size.Height {
		log.Fatalf("Map position is out of bounds. Tried %v, Map size was %v.", mapPos, m.size)
	}
	return utility.Position{X: mapPos.X % m.blockSize.Width, Y: mapPos.Y % m.blockSize.Height}
}

// RandomPositionHighResolution returns a random point in the map using float64s
func (m *GameMap) RandomPositionHighResolution() *utility.PositionHighResolution {
	return &utility.PositionHighResolution{
		X: utility.RandomFloat64InRange(float64(0), float64(m.GetSize().Width)),
		Y: utility.RandomFloat64InRange(float64(0), float64(m.GetSize().Height)),
	}
}

// RandomPosition returns a random point in the map using ints
func (m *GameMap) RandomPosition() *utility.Position {
	return &utility.Position{
		X: utility.RandomIntInRange(0, m.GetSize().Width),
		Y: utility.RandomIntInRange(0, m.GetSize().Height),
	}
}

// RandomCell returns a random cell in the map
func (m *GameMap) RandomCell() *Cell {
	randPos := m.RandomPosition()

	return m.GetCellAt(randPos.X, randPos.Y)
}

// RandomBlock returns a random block in the map
func (m *GameMap) RandomBlock() *Block {
	randPos := m.RandomPosition()

	return m.GetBlockAt(randPos.X, randPos.Y)
}

func noiseWorldGeneration(size utility.Size, elevationNoise, moistureNoise [][]float64, biomeImg *image.RGBA) (worldImg *image.RGBA) {

	// Init world output image
	worldImg = image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))

	// Create a copy of the incoming biome data
	biomeImgTmp := image.NewRGBA(biomeImg.Rect)
	draw.Draw(biomeImgTmp, biomeImg.Rect, biomeImg, biomeImg.Rect.Min, draw.Src)

	// Generate
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			elevation := elevationNoise[x][y]
			moisture := moistureNoise[x][y]
			sampledX, sampledY, biome := determineBiome(
				biomeImgTmp,
				roundToInt(elevation),
				roundToInt(moisture),
				true, // isFuzzed
				15.0, // fuzzFactor
			)
			worldImg.Set(x, y, biome)
			biomeImg.Set(sampledX, sampledY, color.RGBA{255, 0, 0, 255})
		}
	}

	return
}

func voronoiWorldGeneration(size utility.Size, elevationNoise, moistureNoise [][]float64, biomeImg *image.RGBA) (worldImg *image.RGBA) {

	// Generate images to hold data
	worldImg = image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))

	// Leverage noise to sample biome data for use later
	biomeSampledColorData := make([][]*colorSampling, size.Width)
	for x := 0; x < size.Width; x++ {
		biomeSampledColorData[x] = make([]*colorSampling, size.Height)
		for y := 0; y < size.Height; y++ {
			elevation := elevationNoise[x][y]
			moisture := moistureNoise[x][y]
			sampledX, sampledY, biome := determineBiome(
				biomeImg,
				roundToInt(elevation),
				roundToInt(moisture),
				true, // isFuzzed
				15.0, // fuzzFactor
			)
			biomeSampledColorData[x][y] = &colorSampling{color: biome, point: image.Point{X: sampledX, Y: sampledY}}
		}
	}

	// Compute voronoi diagram
	bbox := voronoi.NewBBox(0, float64(size.Width), 0, float64(size.Height))
	sites := utils.RandomSites(bbox, 15000)
	d := voronoi.ComputeDiagram(sites, bbox, true)

	// Relax using Lloyd's algorithm
	relaxationIterations := 1 // TODO: put world generation params into a config JSON
	for i := 0; i < relaxationIterations; i++ {
		sites = utils.LloydRelaxation(d.Cells)
		d = voronoi.ComputeDiagram(sites, bbox, true)
	}

	// Create root diagram
	center := voronoi.Vertex{X: float64(size.Width / 2), Y: float64(size.Height / 2)}
	diagram := &diagram{d, center}

	// Create drawing context
	draw := draw2dimg.NewGraphicContext(worldImg)
	draw.SetLineWidth(2.0)

	// Iterate over cells
	for _, cell := range diagram.Cells {

		// Draw cell edges path
		atFirstPoint := true
		for _, hedge := range cell.Halfedges {
			a := hedge.GetStartpoint()
			b := hedge.GetEndpoint()

			if atFirstPoint {
				atFirstPoint = false
				draw.MoveTo(a.X, a.Y)
			}
			draw.LineTo(b.X, b.Y)
		}
		draw.Close()

		// Color cell using centroid in biome
		center := utils.CellCentroid(cell)
		sampling := biomeSampledColorData[utility.Clamp(roundToInt(center.X), 0, size.Width)][utility.Clamp(roundToInt(center.Y), 0, size.Height)] // TODO: Fixed? ... THIS CAN GO OUT OF BOUND FOR SOME REASON...
		cellColor := sampling.color
		draw.SetFillColor(cellColor)
		draw.SetStrokeColor(cellColor)
		draw.FillStroke()

		// Draw centroid for cell
		// drawCentroid := draw2dimg.NewGraphicContext(worldImg)
		// drawCentroid.SetLineWidth(0.2)
		// // draw2.MoveTo(center.X, center.Y)
		// draw2dkit.Circle(drawCentroid, center.X, center.Y, 0.5)
		// drawCentroid.SetStrokeColor(color.RGBA{0x00, 0x00, 0xff, 0xff})
		// drawCentroid.Stroke()

		// Mark biome data with biome sampling event
		biomeImg.Set(sampling.point.X, sampling.point.Y, color.RGBA{255, 0, 0, 255})
	}

	return
}

func prepareBiomeData(targetSize utility.Size, biomeDataPath string) (biomeImgScaled *image.RGBA) {
	// Load Biome data
	//absPath, _ := filepath.Abs(biomeDataPath)
	biomeImg, err := draw2dimg.LoadFromPngFile(biomeDataPath)
	if err != nil {
		log.Fatalf("Failed to load biome data from file %v. Saw %v", biomeDataPath, err)
	}

	// Transform biome mapping to target size of the world and noise data
	biomeImgSize := biomeImg.Bounds().Size()
	transformMatrix := draw2d.NewMatrixFromRects(
		[4]float64{0, 0, float64(biomeImgSize.X), float64(biomeImgSize.Y)},
		[4]float64{0, 0, float64(targetSize.Width), float64(targetSize.Height)},
	)
	biomeImgScaled = image.NewRGBA(image.Rect(0, 0, targetSize.Width, targetSize.Height))
	draw2dimg.DrawImage(biomeImg, biomeImgScaled, transformMatrix, draw.Src, draw2dimg.BicubicFilter)

	return
}

func prepareOutputImage(size utility.Size, world, biomes, elevation, moisture *image.RGBA) {

	// Create an output target
	img := image.NewRGBA(image.Rect(0, 0, size.Width*2, size.Height*2))

	// Determine translations for output
	baseMatrix := draw2d.NewIdentityMatrix()
	biomeTranslationMatrix := baseMatrix.Copy()
	elevationTranslationMatrix := baseMatrix.Copy()
	moistureTranslationMatrix := baseMatrix.Copy()
	biomeTranslationMatrix.Translate(float64(size.Width), 0)
	elevationTranslationMatrix.Translate(0, float64(size.Height))
	moistureTranslationMatrix.Translate(float64(size.Width), float64(size.Height))

	// Render parts into whole
	draw2dimg.DrawImage(world, img, baseMatrix, draw.Src, draw2dimg.BicubicFilter)
	draw2dimg.DrawImage(biomes, img, biomeTranslationMatrix, draw.Src, draw2dimg.BicubicFilter)
	draw2dimg.DrawImage(elevation, img, elevationTranslationMatrix, draw.Src, draw2dimg.BicubicFilter)
	draw2dimg.DrawImage(moisture, img, moistureTranslationMatrix, draw.Src, draw2dimg.BicubicFilter)

	// Draw text labels
	draw := draw2dimg.NewGraphicContext(img)
	draw2d.SetFontFolder("../assets/font")
	draw2d.SetFontNamer(fontNamer)
	draw.SetLineWidth(4)
	draw.SetFontData(draw2d.FontData{Name: "Rubik"})
	draw.SetFillColor(color.RGBA{255, 255, 255, 255})
	draw.SetStrokeColor(color.RGBA{0, 0, 0, 255})
	draw.SetFontSize(18)
	draw.StrokeStringAt("Generated World", 10, float64(size.Height)-10)
	draw.FillStringAt("Generated World", 10, float64(size.Height)-10)
	draw.StrokeStringAt("Elevation Noise", 10, float64(2*size.Height)-10)
	draw.FillStringAt("Elevation Noise", 10, float64(2*size.Height)-10)
	draw.StrokeStringAt("Moisture Noise", 10+float64(size.Width), float64(size.Height*2)-10)
	draw.FillStringAt("Moisture Noise", 10+float64(size.Width), float64(size.Height*2)-10)
	draw.StrokeStringAt("Biome Mapping", 10+float64(size.Width), float64(size.Height)-10)
	draw.FillStringAt("Biome Mapping", 10+float64(size.Width), float64(size.Height)-10)

	// Save image
	err := draw2dimg.SaveToPngFile("../assets/image/generated/out.png", img)
	if err != nil {
		log.Fatalf("Failed to save generate world map. Saw %v", err)
	}
}

func prepareNoiseData(targetSize utility.Size) (elevationNoise, moistureNoise [][]float64, elevationImg, moistureImg *image.RGBA) {

	// Construct Noise Generators
	noiseGeneratorElevation := noise.NewFBMNoiseGenerator2D(
		16,    // octaveCount
		0.5,   // persistence
		2.0,   // lacunarity
		0.010, // frequency
		0.5,   // scale
		0.5,   // bias
		0.0,   // min
		1.0,   // max
		true,  // isRedistributed
		2.0,   // reDistributionFactor
		false, // isTerraced
		25.0,  // terraceFactor
	)
	noiseGeneratorMoisture := noise.NewFBMNoiseGenerator2D(
		16,    // octaveCount
		0.5,   // persistence
		2.0,   // lacunarity
		0.007, // frequency
		0.5,   // scale
		0.5,   // bias
		0.0,   // min
		1.0,   // max
		false, // isRedistributed
		1.0,   // reDistributionFactor
		true,  // isTerraced
		10.0,  // terraceFactor
	)

	// Generate noise
	moistureNoise = noiseGeneratorMoisture.BuildNoiseMatrix(targetSize.Width, targetSize.Height, 0.0, float64(targetSize.Width))
	elevationNoise = noiseGeneratorElevation.BuildNoiseMatrix(targetSize.Width, targetSize.Height, 0.0, float64(targetSize.Height))
	moistureNoiseRender := noiseGeneratorMoisture.BuildNoiseMatrix(targetSize.Width, targetSize.Height, 0.0, 255.0)
	elevationNoiseRender := noiseGeneratorElevation.BuildNoiseMatrix(targetSize.Width, targetSize.Height, 0.0, 255.0)

	// Generate noise images
	elevationImg = image.NewRGBA(image.Rect(0, 0, targetSize.Width, targetSize.Height))
	moistureImg = image.NewRGBA(image.Rect(0, 0, targetSize.Width, targetSize.Height))
	for x := 0; x < targetSize.Width; x++ {
		for y := 0; y < targetSize.Height; y++ {
			elevationDatum := elevationNoiseRender[x][y]
			moistureDatum := moistureNoiseRender[x][y]

			elevationImg.Set(x, y, color.RGBA{uint8(elevationDatum), uint8(elevationDatum), uint8(elevationDatum), 255})
			moistureImg.Set(x, y, color.RGBA{uint8(moistureDatum), uint8(moistureDatum), uint8(moistureDatum), 255})
		}
	}

	return
}
