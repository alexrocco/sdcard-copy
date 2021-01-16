package asset

import (
	"github.com/alexrocco/sdcard-copy/mock"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"log"
	"testing"
)

func TestAssetProcessor_Process(t *testing.T) {
	type args struct {
		asset Config
	}
	tests := []struct {
		name                 string
		args                 args
		foundAssets          []string
		s3Assets             []string
		wantedUploadedAssets []string
	}{
		{
			name: "It should upload the diff",
			args: args{
				asset: Config{
					Description:    "testDescription",
					SdCardRegex:    "testRegex",
					S3BucketName:   "testBucketName",
					S3BucketPrefix: "/JPG",
				},
			},
			foundAssets:          []string{"/test/123.JPG", "/test/345.JPG"},
			s3Assets:             []string{"/test/123.JPG"},
			wantedUploadedAssets: []string{"/test/345.JPG"},
		},
		{
			name: "It should not upload with no diff",
			args: args{
				asset: Config{
					Description:    "testDescription",
					SdCardRegex:    "testRegex",
					S3BucketName:   "testBucketName",
					S3BucketPrefix: "/JPG",
				},
			},
			foundAssets:          []string{"/test/123.JPG"},
			s3Assets:             []string{"/test/123.JPG"},
			wantedUploadedAssets: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFinder := mock.NewMockFinder(ctrl)
			mockFinder.EXPECT().Find(tt.args.asset.SdCardRegex).Return(tt.foundAssets, nil).Times(1)

			mockS3 := mock.NewMockS3Api(ctrl)
			mockS3.EXPECT().ListAllKeys(tt.args.asset.S3BucketName, tt.args.asset.S3BucketPrefix).Return(tt.s3Assets, nil).Times(1)

			timesCalledUploaded := 0
			if len(tt.wantedUploadedAssets) > 0 {
				timesCalledUploaded = 1
			}
			mockS3.EXPECT().Upload(
				tt.args.asset.S3BucketName,
				tt.args.asset.S3BucketPrefix,
				tt.wantedUploadedAssets,
				gomock.Any()).
				Return(nil).Times(timesCalledUploaded)

			sdCardProcessor := SdCardProcessor{
				Finder: mockFinder,
				S3:     mockS3,
				Log:    log.New(ioutil.Discard, "", log.Ltime),
			}

			gotErr := sdCardProcessor.Process(tt.args.asset)
			if gotErr != nil {
				t.Errorf("No error should be found, but got %v", gotErr)
			}

		})
	}
}
