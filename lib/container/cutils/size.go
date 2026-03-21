package cutils

import (
	"fmt"
	"slices"
)

var SIZE_UNITS = []string{"K", "M", "G", "T", "P", "E"}

func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %sB", float64(size)/float64(div), SIZE_UNITS[exp])
}

func IsValidSizeUnit(str string) bool {
	return slices.Contains(SIZE_UNITS, str)
}

func unitMultiplier(unit string) int64 {
	switch unit {
	case "B":
		return 1
	case "K":
		return 1024
	case "M":
		return 1024 * 1024
	case "G":
		return 1024 * 1024 * 1024
	case "T":
		return 1024 * 1024 * 1024 * 1024
	case "P":
		return 1024 * 1024 * 1024 * 1024 * 1024
	case "E":
		return 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return -1
	}
}

func parseSize(input string) (float64, string, error) {
	var value float64
	var unit string

	// try number + unit (e.g. "10G", "1.5M")
	n, err := fmt.Sscanf(input, "%f%s", &value, &unit)
	if n == 0 || err != nil {
		return 0, "", fmt.Errorf("invalid size format: %s", input)
	}

	if unit == "" {
		unit = "B"
	}

	if unit != "B" && !IsValidSizeUnit(unit) {
		return 0, "", fmt.Errorf("invalid unit: %s", unit)
	}

	return value, unit, nil
}

func AnySizeTo(input string, target string) (int64, error) {
	value, unit, err := parseSize(input)
	if err != nil {
		return 0, err
	}

	fromMul := unitMultiplier(unit)
	toMul := unitMultiplier(target)

	if fromMul < 0 || toMul < 0 {
		return 0, fmt.Errorf("invalid unit conversion: %s → %s", unit, target)
	}

	// convert to bytes first
	bytes := value * float64(fromMul)

	// convert to target
	result := bytes / float64(toMul)

	return int64(result), nil
}
