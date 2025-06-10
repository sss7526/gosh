package main

import (
	"fmt"
	"os"
	"gosh/shell"
)

func main() {
	sh := shell.NewShell()
	if err := sh.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Shell terminated with error: %v\n", err)
		os.Exit(1)
	}
	// shell.Start()
}
