package size

type SIZE_UNIT string

const (
	B SIZE_UNIT = "B"
	K SIZE_UNIT = "K"
	M SIZE_UNIT = "M"
	G SIZE_UNIT = "G"
	T SIZE_UNIT = "T"
	P SIZE_UNIT = "P"
	E SIZE_UNIT = "E"
)

var sizeUnitMap = map[string]SIZE_UNIT{
	"B": B, "K": K, "M": M, "G": G,
	"T": T, "P": P, "E": E,
}

func ParseSizeUnit(s string) (SIZE_UNIT, bool) {
	u, ok := sizeUnitMap[s]
	return u, ok
}

var unitMultipliers = map[SIZE_UNIT]int64{
	B: 1,
	K: 1 << 10,
	M: 1 << 20,
	G: 1 << 30,
	T: 1 << 40,
	P: 1 << 50,
	E: 1 << 60,
}

func UnitMultiplier(unit SIZE_UNIT) int64 {
	m, ok := unitMultipliers[unit]
	if !ok {
		return -1
	}
	return m
}
