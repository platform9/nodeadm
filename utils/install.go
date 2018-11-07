package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Jeffail/gabs"
	log "github.com/platform9/nodeadm/pkg/logrus"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/systemd"
)

func InstallMasterComponents(config *apis.InitConfiguration) {
	PopulateCache()
	placeKubeComponents()
	placeCNIPlugin()
	if err := systemd.StopIfActive("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.DisableIfEnabled("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	placeAndModifyKubeletServiceFile()
	placeAndModifyKubeadmKubeletSystemdDropin()
	if err := systemd.Enable("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.Start("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if config.VIPConfiguration.IP != "" {
		if err := systemd.StopIfActive("keepalived.service"); err != nil {
			log.Fatalf("Failed to install keepalived service: %v", err)
		}
		if err := systemd.DisableIfEnabled("keepalived.service"); err != nil {
			log.Fatalf("Failed to install keepalived service: %v", err)
		}
		if err := writeKeepAlivedServiceFiles(config); err != nil {
			log.Fatalf("Failed to configure keepalived: %v", err)
		}
		if err := systemd.Enable("keepalived.service"); err != nil {
			log.Fatalf("Failed to install keepalived service: %v", err)
		}
		if err := systemd.Start("keepalived.service"); err != nil {
			log.Fatalf("Failed to install keepalived service: %v", err)
		}
	}
}

func InstallNodeComponents() {
	PopulateCache()
	placeKubeComponents()
	placeCNIPlugin()
	if err := systemd.StopIfActive("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.DisableIfEnabled("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	placeAndModifyKubeletServiceFile()
	placeAndModifyKubeadmKubeletSystemdDropin()
	if err := systemd.Enable("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.Start("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
}

func placeAndModifyKubeletServiceFile() {
	serviceFile := filepath.Join(constants.SystemdDir, "kubelet.service")
	_, err := copyFile(filepath.Join(constants.CacheDir, constants.KubeDirName, "kubelet.service"), serviceFile)
	checkError(err, "Unable to copy file")
	ReplaceString(serviceFile, "/usr/bin", constants.BaseInstallDir)
}

func placeAndModifyKubeadmKubeletSystemdDropin() {
	err := os.MkdirAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"), constants.Execute)
	if err != nil {
		log.Fatalf("\nFailed to create dir with error %v", err)
	}
	confFile := filepath.Join(constants.SystemdDir, "kubelet.service.d", constants.KubeadmKubeletSystemdDropinFilename)
	_, err = copyFile(filepath.Join(constants.CacheDir, constants.KubeDirName, constants.KubeadmKubeletSystemdDropinFilename), confFile)
	checkError(err, "Unable to copy file")
	ReplaceString(confFile, "/usr/bin", constants.BaseInstallDir)
}

func placeKubeComponents() {
	_, err := copyFile(filepath.Join(constants.CacheDir, constants.KubeDirName, "kubectl"), filepath.Join(constants.BaseInstallDir, "kubectl"))
	checkError(err, "Unable to copy file")
	_, err = copyFile(filepath.Join(constants.CacheDir, constants.KubeDirName, "kubeadm"), filepath.Join(constants.BaseInstallDir, "kubeadm"))
	checkError(err, "Unable to copy file")
	_, err = copyFile(filepath.Join(constants.CacheDir, constants.KubeDirName, "kubelet"), filepath.Join(constants.BaseInstallDir, "kubelet"))
	checkError(err, "Unable to copy file")
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func copyFile(src string, dst string) ([]byte, error) {
	cmd := exec.Command("cp", src, dst)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}
	return out, err
}

func placeCNIPlugin() {
	tmpFile := fmt.Sprintf("cni-plugins-amd64-%s.tgz", constants.CNIVersion)
	_, err := copyFile(filepath.Join(constants.CacheDir, constants.CNIDirName, tmpFile), filepath.Join("/tmp", tmpFile))
	checkError(err, "Unable to copy file")
	if _, err = os.Stat(constants.CniVersionInstallDir); os.IsNotExist(err) {
		err := os.MkdirAll(constants.CniVersionInstallDir, constants.Execute)
		if err != nil {
			log.Fatalf("\nFailed to create dir %s with error %v", constants.CniVersionInstallDir, err)
		}
		cmd := exec.Command("tar", "-xvf", filepath.Join("/tmp", tmpFile), "-C", constants.CniVersionInstallDir)
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to run %q: %s", strings.Join(cmd.Args, " "), err)
		}
		CreateSymLinks(constants.CniVersionInstallDir, constants.CNIBaseDir, true)
	}

}

func writeTemplateIntoFile(tmpl, name, file string, data interface{}) {
	err := os.MkdirAll(filepath.Dir(file), constants.Read)
	if err != nil {
		log.Fatalf("Failed to create dirs for path %s with error %v", filepath.Dir(file), err)
	}
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Failed to create file %q: %v", file, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	t := template.Must(template.New(name).Parse(tmpl))
	t.Execute(w, data)
	err = w.Flush()
	if err != nil {
		log.Fatalf("Failed to write to file %q: %v", file, err)
	}
}

func writeKeepAlivedServiceFiles(config *apis.InitConfiguration) error {
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse masterConfiguration: %s", err)
	}

	v := p.Path("api.bindPort").Data()
	if v == nil {
		return fmt.Errorf("masterConfiguration.api.bindPort is not defined. This is a bug, please file an issue on github.com/platform9/nodeadm")
	}
	apiBindPort, ok := v.(int)
	if !ok {
		return fmt.Errorf("unable to parse masterConfiguration.api.bindPort")
	}

	configTemplateVals := struct {
		VIPConfiguration   *apis.VIPConfiguration
		APIBindPort        int
		VRRPScriptInterval int
		VRRPScriptRise     int
		VRRPScriptFall     int
		WgetTimeout        int
	}{
		VIPConfiguration:   config.VIPConfiguration,
		APIBindPort:        apiBindPort,
		VRRPScriptInterval: constants.VRRPScriptInterval,
		VRRPScriptRise:     constants.VRRPScriptRise,
		VRRPScriptFall:     constants.VRRPScriptFall,
		WgetTimeout:        constants.WgetTimeout,
	}
	kaConfFileTemplate := `global_defs {
	enable_script_security
}

vrrp_script chk_apiserver {
	script "/usr/bin/wget -T {{.WgetTimeout}} -qO /dev/null https://127.0.0.1:{{.APIBindPort}}/healthz"
	interval {{.VRRPScriptInterval}}
	fall {{.VRRPScriptFall}}
	rise {{.VRRPScriptRise}}
}

vrrp_instance K8S_APISERVER {
	interface {{.VIPConfiguration.NetworkInterface}}
	state BACKUP
	virtual_router_id {{.VIPConfiguration.RouterID}}
	nopreempt
	virtual_ipaddress {
		{{.VIPConfiguration.IP}}
	}
	track_script {
		chk_apiserver
	}
}`
	writeTemplateIntoFile(kaConfFileTemplate, "vipConfFileTemplate", constants.KeepalivedConfigFilename, configTemplateVals)

	kaSvcFileTemplate := `
[Unit]
Description= Keepalived service
After=network.target docker.service
Requires=docker.service
[Service]
Type=simple
ExecStart=/usr/bin/docker run --cap-add=NET_ADMIN \
		--net=host --name vip \
		-v {{.ConfigFile}}:/usr/local/etc/keepalived/keepalived.conf \
		{{.KeepAlivedImg}}
ExecStartPre=-/usr/bin/docker kill vip
ExecStartPre=-/usr/bin/docker rm vip
ExecStop=/usr/bin/docker stop vip
ExecStopPost=/usr/bin/docker rm vip
Restart=on-failure
MemoryLow=10M
[Install]
WantedBy=multi-user.target
	`
	type KaServiceData struct {
		ConfigFile, KeepAlivedImg string
	}
	kaServiceData := KaServiceData{constants.KeepalivedConfigFilename, constants.KeepalivedImage}
	writeTemplateIntoFile(kaSvcFileTemplate, "kaSvcFileTemplate", filepath.Join(constants.SystemdDir, "keepalived.service"), kaServiceData)

	return nil
}
