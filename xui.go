package xui

import (
	"github.com/psilva261/xui/element"
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"fmt"
	"image"
	"io"
	"log"
	"sync"
	"time"
	"github.com/psilva261/xui/events/keyboard"
	"github.com/psilva261/xui/events/mouse"
	//"github.com/psilva261/xui/internal/color"
	"github.com/psilva261/xui/space"
)

const (
	scrollStep = 40
)

type Interface interface {
	R() image.Rectangle
	SetRoot(element.Interface)
	Render()
	Loop()

	Scale(n int) int

	// Create properly scaled points, rectangles and spaces

	Pt(x, y int) image.Point
	Rect(x0, y0, x1, y1 int) image.Rectangle
	Space(vals... int) space.Sp
}

type Xui struct {
	mu sync.Mutex
	keyctl *draw.Keyboardctl
	mousectl *draw.Mousectl
	errch chan error
	root element.Interface
	display *draw.Display
	surface *memdraw.Image
	bgLayer *memdraw.Image

	rootXY image.Point
}

func (x *Xui) SetRoot(el element.Interface) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.root = el
}

var buf = make([]byte, 1024*1024*32)

func (x *Xui) Render() {
	x.mu.Lock()
	defer x.mu.Unlock()

	if x.root != nil {
		//log.Printf("x.Render: x.bgLayer=%v, color.EmptyMask=%v", x.bgLayer, color.EmptyMask)
		//log.Printf("x.Render: x.surface.R=%v, x.bgLayer.R=%v, color.EmptyMask.R=%v", x.surface.R, x.bgLayer.R, color.EmptyMask.R)
		x.surface.Draw(x.bgLayer.R, x.bgLayer, image.ZP, nil, image.ZP, draw.SoverD)

		im := x.root.Render() //x.surface, x.rootXY)
		x.surface.Draw(im.R.Add(x.rootXY), im, image.ZP, nil, image.ZP, draw.SoverD)

		nbuf, err := memdraw.Unload(x.surface, x.surface.R, buf)
		if err != nil { panic(err.Error()) }

		data := buf[:nbuf]
		_, err = x.display.ScreenImage.Load(x.surface.R, data)
		if err != nil {
			panic(err.Error())
		}

		x.display.Flush()
	}
}

func (x *Xui) R() image.Rectangle {
	return x.display.ScreenImage.R
}

func (x *Xui) Scale(n int) int {
	return x.display.Scale(n)
}

func (x *Xui) Pt(x0, y0 int) image.Point {
	return image.Point{
		x.Scale(x0),
		x.Scale(y0),
	}
}

func (x *Xui) Rect(x0, y0, x1, y1 int) image.Rectangle {
	return image.Rectangle{
		x.Pt(x0, y0),
		x.Pt(x1, y1),
	}
}

func (x *Xui) Space(vals... int) space.Sp {
	scaled := make([]int, len(vals))
	for i := 0; i < len(vals); i++ {
		scaled[i] = x.Scale(vals[i])
	}
	return space.New(vals...)
}

func New() (Interface, error) {

	errch := make(chan error, 1)
	display, err := draw.Init(errch, "", "hello", "800x600")
	if err != nil {
		return nil, err
	}
	memdraw.Init()
	_ = display

	x := &Xui{
		errch: errch,
		keyctl: display.InitKeyboard(),
		mousectl: display.InitMouse(),
		display: display,
	}

	x.surface, err = memdraw.AllocImage(x.display.ScreenImage.R, x.display.ScreenImage.Pix)
	if err != nil {
		return nil, fmt.Errorf("alloc image: %w", err)
	}

	x.bgLayer, err = memdraw.AllocImage(x.display.ScreenImage.R, draw.ABGR32)
	if err != nil {
		return nil, fmt.Errorf("alloc image: %w", err)
	}
	memdraw.FillColor(x.bgLayer, draw.White)

	return x, nil
}

func (x *Xui) Loop() {
	go func() {
		for {
			(func() {
				<-time.After(10*time.Millisecond)
				x.Render()
			})()
		}
	}()
	for {
		select {
		case m := <-x.mousectl.C:
			if m.Buttons == 8 {
				x.rootXY.Y += scrollStep
			} else if m.Buttons == 16 {
				x.rootXY.Y -= scrollStep
			}
			go func() {
				x.mu.Lock()
				defer x.mu.Unlock()

				var ev mouse.Event

				if m.Buttons&1 > 0 {
					ev.Type |= mouse.Click
				}


				if x.root != nil {
					ev.Point = m.Point.Sub(x.rootXY)
					ev.Buttons = m.Buttons
					ev.Msec = m.Msec
					x.root.Event(ev)
				}
			}()
		case k := <-x.keyctl.C:
			go func() {
				x.mu.Lock()
				defer x.mu.Unlock()

				log.Printf("KEY %v", k)

				ev := keyboard.Event{}
				ev.Type = keyboard.Pressed
				ev.Key = k
				if x.root != nil {
					x.root.Event(ev)
				}
			}()
		case <-x.mousectl.Resize:
			log.Printf("resize")
		case err := <-x.errch:
			log.Printf("errch: %v", err)
			if err == io.EOF {
				log.Printf("returning")
				return
			}
		case <-time.After(10*time.Millisecond):
			//x.Render()
		}
	}
}
