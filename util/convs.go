package util

import (
	"log/slog"
	"strconv"
)

func String2Float(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		slog.Error("error converting string to float", "error", err)
		return 0
	}
	return f
}
