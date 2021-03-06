package asset

import (
	"github.com/alexrocco/sdcard-copy/aws"
	"github.com/alexrocco/sdcard-copy/slice"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

// Processor process asset
type Processor interface {
	// Process process the AssetProcess
	Process(asset Config) error
}

// SdCardProcessor process asset from the SD card
type SdCardProcessor struct {
	Finder Finder
	S3     aws.S3Api

	Log *logrus.Logger
}

func (a *SdCardProcessor) Process(asset Config) error {
	a.Log.Infof("Processing %s", asset.Description)

	a.Log.Infof("Finding the %s in the Sd Card", asset.Description)

	// Find all the asset in the sd card
	assetPaths, err := a.Finder.Find(asset.SdCardRegex)
	if err != nil {
		return errors.Wrap(err, "error when finding the asset in the sd card")
	}

	a.Log.Infof("Found %d asset(s) in the sd card", len(assetPaths))

	// Map to hold the file name and path to help in the upload part
	assetsSdCard := make(map[string]string)
	for _, jpgPath := range assetPaths {
		assetsSdCard[filepath.Base(jpgPath)] = jpgPath
	}

	// Find all the keys/asset in the S3 bucket
	s3Keys, err := a.S3.ListAllKeys(asset.S3BucketName, asset.S3BucketPrefix)
	if err != nil {
		return errors.Wrap(err, "error when listing all the AWS S3 keys in the bucket")
	}

	// Get the file name from all the S3 keys
	var assetsS3 []string
	for _, assetS3Key := range s3Keys {
		assetsS3 = append(assetsS3, filepath.Base(assetS3Key))
	}

	// Get the file name from the sd card
	var assetFiles []string
	for k := range assetsSdCard {
		assetFiles = append(assetFiles, k)
	}

	// Diff to find which files should be uploaded to S3
	diffs := slice.Diff(assetFiles, assetsS3)

	if len(diffs) > 0 {
		a.Log.Infof("Found %d asset(s) to upload", len(diffs))

		// Get all the paths to upload
		var diffPathsUpload []string
		for _, diff := range diffs {
			diffPathsUpload = append(diffPathsUpload, assetsSdCard[diff])
		}

		// Creates a progress bar to follow what is happening in the s3 upload
		bar := progressbar.Default(int64(len(diffs)))
		err := bar.RenderBlank()
		if err != nil {
			a.Log.Fatalf("Error rendering progress bar blank: %v", err)
		}

		bar.Describe(asset.Description)

		uploaded := make(chan string)
		go func() {
			for {
				select {
				case _, ok := <-uploaded:
					if ok {
						err := bar.Add(1)
						if err != nil {
							a.Log.Fatalf("Error incrementing progress bar: %v", err)
						}
					} else {
						// channel closed, returning function
						return
					}
				}
			}
		}()

		// Upload all the diffs
		err = a.S3.Upload(asset.S3BucketName, asset.S3BucketPrefix, diffPathsUpload, uploaded)
		if err != nil {
			return errors.Wrap(err, "error when uploading all the asset to the AWS S3 bucket")
		}

		// Finish the progress bar
		err = bar.Finish()
		if err != nil {
			a.Log.Fatalf("Error finishing progress bar: %v", err)
		}
	} else {
		a.Log.Info("No asset found to upload")
	}

	// Skip one line for the next iteration
	a.Log.Info()

	return nil
}
