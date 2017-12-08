package object

import (
	bitflag "github.com/mvpninjas/go-bitflag"
)

const (
	// FlagPlayer flags an object as a player controlled object
	FlagPlayer bitflag.Flag = 1 << bitflag.Flag(iota)

	// FlagNPC flags an object as a non-player controlled object
	FlagNPC

	// FlagBuilding flags an object as a building or structure
	FlagBuilding

	// FlagInvulnerable flags an object as unable to take damage
	FlagInvulnerable

	// FlagImmovable flags an object as unable to move in the game map
	FlagImmovable
)

// TypeFlagsBitSet is a bit set for various object types, each a bit flag that is true or false
type TypeFlagsBitSet struct {
	Flags bitflag.Flag
}

// NewTypeFlagsBitSet creates a new bit set comprised of the passed object type bit flags
func NewTypeFlagsBitSet(flags ...bitflag.Flag) *TypeFlagsBitSet {
	var finalFlag bitflag.Flag

	for _, aFlag := range flags {
		finalFlag.Set(aFlag)
	}

	return &TypeFlagsBitSet{Flags: finalFlag}
}

// Examples of calling
// flag.Set(A)
// flag.Set(B, C)
// flag.Set(C | D)
// flag.Clear()
// flag.Set(A, B, C, D)
// flag.Unset(A)
// flag.Unset(B, C)
// flag.Unset(A | C)
// if flag.Isset(A) {
// 	fmt.Println("A")
// }
// if flag.Isset(B) {
// 	fmt.Println("B")
// }
// if flag.Isset(C) {
// 	fmt.Println("C")
// }
// if !flag.Isset(D) {
// 	fmt.Println("D")
// }
// flag.Set(C)
// if !flag.One(D) {
// 	fmt.Println("E")
// }
