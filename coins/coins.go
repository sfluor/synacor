package coins

import (
	"fmt"
	"os"
	"strings"
)

// PrintSolution prints the solution for the coin enigma
func PrintSolution() {
	result := 399

	coins := map[int]string{
		2: "red coin",
		3: "corroded coin",
		5: "shiny coin",
		7: "concave coin",
		9: "blue coin",
	}

	// Ugly solution but the problem is small
	for a := range coins {
		for b := range coins {
			for c := range coins {
				for d := range coins {
					for e := range coins {
						if a+b*c*c+d*d*d-e == result &&
							!(a == b || a == c || a == d || a == e || b == c || b == d || b == e || c == d || c == e || d == e) {
							fmt.Println(strings.Join([]string{coins[a], coins[b], coins[c], coins[d], coins[e]}, " | "))
							// > blue coin | red coin | shiny coin | concave coin | corroded coin
							os.Exit(0)
						}
					}
				}
			}
		}
	}

}
