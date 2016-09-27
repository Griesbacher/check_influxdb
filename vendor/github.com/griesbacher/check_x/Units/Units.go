package Units

import "fmt"

//ByteSize based on: https://golang.org/doc/effective_go.html#constants
type ByteSize float64

const (
	_ = iota // ignore first value by assigning to blank identifier
	KiB ByteSize = 1 << (10 * iota)
	MiB
	GiB
	TiB
	PiB
	EiB
	ZiB
	YiB
)

func (b ByteSize) String() string {
	switch {
	case b >= YiB:
		return fmt.Sprintf("%.2fYiB", b / YiB)
	case b >= ZiB:
		return fmt.Sprintf("%.2fZiB", b / ZiB)
	case b >= EiB:
		return fmt.Sprintf("%.2fEiB", b / EiB)
	case b >= PiB:
		return fmt.Sprintf("%.2fPiB", b / PiB)
	case b >= TiB:
		return fmt.Sprintf("%.2fTiB", b / TiB)
	case b >= GiB:
		return fmt.Sprintf("%.2fGiB", b / GiB)
	case b >= MiB:
		return fmt.Sprintf("%.2fMiB", b / MiB)
	case b >= KiB:
		return fmt.Sprintf("%.2fKiB", b / KiB)
	}
	return fmt.Sprintf("%.2fB", b)
}