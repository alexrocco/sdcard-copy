package asset

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config hold asset configuration to process
type Config struct {
	Description    string `yaml:"description"`
	SdCardRegex    string `yaml:"sdCardRegex"`
	S3BucketName   string `yaml:"s3BucketName"`
	S3BucketPrefix string `yaml:"s3BucketPrefix"`
}

// LoadConfigs loads all the assets config from the yaml
func LoadConfigs() ([]Config, error) {
	content, err := ioutil.ReadFile("asset.yaml")
	if err != nil {
		return []Config{}, errors.Wrap(err, "error reading asset yaml")
	}

	var assets []Config
	err = yaml.Unmarshal(content, &assets)
	if err != nil {
		return []Config{}, errors.Wrap(err, "error parsing asset yaml")
	}

	return assets, nil
}
