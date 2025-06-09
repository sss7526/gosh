package cmd

import "os"

func Exit(args []string) error {
	os.Exit(0)
	return nil
}