package assets

type Asset struct {
	Description    string `yaml:"description"`
	SdCardRegex    string `yaml:"sdCardRegex"`
	S3BucketName   string `yaml:"s3BucketName"`
	S3BucketPrefix string `yaml:"s3BucketPrefix"`
}
