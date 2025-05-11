package label

import (
	"9fans.net/go/draw/memdraw"
	"image"
	"image/color"
	"testing"
	"xui/xuitest"
)

func TestRender(t *testing.T) {
	colors := testRender(t, "aaaa")
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	n4 := colors[black]
	colors = testRender(t, "aaaaaaaa")
	n8 := colors[black]
	if n8/n4 != 2 {
		t.Fail()
	}
}

func testRender(t *testing.T, text string) map[color.Color]int {
	memdraw.Init()
	l := New(image.Pt(0, 0), text)
	img := l.Render()

	info, err := xuitest.Analyze(img)
	if err != nil {
		panic(err.Error())
	}
	colors := info.Colors
	//bbox := info.Bbox

	return colors
}
