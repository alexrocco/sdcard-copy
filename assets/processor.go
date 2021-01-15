package assets

import (
	"github.com/alexrocco/sdcard-copy/aws"
	"github.com/alexrocco/sdcard-copy/slice"
	"github.com/pkg/errors"
	"log"
	"path/filepath"
)

// Processor process assets
type Processor interface {
	// Process process the AssetProcess
	Process(asset Asset) error
}

// AssetProcessor process assets
type AssetProcessor struct {
	Finder Finder
	S3     aws.S3Api

	Log *log.Logger
}

func (a *AssetProcessor) Process(asset Asset) error {
	a.Log.Println("Processing", asset.Description)

	// Find all the assets in the sd card
	assetPaths, err := a.Finder.Find(asset.SdCardRegex)
	if err != nil {
		return errors.Wrap(err, "error when finding the assets in the sd card")
	}

	// Map to hold the file name and path to help in the upload part
	assetsSdCard := make(map[string]string)
	for _, jpgPath := range assetPaths {
		assetsSdCard[filepath.Base(jpgPath)] = jpgPath
	}

	// Find all the keys/assets in the S3 bucket
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
		a.Log.Printf("Found %d assets to upload", len(diffs))

		// Gets all the paths to upload
		var diffPathsUpload []string
		for _, diff := range diffs {
			diffPathsUpload = append(diffPathsUpload, assetsSdCard[diff])
		}

		// Batch upload all the diffs
		err = a.S3.BatchUpload(asset.S3BucketName, asset.S3BucketPrefix, diffPathsUpload, asset.Description)
		if err != nil {
			return errors.Wrap(err, "error when uploading all the assets to the AWS S3 bucket")
		}
	} else {
		a.Log.Println("No assets found to upload")
	}

	return nil
}
