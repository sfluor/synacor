package vm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Return true if we should go to the next operation
func (vm *VM) debug(cmd string) bool {
	if strings.Contains(cmd, "register") {
		vm.printDebug("Register: " + vm.formatRegister() + "\n")
	}

	if strings.Contains(cmd, "stack") {
		vm.printDebug("Stack: " + vm.formatStack() + "\n")
	}

	if strings.Contains(cmd, "cursor") {
		vm.printDebug("Cursor: " + fmt.Sprintf("%d", vm.cursor) + "\n")
	}

	r := regexp.MustCompile(`^\$setreg R([1-8]) (0|[1-9][0-9]*)`)

	// Force a non zero value for R7
	if strings.Contains(cmd, "setreg") {
		match := r.FindStringSubmatch(cmd)
		// Wrong command
		if len(match) < 3 {
			vm.printError("Wrong command ! Should be $setreg R<n> <value>\n")
			return false
		}

		reg, err := strconv.ParseUint(match[1], 10, 16)
		if err != nil {
			vm.printError("Wrong register value\n")
			return false
		}

		val, err := strconv.ParseUint(string(match[2]), 10, 16)
		if err != nil {
			vm.printError("Wrong value for register\n")
			return false
		}

		vm.register[reg-1] = uint16(val)

	}

	if strings.Contains(cmd, "debugon") {
		vm.debugging = true
	}

	if strings.Contains(cmd, "debugoff") {
		vm.debugging = false
	}

	// Avance manually
	if strings.Contains(cmd, "steppingon") {
		vm.stepping = true
	}

	if strings.Contains(cmd, "steppingoff") {
		vm.stepping = false
	}

	if strings.Contains(cmd, "next") {
		return true
	}

	return false
}

func (vm VM) printDebug(str string) {
	// Print debug in light green
	fmt.Print("\033[32m", str, "\033[0m")
}

func (vm VM) printError(str string) {
	// Print error in red
	fmt.Print("\033[31m", str, "\033[0m")
}
