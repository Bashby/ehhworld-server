package player

import (
	bitflag "github.com/mvpninjas/go-bitflag"
)

const (
	// FlagDirtyPosition flags a player as having a position change that needs to be sent to the client
	FlagDirtyPosition bitflag.Flag = 1 << bitflag.Flag(iota)

	// FlagDirtyOrientation flags a player as having a orientation change that needs to be sent to the client
	FlagDirtyOrientation

	// FlagDirtyHealth flags a player as having a health change that needs to be sent to the client
	FlagDirtyHealth

	// FlagDirtyInventory flags a player as having an inventory change that needs to be sent to the client
	FlagDirtyInventory
)

// DirtyFlagsBitSet is a bit set for various player-state-is-dirty flags, each a bit flag that is true or false
type DirtyFlagsBitSet struct {
	Flags bitflag.Flag
}

// NewDirtyFlagsBitSet creates a new bit set comprised of the passed player dirty state bit flags
func NewDirtyFlagsBitSet(flags ...bitflag.Flag) *DirtyFlagsBitSet {
	var finalFlag bitflag.Flag

	for _, aFlag := range flags {
		finalFlag.Set(aFlag)
	}

	return &DirtyFlagsBitSet{Flags: finalFlag}
}
