package asset

import (
	"fmt"
	"github.com/alexrocco/sdcard-copy/shell"
	"strings"
)

//SdCardFinder finds asset in the sd card
type SdCardFinder struct {
	MountedPath string
	Bash        shell.Bash
}

func (sf *SdCardFinder) Find(regex string) ([]string, error) {
	assetsCmd := fmt.Sprintf("find %s -type f -wholename %q", sf.MountedPath, regex)
	output, err := sf.Bash.Execute(assetsCmd)
	if err != nil {
		return []string{}, err
	}

	output = strings.TrimSuffix(output, "\n")

	return strings.Split(output, "\n"), nil
}
