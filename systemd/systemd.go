/**
 *   Copyright 2018 Platform9 Systems, Inc.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package systemd

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

func Start(unit string) error {
	// Before we try to start any unit, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"start", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		return fmt.Errorf("failed to start unit: %v", err)
	}
	return nil
}

func Stop(unit string) error {
	// Before we try to start any unit, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"stop", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		return fmt.Errorf("failed to stop unit: %v", err)
	}
	return nil
}

func Enable(unit string) error {
	// Before we try to enable any unit, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"enable", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		return fmt.Errorf("failed to enable unit: %v", err)
	}
	return nil
}

func Disable(unit string) error {
	// Before we try to disable any unit, make sure that systemd is ready
	if err := reloadSystemd(); err != nil {
		return err
	}
	args := []string{"disable", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		return fmt.Errorf("failed to disable unit: %v", err)
	}
	return nil
}

// EnableAndStartUnit enables and starts the etcd unit
func EnableAndStartUnit(unit string) error {
	if err := Enable(unit); err != nil {
		return err
	}
	return Start(unit)
}

// DisableAndStopUnit disables and stops the etcd unit
func DisableAndStopUnit(unit string) error {
	if err := Disable(unit); err != nil {
		return err
	}
	return Stop(unit)
}

// Active checks if the systemd unit is active
func Active(unit string) (bool, error) {
	args := []string{"is-active", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		switch v := err.(type) {
		case *exec.Error:
			return false, fmt.Errorf("failed to run command %q: %s", v.Name, v.Err)
		case *exec.ExitError:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// Enabled checks if the systemd unit is enabled
func Enabled(unit string) (bool, error) {
	args := []string{"is-enabled", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		switch v := err.(type) {
		case *exec.Error:
			return false, fmt.Errorf("failed to run command %q: %s", v.Name, v.Err)
		case *exec.ExitError:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ResetFailed(unit string) error {
	args := []string{"reset-failed", unit}
	if err := exec.Command("systemctl", args...).Run(); err != nil {
		return fmt.Errorf("failed to reset failed unit: %v", err)
	}
	return nil
}

func Failed(unit string) (bool, error) {
	args := []string{"is-failed", unit}
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

// DisableIfEnabled disables a unit if it is enabled
func DisableIfEnabled(unit string) error {
	enabled, err := Enabled(unit)
	if err != nil {
		return fmt.Errorf("unable to check if unit %s is enabled: %v", unit, err)
	}
	if enabled {
		if err := Disable(unit); err != nil {
			return fmt.Errorf("unable to disable unit %s: %v", unit, err)
		}
	}
	return nil
}

// StopIfActive stops a unit if it is active
func StopIfActive(unit string) error {
	active, err := Active(unit)
	if err != nil {
		return fmt.Errorf("unable to check if unit %s is active: %v", unit, err)
	}
	if active {
		if err := Stop(unit); err != nil {
			return fmt.Errorf("unable to stop unit %s: %v", unit, err)
		}
	}
	return nil
}
