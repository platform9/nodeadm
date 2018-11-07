package kubeadm

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// WriteConfiguration writes the kubeadm configuration to a file.
func WriteConfiguration(filename string, kubeadmConfiguration map[string]interface{}) error {
	y, err := yaml.Marshal(kubeadmConfiguration)
	if err != nil {
		return fmt.Errorf("unable to serialize kubeadm configuration: %s", err)
	}
	if err := ioutil.WriteFile(filename, y, 0600); err != nil {
		return fmt.Errorf("unable to write file: %s", err)
	}
	return nil
}
