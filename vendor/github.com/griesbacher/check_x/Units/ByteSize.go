package Units

import (
	"fmt"
)

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

func castTString(b ByteSize) string {
	return fmt.Sprintf("%0.3f", b)
}

func (b ByteSize) String() string {
	switch {
	case b >= yiB:
		return fmt.Sprintf("%sYiB", castTString(b/yiB))
	case b >= ziB:
		return fmt.Sprintf("%sZiB", castTString(b/ziB))
	case b >= eiB:
		return fmt.Sprintf("%sEiB", castTString(b/eiB))
	case b >= piB:
		return fmt.Sprintf("%sPiB", castTString(b/piB))
	case b >= tiB:
		return fmt.Sprintf("%sTiB", castTString(b/tiB))
	case b >= giB:
		return fmt.Sprintf("%sGiB", castTString(b/giB))
	case b >= miB:
		return fmt.Sprintf("%sMiB", castTString(b/miB))
	case b >= kiB:
		return fmt.Sprintf("%sKiB", castTString(b/kiB))
	}
	return fmt.Sprintf("%sB", castTString(b))
}
