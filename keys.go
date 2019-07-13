package main

import (
	// "fmt"
	"github.com/gotk3/gotk3/gdk"
	// "github.com/gotk3/gotk3/glib"
	// "github.com/gotk3/gotk3/gtk"
	// "log"
	// "reflect"
	// "runtime"
)

type KeySym = uint
type KeyPress struct {
	EventKey gdk.EventKey
}

func (kp *KeyPress) Key() KeySym {
	return kp.EventKey.KeyVal()
}

func (kp *KeyPress) hasModifier(mask uint) bool {
	return kp.EventKey.State()&uint(mask) != 0
}

func (kp *KeyPress) getModMask() uint {
	return kp.EventKey.State() ^ kp.EventKey.KeyVal()
}

func (kp *KeyPress) Alt() bool {
	return kp.hasModifier(gdk.GDK_MOD1_MASK)
}

func (kp *KeyPress) Ctrl() bool {
	return kp.hasModifier(gdk.GDK_CONTROL_MASK)
}

func (kp *KeyPress) Shift() bool {
	return kp.hasModifier(uint(gdk.GDK_SHIFT_MASK))
}

func (kp *KeyPress) State() uint {
	return kp.EventKey.State()
}

type KeyChord struct {
	Key  KeySym
	Alt bool
	Ctrl bool
	Shift bool
}

