package aws

//go:generate mockgen -source=./s3.go -package=mocks -destination=../mocks/mock_s3.go S3Api

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/schollz/progressbar/v3"
)

// S3Api describe the calls to AWS S3 APi
type S3Api interface {
	ListAllKeys(bucketName string, prefix string) ([]string, error)
	BatchUpload(bucketName string, prefix string, uploadPaths []string, description string) error
}

// S3 holds aws configs to use S3 api
type S3 struct {
	AwsRegion string
}

// ListAllKeys list all the keys in a prefix/path of a bucket
func (s *S3) ListAllKeys(bucketName string, prefix string) ([]string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.AwsRegion),
	}))

	s3Svc := s3.New(sess)

	// List all the objects from a prefix
	listObjOut, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return []string{}, err
	}

	var keys []string
	for _, content := range listObjOut.Contents {
		keys = append(keys, *content.Key)
	}

	// ListObjects can only return up to 1000 keys, so if next continuation token is set it means that more keys need to be fetched
	if listObjOut.NextContinuationToken != nil {
		nextContToken := listObjOut.NextContinuationToken

		// Keep fetching util nextContToken is not present
		for nextContToken != nil {
			listObjOut, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:            aws.String(bucketName),
				Prefix:            aws.String(prefix),
				ContinuationToken: nextContToken,
			})
			if err != nil {
				return []string{}, err
			}

			for _, content := range listObjOut.Contents {
				keys = append(keys, *content.Key)
			}

			// Update next cont token if the next "pagination"
			nextContToken = listObjOut.NextContinuationToken
		}
	}

	return keys, nil
}

// BatchUpload batch upload to AWS S3
func (s *S3) BatchUpload(bucketName string, prefix string, uploadPaths []string, description string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.AwsRegion),
	}))

	uploader := s3manager.NewUploader(sess)

	// Create the progress bar in the console
	bar := progressbar.NewOptions(len(uploadPaths),
		progressbar.OptionSetDescription(description))

	for _, upload := range uploadPaths {
		file, err := os.Open(upload)
		if err != nil {
			return errors.Wrapf(err, "can't open file %q", upload)
		}

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fmt.Sprintf("%s/%s", prefix, filepath.Base(file.Name()))),
			Body:   file,
		})
		if err != nil {
			return errors.Wrapf(err, "can't upload %q to AWS S3", upload)
		}

		err = file.Close()
		if err != nil {
			return errors.Wrapf(err, "can't close file %q", upload)
		}

		err = bar.Add(1)
		if err != nil {
			return errors.Wrap(err, "can't increase progress bar")
		}
	}

	return nil
}
