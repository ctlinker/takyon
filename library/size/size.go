package size

import (
	"fmt"
	"os"
)

type Size struct {
	Value float64
	Unit  SIZE_UNIT
}

func (s Size) Bytes() float64 {
	return s.Value * float64(UnitMultiplier(s.Unit))
}

// Returns the value of s in the given unit
func (s Size) In(unit SIZE_UNIT) float64 {
	return s.Bytes() / float64(UnitMultiplier(unit))
}

// Returns a new Size converted to the given unit
func (s Size) ConvertedTo(U SIZE_UNIT) (Size, error) {
	byteUnit := UnitMultiplier(U)
	if byteUnit < 0 {
		return Size{}, fmt.Errorf("invalid conversion: %s → %s", s.Unit, U)
	}

	return Size{
		Value: s.In(U),
		Unit:  U,
	}, nil
}

func (s Size) String() string {
	if s.Value == float64(int64(s.Value)) {
		return fmt.Sprintf("%d%s", int64(s.Value), s.Unit)
	}
	return fmt.Sprintf("%.2f%s", s.Value, s.Unit)
}

func SizeOf(name string) (Size, error) {
	f, err := os.Stat(name)
	if err != nil {
		return Size{}, err
	}

	s := Size{
		Value: float64(f.Size()),
		Unit:  B,
	}

	return s, nil
}
