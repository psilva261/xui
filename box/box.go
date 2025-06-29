package box

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"image"
	//"log"
	"runtime"
	"sync"
	"github.com/psilva261/xui/element"
	"github.com/psilva261/xui/events"
	"github.com/psilva261/xui/events/mouse"
	"github.com/psilva261/xui/internal/color"
	"github.com/psilva261/xui/layout"
	"github.com/psilva261/xui/space"
)

type Interface = element.Interface

type Dir int

const (
	Horizontal Dir = iota + 1
	Vertical
)

type Box struct {
	Elements []element.Interface
	Rs []image.Rectangle
	Dir
	Wrap bool
	boxImg *memdraw.Image

	mouseEntered []bool

	// mouseXY for hover focus
	mouseXY image.Point

	// Optional parameters
	Width int
	Height int
	color.Colorset

	Background *memdraw.Image

	Margin space.Sp
	Border space.Sp
	Padding space.Sp
}

// rMax is optional.
//
// R can otherwise be determined recursively by
// calling els Geometry and applying the layout.
//
// ...can be dynamically resized though
// (similar to Go slices)
func New(els []element.Interface) *Box {
	return &Box{
		Elements: els,
		Rs: make([]image.Rectangle, len(els)),
		mouseEntered: make([]bool, len(els)),
	}
}

func (b *Box) Event(evOrig events.Interface) {
	if mev, ok := evOrig.(mouse.Event); ok {
		b.mouseXY = mev.Point
	}
	for i, el := range b.Elements {
		elR, elMargin := el.Geom()

		switch tevOrig := evOrig.(type) {
		case mouse.Event:
			var ev mouse.Event

			ev.Point = tevOrig.Point.
			  Sub(b.Padding.TopLeft()).
			  Sub(elR.Min).
			  Sub(elMargin.TopLeft())
			ev.Buttons = tevOrig.Buttons
			ev.Msec = tevOrig.Msec

			if tevOrig.Point.In(b.Rs[i]) {
				if !b.mouseEntered[i] {
					b.mouseEntered[i] = true
					ev.Type |= mouse.Enter
				}
				if tevOrig.Type != 0 {
					ev.Type |= tevOrig.Type
				}
			} else {
				if b.mouseEntered[i] {
					b.mouseEntered[i] = false
					ev.Type |= mouse.Leave
				}
			}
			if ev.Type != 0 {
				el.Event(ev)
			}
		default:
			if b.mouseXY.In(b.Rs[i]) {
				el.Event(tevOrig)
			}
		}
	}
}

func (b *Box) Render() *memdraw.Image {
	if b.boxImg == nil {
		b.layoutBoxImg()
	}

	// 2. Render
	ims := make([]*memdraw.Image, len(b.Elements))
	wg := sync.WaitGroup{}
	for i, el := range b.Elements {
		if false /* not working right now */ && runtime.GOARCH != "arm64" {
			wg.Add(1)
			go func(ii int) {
				ims[ii] = el.Render() //b.boxImg, b.Rs[i].Min)
				wg.Done()
			}(i)
		} else {
			// otherwise every item can end up being the same element
			ims[i] = el.Render() //b.boxImg, b.Rs[i].Min)
		}
	}
	wg.Wait()

	for i, im := range ims {
		//log.Printf("Render: ims[%d].R=%+v", i, ims[i].R)
		rIm := im.R.Add(b.Rs[i].Min).Add(b.Padding.TopLeft())
		b.boxImg.Draw(rIm, im, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
	}
	//log.Printf("Box.Render: surface=%v", surface)
	//log.Printf("b.boxImg.R=%v", b.boxImg.R)
	return b.boxImg
}

// Populate
//
// - b.Rs[i]
// - b.boxImg
func (b *Box) layoutBoxImg() {
	// 0. Validations

	// 1. Layout
	var dxy image.Point
	for i, el := range b.Elements {
		rEl, marginEl := el.Geom()
		b.Rs[i] = rEl.Add(dxy)
		b.Rs[i] = b.Rs[i].Add(marginEl.TopLeft())
		switch {
		case b.Dir == Horizontal || (b.Wrap && (b.Width == 0 || rEl.Dx()+dxy.X <= b.Width)):
			//log.Printf("horiz.")
			dxy = dxy.Add(image.Point{X: rEl.Dx()+marginEl.Left.Val})
			if i > 0 {
				_, marginLast := b.Elements[i-1].Geom()
				dxy = dxy.Add(image.Point{X: marginLast.Right.Val})
			}
		case b.Dir == Vertical:
			//log.Printf("vert.")
			fallthrough
		default:
			dxy = dxy.Add(image.Point{Y: rEl.Dy()+marginEl.Top.Val})
			if i > 0 {
				_, marginLast := b.Elements[i-1].Geom()
				dxy = dxy.Add(image.Point{Y: marginLast.Bottom.Val})
			}
			dxy.X = 0
		}
	}

	if b.boxImg == nil {
		var err error
		var r image.Rectangle
		if b.Width != 0 && b.Height != 0 {
			r = image.Rect(0, 0, b.Width, b.Height)
		} else if len(b.Elements) != 0 {
			r = image.Rectangle{
				//Min: b.Rs[0].Min,
				Max: b.Rs[len(b.Rs)-1].Max.Add(b.Rs[0].Min).
				  Add(b.Padding.Size()),
			}

			// Expand outer rectangle if inner element rectangles don't fit
			for _, el := range b.Elements {
				rEl, _ := el.Geom()
				if rEl.Dx() > r.Dx() {
					r.Max.X += rEl.Dx()-r.Dx()
				}
				if rEl.Dy() > r.Dy() {
					r.Max.X += rEl.Dy()-r.Dy()
				}
			}
		}
		// Allocate image
		b.boxImg, err = memdraw.AllocImage(r, draw.ABGR32)
		if err != nil {
			panic(err.Error())
		}
		if b.Colorset.Normal.Background != nil {
			//log.Printf("b.Background.R=%v", b.Background.R)
			b.boxImg.Draw(r, b.Colorset.Normal.Background, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
		} else {
			memdraw.FillColor(b.boxImg, draw.Transparent)
		}
		if b.Background != nil {
			b.boxImg.Draw(r, b.Background, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
		}
	}
}

func (b Box) Focus() {
}

func (b Box) Layout() layout.Interface {
	return layout.Inline{}
}

func (b *Box) Geom() (r image.Rectangle, margin space.Sp) {
	if b.boxImg == nil {
		b.layoutBoxImg()
	}
	return b.boxImg.R, b.Margin
}

