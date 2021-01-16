package main

import (
	"strings"

	"github.com/alexrocco/sdcard-copy/asset"
	"github.com/alexrocco/sdcard-copy/aws"
	"github.com/alexrocco/sdcard-copy/shell"
)

func main() {
	log := NewLogger()

	log.Println("####################################")
	log.Println("### Starting SD Card copy script ###")
	log.Println("####################################")
	defer log.Println("SD Card copy script finished.")

	bash := shell.Bash{}

	// Finds the mounted path of the sd card by looking for a device with exfat type
	sdCardMountPath, err := bash.Execute("mount -t exfat | awk -F' ' '{print $3}'")
	if err != nil {
		log.Fatalf("Error when getting the sd card mounted path: %v", err)
	}

	if len(sdCardMountPath) == 0 {
		log.Fatalln("No SD card found, exiting...")
	}

	// Remove the break line if exists
	sdCardMountPath = strings.TrimSuffix(sdCardMountPath, "\n")

	sdCardFinder := asset.SdCardFinder{
		MountedPath: sdCardMountPath,
		Bash:        bash,
	}

	s3 := aws.S3LocalCred{
		AwsRegion: "us-east-2",
	}

	sdCardProcessor := asset.SdCardProcessor{
		Finder: &sdCardFinder,
		S3:     &s3,
		Log:    log,
	}

	assets, err := asset.LoadConfigs()
	if err != nil {
		log.Fatalf("Error loading asset: %v", err)
	}

	for _, a := range assets {
		err = sdCardProcessor.Process(a)
		if err != nil {
			log.Fatalf("Error processing %q: %v", a.Description, err)
		}
	}
}
