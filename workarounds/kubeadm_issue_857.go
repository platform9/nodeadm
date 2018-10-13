package workarounds

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/platform9/nodeadm/constants"
)

const (
	// patchTemplate has one formatting parameter: the kube-proxy image name, a string
	patchTemplate = `[
    {
        "op": "add",
        "path": "/spec/template/spec/volumes/-",
        "value": {
                "name": "shared-data",
                "mountPath": "/shared-data"
            }
    },
    {
        "op": "add",
        "path": "/spec/template/spec/containers/0/volumeMounts/-",
        "value": {
                "name": "shared-data",
                "mountPath": "/shared-data"
            }
    },
    {
        "op": "replace",
        "path": "/spec/template/spec/containers/0/command/1",
        "value": "--config=/shared-data/config.conf"
    },
    {
        "op": "add",
        "path": "/spec/template/spec/initContainers",
        "value": [
                {
                    "command": [
                        "sh",
                        "-c",
                        "/bin/sed \"s/hostnameOverride: \\\"\\\"/hostnameOverride: $(NODE_NAME)/\" /var/lib/kube-proxy/config.conf > /shared-data/config.conf"
                    ],
                    "env": [
                        {
                            "name": "NODE_NAME",
                            "valueFrom": {
                                "fieldRef": {
                                    "apiVersion": "v1",
                                    "fieldPath": "spec.nodeName"
                                }
                            }
                        }
                    ],
                    "image": "%s",
                    "imagePullPolicy": "IfNotPresent",
                    "name": "update-config-file",
                    "volumeMounts": [
                    {
                        "mountPath": "/var/lib/kube-proxy",
                        "name": "kube-proxy"
                    },
                    {
                        "mountPath": "/shared-data",
                        "name": "shared-data"
                    }
                    ]
                }
            ]
    }
]`

	resultIfPatched = "--config=/shared-data/config.conf"
)

// EnsureKubeProxyRespectsHostoverride patches the kube-proxy daemonset so that
// kube-proxy respects the hostnameOverride setting. The function is idempotent.
// See: https://github.com/kubernetes/kubeadm/issues/857
func EnsureKubeProxyRespectsHostoverride() error {
	log.Infoln("[workarounds] Checking whether kube-proxy daemonset is patched")
	patched, err := isPatchedKubeProxyDaemonSet()
	if err != nil {
		return fmt.Errorf("unable to check if kube-proxy daemonset is patched: %v", err)
	}
	if patched {
		log.Infoln("[workarounds] Kube-proxy daemonset already patched. Continuing. ")
		return nil
	}
	log.Infoln("[workarounds] Patching kube-proxy daemonset")
	err = patchKubeProxyDaemonSet()
	if err != nil {
		return fmt.Errorf("unable to patch kube-proxy daemonset: %v", err)
	}
	log.Infoln("[workarounds] Patched kube-proxy daemonset")
	return nil
}

func isPatchedKubeProxyDaemonSet() (bool, error) {
	name := "/bin/sh"
	// If this field was changed, we assume the patch was applied. Because the
	// entire patch applies or is rejected, we can infer that the other fields
	// were updated as expected. We assume the daemonset was not edited by hand.
	arg := fmt.Sprintf("%s --kubeconfig=%s --namespace=kube-system get daemonset kube-proxy -ojsonpath='{.spec.template.spec.containers[0].command[1]}'", filepath.Join(constants.BaseInstallDir, constants.KubectlFilename), constants.AdminKubeconfigFile)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, "-c", arg)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("error running %q: %v (stdout: %s) (stderr: %s)", strings.Join(cmd.Args, " "), err, string(stdout.Bytes()), string(stderr.Bytes()))
	}

	if strings.Compare(string(stdout.Bytes()), resultIfPatched) == 0 {
		return true, nil
	}
	return false, nil
}

func patchKubeProxyDaemonSet() error {
	patchWithKubeProxyVersion := fmt.Sprintf(patchTemplate, fmt.Sprintf("k8s.gcr.io/kube-proxy-amd64:%s", constants.KubernetesVersion))
	name := "/bin/sh"
	arg := fmt.Sprintf("%s --kubeconfig=%s --namespace=kube-system patch --type=json daemonset kube-proxy --patch='%s'", filepath.Join(constants.BaseInstallDir, constants.KubectlFilename), constants.AdminKubeconfigFile, patchWithKubeProxyVersion)

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
