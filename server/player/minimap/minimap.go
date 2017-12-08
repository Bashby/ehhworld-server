package minimap

import (
	"image/color"

	"bitbucket.org/ehhio/ehhworldserver/server/utility"
)

// Minimap defines a player mini-map
type Minimap struct {
	size     *utility.Size          // The size of the minimap
	known    *utility.BooleanMatrix // Boolean knowledge of which blocks a player is aware of across all sessions
	explored *utility.BooleanMatrix // Boolean knowledge of which blocks a player has visited in their current session
	markers  map[int]*MapMarker     // All markers on the map
}

// NewMinimap creates a new minimap object
func NewMinimap(blockMatrixSize *utility.Size) *Minimap {
	matrixKnown := utility.NewBooleanMatrix(blockMatrixSize)
	matrixExplored := utility.NewBooleanMatrix(blockMatrixSize)

	return &Minimap{
		size:     blockMatrixSize,
		known:    matrixKnown,
		explored: matrixExplored,
		markers:  make(map[int]*MapMarker),
	}
}

// AddMarker inserts a new marker
func (m *Minimap) AddMarker(marker *MapMarker) int {
	// Get keys of map
	keys := make([]int, len(m.markers))
	i := 0
	for k := range m.markers {
		keys[i] = k
		i++
	}

	// Find the max key
	maxIndex := -1
	for _, v := range keys {
		if v > maxIndex {
			maxIndex = v
		}
	}

	// Store
	key := maxIndex + 1
	m.markers[key] = marker

	// Return key
	return key
}

// RemoveMarker deletes a marker from the minimap, by key
func (m *Minimap) RemoveMarker(key int) {
	delete(m.markers, key)
}

// GetMarker returns a marker from the minimap, by key
func (m *Minimap) GetMarker(key int) *MapMarker {
	return m.markers[key]
}

// ExploreUnknownWithin explores all unexplored map indexes within a certain distance of a point.
// Returns a slice of these newly explored blocks
func (m *Minimap) ExploreUnknownWithin(position *utility.Position, dx, dy int) []*utility.Position {
	if !position.IsWithinBounds(m.size) {
		return make([]*utility.Position, 0)
	}

	return m.getAdjacentWithin(m.explored, position, dx, dy, false)
}

// getAdjacentWithin returns adjacent indexs within a given distance that are the requested boolean
func (m *Minimap) getAdjacentWithin(target *utility.BooleanMatrix, position *utility.Position, dx, dy int, getThisBool bool) []*utility.Position {
	xMin := utility.Clamp(position.X-dx, 0, m.size.Width)
	xMax := utility.Clamp(position.X+dx+1, 0, m.size.Width)
	yMin := utility.Clamp(position.Y-dy, 0, m.size.Height)
	yMax := utility.Clamp(position.Y+dy+1, 0, m.size.Height)

	res := make([]*utility.Position, 0, (xMax-xMin)*(yMax-yMin)-1)
	for x := xMin; x < xMax; x++ {
		for y := yMin; y < yMax; y++ {
			// Don't add ourselves
			if x != position.X && y != position.Y {
				// Only add indexs that match the requested boolean state
				if target.Matrix[x][y] == getThisBool {
					target.Matrix[x][y] = !getThisBool
					res = append(res, &utility.Position{X: x, Y: y})
				}
			}
		}
	}

	return res
}

// MapMarker is a object defining a marked point on a map
type MapMarker struct {
	position    *utility.PositionHighResolution
	markerType  MapMarkerType
	color       color.Color
	title       string
	description string
}

// NewMapMarker creates a new marker to be placed on a map
func NewMapMarker(position *utility.PositionHighResolution, markerType MapMarkerType, color color.Color, title, description string) *MapMarker {
	return &MapMarker{
		position:    position,
		markerType:  markerType,
		color:       color,
		title:       title,
		description: description,
	}
}

// SetMarkerPosition changes the position of a marker to a new supplied position
func (m *MapMarker) SetMarkerPosition(position *utility.PositionHighResolution) {
	m.position = position
}

// SetMarkerColor changes the color of a marker to a new supplied color
func (m *MapMarker) SetMarkerColor(color color.Color) {
	m.color = color
}

// SetMarkerTitle changes the title of a marker to a new supplied title, by key
func (m *MapMarker) SetMarkerTitle(title string, key int) {
	m.title = title
}

// SetMarkerDescription changes the description of a marker to a new supplied description, by key
func (m *MapMarker) SetMarkerDescription(description string, key int) {
	m.description = description
}

// MapMarkerType The type of a map marker
type MapMarkerType int

//go:generate stringer -type=MapMarkerType

// MapMarkerType enum
const (
	PlayerMapMarker MapMarkerType = iota
	BattleMapMarker
	QuestAvailableMapMarker
	QuestDestinationMapMarker
	QuestCompleteMapMarker
)
