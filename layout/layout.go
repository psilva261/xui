package layout

import (
	"image"
)

type Interface interface {
	Arrange([]Element)
}

type Element interface {
	Geom() (orig, extent image.Point)
}

type Inline struct {
}

func (inl Inline) Arrange([]Element) {}


