// Package gui provides a basic portable library to construct graphical user
// interfaces.
//
// The main interface is called Window, a draw.Image interface which represents
// a rectangular area of an underlying graphical user interface (GUI). Window
// can thus be manipulated like a regular draw.Image. The Window interface does
// not provide other regular window manipulation functions (i.e. resizing,
// moving) as these should be handled by the window manager. The gui library
// also does not support multi-window programs, as this harms portability and
// may lead to abuse of the capability.
//
// Warning: package gui does not support GPUs and therefore should not be used
// in graphics-intensive solutions.
package gui

import (
	"image/draw"
)

// Conn is a generic, abstracted connection to an underlying GUI.
type Conn interface {
	// Events returns a read-only event channel over which GUI events are sent.
	Events() <-chan Event
	// Each call returns a copy of the same Pointer.
	Pointer() Pointer
	// Each call returns a copy of the same Window.
	Window() Window
}

// Event is a string code for an event.
type Event struct {
	Type   uint32
	Value  string
}

// List of event types.
const (
	KbDown     = 0x001
	KbHold     = 0x002
	KbUp       = 0x003
	Tap        = 0x004
	WinChange  = 0x003
)

// List of common non-linguistic keys.
const (
	KeyAlt        = "alt"
	KeyBackspace  = "backspace"
	KeyBreak      = "break"
	KeyCapsLock   = "capslock"
	KeyCtrl       = "ctrl"
	KeyDel        = "delete"
	KeyDown       = "down"
	KeyEnd        = "end"
	KeyEscape     = "escape"
	KeyF1         = "f1"
	KeyF2         = "f2"
	KeyF3         = "f3"
	KeyF4         = "f4"
	KeyF5         = "f5"
	KeyF6         = "f6"
	KeyF7         = "f7"
	KeyF8         = "f8"
	KeyF9         = "f9"
	KeyF10        = "f10"
	KeyF11        = "f11"
	KeyF12        = "f12"
	KeyFn         = "fn"
	KeyFnLock     = "fnlock"
	KeyHome       = "home"
	KeyInsert     = "insert"
	KeyLeft       = "left"
	KeyNumLock    = "numlock"
	KeyPause      = "pause"
	KeyPgDown     = "pagedown"
	KeyPgUp       = "pageup"
	KeyPrtScr     = "prtscr"
	KeyRight      = "right"
	KeyScrollLock = "scrolllock"
	KeyShift      = "shift"
	KeySuper      = "super"
	KeyUp         = "up"
)

// List of common events.
var (
	// Window close event. Signals a request for the termination of the GUI
	// session.
	Close   = Event{WinChange, "Close"}
	// Underlying GUI error.
	Error   = Event{WinChange, "Error"}
	// Left mouse button for right-handed people, single finger tap on touchscreen.
	Mouse1  = Event{Tap, "Mouse1"}
	// Middle mouse button or equivalent.
	Mouse2  = Event{Tap, "Mouse2"}
	// Right mouse button for right-handed people, or equivalent.
	Mouse3  = Event{Tap, "Mouse3"}
	Resize  = Event{WinChange, "Resize"}
)

// Pointer represents the current position of the pointing device or the last
// position of physical contact.
type Pointer interface {
	Pos() (x, y int)
}

// Window is a draw.Image that represents a rectangular area of an underlying
// two-dimensional GUI.
type Window interface {
	draw.Image
	// Title requests the underlying window manager to set name as the window
	// title.
	Title(name string) error
}

// Dial attempts to establish a connection to the current underlying GUI (if one
// exists).
func Dial() (Conn, error) {
	return dial()
}
