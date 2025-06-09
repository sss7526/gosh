package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"gosh/util"
	// "gosh/shell/execution"
)

func Start() {
	SetupSignalHandling(&fgPid)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(GetPrompt()) // Display prompt
		// Read the command from the user
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("exit") // Print "exit" like bash/zsh does
				os.Exit(0)
			}
			fmt.Println("Error reading input:", err)
			break
		}

		args := strings.Fields(input) // Tokenize input
		if len(args) == 0 {
			continue
		}

		for i, arg := range args {
			args[i] = util.ExpandHomeDirectory(arg)
		}

		// Handlin built-ins
		// if isBuiltInCommand(args) {
		// 	runBuiltInCommand(args)
		// } else {
		// 	executeCommandLowLevel(args, &fgPid)
		// }

		// Check if the command is built-in or external
		if err := HandleBuiltInCommand(args); err == nil {
			continue
		}

		// If not a built-in command, handle execution
		// execution.Execute(args)
		Execute(args, &fgPid)
	}
}