package font

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/plan9font"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	imagedraw "image/draw"
)

func String(text string) (textImg *memdraw.Image, err error) {
	readFile := func(name string) ([]byte, error) {
		return ioutil.ReadFile(filepath.FromSlash(path.Join(FontDir(), name)))
	}
	fontData, err := readFile(FontFile())
	if err != nil {
		log.Fatal(err)
	}
	face, err := plan9font.ParseFont(fontData, readFile)
	if err != nil {
		log.Fatal(err)
	}
	ascent := face.Metrics().Ascent.Ceil()

	w := 0
	for _, r := range text {
		bounds, adv, _ := face.GlyphBounds(r)
		_ = bounds
		w += adv.Ceil()
	}
	r := image.Rect(0, 0, w, 2*ascent)
	if r.Dx() == 0 {
		r.Max.X = r.Min.X + 1
	}
	if r.Dy() == 0 {
		r.Max.Y = r.Min.Y + 1
	}
	dst := image.NewRGBA(r)
	imagedraw.Draw(dst, dst.Bounds(), image.White, image.Point{}, imagedraw.Src)
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(0, ascent),
	}
	d.DrawString(text)

	log.Printf("String: dst.Bounds()=%+v", dst.Bounds())

	textImg, err = memdraw.AllocImage(dst.Bounds(), draw.ABGR32)
	if err != nil {
		panic(fmt.Errorf("allocimage: %s", err).Error())
	}

	_, err = memdraw.Load(textImg, dst.Bounds(), dst.Pix, false)
	if err != nil {
		panic(err.Error())
	}

	return
}
