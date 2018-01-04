package main

import "fmt"
import "io"
import "log"

// formatRegister returns a string reprensentation of the current state of the register
func (vm vm) formatRegister() string {
	res := ""
	for i, v := range vm.register {
		res += fmt.Sprintf("R%d: %6d | ", i+1, v)
	}

	res += "\n"

	return res
}

// log writes the state of the vm in a writer
func (vm vm) log(w io.Writer) {

	_, err := w.Write([]byte(vm.formatRegister()))

	if err != nil {
		log.Fatalf("Could not log vm state: %s", err)
	}
}
