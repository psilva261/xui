package color

import (
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"fmt"
)

var (
	EmptyMask, Border *memdraw.Image

	Hover Colors
)

type Colors struct {
	Text *memdraw.Image
	Background *memdraw.Image
	Border *memdraw.Image
}

type Colorset struct {
	Normal, Hover Colors
}

func init() {
	EmptyMask = MustMake(draw.Opaque)
	Border = MustMake(0xbbbbbbff)
	Hover.Border = MustMake(0xdc7232ff)
}

func MustMake(val draw.Color) (img *memdraw.Image) {
	img, err := Make(val)
	if err != nil {
		panic(err.Error())
	}
	return
}

func Make(val draw.Color) (img *memdraw.Image, err error) {
	img, err = memdraw.AllocImage(draw.Rect(0, 0, 1, 1), draw.ABGR32)
	if err != nil {
		return nil, fmt.Errorf("allocc image: %w", err)
	}
	img.Flags |= memdraw.Frepl
	img.Clipr = draw.Rect(-0x3FFFFFF, -0x3FFFFFF, 0x3FFFFFF, 0x3FFFFFF)
	memdraw.FillColor(img, val)
	return
}
