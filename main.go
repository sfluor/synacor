package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sfluor/synacor/extractor"
	"github.com/sfluor/synacor/vm"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Please give the input binary as parameter: %v challenge.bin", os.Args[0])
		os.Exit(2)
	}

	// Coins solution
	// coins.printSolution()

	// Find R7 value
	// fmt.Println("Correct R7 value: ", vm.FindCorrectR7Value())

	// Read file
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	bin := extractor.Parse(string(b))

	// Extract code
	// extractCode(bin)

	// Initialize VM
	vm := vm.New(bin)

	// Run
	vm.Run()
}

func extractCode(bin []uint16) {
	f, err := os.Create("challenge.extracted")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	extractor.WriteExtractedCode(bin, f)
}
