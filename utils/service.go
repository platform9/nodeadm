package utils

import (
	"fmt"
	"os/exec"
)

func reloadSystemd() error {
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}
	return nil
}

func serviceStart(service string) error {
	// Before we try to start any service, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"start", service}
	return exec.Command("systemctl", args...).Run()
}

func serviceStop(service string) error {
	// Before we try to start any service, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"stop", service}
	return exec.Command("systemctl", args...).Run()
}

func serviceEnable(service string) error {
	// Before we try to enable any service, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"enable", service}
	return exec.Command("systemctl", args...).Run()
}

func serviceDisable(service string) error {
	// Before we try to enable any service, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"disable", service}
	return exec.Command("systemctl", args...).Run()
}

// EnableAndStartService enables and starts the etcd service
func EnableAndStartService(unitFile string) error {
	err := serviceEnable(unitFile)
	if err != nil {
		return err
	}
	err = serviceStart(unitFile)
	if err != nil {
		return err
	}
	return nil
}

// StopAndDisableService stops and disables service
func StopAndDisableService(unitFile string) error {
	err := serviceStop(unitFile)
	if err != nil {
		return err
	}
	err = serviceDisable(unitFile)
	if err != nil {
		return err
	}
	return nil
}
