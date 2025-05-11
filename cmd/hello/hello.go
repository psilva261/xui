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
	"xui"
	"xui/anim"
	"xui/box"
	"xui/button"
	"xui/element"
	"xui/events"
	"xui/events/mouse"
	"xui/field"
	"xui/label"
	"xui/space"
)

func Main() (err error) {
	memdraw.Init()

	x, err := xui.New()

	if err != nil {
		return
	}


	l := label.New(image.Point{10, 5}, "Hello, world!!") //, c)
	l.Margin = space.New(5, 10)

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

	fl := field.New(x, image.ZP, " ", x.Rect(0,0, 150, 50))
	fl.Margin = x.Space(5, 10)

	b := box.New([]element.Interface{
		a, l, btn, fl,
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
