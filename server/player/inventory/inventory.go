package inventory

import (
	"bitbucket.org/ehhio/ehhworldserver/server/object"
)

type Inventory struct {
	slots []*object.Item
}

func (i *Inventory) Add(item *object.Object) *object.Object {
	if item.IsStackable() {
		for idx := range i.slots {
			if i.slots[idx].Type == item.Type && !i.slots[idx].IsMaxStacked() {
				item := i.slots[idx].Stack(item)
				if item.Count == 0 {
					return nil
				}
			}
		}
	}

	if i.IsFull() {
		return item
	}

	i.slots = append(i.slots, item)

	return nil
}

// RemoveAt retrieves all the items stored in the inventory at an index
func (i *Inventory) RemoveAt(index int) *object.Object {
	if index >= len(i.slots) {
		return nil
	}

	return i.slots[index]
}

func (i *Inventory) IsFull() bool {
	return len(i.slots) == cap(i.slots)
}
