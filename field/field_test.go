package field

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"image"
	"testing"
	"xui/events/keyboard"
	"xui/xuitest"
)

func TestRender(t *testing.T) {
	memdraw.Init()
	rScreen := image.Rect(0, 0, 800, 600)
	rField := image.Rect(0, 0, 300, 100)
	img, err := memdraw.AllocImage(rScreen, draw.ABGR32)
	if err != nil {
		t.Fail()
	}
	memdraw.FillColor(img, draw.White)
	f := New(&xuitest.Xui{}, image.ZP, "", rField)

	for i := 0; i < 20; i++ {
		f.Event(keyboard.Event{Key: 'a'})
		img := f.Render()
		info, err := xuitest.Analyze(img)
		if err != nil {
			panic(err.Error())
		}
		t.Logf("info=%v", info.Bbox)
		if !info.Bbox.In(rField) {
			t.Fatalf("bbox outside field rect")
		}
	}
}