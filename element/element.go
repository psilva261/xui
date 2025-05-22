package element

import (
	"9fans.net/go/draw/memdraw"
	"image"
	"github.com/psilva261/xui/events"
	"github.com/psilva261/xui/layout"
	"github.com/psilva261/xui/space"
)

type Interface interface {
	Event(events.Interface)
	Render() *memdraw.Image
	Focus()
	// Layout hints about how to arrange this element
	Layout() layout.Interface
	// Coordinates relative to parent element unless layout expects different coordinate system
	//
	// Other possible coordinates:
	// - global (relative to window)
	// - global scroll (relative to scrolled window)
	// - inline (relative to neighbour)
	// - parent (relative to (scrolled) parent element)
	Geom() (r image.Rectangle, margin space.Sp)
}
