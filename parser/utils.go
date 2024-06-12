package parser

import (
	"golang.org/x/exp/slices"
)

func checkIfPortIsContained(port int, portSlice []int) bool {
	if slices.Contains(portSlice, port) {
		return true
	}

	return false
}
