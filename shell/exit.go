package shell

import (
	"os"
	"strconv"
	"fmt"
)

// Exit terminates the shell with an optional exit code.
func (sh *Shell) Exit(args []string) error {
	exitCode := 0
	if len(args) > 0 {
		var err error
		exitCode, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("exit: %s: numeric argument required\n", args[0])
			exitCode = 1
		}
	}
	os.Exit(exitCode)
	return nil
}