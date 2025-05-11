package font

import (
	//"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"image"
	"image/color"
	"testing"
	"xui/xuitest"
)

func TestString(t *testing.T) {
	memdraw.Init()
	colors, dashBbox := testString(t, "----")
	t.Logf("colors=%+v", colors)
	t.Logf("dashBbox=%+v", dashBbox)
	if dashBbox.Dx() < 5 || dashBbox.Dy() < 1 {
		t.Fail()
	}
	colors, underBbox := testString(t, "____")
	t.Logf("colors=%+v", colors)
	t.Logf("underBbox=%+v", underBbox)
	if underBbox.Dx() < 5 || underBbox.Dy() < 1 {
		t.Fail()
	}
	if underBbox.Min.X >= dashBbox.Min.X {
		t.Fail()
	}
}

func TestStringEmpty(t *testing.T) {
	memdraw.Init()
	_, dashBbox := testString(t, "")
	t.Logf("dashBbox=%v", dashBbox)
}

func testString(t *testing.T, text string) (map[color.Color]int, image.Rectangle) {
	img, err := String(text)
	if err != nil {
		t.Fail()
	}

	info, err := xuitest.Analyze(img)
	if err != nil {
		panic(err.Error())
	}
	colors := info.Colors
	bbox := info.Bbox

	return colors, bbox
}