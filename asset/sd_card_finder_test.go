package asset

import (
	"github.com/alexrocco/sdcard-copy/shell"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestSdCardFinder_Find(t *testing.T) {
	// Creates the temp dir for the DCIM folder to simulate a sd card dir structure
	tmpDir, err := ioutil.TempDir(os.TempDir(), "DCIM")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Creates the temp image
	tmpJpg, err := ioutil.TempFile(tmpDir, "DSC*.JPG")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpJpg.Name())

	sdCardFinder := SdCardFinder{
		Bash:        shell.Bash{},
		MountedPath: tmpDir,
	}

	got, err := sdCardFinder.Find("*.JPG")
	if err != nil {
		t.Errorf("no error should be found, but got %v", err)
		return
	}

	if len(got) > 1 {
		t.Errorf("it should find only one file")
	}

	if got[0] != tmpJpg.Name() {
		t.Errorf("got: %s, wanted: %s", got[0], tmpJpg.Name())
	}
}
