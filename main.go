package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sfluor/synacor/coins"
	"github.com/sfluor/synacor/extractor"
	"github.com/sfluor/synacor/orb"
	"github.com/sfluor/synacor/vm"
)

func main() {

	coinsFlag := flag.Bool("coins", false, "Print the solution for the coin enigma")
	orbFlag := flag.Bool("orb", false, "Print the solution for the orb enigma")
	teleporter := flag.Bool("teleporter", false, "Print the solution for the teleporter enigma")
	file := flag.String("bin", "", "Path to the challenge.bin file")

	flag.Parse()

	if *coinsFlag {
		// Coins solution
		coins.PrintSolution()

	} else if *orbFlag {
		// Orb search
		orb.Search()

	} else if *teleporter {
		// Find R7 value
		fmt.Println("Correct R7 value: ", vm.FindCorrectR7Value())

	} else {
		// Read file
		b, err := ioutil.ReadFile(*file)
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

}

func extractCode(bin []uint16) {
	f, err := os.Create("challenge.extracted")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	extractor.WriteExtractedCode(bin, f)
}
