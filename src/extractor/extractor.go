// Package extractor to handle binary extraction
package extractor

import (
	"fmt"
	"io"
	"strconv"

	"../vm"
)

// op type
type op struct {
	code  uint16 // Code of the operation
	name  string // Name of the operation
	nargs uint16 // Number of arguments
}

var operations map[uint16]op

// Parse Parses the binary as a string and return the list of 16-bits values respecting little-endian convention
func Parse(input string) []uint16 {
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

func init() {
	operations = map[uint16]op{
		0:  op{0, "halt", 0},
		1:  op{1, "set", 2},
		2:  op{2, "push", 1},
		3:  op{3, "pop", 1},
		4:  op{4, "eq", 3},
		5:  op{5, "gt", 3},
		6:  op{6, "jmp", 1},
		7:  op{7, "jt", 2},
		8:  op{8, "jf", 2},
		9:  op{9, "add", 3},
		10: op{10, "mult", 3},
		11: op{11, "mod", 3},
		12: op{12, "and", 3},
		13: op{13, "or", 3},
		14: op{14, "not", 2},
		15: op{15, "rmem", 2},
		16: op{16, "wmem", 2},
		17: op{17, "call", 1},
		18: op{18, "ret", 0},
		19: op{19, "out", 1},
		20: op{20, "in", 1},
		21: op{21, "noop", 0},
	}
}

// WriteExtractedCode writes the "readable" code to an io.Writer
func WriteExtractedCode(binary []uint16, w io.Writer) {
	for cursor := uint16(0); cursor < uint16(len(binary)); {
		op, ok := operations[binary[cursor]]
		if !ok {
			fmt.Printf("Invalid opcode: %v, %v\n", binary[cursor], binary[cursor-5:cursor+5])
			cursor++
		} else {
			row := fmt.Sprintf("(%6d) | %4s: %v", cursor, op.name, convert(binary[cursor+1:cursor+op.nargs+1]))

			if op.name == "out" {
				row += " " + string(binary[cursor+1])
			}

			w.Write([]byte(row + "\n"))

			cursor += op.nargs + 1
		}
	}
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

// transforms a value > M in it's register name
func convert(input []uint16) []string {
	res := []string{}

	for _, v := range input {
		if v >= vm.M {
			res = append(res, fmt.Sprintf("R%d", v-vm.M))
		} else {
			res = append(res, fmt.Sprintf("%d", v))
		}
	}

	return res
}
