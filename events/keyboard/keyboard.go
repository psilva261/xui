package keyboard

type Type int

const (
	Down Type = 1 << iota
	Up
	Pressed
)

type Event struct {
	Type

	Key rune
}
