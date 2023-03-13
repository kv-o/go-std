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
	"io"
)

// Open attempts to open a Window in the current GUI (if one exists), set up a
// Stdin interface, and provides access to a Pointer. Returns error as soon as
// one is encountered.
func Open() (Window, Stdin, Pointer, error) {
	return open()
}

// Window is a draw.Image that represents a rectangular area of an underlying
// GUI.
type Window interface {
	draw.Image
}

// Stdin multiplexes os.Stdin and platform-specific input provided by the
// underlying GUI.
type Stdin interface {
	io.Reader
}

// Pointer represents the current position of the pointing device or the last
// position of physical contact.
type Pointer interface {
	Pos() (x, y int)
}
