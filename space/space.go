package space

import (
	"image"
	"github.com/psilva261/xui/space/offset"
)

type Sp struct {
	Top, Right, Bottom, Left offset.Of
}

func New(vals... int) (s Sp) {
	switch len(vals) {
	case 4:
		s.Top = offset.Of{Val: vals[0]}
		s.Right = offset.Of{Val: vals[1]}
		s.Bottom = offset.Of{Val: vals[2]}
		s.Left = offset.Of{Val: vals[3]}
	case 3:
		s.Top = offset.Of{Val: vals[0]}
		s.Right = offset.Of{Val: vals[1]}
		s.Bottom = offset.Of{Val: vals[2]}
		s.Left = offset.Of{Val: vals[1]}
	case 2:
		s.Top = offset.Of{Val: vals[0]}
		s.Right = offset.Of{Val: vals[1]}
		s.Bottom = offset.Of{Val: vals[0]}
		s.Left = offset.Of{Val: vals[1]}
	case 1:
		s.Top = offset.Of{Val: vals[0]}
		s.Right = offset.Of{Val: vals[0]}
		s.Bottom = offset.Of{Val: vals[0]}
		s.Left = offset.Of{Val: vals[0]}
	}
	return
}

func (s Sp) TopLeft() image.Point {
	return image.Pt(s.Left.Val, s.Top.Val)
}

func (s Sp) Dx() int {
	return s.Left.Val+s.Right.Val
}

func (s Sp) Dy() int {
	return s.Top.Val+s.Bottom.Val
}

func (s Sp) Size() image.Point {
	return image.Pt(s.Dx(), s.Dy())
}
