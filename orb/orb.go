package orb

import (
	"fmt"
)

type room struct {
	operation string
	number    int
	x, y      int
}

// |   *   |   8   |   -   |   1   |
// |   4   |   *   |   11  |   *   |
// |   +   |   4   |   -   |   18  |
// |   22  |   -   |   9   |   *   |
var vault = [][]room{
	[]room{
		room{operation: "*", x: 0, y: 0},
		room{number: 8, x: 1, y: 0},
		room{operation: "-", x: 2, y: 0},
		room{number: 1, x: 3, y: 0},
	},
	[]room{
		room{number: 4, x: 0, y: 1},
		room{operation: "*", x: 1, y: 1},
		room{number: 11, x: 2, y: 1},
		room{operation: "*", x: 3, y: 1},
	},
	[]room{
		room{operation: "+", x: 0, y: 2},
		room{number: 4, x: 1, y: 2},
		room{operation: "-", x: 2, y: 2},
		room{number: 18, x: 3, y: 2},
	},
	[]room{
		room{number: 22, x: 0, y: 3},
		room{operation: "-", x: 1, y: 3},
		room{number: 9, x: 2, y: 3},
		room{operation: "*", x: 3, y: 3},
	},
}

var operations = map[string]func(a, b int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"":  func(a, b int) int { return a },
}

type state struct {
	x, y, orb int
	history   []string
	op        string
}

func Search() {
	queue := []state{
		state{
			0, 3, 22, []string{}, "",
		},
	}

	for len(queue) != 0 {
		state := queue[0]
		queue = queue[1:]
		// Vault Door
		if state.x == 3 && state.y == 0 {
			if state.orb == 30 {
				fmt.Println(state)
				return
			}
		} else if state.orb > 0 && state.orb < 100 {
			queue = append(queue, getNextStates(state)...)
		}
	}
}

func getNextStates(previousState state) []state {
	newStates := []state{}

	x := previousState.x
	y := previousState.y

	if x > 1 || (x == 1 && y != 3) {
		newStates = append(newStates, newState(previousState, x-1, y, "west"))
	}
	if x < 3 {
		newStates = append(newStates, newState(previousState, x+1, y, "east"))
	}
	if y > 0 {
		newStates = append(newStates, newState(previousState, x, y-1, "north"))
	}
	if y < 2 || (y == 2 && x != 0) {
		newStates = append(newStates, newState(previousState, x, y+1, "south"))
	}
	return newStates
}

func newState(previousState state, x, y int, dir string) state {
	// Copy otherwise weird behaviors happens
	history := make([]string, len(previousState.history))
	copy(history, previousState.history)

	return state{
		x:       x,
		y:       y,
		orb:     operations[previousState.op](previousState.orb, vault[y][x].number),
		history: append(history, dir),
		op:      vault[y][x].operation,
	}
}
