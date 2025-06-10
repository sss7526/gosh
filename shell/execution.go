package shell

import (
	"fmt"
	"syscall"
	"os"
	"gosh/util"
)

func (sh *Shell) Execute(args []string) {
	if len(args) == 0 {
		return
	}

	path, err := util.ResolveExecutable(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set up process attributes
	procAttr := &syscall.ProcAttr{
		Dir: "", // Working directory (empty = current)
		Env: os.Environ(),
		Files: []uintptr{
			uintptr(syscall.Stdin),
			uintptr(syscall.Stdout),
			uintptr(syscall.Stderr),
		},
	}

	// Fork and execute the process
	pid, err := syscall.ForkExec(path, args, procAttr)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	sh.fgPid = pid

	// Wait for child proc to complete (fg job)
	var status syscall.WaitStatus
	_, err = syscall.Wait4(pid, &status, 0, nil)
	
	sh.fgPid = 0
	if err != nil {
		fmt.Printf("Error waiting for process: %v\n", err)
		return
	}

	if status.Exited() {
		if status.ExitStatus() != 0 {
			fmt.Printf("Process exited with code %d\n", status.ExitStatus())
		}
	} else if status.Signaled() {
		fmt.Printf("Process killed by signal %d\n", status.Signal())
	}

}