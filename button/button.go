package button

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"image"
	"github.com/psilva261/xui"
	"github.com/psilva261/xui/events"
	"github.com/psilva261/xui/events/mouse"
	"github.com/psilva261/xui/font"
	"github.com/psilva261/xui/internal/color"
	"github.com/psilva261/xui/internal/geom"
	"github.com/psilva261/xui/space"
	"github.com/psilva261/xui/layout"
)

type Interface interface {
}

type Button struct {
	Orig image.Point
	Text string
	x xui.Interface
	textColor *memdraw.Image
	borderColor *memdraw.Image
	borderHoverColor *memdraw.Image

	textImg *memdraw.Image
	btnImg *memdraw.Image
	hoverImg *memdraw.Image

	hover bool

	Margin space.Sp

	cb func(ev events.Interface, userData any)
	cbUserData any
}

func New(x xui.Interface, orig image.Point, text string) (b *Button) {
	b = &Button{}
	b.Orig = orig
	b.Text = text
	b.x = x

	var err error
	b.textImg, err = font.String(text)
	if err != nil {
		panic(err.Error())
	}

	b.btnImg, err = memdraw.AllocImage(b.textImg.R.Inset(-10).Add(image.Point{10,10}), draw.ABGR32)
	if err != nil {
		panic(err.Error())
	}
	memdraw.FillColor(b.btnImg, draw.Opaque)

	b.hoverImg, err = memdraw.AllocImage(b.textImg.R.Inset(-10).Add(image.Point{10,10}), draw.ABGR32)
	if err != nil {
		panic(err.Error())
	}
	memdraw.FillColor(b.hoverImg, draw.Opaque)

	b.borderColor = color.Border
	b.borderHoverColor = color.Hover.Border

	b.btnImg.Draw(b.textImg.R.Add(image.Point{7,13}), b.textImg, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
	geom.DrawRoundedBorder(b.btnImg, b.btnImg.R, b.borderColor)

	b.hoverImg.Draw(b.textImg.R.Add(image.Point{7,13}), b.textImg, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
	geom.DrawRoundedBorder(b.hoverImg, b.hoverImg.R, b.borderHoverColor)

	return
}

func (b *Button) Event(ev events.Interface) {
	mev, ok := ev.(mouse.Event)
	if !ok { return }

	if mev.Type == mouse.Enter {
		b.hover = true
	} else if mev.Type == mouse.Leave {
		b.hover = false
	}

	if b.cb != nil {
		b.cb(mev, b.cbUserData)
	}
}

func (b Button) Render() *memdraw.Image {
	if b.hover {
		return b.hoverImg
	} else {
		return b.btnImg
	}
}

func (b Button) Focus() {
}

func (b Button) Layout() layout.Interface {
	return layout.Inline{}
}

func (b Button) Geom() (r image.Rectangle, margin space.Sp) {
	return b.btnImg.R, b.Margin
}

func (b *Button) SetCallback(cb func(ev events.Interface, userData any), userData any) {
	b.cb = cb
	b.cbUserData = userData
}

