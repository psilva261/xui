package geom

// Original code from github.com/mjl-/duit
//
// Copyright 2018 Mechiel Lukkien mechiel@ueber.net
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
// to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"image"

	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
)

// draw border with rounded corners, on the inside of `r`.
func DrawRoundedBorder(img *memdraw.Image, r image.Rectangle, color *memdraw.Image) {
	radius := 3
	x0 := r.Min.X
	x1 := r.Max.X - 1
	y0 := r.Min.Y
	y1 := r.Max.Y - 1
	tl := image.Pt(x0+radius, y0+radius)
	bl := image.Pt(x0+radius, y1-radius)
	br := image.Pt(x1-radius, y1-radius)
	tr := image.Pt(x1-radius, y0+radius)
	memdraw.Arc(img, tl, radius, radius, 0, color, image.ZP, 90, 90, draw.SoverD)
	memdraw.Arc(img, bl, radius, radius, 0, color, image.ZP, 180, 90, draw.SoverD)
	memdraw.Arc(img, br, radius, radius, 0, color, image.ZP, 270, 90, draw.SoverD)
	memdraw.Arc(img, tr, radius, radius, 0, color, image.ZP, 0, 90, draw.SoverD)
	memdraw.Line(img, image.Pt(x0, y0+radius), image.Pt(x0, y1-radius), 0, 0, 0, color, image.ZP, draw.SoverD)
	memdraw.Line(img, image.Pt(x0+radius, y1), image.Pt(x1-radius, y1), 0, 0, 0, color, image.ZP, draw.SoverD)
	memdraw.Line(img, image.Pt(x1, y1-radius), image.Pt(x1, y0+radius), 0, 0, 0, color, image.ZP, draw.SoverD)
	memdraw.Line(img, image.Pt(x1-radius, y0), image.Pt(x0+radius, y0), 0, 0, 0, color, image.ZP, draw.SoverD)
}

func DrawCursor(dst *memdraw.Image, bounds image.Rectangle, text *memdraw.Image, color *memdraw.Image) {
	h := text.R.Dy()/2
	p1 := text.R.Max.Add(image.Pt(h/5, -h/2))
	p0 := p1.Add(image.Pt(0, -h))
	if p0.X >= bounds.Dx() {
		return
	}
	memdraw.Line(dst, p0, p1, draw.EndSquare, draw.EndSquare, 1, color, image.ZP, draw.SoverD)
}

