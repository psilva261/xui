package xuitest

import (
	"9fans.net/go/draw/memdraw"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"xui/element"
	"xui/space"
)

type Xui struct {
	image.Rectangle
}

func (xui *Xui) R() image.Rectangle {
	return xui.Rectangle
}

func (xui *Xui) Scale(n int) int {
	return n
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

func (xui *Xui) SetRoot(element.Interface) {
}

func (xui *Xui) Render() {
}

func (xui *Xui) Loop() {
}

type Info struct {
	Colors map[color.Color]int
	Bbox image.Rectangle
	BboxByColor map[color.Color]*image.Rectangle
}

func Analyze(img *memdraw.Image) (Info, error) {
	rgba, err := ToRGBA(img)
	if err != nil {
		return Info{}, fmt.Errorf("to rgba: %w", err)
	}

	colors := make(map[color.Color]int)
	black := color.RGBA{R:0, G:0, B:0, A:255}
	newRect := image.Rectangle{
		Min: image.Point{1e6, 1e6},
		Max: image.Point{-1e6, -1e6},
	}
	bbox := newRect
	bboxByColor := make(map[color.Color]*image.Rectangle)
	for x := img.R.Min.X; x < img.R.Max.X; x++ {
		for y := img.R.Min.Y; y < img.R.Max.Y; y++ {
			c := rgba.At(x, y)
			colors[c] = colors[c]+1
			if bboxByColor[c] == nil {
				r := newRect
				bboxByColor[c] = &r
			}
			if c == black {
				if bbox.Min.X > x {
					bbox.Min.X = x
				}
				if bbox.Max.X < x+1 {
					bbox.Max.X = x+1
				}
				if bbox.Min.Y > y {
					bbox.Min.Y = y
				}
				if bbox.Max.Y < y+1 {
					bbox.Max.Y = y+1
				}
			}
			if bboxByColor[c].Min.X > x {
				bboxByColor[c].Min.X = x
			}
			if bboxByColor[c].Max.X < x+1 {
				bboxByColor[c].Max.X = x+1
			}
			if bboxByColor[c].Min.Y > y {
				bboxByColor[c].Min.Y = y
			}
			if bboxByColor[c].Max.Y < y+1 {
				bboxByColor[c].Max.Y = y+1
			}
		}
	}

	info := Info{
		Colors: colors,
		Bbox: bbox,
		BboxByColor: bboxByColor,
	}

	return info, nil
}

func Dump(img *memdraw.Image) {
	rgba, err := ToRGBA(img)
	if err != nil {
		if err != nil { panic(err.Error()) }
	}

	f, err := os.Create("xui.png")
	if err != nil { panic(err.Error()) }
	defer f.Close()
	err = png.Encode(f, rgba)
	if err != nil { panic(err.Error()) }
}

func ToRGBA(img *memdraw.Image) (rgba *image.RGBA, err error) {
	rgba = image.NewRGBA/*32*/(img.R)
	rgba.Pix = make([]byte, 10*1024*1024)
	n, err := memdraw.Unload(img, img.R, rgba.Pix)
	if err != nil {
		return nil, fmt.Errorf("unload: %w", err)
	}
	rgba.Pix = rgba.Pix[:n]

	return
}
