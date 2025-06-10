package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"errors"
	"gosh/util"
)

func (sh *Shell) Start() error {
	
	sh.SetupSignalHandling()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Generate and display prompt
		fmt.Print(GetPrompt()) 

		// Read the command from the user
		input, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("exit") // Print "exit" like bash/zsh does
				return nil			// Graceful exit on EOF
			}
			return fmt.Errorf("error reading input: %w", err)
		}

		args := strings.Fields(input) // Tokenize input
		if len(args) == 0 {
			continue
		}

		for i, arg := range args {
			args[i] = util.ExpandHomeDirectory(arg)
		}

		// Check if the command is built-in or external
		if err = sh.HandleBuiltInCommand(args); err != nil {
			if errors.Is(err, ErrNotBuiltInCommand) {
				// Not a built-in, delegate to external command execution
				sh.Execute(args)
			} else {
				fmt.Fprintf(os.Stderr, "gosh: %s\n", err)
			}
			continue
		}
	}
}