package anim

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"image"
	"time"
	"xui"
	"xui/events"
	"xui/internal/color"
	"xui/layout"
	"xui/space"
)

type Interface interface {
}

type Anim struct {
	Orig image.Point
	x xui.Interface
	images []*memdraw.Image
	current *memdraw.Image
	delay []int
	delaySum []int
	Margin space.Sp
}

func NewFromSeq(x xui.Interface, images []*memdraw.Image, delay []int) *Anim {
	a := &Anim{
		x: x,
		images: images,
		delay: delay,
	}
	a.delaySum = make([]int, len(a.delay))
	for i := 0; i < len(a.delay); i++ {
		a.delaySum[i] = a.delay[i]
		if i > 0 {
			a.delaySum[i] += a.delaySum[i-1]
		}
	}

	var err error
	a.current, err = memdraw.AllocImage(a.images[0].R, draw.ABGR32)
	if err != nil {
		panic(err.Error())
	}
	memdraw.FillColor(a.current, draw.Transparent)

	return a
}

func (a Anim) CurrentFrame(unixNano int64) int {
	m := a.delaySum[len(a.delay)-1]
	aTime := int((unixNano/1e7) % int64(m))
	for i, dc := range a.delaySum {
		if aTime <= dc {
			return i
		}
	}
	return 0
}

func (a Anim) Event(events.Interface) {
}

func (a Anim) Render() *memdraw.Image {
	t := time.Now().UnixNano()
	f := a.CurrentFrame(t)
	//orig := xy

	a.current.Draw(a.images[f].R, a.images[f], image.ZP, color.EmptyMask, image.ZP, draw.SoverD)
	return a.current
}

func (a Anim) Focus() {
}

func (a Anim) Layout() layout.Interface {
	return layout.Inline{}
}

func (a Anim) Geom() (r image.Rectangle, margin space.Sp) {
	return a.images[0].R, a.Margin
}

