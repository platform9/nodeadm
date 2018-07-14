package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/platform9/nodeadm/apis"
)

func InitConfigurationFromFile(path string) (*apis.InitConfiguration, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	}
	config := apis.InitConfiguration{}
	if err := yaml.Unmarshal(f, &config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}
	return &config, nil
}

func JoinConfigurationFromFile(path string) (*apis.JoinConfiguration, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	}
	config := apis.JoinConfiguration{}
	if err := yaml.Unmarshal(f, &config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}
	return &config, nil
}
