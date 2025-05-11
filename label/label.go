package label

import (
	"9fans.net/go/draw/memdraw"
	"image"
	"xui/events"
	"xui/internal/font"
	"xui/layout"
	"xui/space"
)

type Interface interface {
}

type Label struct {
	Orig image.Point
	Text string
	textColor *memdraw.Image

	textImg *memdraw.Image

	Margin space.Sp
}

func New(orig image.Point, text string) (l Label) {
	l.Orig = orig
	l.Text = text

	var err error
	l.textImg, err = font.String(text)
	if err != nil {
		panic(err.Error())
	}

	return
}

func (l Label) Event(events.Interface) {
}

func (l Label) Render() *memdraw.Image {
	return l.textImg
}

func (l Label) Focus() {
}

func (l Label) Layout() layout.Interface {
	return layout.Inline{}
}

func (l Label) Geom() (r image.Rectangle, margin space.Sp) {
	return l.textImg.R, l.Margin
}
