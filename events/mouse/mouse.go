package mouse

import "image"

type Type int

const (
	Enter Type = 1 << iota
	Leave
	Click
	Down
)

const Up = Click

// Mouse Point in local coordinates though
//
// E.g. for click event on Button relative to button
// - thus for mouse leave it can be negative
type Event struct {
	Type

	image.Point

	Buttons int
	// timestamp in ms
	Msec uint32
}
