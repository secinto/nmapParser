package parser

import (
	"golang.org/x/exp/slices"
	"strconv"
)

func checkIfPortIsContained(port string, portSlice []int) bool {
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		log.Errorf("Error converting port number to int: %s", err.Error())
		return false
	}
	if slices.Contains(portSlice, portNumber) {
		return true
	}

	return false
}

func checkIfPortIsNotContained(port Port, portSlice []Port) bool {
	for _, slice := range portSlice {
		if slice.Number == port.Number && slice.Protocol == port.Protocol {
			return false
		}
	}
	return true
}
