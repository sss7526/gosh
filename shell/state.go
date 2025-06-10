package shell

import "sync"

// Job represents a process started by the shell, either in the foreground or background
type Job struct {
	ID 			int		// Unique Job ID
	Command 	string  // The command string being executed
	Pid 		int		// PID of the process
	Status 		string 	// Job status ("Running, "Stopped", "Completed")
	Background 	bool	// If the job is running in the background
}

// Shell encapsulates the state and behavior of the shell instance
type Shell struct {
	CurrentPwd 	string 				// Current working directory
	OldPwd		string 				// Prvious working directory
	Env 		map[string]string 	// Environment variables
	Jobs		map[int]*Job		// List of jobs keyed by JobID
	// lastJobID	int					// Generates unique Job IDs
	fgPid 		int 				// PID of the current foreground job
	LastStatus 	int 				// Exit status of the last executed command
	lock 		sync.Mutex 			// Mutex for thread-safe job and state operations
}