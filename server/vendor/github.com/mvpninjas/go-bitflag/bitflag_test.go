package bitflag

import "testing"

const (
	FLAG_A Flag = 1 << Flag(iota)
	FLAG_B
	FLAG_C
	FLAG_D
)

func TestExample(t *testing.T) {

	var flag Flag

	flag.Set(FLAG_A)
	flag.Set(FLAG_B, FLAG_C)
	flag.Set(FLAG_C | FLAG_D)

	flag.Clear()

	flag.Set(FLAG_A, FLAG_B, FLAG_C, FLAG_D)

	flag.Unset(FLAG_A)

	flag.Unset(FLAG_B, FLAG_C)

	if flag.Isset(FLAG_A) {
		t.Fatal("A")
	}

	if flag.Isset(FLAG_B) {
		t.Fatal("B")
	}

	if flag.Isset(FLAG_C) {
		t.Fatal("C")
	}

	if !flag.Isset(FLAG_D) {
		t.Fatal("D")
	}
}

func TestSet(t *testing.T) {

	var flag Flag

	flag.Set(FLAG_A)
	flag.Set(FLAG_B)

	if !flag.Isset(FLAG_A) {
		t.Fail()
	}

	if !flag.Isset(FLAG_B) {
		t.Fail()
	}

	flag.Set(FLAG_A, FLAG_B)

	if !flag.Isset(FLAG_A, FLAG_B) {
		t.Fail()
	}

	if !flag.One(FLAG_B) {
		t.Fail()
	}
}

func TestUnset(t *testing.T) {

	var flag Flag

	flag.Set(FLAG_A, FLAG_B, FLAG_C)

	if !flag.Isset(FLAG_A, FLAG_B, FLAG_C) {
		t.Fail()
	}

	flag.Unset(FLAG_B)

	if flag.Isset(FLAG_B) {
		t.Fail()
	}

	if !flag.Isset(FLAG_A, FLAG_C) {
		t.Fail()
	}

}

func TestClear(t *testing.T) {

	var flag Flag

	flag.Set(FLAG_A, FLAG_B, FLAG_C)

	if !flag.Isset(FLAG_A, FLAG_B, FLAG_C) {
		t.Fail()
	}

	flag.Clear()

	if flag.One(FLAG_A, FLAG_B, FLAG_C) {
		t.Fail()
	}
}
