package field

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	//"fmt"
	"image"
	"log"
	"sync"
	"github.com/psilva261/xui"
	"github.com/psilva261/xui/events"
	"github.com/psilva261/xui/events/keyboard"
	"github.com/psilva261/xui/events/mouse"
	"github.com/psilva261/xui/font"
	"github.com/psilva261/xui/internal/color"
	"github.com/psilva261/xui/internal/geom"
	"github.com/psilva261/xui/layout"
	"github.com/psilva261/xui/space"
)

const (
	Backspace rune = 8
)

type Interface interface {
}

type Field struct {
	mu sync.RWMutex

	Orig image.Point
	image.Rectangle
	x xui.Interface

	Text string
	color.Colorset

	// Position of the cursor
	Pos int

	textImg *memdraw.Image
	borderImg *memdraw.Image
	hoverImg *memdraw.Image

	hover bool

	cb func(ev events.Interface, userData any)
	cbUserData any

	Margin space.Sp
}

func New(x xui.Interface, orig image.Point, text string, r image.Rectangle) (f *Field) {
	f = &Field{}
	f.Orig = orig
	f.Text = text
	f.Rectangle = r
	f.x = x

	f.Colorset.Normal.Border = color.Border
	f.Colorset.Hover.Border = color.Hover.Border

	return
}

func (f *Field) Event(ev events.Interface) {
	f.mu.Lock()
	defer f.mu.Unlock()

	switch tev := ev.(type) {
	case mouse.Event:
		if tev.Type == mouse.Enter {
			f.hover = true
		} else if tev.Type == mouse.Leave {
			f.hover = false
		}

		if f.cb != nil {
			f.cb(tev, f.cbUserData)
		}
	case keyboard.Event:
		log.Printf("key pressed: %+v %d % x", tev, tev, tev)

		switch tev.Key {
		case Backspace:
			if f.Pos > 0 && len(f.Text) > 0 {
				f.Text = f.Text[:f.Pos-1] + f.Text[f.Pos:]
			}
			f.Pos -= 1
		case draw.KeyDelete:
			if f.Pos < len(f.Text) {
				f.Text = f.Text[:f.Pos] + f.Text[f.Pos+1:]
			}
		case draw.KeyLeft:
			f.Pos -= 1
		case draw.KeyRight:
			f.Pos += 1
		default:
			f.Text += string([]byte{byte(tev.Key)})
			f.Pos = len(f.Text)
		}
		if f.Pos < 0 {
			f.Pos = 0
		}
		if f.Pos > len(f.Text) {
			f.Pos = len(f.Text)
		}
		//log.Printf("event: call updateTextImgs")
		f.updateTextImgs()
	}
}

func (f *Field) Render() *memdraw.Image {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.borderImg == nil {
		//log.Printf("render: call updateTextImgs")
		f.updateTextImgs()
	}

	if f.hover {
		return f.hoverImg
	} else {
		return f.borderImg
	}
}

func (f *Field) updateTextImgs() {

	var err error
	//log.Printf("updateTextImgs: call font.String")
	f.textImg, err = font.String(f.Text)
	if err != nil {
		panic(err.Error())
	}
	r := f.Rectangle//f.textImg.R.Intersect(f.Rectangle)

	f.borderImg, err = memdraw.AllocImage(f.Rectangle.Inset(-f.x.Scale(5)).Add(f.x.Pt(5, 5)), draw.ABGR32)
	if err != nil {
		panic(err.Error())
	}
	memdraw.FillColor(f.borderImg, draw.Opaque)

	f.hoverImg, err = memdraw.AllocImage(f.Rectangle.Inset(-f.x.Scale(5)).Add(f.x.Pt(5, 5)), draw.ABGR32)
	if err != nil {
		panic(err.Error())
	}
	memdraw.FillColor(f.hoverImg, draw.Opaque)

	rr := r
	rr.Min = r.Min.Add(f.x.Pt(3, 7))

	f.borderImg.Draw(rr, f.textImg, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
	geom.DrawRoundedBorder(f.borderImg, f.Rectangle, f.Colorset.Normal.Border)

	f.hoverImg.Draw(rr, f.textImg, image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
	geom.DrawRoundedBorder(f.hoverImg, f.Rectangle, f.Colorset.Hover.Border)

	textPart, err := font.String(f.Text[:f.Pos])
	if err == nil {
		geom.DrawCursor(f.hoverImg, r, textPart.R, f.Colorset.Hover.Border)
	} else {
		log.Printf("font string sub text: %v", err)
	}
}

func (f Field) Focus() {
}

func (f Field) Layout() layout.Interface {
	return layout.Inline{}
}

func (f *Field) Geom() (r image.Rectangle, margin space.Sp) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.borderImg == nil {
		//log.Printf("geom: call updateTextImgs")
		f.updateTextImgs()
	}
	return f.borderImg.R, f.Margin
}

func (f *Field) SetCallback(cb func(ev events.Interface, userData any), userData any) {
	f.cb = cb
	f.cbUserData = userData
}


