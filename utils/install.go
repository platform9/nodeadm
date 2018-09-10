package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/deprecated"
	netutil "k8s.io/apimachinery/pkg/util/net"
)

func InstallMasterComponents(config *apis.InitConfiguration) {
	PopulateCache()
	PlaceComponentsFromCache(config.Networking)
	EnableAndStartService("kubelet.service")
	writeKeepAlivedServiceFiles(config)
	EnableAndStartService("keepalived.service")
}

func InstallNodeComponents(config *apis.JoinConfiguration) {
	PopulateCache()
	PlaceComponentsFromCache(config.Networking)
	EnableAndStartService("kubelet.service")
}

func PlaceComponentsFromCache(netConfig apis.Networking) {
	placeKubeComponents()
	placeCNIPlugin()
	placeAndModifyKubeletServiceFile()
	placeAndModifyKubeadmKubeletSystemdDropin(netConfig)
	placeAndModifyNodeadmKubeletSystemdDropin(netConfig)
	placeNetworkConfig()
}

func placeAndModifyKubeletServiceFile() {
	serviceFile := filepath.Join(constants.SystemdDir, "kubelet.service")
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, "kubelet.service"), serviceFile)
	ReplaceString(serviceFile, "/usr/bin", constants.BaseInstallDir)
}

func placeAndModifyKubeadmKubeletSystemdDropin(netConfig apis.Networking) {
	err := os.MkdirAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"), constants.Execute)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(constants.SystemdDir, "kubelet.service.d", constants.KubeadmKubeletSystemdDropinFilename)
	deprecated.Run("", "cp", filepath.Join(constants.CacheDir, constants.KubeDirName, constants.KubeadmKubeletSystemdDropinFilename), confFile)
	ReplaceString(confFile, "/usr/bin", constants.BaseInstallDir)
}

func placeAndModifyNodeadmKubeletSystemdDropin(netConfig apis.Networking) {
	err := os.MkdirAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"), constants.Execute)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
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
		MaxPods          int
		ClusterDNS       string
		ClusterDomain    string
		HostnameOverride string
		KubeAPIQPS       int
		KubeAPIBurst     int
	}{
		FailSwapOn:       constants.KubeletFailSwapOn,
		MaxPods:          constants.KubeletMaxPods,
		ClusterDNS:       dnsIP.String(),
		ClusterDomain:    netConfig.DNSDomain,
		HostnameOverride: hostnameOverride,
		KubeAPIQPS:       constants.KubeletKubeAPIQPS,
		KubeAPIBurst:     constants.KubeletKubeAPIBurst,
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
			log.Fatalf("Failed to create dir %s with error %v\n", constants.CniVersionInstallDir, err)
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
	//VIPConfig = config.VIPConfiguration
	//masterConf = config.MasterConfiguration
	log.Printf("Vip configuration as parsed from the file %v\n", config)
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
	script "/usr/bin/wget -T {{.WgetTimeout}} -qO - https://127.0.0.1:{{.InitConfig.MasterConfiguration.API.BindPort}}/healthz > /dev/null 2>&1"
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
[Install]
WantedBy=multi-user.target
	`
	type KaServiceData struct {
		ConfigFile, KeepAlivedImg string
	}
	kaServiceData := KaServiceData{confFile, constants.KeepalivedImage}
	writeTemplateIntoFile(kaSvcFileTemplate, "kaSvcFileTemplate", filepath.Join(constants.SystemdDir, "keepalived.service"), kaServiceData)
}
