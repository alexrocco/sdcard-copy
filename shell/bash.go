package shell

import (
	"github.com/pkg/errors"
	"os/exec"
)

// Bash implements Executor and execute it using Bash.
type Bash struct {}

func (b Bash) Execute(cmd string) (string, error) {
	outCmd, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", errors.Wrapf(err, "error executing bash command %s", cmd)
	}

	return string(outCmd), nil
}
