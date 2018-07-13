package utils

import (
	"fmt"
	"os/exec"
	"strings"
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

func resetFailed(service string) error {
	args := []string{"reset-failed", service}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		return fmt.Errorf("failed to reset failed service: %v", err)
	}
	return nil
}

func isFailed(service string) (bool, error) {
	args := []string{"is-failed", service}
	out, err := exec.Command("systemctl", args...).CombinedOutput()

	if err != nil {
		switch err.(type) {
		// systemctl ran and exited
		case *exec.ExitError:
			return false, nil
		// exec encountered a different error
		default:
			return false, err
		}
	}
	// We could return true,nil here but let's check stdout just to make sure
	if strings.TrimSpace(string(out)) == "failed" {
		return true, nil
	}
	// We shouldn't hit this point but this is required.
	return false, nil
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

// ResetFailedService invokes reset-failed on a service if
// the service satisfies the is-failed query, otherwise a noop
func ResetFailedService(unitFile string) error {
	// Check to see if the service actually failed
	failed, err := isFailed(unitFile)
	if err != nil {
		return err
	}
	if failed {
		return resetFailed(unitFile)
	}
	return nil
}
