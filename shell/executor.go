package shell

// Executor will execute a command
type Executor interface {
	// Execute a command and return the stdout
	Execute(cmd string) (string, error)
}
