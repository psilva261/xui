package main

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"fmt"
	"image"
	imagedraw "image/draw"
	"image/gif"
	"log"
	"os"
	"github.com/psilva261/xui"
	"github.com/psilva261/xui/anim"
	"github.com/psilva261/xui/box"
	"github.com/psilva261/xui/button"
	"github.com/psilva261/xui/element"
	"github.com/psilva261/xui/events"
	"github.com/psilva261/xui/events/mouse"
	"github.com/psilva261/xui/field"
	"github.com/psilva261/xui/label"
	"github.com/psilva261/xui/space"
)

func Main() (err error) {
	memdraw.Init()

	x, err := xui.New()

	if err != nil {
		return
	}


	l := label.New(image.Point{10, 5}, "Hello, world!!") //, c)
	l.Margin = space.New(5, 10)
	l2 := label.New(image.Point{10, 5}, "Additional text!!") //, c)
	l2.Margin = space.New(5, 10)
	l3 := label.New(image.Point{10, 5}, "A really long subsequent text to fill the horizontal space!") //, c)
	l3.Margin = space.New(5, 10)
	l4 := label.New(image.Point{10, 5}, "Wrap into the next line for sure!!!") //, c)
	l4.Margin = space.New(5, 10)

	f, err := os.Open("animated-computer-image-0064.gif")
	//f, err := os.Open("g3Ys.gif")
	if err != nil { return }
	gifImg, err := gif.DecodeAll(f)
	if err != nil { return }

	animImg := make([]*memdraw.Image, len(gifImg.Image))
	rGifImg := image.Rectangle{
		Max: image.Point{gifImg.Config.Width, gifImg.Config.Height},
	}
	for i, im := range gifImg.Image {
		ni, err := memdraw.AllocImage(rGifImg, draw.ABGR32)
		if err != nil {
			return fmt.Errorf("allocimage: %s", err)
		}

		// Stolen from duit.ReadImage
		var rgba *image.RGBA
		b := im.Bounds()
		rgba = image.NewRGBA(rGifImg)
		imagedraw.Draw(rgba, rgba.Bounds(), im, b.Min, imagedraw.Src)
		_, err = memdraw.Load(ni, rGifImg, rgba.Pix, false)
		if err != nil {
			return fmt.Errorf("load image_: %w", err)
		}

		animImg[i] = ni
	}

	a := anim.NewFromSeq(x, animImg, gifImg.Delay)
	a.Margin = x.Space(5, 10)

	btn := button.New(x, image.ZP, "click right here")
	btn.SetCallback(func(ev events.Interface, _ any) {
		if mev, ok := ev.(mouse.Event); ok && mev.Type & mouse.Click > 0 {
			log.Printf("clicked!")
		}
	}, nil)
	btn.Margin = x.Space(5, 10)

	fl := field.New(x, image.ZP, "", x.Rect(0,0, 150, 50))
	fl.Margin = x.Space(5, 10)
	b2 := box.New([]element.Interface{
		l,
		l2,
		l3,
		l4,
	})
	b2.Wrap = true
	log.Printf("xui screen image width: %v", x.R().Dx())
	b2.Width = x.R().Dx()/2
	b := box.New([]element.Interface{
		a,
		b2,
		btn, fl,
	})
	//b.Dir = box.Horizontal
	x.SetRoot(b)

	x.Loop()

	return
}

func main() {
	if err := Main(); err != nil {
		log.Fatalf("%v", err)
	}
}
