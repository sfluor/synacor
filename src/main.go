package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Please give the input binary as parameter: %v challenge.bin", os.Args[0])
		os.Exit(2)
	}

	// Read file
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Initialize VM
	vm := &vm{memory: parse(string(b))}

	// Run
	vm.run()
}

func (vm *vm) run() {
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
