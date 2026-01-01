package logic

import (
	"strconv"
	"strings"
)

// ParseFloat handles both dot and comma decimal separators for global compatibility.
func ParseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}