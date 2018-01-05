package vm

import "fmt"
import "io"
import "log"

// formatRegister returns a string reprensentation of the current state of the register
func (vm VM) formatRegister() string {
	res := ""
	for i, v := range vm.register {
		res += fmt.Sprintf("R%d: %6d | ", i+1, v)
	}

	return res
}

// formatStack returns a string representation of the current stack
func (vm VM) formatStack() string {
	return fmt.Sprintf("%v", vm.stack)
}

// log writes the state of the vm in a writer
func (vm VM) log(w io.Writer) {

	_, err := w.Write([]byte(vm.formatRegister()))

	if err != nil {
		log.Fatalf("Could not log vm state: %s", err)
	}
}
