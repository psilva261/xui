package field

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"github.com/psilva261/xui/events/keyboard"
	"github.com/psilva261/xui/xuitest"
	"image"
	"testing"
)

var rField = image.Rect(0, 0, 300, 100)

func newField(t *testing.T) *Field {
	memdraw.Init()
	rScreen := image.Rect(0, 0, 800, 600)
	img, err := memdraw.AllocImage(rScreen, draw.ABGR32)
	if err != nil {
		t.Fail()
	}
	memdraw.FillColor(img, draw.White)
	return New(&xuitest.Xui{}, image.ZP, "", rField)
}

func TestEvent(t *testing.T) {
	f := newField(t)
	f.Text = "Text"
	f.Pos = 2
	f.Event(keyboard.Event{Key: Backspace})
	if f.Text != "Txt" {
		t.Fatalf("Text=%s", f.Text)
	}
	if f.Pos != 1 {
		t.Fatalf("Pos=%d", f.Pos)
	}
}

func TestRender(t *testing.T) {
	f := newField(t)
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
