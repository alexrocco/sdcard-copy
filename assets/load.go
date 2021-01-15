package assets

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Load loads all the assets to process
func Load() ([]Asset, error) {
	content, err := ioutil.ReadFile("assets.yaml")
	if err != nil {
		return []Asset{}, errors.Wrap(err, "error reading assets yaml")
	}

	var assets []Asset
	err = yaml.Unmarshal(content, &assets)
	if err != nil {
		return []Asset{}, errors.Wrap(err, "error parsing assets yaml")
	}

	return assets, nil
}
