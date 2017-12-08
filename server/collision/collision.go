package collision

import (
	"bitbucket.org/ehhio/ehhworldserver/server/object"
	"bitbucket.org/ehhio/ehhworldserver/server/utility"

	"github.com/dhconnelly/rtreego"
)

// Constants
const minRTreeBranchingFactor = 5
const maxRTreeBranchingFactor = 25
const rTreeDimensions = 2 // 2D spatial index

// Trackable is the interface objects must satisfy to be tracked by the collision system
type Trackable interface {
	GetAABBBottomLeftPoint() *utility.PositionHighResolution
	GetAABBSize() *utility.SizeHighResolution
	GetObjectFlags() *object.TypeFlagsBitSet
}

// objectWrapper is a wrapper object that will be used in collision calculations and keep a reference to the original game object
type objectWrapper struct {
	rect *rtreego.Rect
	ref  Trackable // A pointer back to the game object
}

// newobjectWrapper creates a new collision object to be tracked by the collision system
func newObjectWrapper(rect *rtreego.Rect, ref Trackable) *objectWrapper {
	return &objectWrapper{
		rect: rect,
		ref:  ref,
	}
}

// Bounds returns a 2D rect defining the bounds of the object in the collision system
// Implements the interface to be tracked in the collision system
func (co *objectWrapper) Bounds() *rtreego.Rect {
	return co.rect
}

// Collision is the system that tracks and computes collision questions between objects in 2D space
type Collision struct {
	rtree    *rtreego.Rtree
	wrappers map[Trackable]*objectWrapper
}

// NewCollision creates a new collision system for tracking objects in 2D space
func NewCollision() *Collision {
	return &Collision{
		rtree:    rtreego.NewTree(rTreeDimensions, minRTreeBranchingFactor, maxRTreeBranchingFactor),
		wrappers: make(map[Trackable]*objectWrapper),
	}
}

// AddObject adds an object, that implements the Trackable interface, to the collision system.
// If the object has moved, it MUST be updated in the collision system.
func (c *Collision) AddObject(object Trackable) {
	objPoint := object.GetAABBBottomLeftPoint()
	objSize := object.GetAABBSize()
	treePoint := rtreego.Point{objPoint.X, objPoint.Y}
	treeRect, _ := rtreego.NewRect(treePoint, []float64{objSize.Width, objSize.Height})
	colObj := newObjectWrapper(treeRect, object)
	c.wrappers[object] = colObj
	c.rtree.Insert(colObj)
}

// UpdateObject modifies the position / dimensions of the object in the collision system
func (c *Collision) UpdateObject(object Trackable) {
	// You must delete and re-add to "update". Limitation of underlying lib.
	c.DeleteObject(object)
	c.AddObject(object)
}

// DeleteObject removes a tracked object from the collision system
func (c *Collision) DeleteObject(object Trackable) {
	wrapperObject := c.wrappers[object]
	delete(c.wrappers, object)
	c.rtree.Delete(wrapperObject)
}

// GetAllWithinRect returns all game objects within a rect
func (c *Collision) GetAllWithinRect(point *utility.PositionHighResolution, size *utility.SizeHighResolution) []Trackable {
	treePoint := rtreego.Point{point.X, point.Y}
	searchRect, _ := rtreego.NewRect(treePoint, []float64{size.Width, size.Height})
	matches := c.rtree.SearchIntersect(searchRect)

	res := make([]Trackable, len(matches))
	for i := range matches {
		res[i] = matches[i].(*objectWrapper).ref
	}

	return res
}

// GetWithinRect returns game objects within a rect, filtered against object bit flags
func (c *Collision) GetWithinRect(point *utility.PositionHighResolution, size *utility.SizeHighResolution, flags *object.TypeFlagsBitSet) []Trackable {
	treePoint := rtreego.Point{point.X, point.Y}
	searchRect, _ := rtreego.NewRect(treePoint, []float64{size.Width, size.Height})
	matches := c.rtree.SearchIntersect(searchRect, c.objectTypeFlagsBitSetFilterGenerator(flags))

	res := make([]Trackable, len(matches))
	for i := range matches {
		res[i] = matches[i].(*objectWrapper).ref
	}

	return res
}

// objectTypeFlagsBitSetFilter implements a rtree filter function using object type bit flags
func (c *Collision) objectTypeFlagsBitSetFilterGenerator(flags *object.TypeFlagsBitSet) rtreego.Filter {
	return func(results []rtreego.Spatial, object rtreego.Spatial) (refuse, abort bool) {
		refObject := object.(*objectWrapper).ref.(Trackable)
		// Return valid if flags match
		if refObject.GetObjectFlags().Flags.Isset(flags.Flags) {
			return true, false
		}

		return false, false
	}
}
