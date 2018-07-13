package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/platform9/nodeadm/apis"
)

func NodeadmConfigurationFromFile(path string) (*apis.NodeadmConfiguration, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	}
	config := apis.NodeadmConfiguration{}
	if err := yaml.Unmarshal(f, &config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}
	return &config, nil
}
