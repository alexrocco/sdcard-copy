package shell

import "os/exec"

// Bash implements Executer and execute it using Bash.
type Bash struct {}

func (b Bash) Execute(cmd string) (string, error) {
	outCmd, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	return string(outCmd), nil
}
