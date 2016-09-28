package Units

import "fmt"

//ByteSize based on: https://golang.org/doc/effective_go.html#constants
type ByteSize float64

const (
	_            = iota // ignore first value by assigning to blank identifier
	kiB ByteSize = 1 << (10 * iota)
	miB
	giB
	tiB
	piB
	eiB
	ziB
	yiB
)

func (b ByteSize) String() string {
	switch {
	case b >= yiB:
		return fmt.Sprintf("%.2fYiB", b/yiB)
	case b >= ziB:
		return fmt.Sprintf("%.2fZiB", b/ziB)
	case b >= eiB:
		return fmt.Sprintf("%.2fEiB", b/eiB)
	case b >= piB:
		return fmt.Sprintf("%.2fPiB", b/piB)
	case b >= tiB:
		return fmt.Sprintf("%.2fTiB", b/tiB)
	case b >= giB:
		return fmt.Sprintf("%.2fGiB", b/giB)
	case b >= miB:
		return fmt.Sprintf("%.2fMiB", b/miB)
	case b >= kiB:
		return fmt.Sprintf("%.2fKiB", b/kiB)
	}
	return fmt.Sprintf("%.2fB", b)
}
