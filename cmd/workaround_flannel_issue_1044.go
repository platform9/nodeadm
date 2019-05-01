package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/platform9/nodeadm/constants"
)

// ensureFlannelDaemonSetTolerations patches the flannel daemonset so that it
// tolerates all NoSchedule taints. See https://github.com/coreos/flannel/issues/1044
func ensureFlannelDaemonSetToleratesAllNoScheduleTaints() error {
	return patchFlannelDaemonSet()
}

// patchFlannelDaemonSet is idempotent; kubectl patch has a zero exit code if
// the patch has already been applied.
func patchFlannelDaemonSet() error {
	name := "/bin/sh"
	arg := fmt.Sprintf(`%s --kubeconfig=%s --namespace=kube-system patch daemonset kube-flannel-ds --patch='{"spec":{"template":{"spec":{"tolerations":[{"effect":"NoSchedule","operator":"Exists"}]}}}}'`, filepath.Join(constants.BaseInstallDir, constants.KubectlFilename), constants.AdminKubeconfigFile)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, "-c", arg)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %q: %v (stdout: %s) (stderr: %s)", strings.Join(cmd.Args, " "), err, string(stdout.Bytes()), string(stderr.Bytes()))
	}

	return nil
}
