package gui

import (
	"fmt"
	"image"
	"image/color"
	"net"
	"os"
	"path/filepath"

	"codeberg.org/kvo/std"
)

type wlPtr struct {
}

func (p wlPtr) Pos() (x, y int) {
	x = 0; y = 0
	return x, y
}

type wlWin struct {
}

func (w wlWin) At(x, y int) color.Color {
	return color.RGBA{0, 0, 0, 0xff}
}

func (w wlWin) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{0, 0},
		image.Point{0, 0},
	}
}

func (w wlWin) ColorModel() color.Model {
	return color.RGBAModel
}

func (w wlWin) Set(x, y int, c color.Color) {
	return
}

func (w wlWin) Title(name string) error {
	return nil
}

type wlConn struct {
	conn net.Conn
}

func (w wlConn) Events() <-chan Event {
	events := make(chan Event)
	return events
}

func (w wlConn) Pointer() (Pointer, error) {
	return wlPtr{}, nil
}

func (w wlConn) Window() (Window, error) {
	return wlWin{}, nil
}

// Attempts to establish a connection with Wayland, or, if that fails, with X11.
func dial() (Conn, error) {
	var display string
	// Check for WAYLAND_SOCKET
	display = os.Getenv("WAYLAND_DISPLAY")
	if display == "" {
		display = "wayland-0"
	}
	leadChar, err := std.Access([]rune(display), 0)
	if err != nil {
		return nil, err
	}
	if leadChar != '/' {
		xdgRt := os.Getenv("XDG_RUNTIME_DIR")
		if xdgRt == "" {
			return nil, fmt.Errorf("gui: XDG_RUNTIME_DIR not set")
		}
		display = filepath.Join(xdgRt, display)
	}
	netConn, err := net.Dial("unix", display)
	if err != nil {
		return nil, fmt.Errorf("gui: failed to dial display: %v", err)
	}
	// TODO: Add X11 support.
	return wlConn{netConn}, nil
}
