package shell

// Executer will execute a command
type Executer interface {
	// Execute a command and return the stdout
	Execute(cmd string) (string, error)
}
