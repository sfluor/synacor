package vm

import (
	"fmt"
)

// CachedConfirmation is the confirmation function using caching
func CachedConfirmation(R0, R1, R7 uint16, confirmationCache map[string]uint16) uint16 {

	formattedInput := fmt.Sprintf("%d-%d-%d", R0, R1, R7)

	// If already cached return it
	if v, cached := confirmationCache[formattedInput]; cached {
		return v
	}

	// Otherwise compute it and cache it
	var cR0 uint16

	if R0 == 0 {

		cR0 = (R1 + 1) % M
	} else if R1 == 0 {

		cR0 = CachedConfirmation((R0-1)%M, R7, R7, confirmationCache)
	} else {
		tempR1 := CachedConfirmation(R0, (R1-1)%M, R7, confirmationCache)

		cR0 = CachedConfirmation((R0-1)%M, tempR1, R7, confirmationCache)
	}

	confirmationCache[formattedInput] = cR0

	return cR0
}

// FindCorrectR7Value finds the correct R7 value
func FindCorrectR7Value() uint16 {
	answer := make(chan uint16, 1)
	found := make(chan bool)

	search := func(start uint16, end uint16) {
		for R7 := start; R7 < end; R7++ {
			select {
			case <-found:
				return
			default:
				if CachedConfirmation(uint16(4), uint16(1), R7, map[string]uint16{}) == 6 {
					answer <- R7
					close(found)
				}
			}
		}
	}

	// Let's spawn multiple workers
	for i := uint16(0); i < 4; i++ {
		go search(M/4*i, M/4*(i+1))
	}

	R7 := <-answer
	return R7
}
