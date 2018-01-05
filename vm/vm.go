// Package vm that describes the behavior of the VM
package vm

import (
	"bufio"
	"fmt"
	"os"
)

// M is the Mem size
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

// VM type
type VM struct {
	register [8]uint16
	stack    []uint16
	memory   []uint16
	cursor   uint16
}

// New creates a VM instance
func New(memory []uint16) *VM {
	return &VM{
		memory: memory,
	}
}

// Run executes the code in memory
func (vm *VM) Run() {
	// Reader for standard input
	stdinReader := bufio.NewReader(os.Stdin)

	// Log file
	// file, err := os.Create("./vm.log")
	// if err != nil {
	// 	panic(err)
	// }

	// Execute the binary
	for {
		vm.execInstruction(stdinReader)
		// vm.log(file)
	}
}

// execInstruction executes one instruction
func (vm *VM) execInstruction(reader *bufio.Reader) {
	// Our cursor that points to the actual position in the memory

	// Retrieve the operation
	op := vm.memory[vm.cursor]

	// To see what opcodes are called during the confirmation process
	// if vm.cursor > 5500 && vm.cursor < 7000 {
	// 	fmt.Println(vm.cursor, vm.register)
	// }

	// To print the stack
	// fmt.Println(vm.stack)

	switch op {
	case HALT: // Code 0
		fmt.Print("Halt op code !")
		os.Exit(0)

	case SET: // Code 1
		vm.set(vm.b())
		vm.cursor += 3

	case PUSH: // Code 2
		vm.stack = append(vm.stack, vm.a())
		vm.cursor += 2

	case POP: // Code 3
		popped, err := vm.pop()
		if err != nil {
			panic(err)
		}
		vm.set(popped)
		vm.cursor += 2

	case EQ: // Code 4
		if vm.b() == vm.c() {
			vm.set(1)
		} else {
			vm.set(0)
		}
		vm.cursor += 4

	case GT: // Code 5
		if vm.b() > vm.c() {
			vm.set(1)
		} else {
			vm.set(0)
		}
		vm.cursor += 4

	case JMP: // Code 6
		vm.cursor = vm.a()

	case JT: // Code 7
		if vm.a() != 0 {
			vm.cursor = vm.b()
		} else {
			vm.cursor += 3
		}

	case JF: // Code 8
		if vm.a() == 0 {
			vm.cursor = vm.b()
		} else {
			vm.cursor += 3
		}

	case ADD: // Code 9
		vm.set((vm.b() + vm.c()) % M)
		vm.cursor += 4

	case MULT: // Code 10
		vm.set((vm.b() * vm.c()) % M)
		vm.cursor += 4

	case MOD: // Code 11
		vm.set(vm.b() % vm.c())
		vm.cursor += 4

	case AND: // Code 12
		vm.set(vm.b() & vm.c())
		vm.cursor += 4

	case OR: // Code 13
		vm.set(vm.b() | vm.c())
		vm.cursor += 4

	case NOT: // Code 14
		vm.set(0x7fff &^ vm.b())
		vm.cursor += 3

	case RMEM: // Code 15
		vm.set(vm.get(vm.b()))
		vm.cursor += 3

	case WMEM: // Code 16
		vm.memory[vm.a()] = vm.b()
		vm.cursor += 3

	case CALL: // Code 17
		vm.push(vm.cursor + 2)
		vm.cursor = vm.a()

	case RET: // Code 18
		popped, err := vm.pop()
		if err != nil {
			// Halt
			fmt.Print("RET operation resulted in halt !")
			os.Exit(0)
		}
		vm.cursor = popped

	case OUT: // Code 19
		fmt.Print(string(vm.a()))
		vm.cursor += 2

	case IN: // Code 20
		b, _ := reader.ReadByte()
		if string(b) == "$" {
			fmt.Println(vm.formatRegister())
			// Force a non zero value for R7
			vm.register[7] = 7
		} else {
			vm.set(uint16(b))
			vm.cursor += 2
		}
	case NOOP: // Code 21
		vm.cursor++

	default:
		panic(fmt.Errorf("Unrecognized opcode %v", op))
	}
}

// get Retrieves a value by checking the register
func (vm VM) get(addr uint16) uint16 {
	m := vm.memory[addr]
	if m > M+7 {
		panic(fmt.Errorf("Get operation: Invalid address %v", m))
	}

	// Register
	if m >= M {
		return vm.register[m-M]
	}

	return m
}

// set Modify a value in the memory
func (vm *VM) set(value uint16) {
	// We always use set in the first argument < a >
	addr := vm.cursor + 1
	m := vm.memory[addr]
	if m > M+8 {
		panic(fmt.Errorf("Set operation: Invalid address %v", m))
	}

	// Set in register
	vm.register[m-M] = value
}

// Push to stack
func (vm *VM) push(value uint16) {
	vm.stack = append(vm.stack, value)
}

// Pop from stack
func (vm *VM) pop() (uint16, error) {
	if len(vm.stack) > 0 {
		res := vm.stack[len(vm.stack)-1]
		vm.stack = vm.stack[:len(vm.stack)-1]
		return res, nil
	}
	return 0, fmt.Errorf("empty stack ")
}

// a returns the first argument of the current command
func (vm VM) a() uint16 {
	return vm.get(vm.cursor + 1)
}

// b returns the second argument of the current command
func (vm VM) b() uint16 {
	return vm.get(vm.cursor + 2)
}

// c returns the first argument of the current command
func (vm VM) c() uint16 {
	return vm.get(vm.cursor + 3)
}
