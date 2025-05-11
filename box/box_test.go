package box

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"image"
	imagecolor "image/color"
	"testing"
	"xui/element"
	"xui/internal/color"
	"xui/label"
	"xui/space"
	"xui/xuitest"
)

func TestMarginPadding(t *testing.T) {
	memdraw.Init()

	inner := New([]element.Interface{})
	inner.Width = 130
	inner.Height = 70
	inner.Colorset.Normal.Background = color.MustMake(draw.Red)
	inner.Margin = space.New(7)

	t.Logf("inner.Margin=%v", inner.Margin)

	outer := New([]element.Interface{inner})
	outer.Colorset.Normal.Background = color.MustMake(draw.Green)
	outer.Margin = space.New(3)
	outer.Padding = space.New(11)
	img := outer.Render()
	t.Logf("TestMarginPadding: img.R=%v", img.R)
	info, err := xuitest.Analyze(img)
	if err != nil {
		panic(err.Error())
	}
	t.Logf("info=%v", info.BboxByColor)
	xuitest.Dump(img)
	bboxGreen := info.BboxByColor[imagecolor.RGBA{G:255,A:255}]
	if bboxGreen.Dx() != (130+2*7+2*11) || bboxGreen.Dy() != (70+2*7+2*11) {
		t.Fail()
	}
}

func TestRenderRect(t *testing.T) {
	memdraw.Init()

	b := New([]element.Interface{})
	b.Width = 130
	b.Height = 70
	b.Colorset.Normal.Background = color.MustMake(draw.Black)
	img := b.Render()
	info, err := xuitest.Analyze(img)
	if err != nil {
		panic(err.Error())
	}
	t.Logf("info=%v", info.Bbox)
	if info.Bbox.Dx() != 130 || info.Bbox.Dy() != 70 {
		t.Fail()
	}
}

func TestRenderLabel(t *testing.T) {
	memdraw.Init()

	l := label.New(image.ZP, "hello")

	imgTestLabel := l.Render()
	infoTestLabel, err := xuitest.Analyze(imgTestLabel)
	if err != nil {
		panic(err.Error())
	}
	t.Logf("infoTestLabel=%v", infoTestLabel.Bbox)

	rScreen := image.Rect(0, 0, 800, 600)

	b := New([]element.Interface{l})

	img := b.Render()
	info, err := xuitest.Analyze(img)
	if err != nil {
		panic(err.Error())
	}
	t.Logf("info=%v", info.Bbox)
	if !info.Bbox.In(rScreen) {
		t.Fatalf("bbox outside screen rect")
	}

	if !infoTestLabel.Bbox.Eq(info.Bbox) {
		t.Fail()
	}
}

