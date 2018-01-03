package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// Mem size
const M = 32768

// Op codes
const (
	HALT uint16 = iota
	SET
	PUSH
	POP
	EQ
	GT
	JMP
	JT
	JF
	ADD
	MULT
	MOD
	AND
	OR
	NOT
	RMEM
	WMEM
	CALL
	RET
	OUT
	IN
	NOOP
)

// Types
type register [8]uint16
type stack []uint16
type mem []uint16

func main() {
	// Check that file is submitted
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Please give the input binary as parameter: %v challenge.bin", os.Args[0])
		os.Exit(2)
	}

	// Read file
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Initialize memory
	reg := [8]uint16{0, 0, 0, 0, 0, 0, 0, 0}
	stack := []uint16{}
	mem := parse(string(b))

	// Execute
	exec(mem, reg, stack)
}

// exec executes the binary code
func exec(mem mem, reg register, stack stack) {
	// Our cursor that points to the actual position in the memory
	cursor := uint16(0)

	// Reader for opcode IN
	reader := bufio.NewReader(os.Stdin)

	for cursor < uint16(len(mem)) {
		// Retrieve the operation
		op := mem[cursor]
		switch op {
		case HALT: // Code 0
			cursor = uint16(len(mem))

		case SET: // Code 1
			reg[mem[cursor+1]-M] = g(mem[cursor+2], reg)
			cursor += 3

		case PUSH: // Code 2
			stack = append(stack, g(mem[cursor+1], reg))
			cursor += 2

		case POP: // Code 3
			if len(stack) == 0 {
				panic("EMPTY STACK")
			} else {
				reg[mem[cursor+1]-M] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			cursor += 2

		case EQ: // Code 4
			if g(mem[cursor+2], reg) == g(mem[cursor+3], reg) {
				reg[mem[cursor+1]-M] = 1
			} else {
				reg[mem[cursor+1]-M] = 0
			}
			cursor += 4

		case GT: // Code 5
			if g(mem[cursor+2], reg) > g(mem[cursor+3], reg) {
				reg[mem[cursor+1]-M] = 1
			} else {
				reg[mem[cursor+1]-M] = 0
			}
			cursor += 4

		case JMP: // Code 6
			cursor = g(mem[cursor+1], reg)

		case JT: // Code 7
			if g(mem[cursor+1], reg) != 0 {
				cursor = g(mem[cursor+2], reg)
			} else {
				cursor += 3
			}

		case JF: // Code 8
			if g(mem[cursor+1], reg) == 0 {
				cursor = g(mem[cursor+2], reg)
			} else {
				cursor += 3
			}

		case ADD: // Code 9
			reg[mem[cursor+1]-M] = (g(mem[cursor+2], reg) + g(mem[cursor+3], reg)) % M
			cursor += 4

		case MULT: // Code 10
			reg[mem[cursor+1]-M] = (g(mem[cursor+2], reg) * g(mem[cursor+3], reg)) % M
			cursor += 4

		case MOD: // Code 11
			reg[mem[cursor+1]-M] = (g(mem[cursor+2], reg) % g(mem[cursor+3], reg))
			cursor += 4

		case AND: // Code 12
			reg[mem[cursor+1]-M] = g(mem[cursor+2], reg) & g(mem[cursor+3], reg)
			cursor += 4

		case OR: // Code 13
			reg[mem[cursor+1]-M] = g(mem[cursor+2], reg) | g(mem[cursor+3], reg)
			cursor += 4

		case NOT: // Code 14
			reg[mem[cursor+1]-M] = 0x7fff &^ g(mem[cursor+2], reg)
			cursor += 3

		case RMEM: // Code 15
			reg[mem[cursor+1]-M] = mem[g(mem[cursor+2], reg)]
			cursor += 3

		case WMEM: // Code 16
			mem[g(mem[cursor+1], reg)] = g(mem[cursor+2], reg)
			cursor += 3

		case CALL: // Code 17
			stack = append(stack, cursor+2)
			cursor = g(mem[cursor+1], reg)

		case RET: // Code 18
			if len(stack) == 0 {
				cursor = uint16(len(mem))
			} else {
				cursor = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

		case OUT: // Code 19
			fmt.Print(string(g(mem[cursor+1], reg)))
			cursor += 2

		case IN: // Code 20
			b, _ := reader.ReadByte()
			reg[mem[cursor+1]-M] = uint16(b)
			cursor += 2

		case NOOP: // Code 21
			cursor++

		default:
			panic(fmt.Errorf("Unrecognized opcode %v", op))
		}
	}
}

// parse Parses the binary as a string and return the list of 16-bits values respecting little-endian convention
func parse(input string) mem {
	mem := []uint16{}

	for i := 0; i < len(input)-1; i += 2 {
		v, err := strconv.ParseUint(tob(input[i+1])+tob(input[i]), 2, 16)
		if err != nil {
			panic(err)
		}

		mem = append(mem, uint16(v))
	}
	return mem
}

// tob Converts to byte representation of size 8
func tob(c uint8) string {
	res := fmt.Sprintf("%b", c)
	s := len(res)
	for i := 0; i < 8-s; i++ {
		res = "0" + res
	}
	return res
}

// g Retrieves a value by checking the register
func g(nb uint16, reg register) uint16 {
	if nb > 32776 {
		panic(fmt.Errorf("Invalid number %v", nb))
	}

	// Register
	if nb >= M {
		return reg[nb-M]
	}

	return nb
}
