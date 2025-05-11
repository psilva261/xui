package offset

type Opts int

const (
	Auto = 1 << (iota+1)
	Unset
)

type Of struct {
	Val int
	Opts
}
