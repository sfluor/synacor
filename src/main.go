package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"./extractor"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Please give the input binary as parameter: %v challenge.bin", os.Args[0])
		os.Exit(2)
	}

	// Coins solution
	// coins.printSolution()

	// Read file
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	bin := extractor.Parse(string(b))

	// Extract code
	extractCode(bin)

	// Initialize VM
	// vm := vm.New(bin)

	// Run
	// vm.Run()
}

func extractCode(bin []uint16) {
	f, err := os.Create("challenge.extracted")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	extractor.WriteExtractedCode(bin, f)
}
