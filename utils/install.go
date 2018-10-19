package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/platform9/nodeadm/logs"

	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeletconfigv1beta1 "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig/v1beta1"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/deprecated"
	"github.com/platform9/nodeadm/systemd"
	netutil "k8s.io/apimachinery/pkg/util/net"
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
	placeKubeletSystemAndDropinFiles(config.Networking, config.Kubelet)
	placeNetworkConfig()
	if err := systemd.Enable("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.Start("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}

	if err := systemd.StopIfActive("keepalived.service"); err != nil {
		log.Fatalf("Failed to install keepalived service: %v", err)
	}
	if err := systemd.DisableIfEnabled("keepalived.service"); err != nil {
		log.Fatalf("Failed to install keepalived service: %v", err)
	}
	writeKeepAlivedServiceFiles(config)
	if err := systemd.Enable("keepalived.service"); err != nil {
		log.Fatalf("Failed to install keepalived service: %v", err)
	}
	if err := systemd.Start("keepalived.service"); err != nil {
		log.Fatalf("Failed to install keepalived service: %v", err)
	}
}

func InstallNodeComponents(config *apis.JoinConfiguration) {
	PopulateCache()
	placeKubeComponents()
	placeCNIPlugin()
	if err := systemd.StopIfActive("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.DisableIfEnabled("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	placeKubeletSystemAndDropinFiles(config.Networking, config.Kubelet)
	placeNetworkConfig()
	if err := systemd.Enable("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
	if err := systemd.Start("kubelet.service"); err != nil {
		log.Fatalf("Failed to install kubelet service: %v", err)
	}
}

func placeKubeletSystemAndDropinFiles(netConfig apis.Networking, kubeletConfig *kubeletconfigv1beta1.KubeletConfiguration) {
	placeAndModifyKubeletServiceFile()
	placeAndModifyKubeadmKubeletSystemdDropin()
	placeAndModifyNodeadmKubeletSystemdDropin(netConfig, kubeletConfig)
}

func placeAndModifyKubeletServiceFile() {
	serviceFile := filepath.Join(constants.SystemdDir, "kubelet.service")
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, "kubelet.service"), serviceFile)
	ReplaceString(serviceFile, "/usr/bin", constants.BaseInstallDir)
}

func placeAndModifyKubeadmKubeletSystemdDropin() {
	err := os.MkdirAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"), constants.Execute)
	if err != nil {
		log.Fatalf("\nFailed to create dir with error %v", err)
	}
	confFile := filepath.Join(constants.SystemdDir, "kubelet.service.d", constants.KubeadmKubeletSystemdDropinFilename)
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, constants.KubeadmKubeletSystemdDropinFilename), confFile)
	ReplaceString(confFile, "/usr/bin", constants.BaseInstallDir)
}

func placeAndModifyNodeadmKubeletSystemdDropin(netConfig apis.Networking, kubeletConfig *kubeletconfigv1beta1.KubeletConfiguration) {
	err := os.MkdirAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"), constants.Execute)
	if err != nil {
		log.Fatalf("\nFailed to create dir with error %v", err)
	}
	confFile := filepath.Join(constants.SystemdDir, "kubelet.service.d", constants.NodeadmKubeletSystemdDropinFilename)

	dnsIP, err := kubeadmconstants.GetDNSIP(netConfig.ServiceSubnet)
	if err != nil {
		log.Fatalf("Failed to derive DNS IP from service subnet %q: %v", netConfig.ServiceSubnet, err)
	}

	hostnameOverride, err := constants.GetHostnameOverride()
	if err != nil {
		log.Fatalf("Failed to dervice hostname override: %v", err)
	}

	data := struct {
		FailSwapOn       bool
		MaxPods          int32
		ClusterDNS       string
		ClusterDomain    string
		HostnameOverride string
		KubeAPIQPS       int32
		KubeAPIBurst     int32
		EvictionHard     string
		FeatureGates     string
	}{
		FailSwapOn:       *kubeletConfig.FailSwapOn,
		MaxPods:          kubeletConfig.MaxPods,
		ClusterDNS:       dnsIP.String(),
		ClusterDomain:    netConfig.DNSDomain,
		HostnameOverride: hostnameOverride,
		KubeAPIQPS:       *kubeletConfig.KubeAPIQPS,
		KubeAPIBurst:     kubeletConfig.KubeAPIBurst,
		EvictionHard:     constants.KubeletEvictionHard,
		FeatureGates:     constants.FeatureGates,
	}

	writeTemplateIntoFile(constants.NodeadmKubeletSystemdDropinTemplate, "nodeadm-kubelet-systemd-dropin", confFile, data)
}

func placeKubeComponents() {
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, "kubectl"), filepath.Join(constants.BaseInstallDir, "kubectl"))
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, "kubeadm"), filepath.Join(constants.BaseInstallDir, "kubeadm"))
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, "kubelet"), filepath.Join(constants.BaseInstallDir, "kubelet"))
}

func placeCNIPlugin() {
	tmpFile := fmt.Sprintf("cni-plugins-amd64-%s.tgz", constants.CNIVersion)
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.CNIDirName, tmpFile), filepath.Join("/tmp", tmpFile))
	if _, err := os.Stat(constants.CniVersionInstallDir); os.IsNotExist(err) {
		err := os.MkdirAll(constants.CniVersionInstallDir, constants.Execute)
		if err != nil {
			log.Fatalf("\nFailed to create dir %s with error %v", constants.CniVersionInstallDir, err)
		}
		deprecated.Run("", "tar", "-xvf", filepath.Join("/tmp", tmpFile), "-C", constants.CniVersionInstallDir)
		CreateSymLinks(constants.CniVersionInstallDir, constants.CNIBaseDir, true)
	}

}

func placeNetworkConfig() {
	os.MkdirAll(constants.ConfInstallDir, constants.Execute)
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.FlannelDirName, constants.FlannelManifestFilename), filepath.Join(constants.ConfInstallDir, constants.FlannelManifestFilename))
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

func writeKeepAlivedServiceFiles(config *apis.InitConfiguration) {
	log.Infof("\nVip configuration as parsed from the file %v", config)
	if len(config.VIPConfiguration.IP) == 0 {
		ip, err := netutil.ChooseHostInterface()
		if err != nil {
			log.Fatalf("Failed to get default interface with err %v", err)
		}
		config.VIPConfiguration.IP = ip.String()
	}

	if len(config.VIPConfiguration.NetworkInterface) == 0 {
		cmdStr := "route | grep '^default' | grep -o '[^ ]*$'"
		cmd := exec.Command("bash", "-c", cmdStr)
		bytes, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to get default interface with err %v", err)
		}
		config.VIPConfiguration.NetworkInterface = strings.Trim(string(bytes), "\n ")
	}

	if config.VIPConfiguration.RouterID == 0 {
		config.VIPConfiguration.RouterID = constants.DefaultRouterID
	}

	configTemplateVals := struct {
		InitConfig         *apis.InitConfiguration
		VRRPScriptInterval int
		VRRPScriptRise     int
		VRRPScriptFall     int
		WgetTimeout        int
	}{
		InitConfig:         config,
		VRRPScriptInterval: constants.VRRPScriptInterval,
		VRRPScriptRise:     constants.VRRPScriptRise,
		VRRPScriptFall:     constants.VRRPScriptFall,
		WgetTimeout:        constants.WgetTimeout,
	}
	kaConfFileTemplate := `global_defs {
	enable_script_security
}

vrrp_script chk_apiserver {
	script "/usr/bin/wget -T {{.WgetTimeout}} -qO /dev/null https://127.0.0.1:{{.InitConfig.MasterConfiguration.API.BindPort}}/healthz"
	interval {{.VRRPScriptInterval}}
	fall {{.VRRPScriptFall}}
	rise {{.VRRPScriptRise}}
}

vrrp_instance K8S_APISERVER {
	interface {{.InitConfig.VIPConfiguration.NetworkInterface}}
	state BACKUP
	virtual_router_id {{.InitConfig.VIPConfiguration.RouterID}}
	nopreempt
	authentication {
		auth_type AH
		auth_pass ourownpassword
	}
	virtual_ipaddress {
		{{.InitConfig.VIPConfiguration.IP}}
	}
	track_script {
		chk_apiserver
	}
}`
	confFile := filepath.Join(constants.SystemdDir, "keepalived.conf")
	writeTemplateIntoFile(kaConfFileTemplate, "vipConfFileTemplate", confFile, configTemplateVals)

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
Restart=on-failure
MemoryLow=10M
[Install]
WantedBy=multi-user.target
	`
	type KaServiceData struct {
		ConfigFile, KeepAlivedImg string
	}
	kaServiceData := KaServiceData{confFile, constants.KeepalivedImage}
	writeTemplateIntoFile(kaSvcFileTemplate, "kaSvcFileTemplate", filepath.Join(constants.SystemdDir, "keepalived.service"), kaServiceData)
}
