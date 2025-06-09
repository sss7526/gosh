package shell

import (
	"errors"
	"gosh/cmd"
)

// Map to store built-in command with their corresponding handlers
var builtInCommands = map[string]func([]string) error{
	"exit": cmd.Exit,
	"cd": cmd.Cd,
	"pwd": cmd.Pwd,
}

func HandleBuiltInCommand(args []string) error {
	if handler, exists := builtInCommands[args[0]]; exists {
		return handler(args[1:])
	}
	return errors.New("not a built-in command")
}